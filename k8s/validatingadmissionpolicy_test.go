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
	"embed"
	"encoding/json"
	"reflect"
	"testing"
)

//go:embed testdata
var testdata embed.FS

func testfile(name string) string {
	return "testdata/" + name
}

func readTestData(policy, original, updated string) (policyData, originalData, updatedData []byte, err error) {
	policyData, err = testdata.ReadFile(testfile(policy))
	if err == nil {
		if original != "" {
			originalData, err = testdata.ReadFile(testfile(original))
		}
		if err == nil {
			if updated != "" {
				updatedData, err = testdata.ReadFile(testfile(updated))
			}
		}
	}
	return
}

func TestEval(t *testing.T) {
	tests := []struct {
		name     string
		policy   string
		orig     string
		updated  string
		expected EvalResponse
		wantErr  bool
	}{
		{
			name:    "test an expression which should fail",
			policy:  "policy1.yaml",
			orig:    "",
			updated: "updated1.yaml",
			expected: EvalResponse{
				Validations: []EvalValidation{{Result: false}},
			},
		},
		{
			name:    "test an expression which should succeed",
			policy:  "policy2.yaml",
			orig:    "",
			updated: "updated2.yaml",
			expected: EvalResponse{
				Validations: []EvalValidation{{Result: true}},
			},
		},
		{
			name:    "test an expression with variables, expression should fail",
			policy:  "variable1 policy.yaml",
			orig:    "",
			updated: "variable1 updated.yaml",
			expected: EvalResponse{
				Variables: []EvalVariable{{
					Name:  "foo",
					Value: "default",
				}},
				Validations: []EvalValidation{{Result: false}},
			},
		},
		{
			name:    "test an expression with variables, expression should succeed",
			policy:  "variable2 policy.yaml",
			orig:    "",
			updated: "variable2 updated.yaml",
			expected: EvalResponse{
				Variables: []EvalVariable{{
					Name:  "foo",
					Value: "bar",
				}},
				Validations: []EvalValidation{{Result: true}},
			},
		},
		{
			name:    "test an expression with variables evaluating to a map, expression should succeed",
			policy:  "variable3 policy.yaml",
			orig:    "",
			updated: "variable3 updated.yaml",
			expected: EvalResponse{
				Variables: []EvalVariable{{
					Name: "labels",
					Value: map[string]any{
						"app": "kubernetes-bootcamp",
						"foo": "bar",
					},
				}},
				Validations: []EvalValidation{{Result: true}},
			},
		},
		{
			name:    "test an expression with variables evaluating to query parameters in a URL, expression should succeed",
			policy:  "variable4 policy.yaml",
			orig:    "",
			updated: "variable4 updated.yaml",
			expected: EvalResponse{
				Variables: []EvalVariable{{
					Name: "foo",
					Value: map[string]any{
						"query": []any{"val"},
					},
				}},
				Validations: []EvalValidation{{Result: true}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, orig, updated, err := readTestData(tt.policy, tt.orig, tt.updated)
			var results string
			if err == nil {
				results, err = EvalValidatingAdmissionPolicy(policy, orig, updated)
			}
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				evalResponse := EvalResponse{}
				if err := json.Unmarshal([]byte(results), &evalResponse); err != nil {
					t.Errorf("Eval() error = %v", err)
				}
				if !reflect.DeepEqual(tt.expected, evalResponse) {
					t.Errorf("Expected %v\n, received %v", tt.expected, evalResponse)
				}
			}
		})
	}
}
