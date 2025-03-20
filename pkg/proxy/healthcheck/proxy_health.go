/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package healthcheck

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/proxy/metrics"
	"k8s.io/utils/clock"
)

const (
	// ToBeDeletedTaint is a taint used by the CLuster Autoscaler before marking a node for deletion. Defined in
	// https://github.com/kubernetes/autoscaler/blob/e80ab518340f88f364fe3ef063f8303755125971/cluster-autoscaler/utils/deletetaint/delete.go#L36
	ToBeDeletedTaint = "ToBeDeletedByClusterAutoscaler"
)

// ProxyHealthServer allows callers to:
//  1. run a http server with /healthz and /livez endpoint handlers.
//  2. update healthz timestamps before and after synchronizing dataplane.
//  3. sync node status, for reporting unhealthy /healthz response
//     if the node is marked for deletion by autoscaler.
//  4. get proxy health by verifying that the delay between QueuedUpdate()
//     calls and Updated() calls exceeded healthTimeout or not.
type ProxyHealthServer struct {
	listener    listener
	httpFactory httpServerFactory
	clock       clock.Clock

	addr          string
	healthTimeout time.Duration

	lock                   sync.RWMutex
	lastUpdatedMap         map[v1.IPFamily]time.Time
	oldestPendingQueuedMap map[v1.IPFamily]time.Time
	nodeEligible           bool
}

// NewProxyHealthServer returns a proxy health http server.
func NewProxyHealthServer(addr string, healthTimeout time.Duration) *ProxyHealthServer {
	return newProxyHealthServer(stdNetListener{}, stdHTTPServerFactory{}, clock.RealClock{}, addr, healthTimeout)
}

func newProxyHealthServer(listener listener, httpServerFactory httpServerFactory, c clock.Clock, addr string, healthTimeout time.Duration) *ProxyHealthServer {
	return &ProxyHealthServer{
		listener:      listener,
		httpFactory:   httpServerFactory,
		clock:         c,
		addr:          addr,
		healthTimeout: healthTimeout,

		lastUpdatedMap:         make(map[v1.IPFamily]time.Time),
		oldestPendingQueuedMap: make(map[v1.IPFamily]time.Time),
		// The node is eligible (and thus the proxy healthy) while it's starting up
		// and until we've processed the first node event that indicates the
		// contrary.
		nodeEligible: true,
	}
}

// Updated should be called when the proxier of the given IP family has successfully updated
// the service rules to reflect the current state and should be considered healthy now.
func (hs *ProxyHealthServer) Updated(ipFamily v1.IPFamily) {
	hs.lock.Lock()
	defer hs.lock.Unlock()
	delete(hs.oldestPendingQueuedMap, ipFamily)
	hs.lastUpdatedMap[ipFamily] = hs.clock.Now()
}

// QueuedUpdate should be called when the proxier receives a Service or Endpoints event
// from API Server containing information that requires updating service rules. It
// indicates that the proxier for the given IP family has received changes but has not
// yet pushed them to its backend. If the proxier does not call Updated within the
// healthTimeout time then it will be considered unhealthy.
func (hs *ProxyHealthServer) QueuedUpdate(ipFamily v1.IPFamily) {
	hs.lock.Lock()
	defer hs.lock.Unlock()
	// Set oldestPendingQueuedMap[ipFamily] only if it's currently unset
	if _, set := hs.oldestPendingQueuedMap[ipFamily]; !set {
		hs.oldestPendingQueuedMap[ipFamily] = hs.clock.Now()
	}
}

// Health returns only the proxier's health state and last updated time.
func (hs *ProxyHealthServer) Health() (bool, time.Time) {
	hs.lock.RLock()
	defer hs.lock.RUnlock()

	var lastUpdated time.Time
	for _, proxierLastUpdated := range hs.lastUpdatedMap {
		if proxierLastUpdated.After(lastUpdated) {
			lastUpdated = proxierLastUpdated
		}
	}

	currentTime := hs.clock.Now()
	for ipFamily, _ := range hs.lastUpdatedMap {
		if _, set := hs.oldestPendingQueuedMap[ipFamily]; !set {
			// the proxier is healthy while it's starting up
			// or the proxier is fully synced.
			continue
		}

		if currentTime.Sub(hs.oldestPendingQueuedMap[ipFamily]) < hs.healthTimeout {
			// there's an unprocessed update queued for this proxier, but it's not late yet.
			continue
		}
		return false, lastUpdated
	}
	return true, lastUpdated
}

