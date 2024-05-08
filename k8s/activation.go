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

import "github.com/google/cel-go/interpreter"

type k8sActivation struct {
	inputData                 map[string]any
	authorizer                any
	authorizerRequestResource any
}

func NewActivation(inputData map[string]any, authorizer any, authorizerRequestResource any) interpreter.Activation {
	return &k8sActivation{
		inputData:                 inputData,
		authorizer:                authorizer,
		authorizerRequestResource: authorizerRequestResource,
	}
}

func (a *k8sActivation) ResolveName(name string) (interface{}, bool) {
	switch name {
	case "authorizer":
		return a.authorizer, a.authorizer != nil
	case "authorizer.requestResource":
		return a.authorizerRequestResource, a.authorizerRequestResource != nil
	default:
		val, ok := a.inputData[name]
		return val, ok
	}
}

func (a *k8sActivation) Parent() interpreter.Activation {
	return nil
}
