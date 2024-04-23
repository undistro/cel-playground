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

package k8s_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/undistro/cel-playground/k8s"
)

func vapTestfile(file string) string {
	return testfile("vap/" + file)
}

func readValidationTestData(policy, original, updated, namespace, request, authorizer string) (policyData, originalData, updatedData, namespaceData, requestData, authorizerData []byte, err error) {
	policyData, err = testdata.ReadFile(vapTestfile(policy))
	if err == nil && original != "" {
		originalData, err = testdata.ReadFile(vapTestfile(original))
	}
	if err == nil && updated != "" {
		updatedData, err = testdata.ReadFile(vapTestfile(updated))
	}
	if err == nil && namespace != "" {
		namespaceData, err = testdata.ReadFile(vapTestfile(namespace))
	}
	if err == nil && request != "" {
		requestData, err = testdata.ReadFile(vapTestfile(request))
	}
	if err == nil && authorizer != "" {
		authorizerData, err = testdata.ReadFile(vapTestfile(authorizer))
	}
	return
}

func TestValidationEval(t *testing.T) {
	tests := []struct {
		name       string
		policy     string
		orig       string
		updated    string
		namespace  string
		request    string
		authorizer string
		expected   k8s.EvalResponse
		wantErr    bool
	}{{
		name:    "test an expression which should fail",
		policy:  "policy1.yaml",
		orig:    "",
		updated: "updated1.yaml",
		expected: k8s.EvalResponse{
			Validations: []*k8s.EvalResult{{Message: "All production deployments should be HA with at least three replicas", Result: false, Cost: uint64ptr(4)}},
			Cost:        uint64ptr(4),
		},
	}, {
		name:    "test an expression which should succeed",
		policy:  "policy2.yaml",
		orig:    "",
		updated: "updated2.yaml",
		expected: k8s.EvalResponse{
			Validations: []*k8s.EvalResult{{Result: true, Cost: uint64ptr(4)}},
			Cost:        uint64ptr(4),
		},
	}, {
		name:    "test an expression with variables, expression should fail with no audit annotation",
		policy:  "variable1 policy.yaml",
		orig:    "",
		updated: "variable1 updated.yaml",
		expected: k8s.EvalResponse{
			ValidationVariables: []*k8s.EvalVariable{{
				Name:  "foo",
				Value: "default",
				Cost:  uint64ptr(6),
			}},
			Validations: []*k8s.EvalResult{{Result: false, Cost: uint64ptr(2)}},
			Cost:        uint64ptr(8),
		},
	}, {
		name:    "test an expression with variables, expression should succeed with audit annotation",
		policy:  "variable2 policy.yaml",
		orig:    "",
		updated: "variable2 updated.yaml",
		expected: k8s.EvalResponse{
			ValidationVariables: []*k8s.EvalVariable{{
				Name:  "foo",
				Value: "bar",
				Cost:  uint64ptr(11),
			}},
			Validations: []*k8s.EvalResult{{
				Result: true,
				Cost:   uint64ptr(2),
			}},
			AuditAnnotations: []*k8s.EvalResult{{
				Name:    strptr("foo-label"),
				Message: "Label for foo is set to bar",
				Cost:    uint64ptr(2),
			}},
			Cost: uint64ptr(15),
		},
	}, {
		name:    "test an expression with variables evaluating to a map, expression should succeed",
		policy:  "variable3 policy.yaml",
		orig:    "",
		updated: "variable3 updated.yaml",
		expected: k8s.EvalResponse{
			ValidationVariables: []*k8s.EvalVariable{{
				Name: "labels",
				Value: map[string]any{
					"app": "kubernetes-bootcamp",
					"foo": "bar",
				},
				Cost: uint64ptr(5),
			}},
			Validations: []*k8s.EvalResult{{Result: true, Cost: uint64ptr(2)}},
			Cost:        uint64ptr(7),
		},
	}, {
		name:    "test an expression with variables evaluating to query parameters in a URL, expression should succeed",
		policy:  "variable4 policy.yaml",
		orig:    "",
		updated: "variable4 updated.yaml",
		expected: k8s.EvalResponse{
			ValidationVariables: []*k8s.EvalVariable{{
				Name: "foo",
				Value: map[string]any{
					"query": []any{"val"},
				},
				Cost: uint64ptr(14),
			}},
			Validations: []*k8s.EvalResult{{Result: true, Cost: uint64ptr(2)}},
			Cost:        uint64ptr(16),
		},
	}, {
		name:    "test valid matchConditions, should see validations and auditAnnotations",
		policy:  "match1 policy.yaml",
		orig:    "",
		updated: "match1 updated.yaml",
		request: "match1 request.yaml",
		expected: k8s.EvalResponse{
			MatchConditions: []*k8s.EvalResult{{
				Name:   strptr("exclude-leases"),
				Result: true,
				Cost:   uint64ptr(5),
			}, {
				Name:   strptr("exclude-kubelet-requests"),
				Result: true,
				Cost:   uint64ptr(5),
			}},
			Validations: []*k8s.EvalResult{{Result: true, Cost: uint64ptr(5)}},
			AuditAnnotations: []*k8s.EvalResult{{
				Name:    strptr("test-annotation"),
				Message: "Name is kubernetes-bootcamp, namespace is default",
				Cost:    uint64ptr(9),
			}},
			Cost: uint64ptr(24),
		},
	}, {
		name:    "test invalid matchConditions, should not see validations and auditAnnotations",
		policy:  "match2 policy.yaml",
		orig:    "",
		updated: "match2 updated.yaml",
		request: "match2 request.yaml",
		expected: k8s.EvalResponse{
			MatchConditionsVariables: []*k8s.EvalVariable{{
				Name:  "isLease",
				Value: false,
				Cost:  uint64ptr(4),
			}},
			MatchConditions: []*k8s.EvalResult{{
				Name:   strptr("exclude-leases"),
				Result: true,
				Cost:   uint64ptr(2),
			}, {
				Name:   strptr("exclude-kubelet-requests"),
				Result: false,
				Cost:   uint64ptr(5),
			}},
			Cost: uint64ptr(11),
		},
	}, {
		name:      "test an expression using namespace attributes",
		policy:    "namespace1 policy.yaml",
		orig:      "",
		updated:   "namespace1 updated.yaml",
		namespace: "namespace1 namespace.yaml",
		expected: k8s.EvalResponse{
			ValidationVariables: []*k8s.EvalVariable{{
				Name:  "environment",
				Value: "prod",
				Cost:  uint64ptr(7),
			}, {
				Name:  "exempt",
				Value: false,
				Cost:  uint64ptr(9),
			}, {
				Name: "containers",
				Value: []any{
					map[string]any{
						"image":                    "prod.policy.example.com/google-samples/kubernetes-bootcamp:v1",
						"imagePullPolicy":          "IfNotPresent",
						"name":                     "kubernetes-bootcamp",
						"resources":                map[string]any{},
						"terminationMessagePath":   "/dev/termination-log",
						"terminationMessagePolicy": "File",
					},
				},
				Cost: uint64ptr(5),
			}, {
				Name: "containersToCheck",
				Value: []any{
					map[string]any{
						"image":                    "prod.policy.example.com/google-samples/kubernetes-bootcamp:v1",
						"imagePullPolicy":          "IfNotPresent",
						"name":                     "kubernetes-bootcamp",
						"resources":                map[string]any{},
						"terminationMessagePath":   "/dev/termination-log",
						"terminationMessagePolicy": "File",
					},
				},
				Cost: uint64ptr(18),
			}},
			Validations: []*k8s.EvalResult{{
				Result: true,
				Cost:   uint64ptr(11),
			}},
			Cost: uint64ptr(50),
		},
	}, {
		name:    "test an expression using request attributes",
		policy:  "request1 policy.yaml",
		orig:    "",
		updated: "request1 updated.yaml",
		request: "request1 request.yaml",
		expected: k8s.EvalResponse{
			Validations: []*k8s.EvalResult{{Result: true, Cost: uint64ptr(12)}},
			Cost:        uint64ptr(12),
		},
	}, {
		name:       "test an expression using allowed authorizer checks",
		policy:     "authorizer1 policy.yaml",
		orig:       "",
		updated:    "authorizer1 updated.yaml",
		namespace:  "authorizer1 namespace.yaml",
		authorizer: "authorizer1 authorizer.yaml",
		expected: k8s.EvalResponse{
			ValidationVariables: []*k8s.EvalVariable{{
				Name:  "environment",
				Value: "prod",
				Cost:  uint64ptr(7),
			}, {
				Name:  "isProd",
				Value: true,
				Cost:  uint64ptr(2),
			}},
			Validations: []*k8s.EvalResult{{
				Result: true,
				Cost:   uint64ptr(10),
			}},
			AuditAnnotations: []*k8s.EvalResult{{
				Name:    strptr("test-annotation"),
				Message: "Deployment is allowed in namespace default",
				Cost:    uint64ptr(4),
			}},
			Cost: uint64ptr(23),
		},
	}, {
		name:       "test an expression using disallowed authorizer checks",
		policy:     "authorizer2 policy.yaml",
		orig:       "",
		updated:    "authorizer2 updated.yaml",
		namespace:  "authorizer2 namespace.yaml",
		authorizer: "authorizer2 authorizer.yaml",
		expected: k8s.EvalResponse{
			ValidationVariables: []*k8s.EvalVariable{{
				Name:  "environment",
				Value: "prod",
				Cost:  uint64ptr(7),
			}, {
				Name:  "isProd",
				Value: true,
				Cost:  uint64ptr(2),
			}},
			Validations: []*k8s.EvalResult{{
				Result: false,
				Cost:   uint64ptr(10),
			}},
			Cost: uint64ptr(19),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			policy, orig, updated, namespace, request, authorizer, err := readValidationTestData(tt.policy, tt.orig, tt.updated, tt.namespace, tt.request, tt.authorizer)
			var results string
			if err == nil {
				results, err = k8s.EvalValidatingAdmissionPolicy(policy, orig, updated, namespace, request, authorizer)
			}
			if err != nil {
				if !tt.wantErr {
					t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				evalResponse := k8s.EvalResponse{}
				if err := json.Unmarshal([]byte(results), &evalResponse); err != nil {
					t.Errorf("Eval() error = %v", err)
				}
				if !reflect.DeepEqual(tt.expected, evalResponse) {
					expected, expErr := json.Marshal(tt.expected)
					response, respErr := json.Marshal(evalResponse)
					if expErr != nil || respErr != nil {
						t.Errorf("Error marshalling expected results or evaluated responses: %v, %v", expErr, respErr)
					} else {
						t.Errorf("Expected %s\n, received %s", expected, response)
					}
				}
			}
		})
	}
}
