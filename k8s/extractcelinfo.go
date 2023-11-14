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
	"fmt"
	"reflect"

	"k8s.io/api/admissionregistration/v1alpha1"
	"k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

type MatchCondition struct {
	name       string
	expression string
}

type CelVariableInfo struct {
	name       string
	expression string
}

type CelValidationInfo struct {
	expression        string
	message           string
	messageExpression string
}

type CelInformation struct {
	name        string
	namespace   string
	variables   []CelVariableInfo
	validations []CelValidationInfo
}

func deserialize(data []byte) (runtime.Object, error) {
	decoder := scheme.Codecs.UniversalDeserializer()

	runtimeObject, _, err := decoder.Decode(data, nil, nil)
	if err != nil {
		return nil, err
	}

	return runtimeObject, nil
}

func extractCelInformation(policyInput []byte) (*CelInformation, error) {
	if deser, err := deserialize(policyInput); err != nil {
		return nil, fmt.Errorf("failed to decode ValidatingAdmissionPolicy: %w", err)
	} else {
		switch policy := deser.(type) {
		case *v1alpha1.ValidatingAdmissionPolicy:
			return extractV1Alpha1CelInformation(policy), nil
		case *v1beta1.ValidatingAdmissionPolicy:
			return extractV1Beta1CelInformation(policy), nil
		default:
			policyType := reflect.TypeOf(deser)
			return nil, fmt.Errorf("expected ValidatingAdmissionPolicy, received %s", policyType.Kind())
		}
	}
}

func extractV1Alpha1CelInformation(policy *v1alpha1.ValidatingAdmissionPolicy) *CelInformation {
	namespace := policy.ObjectMeta.GetNamespace()
	name := policy.ObjectMeta.GetName()

	variables := []CelVariableInfo{}
	for _, variable := range policy.Spec.Variables {
		variables = append(variables, CelVariableInfo{
			name:       variable.Name,
			expression: variable.Expression,
		})
	}

	validations := []CelValidationInfo{}
	for _, validation := range policy.Spec.Validations {

		validations = append(validations, CelValidationInfo{
			expression:        validation.Expression,
			message:           validation.Message,
			messageExpression: validation.MessageExpression,
		})
	}

	return &CelInformation{
		name, namespace, variables, validations,
	}
}

func extractV1Beta1CelInformation(policy *v1beta1.ValidatingAdmissionPolicy) *CelInformation {
	namespace := policy.ObjectMeta.GetNamespace()
	name := policy.ObjectMeta.GetName()

	variables := []CelVariableInfo{}
	for _, variable := range policy.Spec.Variables {
		variables = append(variables, CelVariableInfo{
			name:       variable.Name,
			expression: variable.Expression,
		})
	}

	validations := []CelValidationInfo{}
	for _, validation := range policy.Spec.Validations {

		validations = append(validations, CelValidationInfo{
			expression:        validation.Expression,
			message:           validation.Message,
			messageExpression: validation.MessageExpression,
		})
	}

	return &CelInformation{
		name, namespace, variables, validations,
	}
}
