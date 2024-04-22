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

func webhookTestfile(file string) string {
	return testfile("webhook/" + file)
}

func readWebhookTestData(webhook, original, updated, request, authorizer string) (webhookData, originalData, updatedData, requestData, authorizerData []byte, err error) {
	webhookData, err = testdata.ReadFile(webhookTestfile(webhook))
	if err == nil && original != "" {
		originalData, err = testdata.ReadFile(webhookTestfile(original))
	}
	if err == nil && updated != "" {
		updatedData, err = testdata.ReadFile(webhookTestfile(updated))
	}
	if err == nil && request != "" {
		requestData, err = testdata.ReadFile(webhookTestfile(request))
	}
	if err == nil && authorizer != "" {
		authorizerData, err = testdata.ReadFile(webhookTestfile(authorizer))
	}
	return
}

func TestWebhookEval(t *testing.T) {
	tests := []struct {
		name       string
		webhook    string
		orig       string
		updated    string
		request    string
		authorizer string
		expected   k8s.EvalResponse
		wantErr    bool
	}{{
		name:    "test a single webhook, match conditions will be successful",
		webhook: "webhook1.yaml",
		updated: "updated1.yaml",
		expected: k8s.EvalResponse{
			WebhookMatchConditions: [][]*k8s.EvalResult{{{Name: strptr("include-bootcamp"), Result: true, Cost: uint64ptr(6)}}},
			Cost:                   uint64ptr(6),
		},
	}, {
		name:    "test single webhook, match conditions will not be successful",
		webhook: "webhook2.yaml",
		updated: "updated2.yaml",
		expected: k8s.EvalResponse{
			WebhookMatchConditions: [][]*k8s.EvalResult{{{Name: strptr("exclude-bootcamp"), Result: false, Cost: uint64ptr(7)}}},
			Cost:                   uint64ptr(7),
		},
	}, {
		name:    "test a single webhook, match conditions will rely on request information",
		webhook: "webhook3.yaml",
		updated: "updated3.yaml",
		request: "request3.yaml",
		expected: k8s.EvalResponse{
			WebhookMatchConditions: [][]*k8s.EvalResult{{
				{Name: strptr("exclude-leases"), Result: true, Cost: uint64ptr(5)},
				{Name: strptr("exclude-kubelet-requests"), Result: true, Cost: uint64ptr(5)},
			}},
			Cost: uint64ptr(10),
		},
	}, {
		name:       "test a single webhook, match conditions will rely on authorizer information",
		webhook:    "webhook4.yaml",
		updated:    "updated4.yaml",
		request:    "request4.yaml",
		authorizer: "authorizer4.yaml",
		expected: k8s.EvalResponse{
			WebhookMatchConditions: [][]*k8s.EvalResult{{
				{Name: strptr("breakglass"), Result: true, Cost: uint64ptr(7)},
			}},
			Cost: uint64ptr(7),
		},
	}, {
		name:       "test multiple webhooks, match conditions will rely on request and authorizer information and will be successful",
		webhook:    "multi webhook1.yaml",
		updated:    "multi updated1.yaml",
		request:    "multi request1.yaml",
		authorizer: "multi authorizer1.yaml",
		expected: k8s.EvalResponse{
			WebhookMatchConditions: [][]*k8s.EvalResult{
				{{Name: strptr("breakglass"), Result: true, Cost: uint64ptr(7)}},
				{{Name: strptr("exclude-leases"), Result: true, Cost: uint64ptr(5)}, {Name: strptr("exclude-kubelet-requests"), Result: true, Cost: uint64ptr(5)}},
			},
			Cost: uint64ptr(17),
		},
	}, {
		name:       "test multiple webhooks, match conditions will rely on request and authorizer information and will not be successful",
		webhook:    "multi webhook2.yaml",
		updated:    "multi updated2.yaml",
		request:    "multi request2.yaml",
		authorizer: "multi authorizer2.yaml",
		expected: k8s.EvalResponse{
			WebhookMatchConditions: [][]*k8s.EvalResult{
				{{Name: strptr("breakglass"), Result: false, Cost: uint64ptr(7)}},
				{{Name: strptr("exclude-bootcamp"), Result: false, Cost: uint64ptr(7)}},
			},
			Cost: uint64ptr(14),
		},
	}, {
		name:       "test multiple webhooks, match conditions will rely on request and authorizer information with mixed responses",
		webhook:    "multi webhook3.yaml",
		updated:    "multi updated3.yaml",
		request:    "multi request3.yaml",
		authorizer: "multi authorizer3.yaml",
		expected: k8s.EvalResponse{
			WebhookMatchConditions: [][]*k8s.EvalResult{
				{{Name: strptr("breakglass"), Result: false, Cost: uint64ptr(7)}},
				{{Name: strptr("exclude-leases"), Result: true, Cost: uint64ptr(5)}, {Name: strptr("exclude-kubelet-requests"), Result: true, Cost: uint64ptr(5)}},
			},
			Cost: uint64ptr(17),
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			webhook, orig, updated, request, authorizer, err := readWebhookTestData(tt.webhook, tt.orig, tt.updated, tt.request, tt.authorizer)
			var results string
			if err == nil {
				results, err = k8s.EvalWebhook(webhook, orig, updated, request, authorizer)
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