// SyncNode syncs the node and determines if it is eligible or not. Eligible is
// defined as being: not tainted by ToBeDeletedTaint and not deleted.
func (hs *ProxyHealthServer) SyncNode(node *v1.Node) {
	hs.lock.Lock()
	defer hs.lock.Unlock()

	if !node.DeletionTimestamp.IsZero() {
		hs.nodeEligible = false
		return
	}
	for _, taint := range node.Spec.Taints {
		if taint.Key == ToBeDeletedTaint {
			hs.nodeEligible = false
			return
		}
	}
	hs.nodeEligible = true
}

// NodeEligible returns nodeEligible field of ProxyHealthServer.
func (hs *ProxyHealthServer) NodeEligible() bool {
	hs.lock.RLock()
	defer hs.lock.RUnlock()
	return hs.nodeEligible
}

// Run starts the healthz HTTP server and blocks until it exits.
func (hs *ProxyHealthServer) Run(ctx context.Context) error {
	serveMux := http.NewServeMux()
	serveMux.Handle("/healthz", healthzHandler{hs: hs})
	serveMux.Handle("/livez", livezHandler{hs: hs})
	server := hs.httpFactory.New(serveMux)

	listener, err := hs.listener.Listen(ctx, hs.addr)
	if err != nil {
		return fmt.Errorf("failed to start proxy healthz on %s: %w", hs.addr, err)
	}

	klog.V(3).InfoS("Starting healthz HTTP server", "address", hs.addr)

	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("proxy healthz closed with error: %w", err)
	}
	return nil
}

type healthzHandler struct {
	hs *ProxyHealthServer
}

func (h healthzHandler) ServeHTTP(resp http.ResponseWriter, _ *http.Request) {
	nodeEligible := h.hs.NodeEligible()
	healthy, lastUpdated := h.hs.Health()
	currentTime := h.hs.clock.Now()

	healthy = healthy && nodeEligible
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("X-Content-Type-Options", "nosniff")
	if !healthy {
		metrics.ProxyHealthzTotal.WithLabelValues("503").Inc()
		resp.WriteHeader(http.StatusServiceUnavailable)
	} else {
		metrics.ProxyHealthzTotal.WithLabelValues("200").Inc()
		resp.WriteHeader(http.StatusOK)
		// In older releases, the returned "lastUpdated" time indicated the last
		// time the proxier sync loop ran, even if nothing had changed. To
		// preserve compatibility, we use the same semantics: the returned
		// lastUpdated value is "recent" if the server is healthy. The kube-proxy
		// metrics provide more detailed information.
		lastUpdated = currentTime
	}
	fmt.Fprintf(resp, `{"lastUpdated": %q,"currentTime": %q, "nodeEligible": %v}`, lastUpdated, currentTime, nodeEligible)
}

type livezHandler struct {
	hs *ProxyHealthServer
}

func (h livezHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	healthy, lastUpdated := h.hs.Health()
	currentTime := h.hs.clock.Now()
	resp.Header().Set("Content-Type", "application/json")
	resp.Header().Set("X-Content-Type-Options", "nosniff")
	if !healthy {
		metrics.ProxyLivezTotal.WithLabelValues("503").Inc()
		resp.WriteHeader(http.StatusServiceUnavailable)
	} else {
		metrics.ProxyLivezTotal.WithLabelValues("200").Inc()
		resp.WriteHeader(http.StatusOK)
		// In older releases, the returned "lastUpdated" time indicated the last
		// time the proxier sync loop ran, even if nothing had changed. To
		// preserve compatibility, we use the same semantics: the returned
		// lastUpdated value is "recent" if the server is healthy. The kube-proxy
		// metrics provide more detailed information.
		lastUpdated = currentTime
	}
	fmt.Fprintf(resp, `{"lastUpdated": %q,"currentTime": %q}`, lastUpdated, currentTime)
}
