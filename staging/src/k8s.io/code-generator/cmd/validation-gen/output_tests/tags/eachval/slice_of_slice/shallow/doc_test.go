/*
Copyright 2024 The Kubernetes Authors.

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

package shallow

import (
	"testing"
)

func Test(t *testing.T) {
	st := localSchemeBuilder.Test(t)

	st.Value(&Struct{
		// All zero values.
	}).ExpectValid()

	st.Value(&Struct{
		ListField:        [][]string{make([]string, 0), make([]string, 0)},
		ListPtrField:     [][]*string{make([]*string, 0), make([]*string, 0)},
		ListTypedefField: []SliceType{make(SliceType, 0), make(SliceType, 0)},
	}).ExpectValidateFalseByPath(map[string][]string{
		"listField[0]":        {"field Struct.ListField[*]"},
		"listField[1]":        {"field Struct.ListField[*]"},
		"listPtrField[0]":     {"field Struct.ListPtrField[*]"},
		"listPtrField[1]":     {"field Struct.ListPtrField[*]"},
		"listTypedefField[0]": {"field Struct.ListTypedefField[*]"},
		"listTypedefField[1]": {"field Struct.ListTypedefField[*]"},
	})
}
