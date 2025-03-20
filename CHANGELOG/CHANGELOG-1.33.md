<!-- BEGIN MUNGE: GENERATED_TOC -->

- [v1.33.0-beta.0](#v1330-beta0)
  - [Downloads for v1.33.0-beta.0](#downloads-for-v1330-beta0)
    - [Source Code](#source-code)
    - [Client Binaries](#client-binaries)
    - [Server Binaries](#server-binaries)
    - [Node Binaries](#node-binaries)
    - [Container Images](#container-images)
  - [Changelog since v1.33.0-alpha.3](#changelog-since-v1330-alpha3)
  - [Changes by Kind](#changes-by-kind)
    - [API Change](#api-change)
    - [Feature](#feature)
    - [Bug or Regression](#bug-or-regression)
    - [Other (Cleanup or Flake)](#other-cleanup-or-flake)
  - [Dependencies](#dependencies)
    - [Added](#added)
    - [Changed](#changed)
    - [Removed](#removed)
- [v1.33.0-alpha.3](#v1330-alpha3)
  - [Downloads for v1.33.0-alpha.3](#downloads-for-v1330-alpha3)
    - [Source Code](#source-code-1)
    - [Client Binaries](#client-binaries-1)
    - [Server Binaries](#server-binaries-1)
    - [Node Binaries](#node-binaries-1)
    - [Container Images](#container-images-1)
  - [Changelog since v1.33.0-alpha.2](#changelog-since-v1330-alpha2)
  - [Urgent Upgrade Notes](#urgent-upgrade-notes)
    - [(No, really, you MUST read this before you upgrade)](#no-really-you-must-read-this-before-you-upgrade)
  - [Changes by Kind](#changes-by-kind-1)
    - [Deprecation](#deprecation)
    - [API Change](#api-change-1)
    - [Feature](#feature-1)
    - [Bug or Regression](#bug-or-regression-1)
    - [Other (Cleanup or Flake)](#other-cleanup-or-flake-1)
  - [Dependencies](#dependencies-1)
    - [Added](#added-1)
    - [Changed](#changed-1)
    - [Removed](#removed-1)
- [v1.33.0-alpha.2](#v1330-alpha2)
  - [Downloads for v1.33.0-alpha.2](#downloads-for-v1330-alpha2)
    - [Source Code](#source-code-2)
    - [Client Binaries](#client-binaries-2)
    - [Server Binaries](#server-binaries-2)
    - [Node Binaries](#node-binaries-2)
    - [Container Images](#container-images-2)
  - [Changelog since v1.33.0-alpha.1](#changelog-since-v1330-alpha1)
  - [Changes by Kind](#changes-by-kind-2)
    - [Deprecation](#deprecation-1)
    - [API Change](#api-change-2)
    - [Feature](#feature-2)
    - [Bug or Regression](#bug-or-regression-2)
    - [Other (Cleanup or Flake)](#other-cleanup-or-flake-2)
  - [Dependencies](#dependencies-2)
    - [Added](#added-2)
    - [Changed](#changed-2)
    - [Removed](#removed-2)
- [v1.33.0-alpha.1](#v1330-alpha1)
  - [Downloads for v1.33.0-alpha.1](#downloads-for-v1330-alpha1)
    - [Source Code](#source-code-3)
    - [Client Binaries](#client-binaries-3)
    - [Server Binaries](#server-binaries-3)
    - [Node Binaries](#node-binaries-3)
    - [Container Images](#container-images-3)
  - [Changelog since v1.32.0](#changelog-since-v1320)
  - [Urgent Upgrade Notes](#urgent-upgrade-notes-1)
    - [(No, really, you MUST read this before you upgrade)](#no-really-you-must-read-this-before-you-upgrade-1)
  - [Changes by Kind](#changes-by-kind-3)
    - [API Change](#api-change-3)
    - [Feature](#feature-3)
    - [Documentation](#documentation)
    - [Bug or Regression](#bug-or-regression-3)
    - [Other (Cleanup or Flake)](#other-cleanup-or-flake-3)
  - [Dependencies](#dependencies-3)
    - [Added](#added-3)
    - [Changed](#changed-3)
    - [Removed](#removed-3)

<!-- END MUNGE: GENERATED_TOC -->

# v1.33.0-beta.0


## Downloads for v1.33.0-beta.0



### Source Code

filename | sha512 hash
-------- | -----------
[kubernetes.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes.tar.gz) | 53a7e0e0ad351ca0cfb99ca3258835cd9356dd10df3dc9737dc3ef08510b8afc0eafcac503b6168c24c13bbd1a93f9a06508b5b5c5c5ec2f45e31f86012409e0
[kubernetes-src.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-src.tar.gz) | 56d380d07e265c18f4b86e294b3944f330892588bd62301f8827ce726afd1e9d5e7335bc0c939c3a6297d2e4f5132c82d048858024718b10ff11c6e8d2c40cc9

### Client Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-client-darwin-amd64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-darwin-amd64.tar.gz) | 3ae3b1bc58812ce8a2a1e3ca0014c15b00e3f4edcc72d7cbaa50cede697384d8765f5bb49aac0d5b786295528dc1b07f0135931cfda4f48e33022640fb3c6b7a
[kubernetes-client-darwin-arm64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-darwin-arm64.tar.gz) | 4cb64d8f647c454f1c2c640077724237dcf056b5c2c461e0c5022650667ceb34e013d0727d4ef6d129502c094776501dbab68e993e96a2c6697cde212b42723e
[kubernetes-client-linux-386.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-linux-386.tar.gz) | ff44a2e6622c25c89430006dba64d1e40b78dfe88445d7af35ec7e92979fe880c009130694a9163266fda771fe1da9e0ebcbe9735b3591b0344c2e67513d996a
[kubernetes-client-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-linux-amd64.tar.gz) | 429c1b3838e0a7ce0c6f89b36589c8de64ca83ae1a197097feb1c19dbd9241f03648b590c8a57fa0c1ab1bcda769c46c2c562846bfe924317e86dba117f422b2
[kubernetes-client-linux-arm.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-linux-arm.tar.gz) | aaa2d51b539d269e2b1ec89f5c6308afe23bd13f766fad6e949f424c5db2002f2400dedab8cea6922339c920414de66a16fbd5d752e518982ec501cb803c0339
[kubernetes-client-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-linux-arm64.tar.gz) | 4591f8bdb027fe2eb52652335834777cb0ce509ff5643877724746210747ff3ca1d3b62b41d7c93e05dc9e30923e32e3cbe8fac856deccda2e958ed638b60e0f
[kubernetes-client-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-linux-ppc64le.tar.gz) | 55c3d7b37929af918bb29f304bb94dd21e21cd50b920290b4309a11d1507c9f3ca7c0506e6e23f94b9503da593d727fb136bdcb12d2da8766b993655107cadeb
[kubernetes-client-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-linux-s390x.tar.gz) | c7c1ba2071957e963d9d0824e061af0752489a9fcf9c2a601ce6d66dedbbc5f0f02c14c72cd16c48092024d71996a83bd59a2d01377fa537ff6093bef518e3fa
[kubernetes-client-windows-386.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-windows-386.tar.gz) | b2823da3f55a47940b3238e5b8d276e1ee6af6b10ffbd972ae1299da38f39e36b45772234f57d81c27759a4add33d486c29d8efc32deb779aba703fe3230a5cd
[kubernetes-client-windows-amd64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-windows-amd64.tar.gz) | f5de40e2b5596f40cf59422ed22a448e1d397a11c73de4b1694b04c235bbb2538bd5bf705efb99dc5d7cf24da3d935f7530bbc8680180a9967de4bc341d745d3
[kubernetes-client-windows-arm64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-client-windows-arm64.tar.gz) | 3c24e08b0634465bd7910389be63a09b9b750529c076e7759c9100c07dbeb9dd4d0caa0f871bd776261cec88059aecdf29e6dbbbabaf357b79cbcf620ee1b0d8

### Server Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-server-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-server-linux-amd64.tar.gz) | 8e2c99d48ecc0b806208a983837026943916580ccd2911362b4af5b3aa45e16703f18262fd64d81844854b06f7025c543a6964cc0b3455b5c300e099773c2847
[kubernetes-server-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-server-linux-arm64.tar.gz) | 61688ad68057dd4c7f7f41206dfb558407a7317cdfdb33d305d81c08e93a0f5e11efbd68e10e12aeb7cc873550bbd822f270167ef2208f877bdb8db58f12f14b
[kubernetes-server-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-server-linux-ppc64le.tar.gz) | 4349666887f862bd45bca0d0488128b33107e3f3ebc69cf9c67dbefbd3539431c4e3ff944b8e6948ad0896bbc9c7a99ac9242d478ee44d052a49d8519e9cc017
[kubernetes-server-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-server-linux-s390x.tar.gz) | 351c6345cf88079c124d4f1a1401528d2c5ba1d8bc24ad16c49ef237890ff7090436224c631f5e4fe9a0a9a0a439e983f883854f24a1f96caeed5e9f12522e11

### Node Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-node-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-node-linux-amd64.tar.gz) | e59bf9f26252f94bca19967b1816db3071f0936c1716c58381b4ec0601b16aa050b07404fe40accf17b978d0f08ceddf859e333ff0ba3982a9c161e5b165526a
[kubernetes-node-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-node-linux-arm64.tar.gz) | a407c77a47a7fa38dd1cbf926e9b3b43a46e4dc46598b55cbc7dfa6073b7f6e429f2ba3c2e2d0c2cf8dcf51afb98c108d46986ed06c41a53973e8792d32c20a3
[kubernetes-node-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-node-linux-ppc64le.tar.gz) | 33312107574b1a6403b63851c65b660f37a0086a7f143843d3877f8974081fb4076063b48bc2e0b0f6733690a2edf00c3f704fabc80903fd8f07690a9d86f52d
[kubernetes-node-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-node-linux-s390x.tar.gz) | c253042a95cac403026ac69a304d0a41c36fa210d89c164f81b6388bd695720cf2b143b9543d79c965a5939a116aecffef2476b3f4888f6ab8da27bcd37529e3
[kubernetes-node-windows-amd64.tar.gz](https://dl.k8s.io/v1.33.0-beta.0/kubernetes-node-windows-amd64.tar.gz) | d6ef20e8f5fd6378065354c461221e879f16f90de58ea7c5662efe7981d10949031e986ff01dd719ad8d7e267491d8dba9fdfd2166265a87b863f9771241000f

### Container Images

All container images are available as manifest lists and support the described
architectures. It is also possible to pull a specific architecture directly by
adding the "-$ARCH" suffix  to the container image name.

name | architectures
---- | -------------
[registry.k8s.io/conformance:v1.33.0-beta.0](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-s390x)
[registry.k8s.io/kube-apiserver:v1.33.0-beta.0](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-s390x)
[registry.k8s.io/kube-controller-manager:v1.33.0-beta.0](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-s390x)
[registry.k8s.io/kube-proxy:v1.33.0-beta.0](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-s390x)
[registry.k8s.io/kube-scheduler:v1.33.0-beta.0](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-s390x)
[registry.k8s.io/kubectl:v1.33.0-beta.0](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-s390x)

## Changelog since v1.33.0-alpha.3

## Changes by Kind

### API Change

- DRA support for a "one-of" prioritized list of selection criteria to satisfy a device request in a resource claim. ([#128586](https://github.com/kubernetes/kubernetes/pull/128586), [@mortent](https://github.com/mortent)) [SIG API Machinery, Apps, Etcd, Node, Scheduling and Testing]
- For the InPlacePodVerticalScaling feature, the API server will no longer set the resize status to `Proposed` upon receiving a resize request. ([#130574](https://github.com/kubernetes/kubernetes/pull/130574), [@natasha41575](https://github.com/natasha41575)) [SIG Apps, Node and Testing]
- The apiserver will now return warnings if you create objects with "invalid" IP or
  CIDR values (like "192.168.000.005", which should not have the extra zeros).
  Values with non-standard formats can introduce security problems, and will
  likely be forbidden in a future Kubernetes release. ([#128786](https://github.com/kubernetes/kubernetes/pull/128786), [@danwinship](https://github.com/danwinship)) [SIG Apps, Network and Node]
- When the `PodObservedGenerationTracking` feature gate is set, the kubelet will populate `status.observedGeneration` to reflect the pod's latest `metadata.generation` that it has observed. ([#130352](https://github.com/kubernetes/kubernetes/pull/130352), [@natasha41575](https://github.com/natasha41575)) [SIG API Machinery, Apps, CLI, Node, Release, Scheduling, Storage, Testing and Windows]

### Feature

- Add mechanism that every 5 minutes calculates a digest of etcd and watch cache and exposes it as `apiserver_storage_digest` metric ([#130475](https://github.com/kubernetes/kubernetes/pull/130475), [@serathius](https://github.com/serathius)) [SIG API Machinery, Instrumentation and Testing]
- Adds apiserver.latency.k8s.io/authentication annotation to the audit log to record the time spent authenticating slow requests.
  - Adds apiserver.latency.k8s.io/authorization annotation to the audit log to record the time spent authorizing slow requests. ([#130571](https://github.com/kubernetes/kubernetes/pull/130571), [@hakuna-matatah](https://github.com/hakuna-matatah)) [SIG Auth]
- Allow for dynamic configuration of service account name and audience kubelet can request a token for as part of the node audience restriction feature. ([#130485](https://github.com/kubernetes/kubernetes/pull/130485), [@aramase](https://github.com/aramase)) [SIG Auth and Testing]
- Endpoints resources created by the Endpoints controller now have a label indicating this.
  (Users creating Endpoints by hand _can_ also add this label themselves, but they ought
  to switch to creating EndpointSlices rather than Endpoints anyway.) ([#130564](https://github.com/kubernetes/kubernetes/pull/130564), [@danwinship](https://github.com/danwinship)) [SIG Apps and Network]
- Pod resource checkpointing is now tracked by the `allocated_pods_state` and `actuated_pods_state` files, and  `pod_status_manager_state` is no longer used. ([#130599](https://github.com/kubernetes/kubernetes/pull/130599), [@tallclair](https://github.com/tallclair)) [SIG Node]
- Scheduling Framework exposes NodeInfo to the ScorePlugin. ([#130537](https://github.com/kubernetes/kubernetes/pull/130537), [@saintube](https://github.com/saintube)) [SIG Scheduling, Storage and Testing]
- Set feature gate `OrderedNamespaceDeletion` on by default. ([#130507](https://github.com/kubernetes/kubernetes/pull/130507), [@cici37](https://github.com/cici37)) [SIG API Machinery and Apps]

### Bug or Regression

- Fix a bug on InPlacePodVerticalScalingExclusiveCPUs feature gate exclusive assignment availability check. ([#130559](https://github.com/kubernetes/kubernetes/pull/130559), [@esotsal](https://github.com/esotsal)) [SIG Node]
- Fix kubelet restart unmounts volumes of running pods if the referenced PVC is being deleted by the user ([#130335](https://github.com/kubernetes/kubernetes/pull/130335), [@carlory](https://github.com/carlory)) [SIG Node, Storage and Testing]
- Removed a warning around Linux user namespaces and kernel version. If the feature gate `UserNamespacesSupport` was enabled, the kubelet previously warned when detecting a Linux kernel version earlier than 6.3.0. User namespace support on Linux typically does still need kernel 6.3 or newer, but it can work in older kernels too. ([#130243](https://github.com/kubernetes/kubernetes/pull/130243), [@rata](https://github.com/rata)) [SIG Node]
- The BalancedAllocation plugin will skip all best-effort (zero-requested) pod. ([#130260](https://github.com/kubernetes/kubernetes/pull/130260), [@Bowser1704](https://github.com/Bowser1704)) [SIG Scheduling]
- YAML input which might previously have been confused for JSON is now accepted. ([#130666](https://github.com/kubernetes/kubernetes/pull/130666), [@thockin](https://github.com/thockin)) [SIG API Machinery]

### Other (Cleanup or Flake)

- Changed the error message displayed when a pod is trying to attach a volume that does not match the label/selector from "x node(s) had volume node affinity conflict" to "x node(s) didn't match PersistentVolume's node affinity". ([#129887](https://github.com/kubernetes/kubernetes/pull/129887), [@rhrmo](https://github.com/rhrmo)) [SIG Scheduling and Storage]
- Client-gen now sorts input group/versions to ensure stable output generation even with unsorted inputs ([#130626](https://github.com/kubernetes/kubernetes/pull/130626), [@BenTheElder](https://github.com/BenTheElder)) [SIG API Machinery]
- E2e.test: [Feature:OffByDefault] is added to test names when specifying a featuregate which is not on by default ([#130655](https://github.com/kubernetes/kubernetes/pull/130655), [@BenTheElder](https://github.com/BenTheElder)) [SIG Auth and Testing]
- Kubelet no longer logs multiple errors when running on a system with no iptables binaries installed. ([#129826](https://github.com/kubernetes/kubernetes/pull/129826), [@danwinship](https://github.com/danwinship)) [SIG Network and Node]

## Dependencies

### Added
- github.com/containerd/errdefs/pkg: [v0.3.0](https://github.com/containerd/errdefs/tree/pkg/v0.3.0)
- github.com/klauspost/compress: [v1.18.0](https://github.com/klauspost/compress/tree/v1.18.0)
- github.com/kylelemons/godebug: [v1.1.0](https://github.com/kylelemons/godebug/tree/v1.1.0)
- github.com/opencontainers/cgroups: [v0.0.1](https://github.com/opencontainers/cgroups/tree/v0.0.1)
- github.com/russross/blackfriday: [v1.6.0](https://github.com/russross/blackfriday/tree/v1.6.0)
- github.com/santhosh-tekuri/jsonschema/v5: [v5.3.1](https://github.com/santhosh-tekuri/jsonschema/tree/v5.3.1)
- sigs.k8s.io/randfill: v1.0.0

### Changed
- cloud.google.com/go/compute: v1.25.1 → v1.23.3
- github.com/cilium/ebpf: [v0.16.0 → v0.17.3](https://github.com/cilium/ebpf/compare/v0.16.0...v0.17.3)
- github.com/containerd/containerd/api: [v1.7.19 → v1.8.0](https://github.com/containerd/containerd/compare/api/v1.7.19...api/v1.8.0)
- github.com/containerd/errdefs: [v0.1.0 → v1.0.0](https://github.com/containerd/errdefs/compare/v0.1.0...v1.0.0)
- github.com/containerd/ttrpc: [v1.2.5 → v1.2.6](https://github.com/containerd/ttrpc/compare/v1.2.5...v1.2.6)
- github.com/containerd/typeurl/v2: [v2.2.0 → v2.2.2](https://github.com/containerd/typeurl/compare/v2.2.0...v2.2.2)
- github.com/cyphar/filepath-securejoin: [v0.3.5 → v0.4.1](https://github.com/cyphar/filepath-securejoin/compare/v0.3.5...v0.4.1)
- github.com/go-logfmt/logfmt: [v0.5.1 → v0.4.0](https://github.com/go-logfmt/logfmt/compare/v0.5.1...v0.4.0)
- github.com/google/cadvisor: [v0.51.0 → v0.52.1](https://github.com/google/cadvisor/compare/v0.51.0...v0.52.1)
- github.com/google/go-cmp: [v0.6.0 → v0.7.0](https://github.com/google/go-cmp/compare/v0.6.0...v0.7.0)
- github.com/google/gofuzz: [v1.2.0 → v1.0.0](https://github.com/google/gofuzz/compare/v1.2.0...v1.0.0)
- github.com/matttproud/golang_protobuf_extensions: [v1.0.2 → v1.0.1](https://github.com/matttproud/golang_protobuf_extensions/compare/v1.0.2...v1.0.1)
- github.com/opencontainers/image-spec: [v1.1.0 → v1.1.1](https://github.com/opencontainers/image-spec/compare/v1.1.0...v1.1.1)
- github.com/opencontainers/runc: [v1.2.1 → v1.2.5](https://github.com/opencontainers/runc/compare/v1.2.1...v1.2.5)
- github.com/prometheus/client_golang: [v1.19.1 → v1.22.0-rc.0](https://github.com/prometheus/client_golang/compare/v1.19.1...v1.22.0-rc.0)
- github.com/prometheus/common: [v0.55.0 → v0.62.0](https://github.com/prometheus/common/compare/v0.55.0...v0.62.0)
- golang.org/x/time: v0.7.0 → v0.9.0
- google.golang.org/appengine: v1.6.7 → v1.4.0
- google.golang.org/protobuf: v1.35.2 → v1.36.5
- k8s.io/kube-openapi: 2c72e55 → e5f78fe
- sigs.k8s.io/structured-merge-diff/v4: v4.4.2 → v4.6.0

### Removed
- github.com/checkpoint-restore/go-criu/v6: [v6.3.0](https://github.com/checkpoint-restore/go-criu/tree/v6.3.0)
- github.com/containerd/console: [v1.0.4](https://github.com/containerd/console/tree/v1.0.4)
- github.com/go-kit/log: [v0.2.1](https://github.com/go-kit/log/tree/v0.2.1)
- github.com/moby/sys/user: [v0.3.0](https://github.com/moby/sys/tree/user/v0.3.0)
- github.com/seccomp/libseccomp-golang: [v0.10.0](https://github.com/seccomp/libseccomp-golang/tree/v0.10.0)
- github.com/syndtr/gocapability: [42c35b4](https://github.com/syndtr/gocapability/tree/42c35b4)
- github.com/urfave/cli: [v1.22.14](https://github.com/urfave/cli/tree/v1.22.14)



# v1.33.0-alpha.3


## Downloads for v1.33.0-alpha.3



### Source Code

filename | sha512 hash
-------- | -----------
[kubernetes.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes.tar.gz) | 52751abcbaac8786aa52a8687c6c7d72c6aaa1a8e837ce873ecd66503a92a35c09fd01e85543240c5b51d0c9f97fd374a0dec64fea4acdda6e65b0b6bc183202
[kubernetes-src.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-src.tar.gz) | cbd9967ec5bc31c509f8f9f09a6b40d977c30454554eb743e4c2da92382451fd1d8fae6f011738bccb11fc158fe8f31cc0ddf8d91be298ffa69da8a569c7ef3e

### Client Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-client-darwin-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-darwin-amd64.tar.gz) | 849b061df1d8cc4a727977329504476023bf4c4f4c4ea4b274e914874e3e960eb7b96f96d6377b582e0f17c44bb87aee72b494040dac0b3316f5761c0ad0f227
[kubernetes-client-darwin-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-darwin-arm64.tar.gz) | 910c2f6df7bb8fb901db21f6ddd7e8ec3831cd9f195282139443535e8d590505a56bdf28a3414b3e8438b1ecf716b06b660713e6ed61a907bb381dee1b1391a7
[kubernetes-client-linux-386.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-linux-386.tar.gz) | 202d4af420be89c295facaf5604da05aed1716c8b8f4b81d9e408aaf56cb74593a652fa6b0d45c9c491e12a5c3fadd5eb1aa5988ec5b2d4975e2feb864df5327
[kubernetes-client-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-linux-amd64.tar.gz) | d5bb5bea82ff07540188e0844454a40313752aae99c1dccba54673cb9155b22d8734b3500d83e93b5d59e44b173be666f40a5471927037fa90653b9f7e11f725
[kubernetes-client-linux-arm.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-linux-arm.tar.gz) | 3cd086b710581dd40a0f6a3449b820a9a98a0721099d43844e2764a1d05ef4f62f3232bbbae7d65b63d9c0c994d8bdbba033c1042406beefea483c8358a9c29a
[kubernetes-client-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-linux-arm64.tar.gz) | d3a65059addbf899bab1551b3eed78e2b3ef5222b01d041e6bad31454e8e7e05f7f3ae5650691b726bd2bbf8896fd9699f788aa939e1f27089d3fe4cfcccf8cc
[kubernetes-client-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-linux-ppc64le.tar.gz) | 4fdbc47cbae8fe3f23cc0b42157542e98cadcf82056d9a36c239d6fd720afcbf307b3b01734893c62235ee39618a76a947ae821e12de87d4eb18d22b4bd93bfb
[kubernetes-client-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-linux-s390x.tar.gz) | 5202a5b8afbb0685b370cc0d6866b7a8fe9ece8cc586052af538e438a38019b648c0c4f7b30834529c74f04f9ce740d057d7af77c144ba8b71755895dccd9866
[kubernetes-client-windows-386.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-windows-386.tar.gz) | 6841c68ae7281e7d0352a14123e9ffa06ea70b3d467184718f2894a2062c986b8d42f0d3446508ebe5c3128148592119664bccaad06f8f8dd04424185a7e8911
[kubernetes-client-windows-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-windows-amd64.tar.gz) | da49d82906d57efb03268a2df299eca13e73e33936bb0b15fe5ff6f93037e06818d7710200d5a69e911b361db3009094a05a022beee2fabe05cae744d13e62b3
[kubernetes-client-windows-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-client-windows-arm64.tar.gz) | 2e047a38d94083c2a89b848fa8b9877ee083ac973cd97fbaaa0a6cc05f46a9ca21b6b5769478843f3239c5d4f8ed343b77a30ab6c4d8f84e5f60569b754a93a9

### Server Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-server-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-server-linux-amd64.tar.gz) | 5c4849fb85141d8cc1e327a567c74650914cdee92d39e5a8513dacc0afe4424986e899eae6fbe6160eedf9bb5102921634330598a10ab41c20f370e07b5d8de0
[kubernetes-server-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-server-linux-arm64.tar.gz) | a3db2fb73e65237181a92a908d713e033bd9c8be98ab538aea6f86945ad4b402c9b36b8496f69815667c50b3ab44b5c6a5f50a91d253a3b6e7964939cb59e8c1
[kubernetes-server-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-server-linux-ppc64le.tar.gz) | f7593d0e205e4797b635cdff3b20e25b90981dd5403930fb6c4be3c99143bc37eda5f7c45bdf6088e4d61625f10b25fd3b5d0b4d1b2d460b47764b28645f395f
[kubernetes-server-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-server-linux-s390x.tar.gz) | f0105558d1f31f710482e367d61a76e2c97207d36a74d16e3bf7b94c0482a867f551fcc505cd228768606a42a97e4620df2863fbeca0a737387a34734e7ae553

### Node Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-node-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-node-linux-amd64.tar.gz) | fe4fdad8e3f0bb159fa8ad1d1fc2952d650d49e2e517b053636cc2969ba605f2c5258d468bfd2ac02f580d6d462f17aa98136a198c58dc56cb0fbd4ca53745ec
[kubernetes-node-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-node-linux-arm64.tar.gz) | faaa6f5bcd238729b12557ac27f99741557921daa6bdbfca6c78f8f84390117cd6f378e41ed4e7afa57bba08b1f5e361a6984d299b8f4ed88c8c39890a0a03cd
[kubernetes-node-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-node-linux-ppc64le.tar.gz) | 6a46ece235a5496c82ceaaf052bc07126ec6c3154cffa0a5c1ade77fd7edd445ed75692e1937a46384cad67800e9ea37b09ae8d342be0625155e4a5f0451f569
[kubernetes-node-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-node-linux-s390x.tar.gz) | 06450e76910f8342cd2d48cdba5c7d5b1f2b5e4faa744d6b43184df2a127701bd626e90d15e7621bd1a0749ffa6581dae13eb651798f3023d0fae7f30907e9d8
[kubernetes-node-windows-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.3/kubernetes-node-windows-amd64.tar.gz) | 852b94df5a79fdf7d96e3fccfde67cb238bae90aa56b70271d545c158693c0437777caa62b2381e98520f510e9124ace9c9126cafd225e5cb5ca30c9099867ce

### Container Images

All container images are available as manifest lists and support the described
architectures. It is also possible to pull a specific architecture directly by
adding the "-$ARCH" suffix  to the container image name.

name | architectures
---- | -------------
[registry.k8s.io/conformance:v1.33.0-alpha.3](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-s390x)
[registry.k8s.io/kube-apiserver:v1.33.0-alpha.3](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-s390x)
[registry.k8s.io/kube-controller-manager:v1.33.0-alpha.3](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-s390x)
[registry.k8s.io/kube-proxy:v1.33.0-alpha.3](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-s390x)
[registry.k8s.io/kube-scheduler:v1.33.0-alpha.3](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-s390x)
[registry.k8s.io/kubectl:v1.33.0-alpha.3](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-s390x)

## Changelog since v1.33.0-alpha.2

## Urgent Upgrade Notes

### (No, really, you MUST read this before you upgrade)

 - The behavior of the KUBE_PROXY_NFTABLES_SKIP_KERNEL_VERSION_CHECK environment variable has been fixed in the nftables proxier. The kernel version check is only skipped when this variable is explicitly set to a non-empty value. If you need to skip the check, set the KUBE_PROXY_NFTABLES_SKIP_KERNEL_VERSION_CHECK environment variable. ([#130401](https://github.com/kubernetes/kubernetes/pull/130401), [@ryota-sakamoto](https://github.com/ryota-sakamoto)) [SIG Network]
 
## Changes by Kind

### Deprecation

- The v1 Endpoints API is now officially deprecated (though still fully supported). The API will not be removed, but all users should use the EndpointSlice API instead. ([#130098](https://github.com/kubernetes/kubernetes/pull/130098), [@danwinship](https://github.com/danwinship)) [SIG API Machinery and Network]

### API Change

- InPlacePodVerticalScaling: Memory limits cannot be decreased unless the memory resize restart policy is set to `RestartContainer`. Container resizePolicy is no longer mutable. ([#130183](https://github.com/kubernetes/kubernetes/pull/130183), [@tallclair](https://github.com/tallclair)) [SIG Apps and Node]
- Introduced API type coordination.k8s.io/v1beta1/LeaseCandidate ([#130291](https://github.com/kubernetes/kubernetes/pull/130291), [@Jefftree](https://github.com/Jefftree)) [SIG API Machinery, Etcd and Testing]
- KEP-3857: Recursive Read-only (RRO) mounts: promote to GA ([#130116](https://github.com/kubernetes/kubernetes/pull/130116), [@AkihiroSuda](https://github.com/AkihiroSuda)) [SIG Apps, Node and Testing]
- MergeDefaultEvictionSettings indicates that defaults for the evictionHard, evictionSoft, evictionSoftGracePeriod, and evictionMinimumReclaim fields should be merged into values specified for those fields in this configuration. Signals specified in this configuration take precedence. Signals not specified in this configuration inherit their defaults. ([#127577](https://github.com/kubernetes/kubernetes/pull/127577), [@vaibhav2107](https://github.com/vaibhav2107)) [SIG API Machinery and Node]
- Promote the Job's JobBackoffLimitPerIndex feature-gate to stable. ([#130061](https://github.com/kubernetes/kubernetes/pull/130061), [@mimowo](https://github.com/mimowo)) [SIG API Machinery, Apps, Architecture and Testing]
- Promoted the feature gate `AnyVolumeDataSource` to GA. ([#129770](https://github.com/kubernetes/kubernetes/pull/129770), [@sunnylovestiramisu](https://github.com/sunnylovestiramisu)) [SIG Apps, Storage and Testing]

### Feature

- Added a `/statusz` endpoint for kube-scheduler ([#128987](https://github.com/kubernetes/kubernetes/pull/128987), [@Henrywu573](https://github.com/Henrywu573)) [SIG Instrumentation, Scheduling and Testing]
- Added a alpha feature gate `OrderedNamespaceDeletion`. When enabled, the pods resources are deleted before all other resources while namespace deletion to ensure workload security. ([#130035](https://github.com/kubernetes/kubernetes/pull/130035), [@cici37](https://github.com/cici37)) [SIG API Machinery, Apps and Testing]
- Allow ImageVolume for Restricted PSA profiles ([#130394](https://github.com/kubernetes/kubernetes/pull/130394), [@Barakmor1](https://github.com/Barakmor1)) [SIG Auth]
- Changed metadata management for Pods to populate `.metadata.generation` on writes. New pods will have a `metadata.generation` of 1; updates to mutable fields in the Pod `.spec` will result in `metadata.generation` being incremented by 1. ([#130181](https://github.com/kubernetes/kubernetes/pull/130181), [@natasha41575](https://github.com/natasha41575)) [SIG Apps, Node and Testing]
- Extended the kube-apiserver loopback client certificate validity to 14 months to align with the updated Kubernetes support lifecycle. ([#130047](https://github.com/kubernetes/kubernetes/pull/130047), [@HirazawaUi](https://github.com/HirazawaUi)) [SIG API Machinery and Auth]
- Improved how the API server responds to **list** requests where the response format negotiates to JSON. List responses in JSON are marshalled one element at the time, drastically reducing memory needed to serve large collections. Streaming list responses can be disabled via the `StreamingJSONListEncoding` feature gate. ([#129334](https://github.com/kubernetes/kubernetes/pull/129334), [@serathius](https://github.com/serathius)) [SIG API Machinery, Architecture and Release]
- Kubernetes is now built with go 1.24.0 ([#129688](https://github.com/kubernetes/kubernetes/pull/129688), [@cpanato](https://github.com/cpanato)) [SIG API Machinery, Architecture, Auth, CLI, Cloud Provider, Cluster Lifecycle, Instrumentation, Network, Node, Release, Scheduling, Storage and Testing]
- Promote RelaxedDNSSearchValidation to beta, allowing for Pod search domains to be a single dot "." or contain an underscore "_" ([#130128](https://github.com/kubernetes/kubernetes/pull/130128), [@adrianmoisey](https://github.com/adrianmoisey)) [SIG Apps and Network]
- Promoted the `CRDValidationRatcheting` feature gate to GA in 1.33 ([#130013](https://github.com/kubernetes/kubernetes/pull/130013), [@yongruilin](https://github.com/yongruilin)) [SIG API Machinery]
- Promoted the feature gate `HonorPVReclaimPolicy` to GA. ([#129583](https://github.com/kubernetes/kubernetes/pull/129583), [@carlory](https://github.com/carlory)) [SIG Apps, Storage and Testing]
- Promotes kubectl --subresource flag to stable. ([#130238](https://github.com/kubernetes/kubernetes/pull/130238), [@soltysh](https://github.com/soltysh)) [SIG CLI]
- Various controllers that write out IP address or CIDR values to API objects now
  ensure that they always write out the values in canonical form. ([#130101](https://github.com/kubernetes/kubernetes/pull/130101), [@danwinship](https://github.com/danwinship)) [SIG Apps, Network and Node]

### Bug or Regression

- Add progress tracking to volumes permission and ownership change ([#130398](https://github.com/kubernetes/kubernetes/pull/130398), [@gnufied](https://github.com/gnufied)) [SIG Node and Storage]
- Bugfix for Events that fail to be created when the referenced object name is not a valid Event name, using an UUID as name instead of the referenced object name and the timestamp suffix. ([#129790](https://github.com/kubernetes/kubernetes/pull/129790), [@aojea](https://github.com/aojea)) [SIG API Machinery]
- CSI drivers that calls IsLikelyNotMountPoint should not assume false means that the path is a mount point. Each CSI driver needs to make sure correct usage of return value of IsLikelyNotMountPoint because if the file is an irregular file but not a mount point is acceptable ([#129370](https://github.com/kubernetes/kubernetes/pull/129370), [@andyzhangx](https://github.com/andyzhangx)) [SIG Storage and Windows]
- Fix very rare and sporadic network issues when the host is under heavy load by adding retries for interrupted netlink calls ([#130256](https://github.com/kubernetes/kubernetes/pull/130256), [@adrianmoisey](https://github.com/adrianmoisey)) [SIG Network]
- Fixed an issue in register-gen where imports for k8s.io/apimachinery/pkg/runtime and k8s.io/apimachinery/pkg/runtime/schema were missing. ([#129307](https://github.com/kubernetes/kubernetes/pull/129307), [@LionelJouin](https://github.com/LionelJouin)) [SIG API Machinery]
- Fixes a 1.32 regression starting pods with postStart hooks specified ([#129946](https://github.com/kubernetes/kubernetes/pull/129946), [@alex-petrov-vt](https://github.com/alex-petrov-vt)) [SIG API Machinery]
- Fixes a 1.32 regression where nodes may fail to report status and renew serving certificates after the kubelet restarts ([#130348](https://github.com/kubernetes/kubernetes/pull/130348), [@aojea](https://github.com/aojea)) [SIG Node]
- Fixes an issue in the CEL CIDR library where subnets contained within another CIDR were incorrectly rejected as not contained ([#130450](https://github.com/kubernetes/kubernetes/pull/130450), [@JoelSpeed](https://github.com/JoelSpeed)) [SIG API Machinery]
- Kube-apiserver: Fix a bug where the `ResourceQuota` admission plugin does not respect ANY scope change when a resource is being updated. i.e., to set/unset an existing pod's `terminationGracePeriodSeconds` field. ([#130060](https://github.com/kubernetes/kubernetes/pull/130060), [@carlory](https://github.com/carlory)) [SIG API Machinery, Scheduling and Testing]
- Kube-apiserver: shortening the grace period during a pod deletion no longer moves the metadata.deletionTimestamp into the past ([#122646](https://github.com/kubernetes/kubernetes/pull/122646), [@liggitt](https://github.com/liggitt)) [SIG API Machinery]
- Kube-proxy, when using a Service with External or LoadBalancer IPs on UDP services , was consuming a large amount of CPU because it was not filtering by the Service destination port and trying to delete all the UDP entries associated to the service. ([#130484](https://github.com/kubernetes/kubernetes/pull/130484), [@aojea](https://github.com/aojea)) [SIG Network]
- Kubeadm: fix panic when no UpgradeConfiguration was found in the config file ([#130202](https://github.com/kubernetes/kubernetes/pull/130202), [@SataQiu](https://github.com/SataQiu)) [SIG Cluster Lifecycle]
- The following roles have had `Watch` added to them (prefixed with `system:controller:`):
  
  - `cronjob-controller`
  - `endpoint-controller`
  - `endpointslice-controller`
  - `endpointslicemirroring-controller`
  - `horizontal-pod-autoscaler`
  - `node-controller`
  - `pod-garbage-collector`
  - `storage-version-migrator-controller` ([#130405](https://github.com/kubernetes/kubernetes/pull/130405), [@kariya-mitsuru](https://github.com/kariya-mitsuru)) [SIG Auth]
- The response from kube-apiserver /flagz endpoint would respond correctly with parsed flags value when the feature-gate ComponentFlagz is enabled ([#130328](https://github.com/kubernetes/kubernetes/pull/130328), [@richabanker](https://github.com/richabanker)) [SIG API Machinery and Instrumentation]
- When using the Alpha DRAResourceClaimDeviceStatus feature, IP address values
  in the NetworkDeviceData are now validated more strictly. ([#129219](https://github.com/kubernetes/kubernetes/pull/129219), [@danwinship](https://github.com/danwinship)) [SIG Network]

### Other (Cleanup or Flake)

- 1. kube-apiserver: removed the deprecated the `--cloud-provider` and `--cloud-config` CLI parameters.
  2. removed generally available feature-gate `DisableCloudProviders` and `DisableKubeletCloudCredentialProviders` ([#130162](https://github.com/kubernetes/kubernetes/pull/130162), [@carlory](https://github.com/carlory)) [SIG API Machinery, Cloud Provider, Node and Testing]
- Changed the error message displayed when a pod is trying to attach a volume that does not match the label/selector from "x node(s) had volume node affinity conflict" to "x node(s) didn't match PersistentVolume's node affinity". ([#129887](https://github.com/kubernetes/kubernetes/pull/129887), [@rhrmo](https://github.com/rhrmo)) [SIG Scheduling and Storage]
- Kubeadm: Use generic terminology in logs instead of direct mentions of yaml/json. ([#130345](https://github.com/kubernetes/kubernetes/pull/130345), [@HirazawaUi](https://github.com/HirazawaUi)) [SIG Cluster Lifecycle]
- Remove the JobPodFailurePolicy feature gate that graduated to GA in 1.31 and was unconditionally enabled. ([#129498](https://github.com/kubernetes/kubernetes/pull/129498), [@carlory](https://github.com/carlory)) [SIG Apps]
- Removed general available feature-gate `AppArmor`. ([#129375](https://github.com/kubernetes/kubernetes/pull/129375), [@carlory](https://github.com/carlory)) [SIG Auth and Node]
- Removed generally available feature-gate `AppArmorFields`. ([#129497](https://github.com/kubernetes/kubernetes/pull/129497), [@carlory](https://github.com/carlory)) [SIG Node]

## Dependencies

### Added
- github.com/planetscale/vtprotobuf: [0393e58](https://github.com/planetscale/vtprotobuf/tree/0393e58)
- go.opentelemetry.io/auto/sdk: v1.1.0

### Changed
- cloud.google.com/go/compute/metadata: v0.3.0 → v0.5.0
- github.com/cncf/xds/go: [555b57e → b4127c9](https://github.com/cncf/xds/compare/555b57e...b4127c9)
- github.com/envoyproxy/go-control-plane: [v0.12.0 → v0.13.0](https://github.com/envoyproxy/go-control-plane/compare/v0.12.0...v0.13.0)
- github.com/envoyproxy/protoc-gen-validate: [v1.0.4 → v1.1.0](https://github.com/envoyproxy/protoc-gen-validate/compare/v1.0.4...v1.1.0)
- github.com/golang/glog: [v1.2.1 → v1.2.2](https://github.com/golang/glog/compare/v1.2.1...v1.2.2)
- github.com/gorilla/websocket: [v1.5.0 → v1.5.3](https://github.com/gorilla/websocket/compare/v1.5.0...v1.5.3)
- github.com/grpc-ecosystem/grpc-gateway/v2: [v2.20.0 → v2.24.0](https://github.com/grpc-ecosystem/grpc-gateway/compare/v2.20.0...v2.24.0)
- github.com/rogpeppe/go-internal: [v1.12.0 → v1.13.1](https://github.com/rogpeppe/go-internal/compare/v1.12.0...v1.13.1)
- github.com/stretchr/testify: [v1.9.0 → v1.10.0](https://github.com/stretchr/testify/compare/v1.9.0...v1.10.0)
- go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc: v0.53.0 → v0.58.0
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp: v0.53.0 → v0.58.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc: v1.27.0 → v1.33.0
- go.opentelemetry.io/otel/exporters/otlp/otlptrace: v1.28.0 → v1.33.0
- go.opentelemetry.io/otel/metric: v1.28.0 → v1.33.0
- go.opentelemetry.io/otel/sdk: v1.28.0 → v1.33.0
- go.opentelemetry.io/otel/trace: v1.28.0 → v1.33.0
- go.opentelemetry.io/otel: v1.28.0 → v1.33.0
- go.opentelemetry.io/proto/otlp: v1.3.1 → v1.4.0
- golang.org/x/crypto: v0.31.0 → v0.35.0
- golang.org/x/oauth2: v0.23.0 → v0.27.0
- golang.org/x/sync: v0.10.0 → v0.11.0
- golang.org/x/sys: v0.28.0 → v0.30.0
- golang.org/x/term: v0.27.0 → v0.29.0
- golang.org/x/text: v0.21.0 → v0.22.0
- google.golang.org/genproto/googleapis/api: f6391c0 → e6fa225
- google.golang.org/genproto/googleapis/rpc: f6391c0 → e6fa225
- google.golang.org/grpc: v1.65.0 → v1.68.1
- google.golang.org/protobuf: v1.35.1 → v1.35.2
- k8s.io/gengo/v2: 2b36238 → 1244d31
- sigs.k8s.io/apiserver-network-proxy/konnectivity-client: v0.31.1 → v0.31.2

### Removed
_Nothing has changed._



# v1.33.0-alpha.2


## Downloads for v1.33.0-alpha.2



### Source Code

filename | sha512 hash
-------- | -----------
[kubernetes.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes.tar.gz) | ee13af765b25d466423e51cea5359effb1a095b9033032040bca8569a372656ab27ec38b8b9a4a85a7256f6390c33c0cb7d145ce876ccf282cdf5b3224560724
[kubernetes-src.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-src.tar.gz) | bc32551357ae67573ac9ab4c650bcd547f46a29848e20fc3db286d0e45a22ed254ee2c8d6fe84c4288ebc3df6c3acb118435a532c9cf9f3f5e8d33f4512de806

### Client Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-client-darwin-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-darwin-amd64.tar.gz) | aab9eac3bc604831cfdc926f6d3f12afe6266a2c3808503141ad5780ffcd188f08db3fbad4fedc73da1c612d19bd2e55ba13031fef22ea4839cb294eb54b5767
[kubernetes-client-darwin-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-darwin-arm64.tar.gz) | 373fa812af4ed11b9a3b278c44335fd3618c9fb77aa789311e07e37c4bad81e08b066528dd086356e0bb1e116fa807f0015bc71f225afd5bef4dbbe3079034e1
[kubernetes-client-linux-386.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-linux-386.tar.gz) | e9f8a8925b2b7d3cf89dbaad251f0224945be354ae62c7736b891c73e19334039e68ac7b2dda99f26df0d7028127ccb630de085d2ad45255e263cb03f1f1e552
[kubernetes-client-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-linux-amd64.tar.gz) | 305ea43a314586911f32ae43b16f7a29274fe2a7d87b00b9fb57a4c5c885187a317272c731ddf9d41335905ff5f3640d7a4df7e68d070076e20ff1b2a32a78cd
[kubernetes-client-linux-arm.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-linux-arm.tar.gz) | f012b9e7d46874748655782e125a1a9b7d22c9bee77226eea9c789bc67f5644a9c8380d5fa5d7cc161659011266b9be060dd663603d85b7256deaab4866697c2
[kubernetes-client-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-linux-arm64.tar.gz) | 6952882b71ccc27412fce180844f2a5f9c147b5fb59c4b684d338b3cc767c6e0257f8edde1d1874acda0299ac7c22dba3788292dcbb083fdcc5e61387e8a16a8
[kubernetes-client-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-linux-ppc64le.tar.gz) | d4138ece8741e29c4d4fce07cd9cda38f622b5133a8757334cf5992e3242791213391c2a7ae7db95fee1d70d31b17fda3215d591fb8c9788e0e7d606fcc3a87f
[kubernetes-client-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-linux-s390x.tar.gz) | 511c4c53b20ecff1fc200e85a14211781e0d887a5536a3343a6a0c8ce05c175d073b810945fd1ddd2389318ea26e0ca412b7025ce9f168b76ad24a7ee85213a7
[kubernetes-client-windows-386.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-windows-386.tar.gz) | 68b781adad28a0ac8e19a624e6811f4e593ad4a1422294a40aa356f8ac05dfc5978f90b55a8716059b4a613caad8904961e9c7e74a4a803fed76c98739b126dd
[kubernetes-client-windows-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-windows-amd64.tar.gz) | 009f05ff583c6b43ffea01e9ff2f7e3cc13184646ce358338a2a1188f4750b02a9253a250c977576664d4d173ce8469a0d1be9a3968890a99969292ad1e001ec
[kubernetes-client-windows-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-client-windows-arm64.tar.gz) | 88dcf4ee3f86484d882632a10e63b7b6e64b844b17c3cc674a49e5ddab9cea091710e4503c46ee59d70fcf762dd1c4e954f5091154d23747a528ffa31d593273

### Server Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-server-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-server-linux-amd64.tar.gz) | 8023512c58f639b20bca94aa7bc3e908cd9fe2e213b655d1ad63da1507223651c6eb61ddf0d6670d664080e19e714640e3cf5aab4b9c6eb62fc0166cceabd3fd
[kubernetes-server-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-server-linux-arm64.tar.gz) | 7bb2a4530294bafb8f43ddfcfeefdd3fc8629c8dbfd11c2e789a59a930fe624262698311ed149e2c98cdde9bbf321b8c77213b4f562a5120a35ae645d1abf1ce
[kubernetes-server-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-server-linux-ppc64le.tar.gz) | 2f0071550e98d58b87dc56e5d27a1832827b256aa77ad4f68c3713ecd9e81fa66822d7604988c617c139d7e131e05664409f48f94f450cef467ab63727527e14
[kubernetes-server-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-server-linux-s390x.tar.gz) | 620241063ca4f09b4c71a3659e301246e82d841921e7956759d4a3a74bae7dff1d0951f5aea6928039714569ffbb5040f1ca73633bd90123000f4e18e9f196df

### Node Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-node-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-node-linux-amd64.tar.gz) | d54a8d3406df58a6941837e988e32cdc93bd5025dca1910dbcc1c89d8fa29dc09375c24d7f109fcf4d72c977933c091c225241a0988893a642a35edac04ee38d
[kubernetes-node-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-node-linux-arm64.tar.gz) | ddbf090dc9be5c30a968b655d2007485b8c94e5d95b7cd7e29bbb47ba562ae3ed5c15b965acd81acb715a8d706d967595601c5f0f8f5d6c0181626dcbe156c02
[kubernetes-node-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-node-linux-ppc64le.tar.gz) | c1dd2e061b7b305d481791be17234a5ca02f9c0c302a6044ac2b87940b10c5fc9c2817e00f59adeaab8b564181f8ccda4640dcfde67784daea38361f6faa4b2a
[kubernetes-node-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-node-linux-s390x.tar.gz) | 90974009d003cb911a54cad11bcca6805ceca64ed39120ce70029ece9c8e9a33d89803e92b5d251dce9f16267143914c1ed8542d9507cb3a020823a35b42cfdb
[kubernetes-node-windows-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.2/kubernetes-node-windows-amd64.tar.gz) | cc82205db3e6b6e1640ddbb4fbf8e1d81409c894c92aec1e2d5941c6a282414ada136d1f95403e25cb1f739095f838f6d40c97e65d2fa1dc2f3e6205bfb67249

### Container Images

All container images are available as manifest lists and support the described
architectures. It is also possible to pull a specific architecture directly by
adding the "-$ARCH" suffix  to the container image name.

name | architectures
---- | -------------
[registry.k8s.io/conformance:v1.33.0-alpha.2](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-s390x)
[registry.k8s.io/kube-apiserver:v1.33.0-alpha.2](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-s390x)
[registry.k8s.io/kube-controller-manager:v1.33.0-alpha.2](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-s390x)
[registry.k8s.io/kube-proxy:v1.33.0-alpha.2](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-s390x)
[registry.k8s.io/kube-scheduler:v1.33.0-alpha.2](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-s390x)
[registry.k8s.io/kubectl:v1.33.0-alpha.2](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-s390x)

## Changelog since v1.33.0-alpha.1

## Changes by Kind

### Deprecation

- The WatchFromStorageWithoutResourceVersion feature flag is deprecated and can no longer be enabled ([#129930](https://github.com/kubernetes/kubernetes/pull/129930), [@serathius](https://github.com/serathius)) [SIG API Machinery]

### API Change

- Added support for in-place vertical scaling of Pods with sidecars (containers defined within `initContainers` where the `restartPolicy` is Always). ([#128367](https://github.com/kubernetes/kubernetes/pull/128367), [@vivzbansal](https://github.com/vivzbansal)) [SIG API Machinery, Apps, CLI, Node, Scheduling and Testing]
- Kubectl: added alpha support for customizing kubectl behavior using preferences from a `kuberc` file (separate from kubeconfig). ([#125230](https://github.com/kubernetes/kubernetes/pull/125230), [@ardaguclu](https://github.com/ardaguclu)) [SIG API Machinery, CLI and Testing]

### Feature

- Added a `/statusz` endpoint for kube-controller-manager ([#128991](https://github.com/kubernetes/kubernetes/pull/128991), [@Henrywu573](https://github.com/Henrywu573)) [SIG API Machinery, Cloud Provider, Instrumentation and Testing]
- Fixed SELinuxWarningController defaults when running kube-controller-manager in a container. ([#130037](https://github.com/kubernetes/kubernetes/pull/130037), [@jsafrane](https://github.com/jsafrane)) [SIG Apps and Storage]
- Graduate BtreeWatchCache feature gate to GA ([#129934](https://github.com/kubernetes/kubernetes/pull/129934), [@serathius](https://github.com/serathius)) [SIG API Machinery]
- Introduced the `LegacySidecarContainers` feature gate enabling the legacy code path that predates the `SidecarContainers` feature. This temporary feature gate is disabled by default, only available in v1.33, and will be removed in v1.34. ([#130058](https://github.com/kubernetes/kubernetes/pull/130058), [@gjkim42](https://github.com/gjkim42)) [SIG Node]
- Kubeadm: 'kubeadm upgrade plan' now supports '--etcd-upgrade' flag to control whether the etcd upgrade plan should be displayed. Add an `EtcdUpgrade` field into `UpgradeConfiguration.Plan` for v1beta4. ([#130023](https://github.com/kubernetes/kubernetes/pull/130023), [@SataQiu](https://github.com/SataQiu)) [SIG Cluster Lifecycle]
- Kubeadm: added preflight check for `cp` on Linux nodes and `xcopy` on Windows nodes. These binaries are required for kubeadm to work properly. ([#130045](https://github.com/kubernetes/kubernetes/pull/130045), [@carlory](https://github.com/carlory)) [SIG Cluster Lifecycle]
- Kubeadm: improved `kubeadm init` and `kubeadm join` to provide consistent error messages when the kubelet failed or when failed to wait for control plane components. ([#130040](https://github.com/kubernetes/kubernetes/pull/130040), [@HirazawaUi](https://github.com/HirazawaUi)) [SIG Cluster Lifecycle]
- Kubeadm: promoted the feature gate `ControlPlaneKubeletLocalMode` to Beta. Kubeadm will per default use the local kube-apiserver endpoint for the kubelet when creating a cluster with "kubeadm init" or when joining control plane nodes with "kubeadm join". Enabling the feature gate also affects the `kubeadm init phase kubeconfig kubelet` phase, where the flag `--control-plane-endpoint` no longer affects the generated kubeconfig `Server` field, but the flag `--apiserver-advertise-address` can now be used for the same purpose. ([#129956](https://github.com/kubernetes/kubernetes/pull/129956), [@chrischdi](https://github.com/chrischdi)) [SIG Cluster Lifecycle]
- Kubernetes is now built with go 1.23.5 ([#129962](https://github.com/kubernetes/kubernetes/pull/129962), [@cpanato](https://github.com/cpanato)) [SIG Release and Testing]
- Kubernetes is now built with go 1.23.6 ([#130074](https://github.com/kubernetes/kubernetes/pull/130074), [@cpanato](https://github.com/cpanato)) [SIG Release and Testing]
- NodeRestriction admission now validates the audience value that kubelet is requesting a service account token for is part of the pod spec volume. The kube-apiserver featuregate `ServiceAccountNodeAudienceRestriction` is enabled by default in 1.33. ([#130017](https://github.com/kubernetes/kubernetes/pull/130017), [@aramase](https://github.com/aramase)) [SIG Auth]
- The nftables mode of kube-proxy is now GA. (The iptables mode remains the
  default; you can select the nftables mode by passing `--proxy-mode nftables`
  or using a config file with `mode: nftables`. See the kube-proxy documentation
  for more details.) ([#129653](https://github.com/kubernetes/kubernetes/pull/129653), [@danwinship](https://github.com/danwinship)) [SIG Network]
- `kubeproxy_conntrack_reconciler_deleted_entries_total` metric can be used to track cumulative sum of conntrack flows cleared by reconciler ([#130204](https://github.com/kubernetes/kubernetes/pull/130204), [@aroradaman](https://github.com/aroradaman)) [SIG Network]
- `kubeproxy_conntrack_reconciler_sync_duration_seconds` metric can be used to track conntrack reconciliation latency ([#130200](https://github.com/kubernetes/kubernetes/pull/130200), [@aroradaman](https://github.com/aroradaman)) [SIG Network]

### Bug or Regression

- Fix: adopt go1.23 behavior change in mount point parsing on Windows ([#129368](https://github.com/kubernetes/kubernetes/pull/129368), [@andyzhangx](https://github.com/andyzhangx)) [SIG Storage and Windows]
- Fixes a regression with the ServiceAccountNodeAudienceRestriction feature where `azureFile` volumes encounter "failed to get service accoount token attributes" errors ([#129993](https://github.com/kubernetes/kubernetes/pull/129993), [@aramase](https://github.com/aramase)) [SIG Auth and Testing]
- Kube-proxy: fixes a potential memory leak which can occur in clusters with high volume of UDP workflows ([#130032](https://github.com/kubernetes/kubernetes/pull/130032), [@aroradaman](https://github.com/aroradaman)) [SIG Network]
- Resolves a performance regression in default 1.31+ configurations, related to the ConsistentListFromCache feature, where rapid create / update API requests across different namespaces encounter increased latency. ([#130113](https://github.com/kubernetes/kubernetes/pull/130113), [@AwesomePatrol](https://github.com/AwesomePatrol)) [SIG API Machinery]
- The response from kube-apiserver /flagz endpoint would respond correctly with parsed flags value. ([#129996](https://github.com/kubernetes/kubernetes/pull/129996), [@yongruilin](https://github.com/yongruilin)) [SIG API Machinery, Architecture, Instrumentation and Testing]
- When cpu-manager-policy=static is configured containers meeting the qualifications for static cpu assignment (i.e. Containers with integer CPU `requests` in pods with `Guaranteed` QOS) will not have cfs quota enforced. Because this fix changes a long-established behavior, users observing a regressions can use the DisableCPUQuotaWithExclusiveCPUs feature gate (default on) to restore the old behavior. Please file an issue if you encounter problems and have to use the Feature Gate. ([#127525](https://github.com/kubernetes/kubernetes/pull/127525), [@scott-grimes](https://github.com/scott-grimes)) [SIG Node and Testing]

### Other (Cleanup or Flake)

- Flip StorageNamespaceIndex feature gate to false and deprecate it ([#129933](https://github.com/kubernetes/kubernetes/pull/129933), [@serathius](https://github.com/serathius)) [SIG Node]
- The SeparateCacheWatchRPC  feature gate is deprecated and disabled by default. ([#129929](https://github.com/kubernetes/kubernetes/pull/129929), [@serathius](https://github.com/serathius)) [SIG API Machinery]

## Dependencies

### Added
_Nothing has changed._

### Changed
- github.com/vishvananda/netlink: [b1ce50c → 62fb240](https://github.com/vishvananda/netlink/compare/b1ce50c...62fb240)

### Removed
_Nothing has changed._



# v1.33.0-alpha.1


## Downloads for v1.33.0-alpha.1



### Source Code

filename | sha512 hash
-------- | -----------
[kubernetes.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes.tar.gz) | 809c3565365eccf43761888113fe63c37a700edb6c662f4a29b93768d8d49d6c8ef052a6ffc41f61e9eecb22e006dc03c4399ad05886dc6a7635b2e573d0097d
[kubernetes-src.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-src.tar.gz) | 204a8f6723e8c0b0350994174b43f3a9272dacbd4f2992919b8ec95748df6af53dea385210b89417f1eeaa733732fee6c80559f0779f02f7cb73ccde6384bc9b

### Client Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-client-darwin-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-darwin-amd64.tar.gz) | 7762f1e33b94102a7fb943dfda3067e69ac534aeca040e95462781bd5973ee2436fe60c4ca2eeaea79f210a07c91167629d620bafc5b108839c02a4865ee0b64
[kubernetes-client-darwin-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-darwin-arm64.tar.gz) | ece5bda2f89981659957cc7bc40cd7db20283778c8f1755b9a21499057ec808708eeb7db3f195c0231ba43a0fd9165fb4bf6367183a486d82145414db2327790
[kubernetes-client-linux-386.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-linux-386.tar.gz) | 559689427abb113695ea3a1a1b3cbd388c0887dc8f775878337c1d413c1eb0fccfad161c9af23d7a40a0536b438bd800078fae182fcfde2905568ef4079b1062
[kubernetes-client-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-linux-amd64.tar.gz) | ba65065523407b5596a9efc53f7dd2e5e37b39c3968bbdb13a50944a80635dfc5903395741b5cb0f5f24482384788271fa1354b56f7f6b0b2f7482237aea8cc8
[kubernetes-client-linux-arm.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-linux-arm.tar.gz) | 585edd8319aec86378c16da7515f42fdcae5c618fba5dfba4af1455d5db8f5433fe16b95ff7193a2e648a847261ea51d3b412133459d33b48159ddf695a76f26
[kubernetes-client-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-linux-arm64.tar.gz) | 5d228232661dd237df57181920ee73008e1b28eda0366a85d125f569b15a21ebae8f9e2536b244908f9f82184e097b4ac9722863eed352cd0c957b7444bcc5fa
[kubernetes-client-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-linux-ppc64le.tar.gz) | 59e93927f46aff4f304ccad25a0d6220fa643c42c81b65015bd450d7615a809a8b4912efba0e66fe37f33def4b9fe77785ce43688582003c849377bde3277006
[kubernetes-client-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-linux-s390x.tar.gz) | 7c3bd8c464b0a46a216deb1144e3b042cc218464de6e418345a644024de09a04ec78e13a7c5a3f17d90ad9fda254482dd17d05ae67cd267ee2e0504da8258cf2
[kubernetes-client-windows-386.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-windows-386.tar.gz) | 0ea8503268858c551f9b9e51eb360cc160c76cb19c72c434df79ed421766bcb9addd33e6092525ab8e3556f217ae55dfc13f4506afd27585b5031118a6005403
[kubernetes-client-windows-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-windows-amd64.tar.gz) | f811e3c8e5b4fa31f9ae3493d757b4511de6cf0fc37a161da3c25f1503cf11149af6b79b9abf11314abf2e4cf410f1e41b10414981c141f702bec297a2beeae7
[kubernetes-client-windows-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-client-windows-arm64.tar.gz) | a8dfbb963a5d719dc8890ef14340ce35880e006955a229ff9204bb35da2a29df41b6797dc02269f2cc8de361014f8dd6b2535a9414359b48d820ff2cf536c4e1

### Server Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-server-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-server-linux-amd64.tar.gz) | daf5f5f38ab4357a724d688bfc33f3344f340fc4896d6d0c3da777beb76abe133707bbb6bd47cb954cd46bd62d5f4a7311fcaa5cd99f3389472d846c15d2e604
[kubernetes-server-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-server-linux-arm64.tar.gz) | 28d03d130e28eb7e812db35ca387eb515dfe8c21bbb2e7690285343d381ecd87828c0362ad19b3d13ec8d1d37763924cf9fdb1d814eb75d6e695322c27db06b4
[kubernetes-server-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-server-linux-ppc64le.tar.gz) | b479688f8aaa93d48d5809d21f21837b67144a5c115370f5154b9a13005f47e579f9f54b8f6d371e97165bd4f1a3d8eda85d2a37c83ac1615ca4dad7155d9a6e
[kubernetes-server-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-server-linux-s390x.tar.gz) | ed02308911595375b313b7df2fc6ad94b7dbcfc6f57fb0b9ced5512c4eca8f086852ea24bbfa7f3c146dc9cb98a1e5964dfc911dd46e41f815eeb884b82efdab

### Node Binaries

filename | sha512 hash
-------- | -----------
[kubernetes-node-linux-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-node-linux-amd64.tar.gz) | 846d0079fe2c53bdec279d6cc185f968cfed908762ce63c053830fdaeda78da4856f19253f98b908406694179da82dd2c387a4a08ad01d2522dc67832c7e2ac5
[kubernetes-node-linux-arm64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-node-linux-arm64.tar.gz) | c6b35f71acf7e9009ba1c6d274f1d2655039a0de59c0dd3f544bf240a8e74c43fa7bf830377f7d87dc14ce271e2f312a85930804ddd236a6877d13410131028e
[kubernetes-node-linux-ppc64le.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-node-linux-ppc64le.tar.gz) | c67735374d4f9062c495040c1bb28fc7f15362908d116542e663c58c900fc5e7939468118603d2233c8a951175484d839039f9d2ee1e0473e227fa994a391480
[kubernetes-node-linux-s390x.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-node-linux-s390x.tar.gz) | 2161369d2590959d8d28f81fa1d642028c816a4ce761d7af3d3edae369cda2a58fe8fa466d16e071d34148331ae572512421296ec53a1f5a1312a00376d67a01
[kubernetes-node-windows-amd64.tar.gz](https://dl.k8s.io/v1.33.0-alpha.1/kubernetes-node-windows-amd64.tar.gz) | f8051a237f06566e6bfd51881e1ae50a359b76dd5c8865ba6f3bf936e8be327a9a71d22192e252d49a2fb243be601fd2ceb17ea989b21e57c35f833e7b977341

### Container Images

All container images are available as manifest lists and support the described
architectures. It is also possible to pull a specific architecture directly by
adding the "-$ARCH" suffix  to the container image name.

name | architectures
---- | -------------
[registry.k8s.io/conformance:v1.33.0-alpha.1](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/conformance-s390x)
[registry.k8s.io/kube-apiserver:v1.33.0-alpha.1](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-apiserver-s390x)
[registry.k8s.io/kube-controller-manager:v1.33.0-alpha.1](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-controller-manager-s390x)
[registry.k8s.io/kube-proxy:v1.33.0-alpha.1](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-proxy-s390x)
[registry.k8s.io/kube-scheduler:v1.33.0-alpha.1](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kube-scheduler-s390x)
[registry.k8s.io/kubectl:v1.33.0-alpha.1](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl) | [amd64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-amd64), [arm64](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-arm64), [ppc64le](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-ppc64le), [s390x](https://console.cloud.google.com/artifacts/docker/k8s-artifacts-prod/southamerica-east1/images/kubectl-s390x)

## Changelog since v1.32.0

## Urgent Upgrade Notes

### (No, really, you MUST read this before you upgrade)

 - Action required for custom plugin developers.
  The `UpdatePodTolerations` action type is renamed to `UpdatePodToleration`, you have to follow the renaming if you're using it. ([#129023](https://github.com/kubernetes/kubernetes/pull/129023), [@zhifei92](https://github.com/zhifei92)) [SIG Scheduling and Testing]
 
## Changes by Kind

### API Change

- A new status field `.status.terminatingReplicas` is added to Deployments and ReplicaSets to allow tracking of terminating pods when the DeploymentPodReplacementPolicy feature-gate is enabled. ([#128546](https://github.com/kubernetes/kubernetes/pull/128546), [@atiratree](https://github.com/atiratree)) [SIG API Machinery, Apps and Testing]
- DRA API: the maximum number of pods which can use the same ResourceClaim is now 256 instead of 32. Beware that downgrading a cluster where this relaxed limit is in use to Kubernetes 1.32.0 is not supported because 1.32.0 would refuse to update ResourceClaims with more than 32 entries in the status.reservedFor field. ([#129543](https://github.com/kubernetes/kubernetes/pull/129543), [@pohly](https://github.com/pohly)) [SIG API Machinery, Node and Testing]
- DRA: CEL expressions using attribute strings exceeded the cost limit because their cost estimation was incomplete. ([#129661](https://github.com/kubernetes/kubernetes/pull/129661), [@pohly](https://github.com/pohly)) [SIG Node]
- DRA: when asking for "All" devices on a node, Kubernetes <= 1.32 proceeded to schedule pods onto nodes with no devices by not allocating any devices for those pods. Kubernetes 1.33 changes that to only picking nodes which have at least one device. Users who want the "proceed with scheduling also without devices" semantic can use the upcoming prioritized list feature with one sub-request for "all" devices and a second alternative with "count: 0". ([#129560](https://github.com/kubernetes/kubernetes/pull/129560), [@bart0sh](https://github.com/bart0sh)) [SIG API Machinery and Node]
- Graduate MultiCIDRServiceAllocator to stable and DisableAllocatorDualWrite to beta (disabled by default).
  Action required for Kubernetes distributions that manage the cluster Service CIDR.
  This feature allows users to define the cluster Service CIDR via a new API object: ServiceCIDR.
  Distributions or administrators of Kubernetes may want to control that new Service CIDRs added to the cluster does not overlap with other networks on the cluster, that only belong to a specific range of IPs or just simple retain the existing behavior of only having one ServiceCIDR per cluster.  An example of a Validation Admission Policy to achieve this is:
  
  ---
  apiVersion: admissionregistration.k8s.io/v1
  kind: ValidatingAdmissionPolicy
  metadata:
    name: "servicecidrs.default"
  spec:
    failurePolicy: Fail
    matchConstraints:
      resourceRules:
      - apiGroups:   ["networking.k8s.io"]
        apiVersions: ["v1","v1beta1"]
        operations:  ["CREATE", "UPDATE"]
        resources:   ["servicecidrs"]
    matchConditions:
    - name: 'exclude-default-servicecidr'
      expression: "object.metadata.name != 'kubernetes'"
    variables:
    - name: allowed
      expression: "['10.96.0.0/16','2001:db8::/64']"
    validations:
    - expression: "object.spec.cidrs.all(i , variables.allowed.exists(j , cidr(j).containsCIDR(i)))"
  ---
  apiVersion: admissionregistration.k8s.io/v1
  kind: ValidatingAdmissionPolicyBinding
  metadata:
    name: "servicecidrs-binding"
  spec:
    policyName: "servicecidrs.default"
    validationActions: [Deny,Audit]
  --- ([#128971](https://github.com/kubernetes/kubernetes/pull/128971), [@aojea](https://github.com/aojea)) [SIG Apps, Architecture, Auth, CLI, Etcd, Network, Release and Testing]
- Kubenetes starts validating NodeSelectorRequirement's values when creating pods. ([#128212](https://github.com/kubernetes/kubernetes/pull/128212), [@AxeZhan](https://github.com/AxeZhan)) [SIG Apps and Scheduling]
- Kubernetes components that accept x509 client certificate authentication now read the user UID from a certificate subject name RDN with object id 1.3.6.1.4.1.57683.2. An RDN with this object id must contain a string value, and appear no more than once in the certificate subject. Reading the user UID from this RDN can be disabled by setting the beta feature gate `AllowParsingUserUIDFromCertAuth` to false (until the feature gate graduates to GA). ([#127897](https://github.com/kubernetes/kubernetes/pull/127897), [@modulitos](https://github.com/modulitos)) [SIG API Machinery, Auth and Testing]
- Removed general available feature-gate `PDBUnhealthyPodEvictionPolicy`. ([#129500](https://github.com/kubernetes/kubernetes/pull/129500), [@carlory](https://github.com/carlory)) [SIG API Machinery, Apps and Auth]
- `kubectl apply` now coerces `null` values for labels and annotations in manifests to empty string values, consistent with typed JSON metadata decoding, rather than dropping all labels and annotations ([#129257](https://github.com/kubernetes/kubernetes/pull/129257), [@liggitt](https://github.com/liggitt)) [SIG API Machinery]

### Feature

- Add unit test helpers to validate CEL and patterns in CustomResourceDefinitions. ([#129028](https://github.com/kubernetes/kubernetes/pull/129028), [@sttts](https://github.com/sttts)) [SIG API Machinery]
- Added a `/flagz` endpoint for kube-proxy ([#128985](https://github.com/kubernetes/kubernetes/pull/128985), [@yongruilin](https://github.com/yongruilin)) [SIG Instrumentation and Network]
- Added a `/status` endpoint for kube-proxy ([#128989](https://github.com/kubernetes/kubernetes/pull/128989), [@Henrywu573](https://github.com/Henrywu573)) [SIG Instrumentation and Network]
- Added e2e tests for volume group snapshots. ([#128972](https://github.com/kubernetes/kubernetes/pull/128972), [@manishym](https://github.com/manishym)) [SIG Cloud Provider, Storage and Testing]
- Adds a /flagz endpoint for kube-scheduler endpoint ([#128818](https://github.com/kubernetes/kubernetes/pull/128818), [@yongruilin](https://github.com/yongruilin)) [SIG Architecture, Instrumentation, Scheduling and Testing]
- Adds a /statusz endpoint for kubelet endpoint ([#128811](https://github.com/kubernetes/kubernetes/pull/128811), [@zhifei92](https://github.com/zhifei92)) [SIG Architecture, Instrumentation and Node]
- Bugfix: Ensure container-level swap metrics are collected ([#129486](https://github.com/kubernetes/kubernetes/pull/129486), [@iholder101](https://github.com/iholder101)) [SIG Node and Testing]
- Calculated pod resources are now cached when adding pods to NodeInfo in the scheduler framework, improving performance when processing unschedulable pods. ([#129635](https://github.com/kubernetes/kubernetes/pull/129635), [@macsko](https://github.com/macsko)) [SIG Scheduling]
- Cel-go has been bumped to v0.23.2. ([#129844](https://github.com/kubernetes/kubernetes/pull/129844), [@cici37](https://github.com/cici37)) [SIG API Machinery, Auth, Cloud Provider and Node]
- Client-go/rest: fully supports contextual logging. BackoffManagerWithContext should be used instead of BackoffManager to ensure that the caller can interrupt the sleep. ([#127709](https://github.com/kubernetes/kubernetes/pull/127709), [@pohly](https://github.com/pohly)) [SIG API Machinery, Architecture, Auth, Cloud Provider, Instrumentation, Network and Node]
- Graduated the `KubeletFineGrainedAuthz` feature gate to beta; the gate is now enabled by default. ([#129656](https://github.com/kubernetes/kubernetes/pull/129656), [@vinayakankugoyal](https://github.com/vinayakankugoyal)) [SIG Auth, CLI, Node, Storage and Testing]
- Improved scheduling performance of pods with required topology spreading. ([#129119](https://github.com/kubernetes/kubernetes/pull/129119), [@macsko](https://github.com/macsko)) [SIG Scheduling]
- Kube-apiserver: Promoted the  `ServiceAccountTokenNodeBinding` feature gate general availability. It is now locked to enabled. ([#129591](https://github.com/kubernetes/kubernetes/pull/129591), [@liggitt](https://github.com/liggitt)) [SIG Auth and Testing]
- Kube-proxy extends the schema of its healthz/ and livez/ endpoints to incorporate information about the corresponding IP family ([#129271](https://github.com/kubernetes/kubernetes/pull/129271), [@aroradaman](https://github.com/aroradaman)) [SIG Network and Windows]
- Kubeadm: graduated the WaitForAllControlPlaneComponents feature gate to Beta. When checking the health status of a control plane component, make sure that the address and port defined as arguments in the respective component's static Pod manifest are used. ([#129620](https://github.com/kubernetes/kubernetes/pull/129620), [@neolit123](https://github.com/neolit123)) [SIG Cluster Lifecycle]
- Kubeadm: if the `NodeLocalCRISocket` feature gate is enabled, remove the `kubeadm.alpha.kubernetes.io/cri-socket` annotation from a given node on `kubeadm upgrade`. ([#129279](https://github.com/kubernetes/kubernetes/pull/129279), [@HirazawaUi](https://github.com/HirazawaUi)) [SIG Cluster Lifecycle and Testing]
- Kubeadm: if the `NodeLocalCRISocket` feature gate is enabled, remove the flag `--container-runtime-endpoint` from the `/var/lib/kubelet/kubeadm-flags.env` file on `kubeadm upgrade`. ([#129278](https://github.com/kubernetes/kubernetes/pull/129278), [@HirazawaUi](https://github.com/HirazawaUi)) [SIG Cluster Lifecycle]
- Kubeadm: promoted the feature gate `ControlPlaneKubeletLocalMode` to Beta. Kubeadm will per default use the local kube-apiserver endpoint for the kubelet when creating a cluster with "kubeadm init" or when joining control plane nodes with "kubeadm join". Enabling the feature gate also affects the `kubeadm init phase kubeconfig kubelet` phase, where the flag `--control-plane-endpoint` no longer affects the generated kubeconfig `Server` field, but the flag `--apiserver-advertise-address` can now be used for the same purpose. ([#129956](https://github.com/kubernetes/kubernetes/pull/129956), [@chrischdi](https://github.com/chrischdi)) [SIG Cluster Lifecycle]
- Kubeadm: removed preflight check for nsenter on Linux nodes
  kubeadm: added preflight check for `losetup` on Linux nodes. It's required by kubelet for keeping a block device opened. ([#129450](https://github.com/kubernetes/kubernetes/pull/129450), [@carlory](https://github.com/carlory)) [SIG Cluster Lifecycle]
- Kubeadm: removed the feature gate EtcdLearnerMode which graduated to GA in 1.32. ([#129589](https://github.com/kubernetes/kubernetes/pull/129589), [@neolit123](https://github.com/neolit123)) [SIG Cluster Lifecycle]
- Kubernetes is now built with go 1.23.4 ([#129422](https://github.com/kubernetes/kubernetes/pull/129422), [@cpanato](https://github.com/cpanato)) [SIG Release and Testing]
- Kubernetes is now built with go 1.23.5 ([#129962](https://github.com/kubernetes/kubernetes/pull/129962), [@cpanato](https://github.com/cpanato)) [SIG Release and Testing]
- Promoted the feature gate `CSIMigrationPortworx` to GA. If your applications are using Portworx volumes, please make sure that the corresponding Portworx CSI driver is installed on your cluster **before** upgrading to 1.31 or later because all operations for the in-tree `portworxVolume` type are redirected to the pxd.portworx.com CSI driver when the feature gate is enabled. ([#129297](https://github.com/kubernetes/kubernetes/pull/129297), [@gohilankit](https://github.com/gohilankit)) [SIG Storage]
- The `SidecarContainers` feature has graduated to GA. 'SidecarContainers' feature gate was locked to default value and will be removed in v1.36. If you were setting this feature gate explicitly, please remove it now. ([#129731](https://github.com/kubernetes/kubernetes/pull/129731), [@gjkim42](https://github.com/gjkim42)) [SIG Apps, Node, Scheduling and Testing]
- Upgrade autoscalingv1 to autoscalingv2 in kubectl autoscale cmd, The cmd will attempt to use the autoscaling/v2 API first. If the autoscaling/v2 API is not available or an error occurs, it will fall back to the autoscaling/v1 API. ([#128950](https://github.com/kubernetes/kubernetes/pull/128950), [@googs1025](https://github.com/googs1025)) [SIG Autoscaling and CLI]
- Validate ContainerLogMaxFiles in kubelet config validation ([#129072](https://github.com/kubernetes/kubernetes/pull/129072), [@kannon92](https://github.com/kannon92)) [SIG Node]

### Documentation

- Give example of set-based requirement for -l/--selector flag ([#129106](https://github.com/kubernetes/kubernetes/pull/129106), [@rotsix](https://github.com/rotsix)) [SIG CLI]
- Kubeadm: improved the `kubeadm reset` message for manual cleanups and referenced https://k8s.io/docs/reference/setup-tools/kubeadm/kubeadm-reset/. ([#129644](https://github.com/kubernetes/kubernetes/pull/129644), [@neolit123](https://github.com/neolit123)) [SIG Cluster Lifecycle]

### Bug or Regression

- --feature-gate=InOrderInformers (default on), causes informers to process watch streams in order as opposed to grouping updates for the same item close together.  Binaries embedding client-go, but not wiring the featuregates can disable by setting the `KUBE_FEATURE_InOrderInformers=false`. ([#129568](https://github.com/kubernetes/kubernetes/pull/129568), [@deads2k](https://github.com/deads2k)) [SIG API Machinery]
- Adding a validation for revisionHistoryLimit field in statefulset.spec to prevent it being set to negative value. ([#129017](https://github.com/kubernetes/kubernetes/pull/129017), [@ardaguclu](https://github.com/ardaguclu)) [SIG Apps]
- DRA: the explanation for why a pod which wasn't using ResourceClaims was unscheduleable included a useless "no new claims to deallocate" when it was unscheduleable for some other reasons. ([#129823](https://github.com/kubernetes/kubernetes/pull/129823), [@googs1025](https://github.com/googs1025)) [SIG Node and Scheduling]
- Enables ratcheting validation on status subresources for CustomResourceDefinitions ([#129506](https://github.com/kubernetes/kubernetes/pull/129506), [@JoelSpeed](https://github.com/JoelSpeed)) [SIG API Machinery]
- Fix the issue where the named ports exposed by restartable init containers (a.k.a. sidecar containers) cannot be accessed using a Service. ([#128850](https://github.com/kubernetes/kubernetes/pull/128850), [@toVersus](https://github.com/toVersus)) [SIG Network and Testing]
- Fixed `kubectl wait --for=create` behavior with label selectors, to properly wait for resources with matching labels to appear. ([#128662](https://github.com/kubernetes/kubernetes/pull/128662), [@omerap12](https://github.com/omerap12)) [SIG CLI and Testing]
- Fixed a bug where adding an ephemeral container to a pod which references a new secret or config map doesn't give the pod access to that new secret or config map. (#114984, @cslink) ([#129670](https://github.com/kubernetes/kubernetes/pull/129670), [@cslink](https://github.com/cslink)) [SIG Auth]
- Fixed a data race that could occur when a single Go type was serialized to CBOR concurrently for the first time within a program. ([#129170](https://github.com/kubernetes/kubernetes/pull/129170), [@benluddy](https://github.com/benluddy)) [SIG API Machinery]
- Fixed a storage bug around multipath. iSCSI and Fibre Channel devices attached to nodes via multipath now resolve correctly if partitioned. ([#128086](https://github.com/kubernetes/kubernetes/pull/128086), [@RomanBednar](https://github.com/RomanBednar)) [SIG Storage]
- Fixed in-tree to CSI migration for Portworx volumes, in clusters where Portworx security feature is enabled (it's a Portworx feature, not Kubernetes feature). It required secret data from the secret mentioned in-tree SC, to be passed in CSI requests which was not happening before this fix. ([#129630](https://github.com/kubernetes/kubernetes/pull/129630), [@gohilankit](https://github.com/gohilankit)) [SIG Storage]
- Fixed: kube-proxy EndpointSliceCache memory is leaked ([#128929](https://github.com/kubernetes/kubernetes/pull/128929), [@orange30](https://github.com/orange30)) [SIG Network]
- Fixes CVE-2024-51744 ([#128621](https://github.com/kubernetes/kubernetes/pull/128621), [@kmala](https://github.com/kmala)) [SIG Auth, Cloud Provider and Node]
- Fixes a panic in kube-controller-manager handling StatefulSet objects when revisionHistoryLimit is negative ([#129301](https://github.com/kubernetes/kubernetes/pull/129301), [@ardaguclu](https://github.com/ardaguclu)) [SIG Apps]
- HPA's with ContainerResource metrics will no longer error when container metrics are missing, instead they will use the same logic Resource metrics are using to make calculations ([#127193](https://github.com/kubernetes/kubernetes/pull/127193), [@DP19](https://github.com/DP19)) [SIG Apps and Autoscaling]
- Implemented logging and event recording for probe results with an `Unknown` status in the kubelet's prober module. This helps in better diagnosing and monitoring cases where container probes return an `Unknown` result, improving the observability and reliability of health checks. ([#125901](https://github.com/kubernetes/kubernetes/pull/125901), [@jralmaraz](https://github.com/jralmaraz)) [SIG Node]
- Improved reboot event reporting. The kubelet will only emit one reboot Event when a server-level reboot is detected, even if the kubelet cannot write its status to the associated Node (which triggers a retry). ([#129151](https://github.com/kubernetes/kubernetes/pull/129151), [@rphillips](https://github.com/rphillips)) [SIG Node]
- Kube-apiserver: --service-account-max-token-expiration can now be used in combination with an external token signer --service-account-signing-endpoint, as long as the --service-account-max-token-expiration is not longer than the external token signer's max expiration. ([#129816](https://github.com/kubernetes/kubernetes/pull/129816), [@sambdavidson](https://github.com/sambdavidson)) [SIG API Machinery and Auth]
- Kubeadm: avoid loading the file passed to `--kubeconfig` during `kubeadm init` phases more than once. ([#129006](https://github.com/kubernetes/kubernetes/pull/129006), [@kokes](https://github.com/kokes)) [SIG Cluster Lifecycle]
- Kubeadm: fix a bug where the 'node.skipPhases' in UpgradeConfiguration is not respected by 'kubeadm upgrade node' command ([#129452](https://github.com/kubernetes/kubernetes/pull/129452), [@SataQiu](https://github.com/SataQiu)) [SIG Cluster Lifecycle]
- Kubeadm: fixed a bug where an image is not pulled if there is an error with the sandbox image from CRI. ([#129594](https://github.com/kubernetes/kubernetes/pull/129594), [@neolit123](https://github.com/neolit123)) [SIG Cluster Lifecycle]
- Kubeadm: fixed the bug where the v1beta4 Timeouts.EtcdAPICall field was not respected in etcd client operations, and the default timeout of 2 minutes was always used. ([#129859](https://github.com/kubernetes/kubernetes/pull/129859), [@neolit123](https://github.com/neolit123)) [SIG Cluster Lifecycle]
- Kubeadm: if an addon is disabled in the ClusterConfiguration, skip it during upgrade. ([#129418](https://github.com/kubernetes/kubernetes/pull/129418), [@neolit123](https://github.com/neolit123)) [SIG Cluster Lifecycle]
- Kubeadm: run kernel version and OS version preflight checks on `kubeadm upgrade`. ([#129401](https://github.com/kubernetes/kubernetes/pull/129401), [@pacoxu](https://github.com/pacoxu)) [SIG Cluster Lifecycle]
- Provides an additional function argument to directly specify the version for the tools that the consumers wishes to use ([#129658](https://github.com/kubernetes/kubernetes/pull/129658), [@unmarshall](https://github.com/unmarshall)) [SIG API Machinery]
- Remove the limitation on exposing port 10250 externally in service. ([#129174](https://github.com/kubernetes/kubernetes/pull/129174), [@RyanAoh](https://github.com/RyanAoh)) [SIG Apps and Network]
- This PR changes the signature of the `PublishResources` to accept a `resourceslice.DriverResources` parameter instead of a `Resources` parameter. ([#129142](https://github.com/kubernetes/kubernetes/pull/129142), [@googs1025](https://github.com/googs1025)) [SIG Node and Testing]
- [kubectl] Improved the describe output for projected volume sources to clearly indicate whether Secret and ConfigMap entries are optional. ([#129457](https://github.com/kubernetes/kubernetes/pull/129457), [@gshaibi](https://github.com/gshaibi)) [SIG CLI]

### Other (Cleanup or Flake)

- Implemented scheduler_cache_size metric.
  Also, scheduler_scheduler_cache_size metric is deprecated in favor of scheduler_cache_size, 
  and will be removed at v1.34. ([#128810](https://github.com/kubernetes/kubernetes/pull/128810), [@googs1025](https://github.com/googs1025)) [SIG Scheduling]
- Kube-apiserver: inactive serving code is removed for authentication.k8s.io/v1alpha1 APIs ([#129186](https://github.com/kubernetes/kubernetes/pull/129186), [@liggitt](https://github.com/liggitt)) [SIG Auth and Testing]
- Kube-proxy extends the schema of its metrics/ endpoints to incorporate information about the corresponding IP family ([#129173](https://github.com/kubernetes/kubernetes/pull/129173), [@aroradaman](https://github.com/aroradaman)) [SIG Network and Windows]
- Kube-proxy nftables logs the failed transactions and the full table when using log level 4 or higher. Logging is rate limited to one entry every 24 hours to avoid performance issues. ([#128886](https://github.com/kubernetes/kubernetes/pull/128886), [@npinaeva](https://github.com/npinaeva)) [SIG Network]
- Kubeadm: removed preflight check for `ip`, `iptables`, `ethtool` and `tc` on Linux nodes. kubelet and kube-proxy will continue to report `iptables` errors if its usage is required. The tools `ip`, `ethtool` and `tc` had legacy usage in the kubelet but are no longer required. ([#129131](https://github.com/kubernetes/kubernetes/pull/129131), [@pacoxu](https://github.com/pacoxu)) [SIG Cluster Lifecycle]
- Kubeadm: removed preflight check for `touch` on Linux nodes. ([#129317](https://github.com/kubernetes/kubernetes/pull/129317), [@carlory](https://github.com/carlory)) [SIG Cluster Lifecycle]
- NOE ([#128856](https://github.com/kubernetes/kubernetes/pull/128856), [@adrianmoisey](https://github.com/adrianmoisey)) [SIG Apps and Network]
- Removed generally available feature gate `KubeProxyDrainingTerminatingNodes`. ([#129692](https://github.com/kubernetes/kubernetes/pull/129692), [@alexanderConstantinescu](https://github.com/alexanderConstantinescu)) [SIG Network]
- Removed support for v1alpha1 version of ValidatingAdmissionPolicy and ValidatingAdmissionPolicyBinding API kinds. ([#129207](https://github.com/kubernetes/kubernetes/pull/129207), [@Jefftree](https://github.com/Jefftree)) [SIG Etcd and Testing]
- The deprecated pod_scheduling_duration_seconds metric is removed.
  You can migrate to pod_scheduling_sli_duration_seconds. ([#128906](https://github.com/kubernetes/kubernetes/pull/128906), [@sanposhiho](https://github.com/sanposhiho)) [SIG Instrumentation and Scheduling]
- This renames some coredns metrics, see https://github.com/coredns/coredns/blob/v1.11.0/plugin/forward/README.md#metrics. ([#129175](https://github.com/kubernetes/kubernetes/pull/129175), [@DamianSawicki](https://github.com/DamianSawicki)) [SIG Cloud Provider]
- This renames some coredns metrics, see https://github.com/coredns/coredns/blob/v1.11.0/plugin/forward/README.md#metrics. ([#129232](https://github.com/kubernetes/kubernetes/pull/129232), [@DamianSawicki](https://github.com/DamianSawicki)) [SIG Cloud Provider]
- Updated CNI plugins to v1.6.2. ([#129776](https://github.com/kubernetes/kubernetes/pull/129776), [@saschagrunert](https://github.com/saschagrunert)) [SIG Cloud Provider, Node and Testing]
- Updated cri-tools to v1.32.0. ([#129116](https://github.com/kubernetes/kubernetes/pull/129116), [@saschagrunert](https://github.com/saschagrunert)) [SIG Cloud Provider]
- Upgrade CoreDNS to v1.12.0 ([#128926](https://github.com/kubernetes/kubernetes/pull/128926), [@bzsuni](https://github.com/bzsuni)) [SIG Cloud Provider and Cluster Lifecycle]

## Dependencies

### Added
- gopkg.in/go-jose/go-jose.v2: v2.6.3

### Changed
- cel.dev/expr: v0.18.0 → v0.19.1
- github.com/coredns/corefile-migration: [v1.0.24 → v1.0.25](https://github.com/coredns/corefile-migration/compare/v1.0.24...v1.0.25)
- github.com/coreos/go-oidc: [v2.2.1+incompatible → v2.3.0+incompatible](https://github.com/coreos/go-oidc/compare/v2.2.1...v2.3.0)
- github.com/cyphar/filepath-securejoin: [v0.3.4 → v0.3.5](https://github.com/cyphar/filepath-securejoin/compare/v0.3.4...v0.3.5)
- github.com/davecgh/go-spew: [d8f796a → v1.1.1](https://github.com/davecgh/go-spew/compare/d8f796a...v1.1.1)
- github.com/golang-jwt/jwt/v4: [v4.5.0 → v4.5.1](https://github.com/golang-jwt/jwt/compare/v4.5.0...v4.5.1)
- github.com/google/btree: [v1.0.1 → v1.1.3](https://github.com/google/btree/compare/v1.0.1...v1.1.3)
- github.com/google/cel-go: [v0.22.0 → v0.23.2](https://github.com/google/cel-go/compare/v0.22.0...v0.23.2)
- github.com/google/gnostic-models: [v0.6.8 → v0.6.9](https://github.com/google/gnostic-models/compare/v0.6.8...v0.6.9)
- github.com/pmezard/go-difflib: [5d4384e → v1.0.0](https://github.com/pmezard/go-difflib/compare/5d4384e...v1.0.0)
- golang.org/x/crypto: v0.28.0 → v0.31.0
- golang.org/x/net: v0.30.0 → v0.33.0
- golang.org/x/sync: v0.8.0 → v0.10.0
- golang.org/x/sys: v0.26.0 → v0.28.0
- golang.org/x/term: v0.25.0 → v0.27.0
- golang.org/x/text: v0.19.0 → v0.21.0
- k8s.io/kube-openapi: 32ad38e → 2c72e55
- sigs.k8s.io/apiserver-network-proxy/konnectivity-client: v0.31.0 → v0.31.1
- sigs.k8s.io/kustomize/api: v0.18.0 → v0.19.0
- sigs.k8s.io/kustomize/cmd/config: v0.15.0 → v0.19.0
- sigs.k8s.io/kustomize/kustomize/v5: v5.5.0 → v5.6.0
- sigs.k8s.io/kustomize/kyaml: v0.18.1 → v0.19.0

### Removed
- github.com/asaskevich/govalidator: [f61b66f](https://github.com/asaskevich/govalidator/tree/f61b66f)
- gopkg.in/square/go-jose.v2: v2.6.0