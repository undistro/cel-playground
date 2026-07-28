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
				Cost: uint64ptr(19),
			}},
			Validations: []*k8s.EvalResult{{Result: true, Cost: uint64ptr(2)}},
			Cost:        uint64ptr(21),
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
				Cost:   uint64ptr(350009),
			}},
			AuditAnnotations: []*k8s.EvalResult{{
				Name:    strptr("test-annotation"),
				Message: "Deployment is allowed in namespace default",
				Cost:    uint64ptr(4),
			}},
			Cost: uint64ptr(350022),
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
				Cost:   uint64ptr(350009),
			}},
			Cost: uint64ptr(350018),
		},
	}, {
		name:    "test a broken expression within variables, expression should fail with no audit annotation",
		policy:  "broken1 policy.yaml",
		orig:    "",
		updated: "broken1 updated.yaml",
		expected: k8s.EvalResponse{
			ValidationVariables: []*k8s.EvalVariable{{
				Name:  "foo",
				Value: "default",
				Cost:  uint64ptr(6),
			}, {
				Name:    "containers",
				IsError: true,
				Error:   strptr("unexpected error evaluating expression containers: no such key: spc"),
			}},
			Validations: []*k8s.EvalResult{{
				IsError: true,
				Error:   strptr("unexpected error evaluating expression 'variables.foo == 'default' && variables.containers.all(c, c.image.startsWith(\"test\"))', caused by nested exception: 'no such key: spc'"),
			}},
			AuditAnnotations: []*k8s.EvalResult{{
				Name:    strptr("foo-label"),
				Message: "Label for foo is set to default",
				Cost:    uint64ptr(2),
			}},
			Cost: uint64ptr(8),
		},
	}, {
		name:    "test optional.none() dereference",
		policy:  "optional_none_dereference policy.yaml",
		orig:    "",
		updated: "optional_none_dereference updated.yaml",

		expected: k8s.EvalResponse{
			ValidationVariables: []*k8s.EvalVariable{{
				Name: "containers",
				Value: []any{
					map[string]any{
						"image":                    "gcr.io/google-samples/kubernetes-bootcamp:v1",
						"imagePullPolicy":          "IfNotPresent",
						"name":                     "kubernetes-bootcamp",
						"resources":                map[string]any{},
						"terminationMessagePath":   "/dev/termination-log",
						"terminationMessagePolicy": "File",
					},
				},
				Cost: uint64ptr(5),
			}, {
				Name:  "securityContexts",
				Value: []any{nil},
				Cost:  uint64ptr(15),
			}, {
				Name: "namedSecurityContexts",
				Value: []any{
					map[string]any{
						"kubernetes-bootcamp": nil,
					},
				},
				Cost: uint64ptr(47),
			}},
			Validations: []*k8s.EvalResult{{
				Result:  false,
				Message: "all containers must set runAsNonRoot to true",
				Cost:    uint64ptr(8),
			}, {
				Result:  false,
				Message: "all containers must set readOnlyRootFilesystem to true",
				Cost:    uint64ptr(8),
			}, {
				Result: true,
				Cost:   uint64ptr(8),
			}, {
				Result: true,
				Cost:   uint64ptr(8),
			}, {
				Result: true,
				Cost:   uint64ptr(8),
			}},
			Cost: uint64ptr(107),
		},
	}, {
		// The object is decoded the way a cluster decodes it. `eviction: on` is
		// YAML 1.1, so it reaches CEL as a bool rather than the string "on" --
		// the first validation fails without that coercion. The second guards
		// the other half of the decoder: `replicas: 3` must reach CEL as an int
		// rather than a double, or `%` has no overload.
		name:    "test an object whose values depend on a cluster-faithful decode",
		policy:  "decode1 policy.yaml",
		orig:    "",
		updated: "decode1 updated.yaml",
		expected: k8s.EvalResponse{
			Validations: []*k8s.EvalResult{
				{Result: true, Cost: uint64ptr(4)},
				{Result: true, Cost: uint64ptr(5)},
			},
			Cost: uint64ptr(9),
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

// TestValidationCELLibraries exercises the CEL libraries added to k8s/cel.go so
// they stay wired up. Each case is a validation expression that uses one
// library and must evaluate to true. It asserts availability, not cost, so it
// checks the result and the absence of an error and does not pin a cost.
func TestValidationCELLibraries(t *testing.T) {
	object := `{"apiVersion":"apps/v1","kind":"Deployment","metadata":{"name":"x"}}`
	cases := []struct {
		lib  string
		expr string
	}{
		{"semver", `semver('1.2.3').minor() == 2`},
		{"ip", `isIP('10.0.0.1') && ip('10.0.0.1').family() == 4`},
		{"cidr", `isCIDR('10.0.0.0/8') && cidr('10.0.0.0/8').containsIP(ip('10.0.0.1'))`},
		{"format", `format.named('dns1123Label').hasValue()`},
		{"sets", `sets.contains([1, 2, 3], [2, 3])`},
		{"twoVarComprehensions", `[10, 20, 30].all(i, v, v >= 10)`},
		{"extLists", `[3, 1, 2].sort() == [1, 2, 3]`},
	}
	for _, tc := range cases {
		t.Run(tc.lib, func(t *testing.T) {
			policy := `apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingAdmissionPolicy
metadata:
  name: lib-check
spec:
  failurePolicy: Fail
  matchConstraints:
    resourceRules:
    - apiGroups: ["apps"]
      apiVersions: ["v1"]
      operations: ["CREATE"]
      resources: ["deployments"]
  validations:
    - expression: "` + tc.expr + `"
`
			out, err := k8s.EvalValidatingAdmissionPolicy([]byte(policy), nil, []byte(object), nil, nil, nil)
			if err != nil {
				t.Fatalf("Eval() error = %v", err)
			}
			var resp k8s.EvalResponse
			if err := json.Unmarshal([]byte(out), &resp); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if len(resp.Validations) != 1 {
				t.Fatalf("got %d validations, want 1: %s", len(resp.Validations), out)
			}
			v := resp.Validations[0]
			if v.IsError {
				msg := "<nil>"
				if v.Error != nil {
					msg = *v.Error
				}
				t.Fatalf("%s expression errored (library not registered?): %s", tc.lib, msg)
			}
			if v.Result != true {
				t.Errorf("%s expression = %v, want true", tc.lib, v.Result)
			}
		})
	}
}
