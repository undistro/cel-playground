// Copyright 2023 Undistro Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s

import "testing"

func TestDecodeObjectInputIsClusterFaithful(t *testing.T) {
	input := []byte(`
metadata:
  annotations:
    on: keepme
enabled: yes
replicas: 5
`)
	decoded, err := decodeObjectInput(input)
	if err != nil {
		t.Fatalf("decodeObjectInput() error: %v", err)
	}

	annotations := decoded["metadata"].(map[string]any)["annotations"].(map[string]any)
	// The bare key `on` is coerced to the boolean key `true`, exactly as a
	// cluster does; yaml.v3 would have preserved the string key "on".
	if _, ok := annotations["on"]; ok {
		t.Errorf(`annotation key "on" survived as a string; want it coerced to "true"`)
	}
	if got, ok := annotations["true"]; !ok || got != "keepme" {
		t.Errorf(`annotations["true"] = %v (present=%v), want "keepme"`, got, ok)
	}

	// The bare value `yes` is coerced to a bool, not kept as the string "yes".
	if got := decoded["enabled"]; got != true {
		t.Errorf(`decoded["enabled"] = %#v (%T), want bool true`, got, got)
	}

	// Integral numbers stay int64 (CEL int), not float64, so integer
	// arithmetic and `%` keep working as they do on a cluster.
	if got := decoded["replicas"]; got != int64(5) {
		t.Errorf(`decoded["replicas"] = %#v (%T), want int64(5)`, got, got)
	}
}
