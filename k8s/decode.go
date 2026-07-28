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

import (
	"strings"

	apimachineryjson "k8s.io/apimachinery/pkg/util/json"
	sigsyaml "sigs.k8s.io/yaml"
)

// decodeObjectInput parses user-supplied resource YAML (the object/oldObject
// tabs) into the map[string]any the CEL activation binds, reproducing exactly
// what a real cluster does to the document before admission control ever sees
// it. kubectl converts YAML to JSON client-side with
// sigs.k8s.io/yaml.YAMLToJSON -- which is YAML 1.1, so bare `on`/`off`/`yes`/`no`
// keys and values coerce to booleans (and, because `on` and `yes` both become
// the key `true`, one such annotation is silently dropped, just as a cluster
// drops it). The apiserver then decodes that JSON with apimachinery's util/json,
// which keeps integral numbers as int64 rather than float64, so CEL sees them as
// `int` and integer arithmetic, `%` and `type(x) == int` behave as they do on a
// cluster. An empty or whitespace-only input yields a nil map, mirroring an
// empty editor tab in the playground UI.
//
// The vap and webhook modes share this decoder so that both modes show what a
// cluster would do with the same document.
func decodeObjectInput(input []byte) (map[string]any, error) {
	if len(strings.TrimSpace(string(input))) == 0 {
		return nil, nil
	}
	jsonBytes, err := sigsyaml.YAMLToJSON(input)
	if err != nil {
		return nil, err
	}
	var decoded map[string]any
	if err := apimachineryjson.Unmarshal(jsonBytes, &decoded); err != nil {
		return nil, err
	}
	return decoded, nil
}
