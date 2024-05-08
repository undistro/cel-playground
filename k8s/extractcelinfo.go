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

	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/api/admissionregistration/v1alpha1"
	"k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

type CelVariableInfo struct {
	name       string
	expression string
}

type CelValidationInfo struct {
	expression        string
	message           string
	messageExpression string
}

type CelAuditAnnotationsInfo struct {
	key        string
	expression string
}

type CelMatchConditionsInfo struct {
	name       string
	expression string
}

type CelInformation struct {
	name                   string
	namespace              string
	variables              []CelVariableInfo
	validations            []CelValidationInfo
	auditAnnotations       []CelAuditAnnotationsInfo
	matchConditions        []CelMatchConditionsInfo
	webhookMatchConditions [][]CelMatchConditionsInfo
}

func deserializeCelInformation(data []byte) (runtime.Object, error) {
	decoder := scheme.Codecs.UniversalDeserializer()

	runtimeObject, _, err := decoder.Decode(data, nil, nil)
	if err != nil {
		return nil, err
	}

	return runtimeObject, nil
}

func extractCelInformation(input []byte) (*CelInformation, error) {
	if deser, err := deserializeCelInformation(input); err != nil {
		return nil, fmt.Errorf("failed to decode input: %w", err)
	} else {
		switch resource := deser.(type) {
		case *v1alpha1.ValidatingAdmissionPolicy:
			return extractVAPV1Alpha1CelInformation(resource), nil
		case *v1beta1.ValidatingAdmissionPolicy:
			return extractVAPV1Beta1CelInformation(resource), nil
		case *v1.ValidatingAdmissionPolicy:
			return extractVAPV1CelInformation(resource), nil
		case *v1beta1.ValidatingWebhookConfiguration:
			return extractVWV1Beta1CelInformation(resource), nil
		case *v1.ValidatingWebhookConfiguration:
			return extractVWV1CelInformation(resource), nil
		case *v1beta1.MutatingWebhookConfiguration:
			return extractMWV1Beta1CelInformation(resource), nil
		case *v1.MutatingWebhookConfiguration:
			return extractMWV1CelInformation(resource), nil
		default:
			deserType := reflect.TypeOf(deser)
			return nil, fmt.Errorf("unexpected input type %s", deserType.Kind())
		}
	}
}

func extractVAPV1Alpha1CelInformation(policy *v1alpha1.ValidatingAdmissionPolicy) *CelInformation {
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

	auditAnnotations := []CelAuditAnnotationsInfo{}
	for _, auditAnnotation := range policy.Spec.AuditAnnotations {
		auditAnnotations = append(auditAnnotations, CelAuditAnnotationsInfo{
			key:        auditAnnotation.Key,
			expression: auditAnnotation.ValueExpression,
		})
	}

	matchConditions := []CelMatchConditionsInfo{}
	for _, matchCondition := range policy.Spec.MatchConditions {
		matchConditions = append(matchConditions, CelMatchConditionsInfo{
			name:       matchCondition.Name,
			expression: matchCondition.Expression,
		})
	}

	return &CelInformation{
		name, namespace, variables, validations, auditAnnotations, matchConditions, nil,
	}
}

func extractVAPV1Beta1CelInformation(policy *v1beta1.ValidatingAdmissionPolicy) *CelInformation {
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

	auditAnnotations := []CelAuditAnnotationsInfo{}
	for _, auditAnnotation := range policy.Spec.AuditAnnotations {
		auditAnnotations = append(auditAnnotations, CelAuditAnnotationsInfo{
			key:        auditAnnotation.Key,
			expression: auditAnnotation.ValueExpression,
		})
	}

	matchConditions := []CelMatchConditionsInfo{}
	for _, matchCondition := range policy.Spec.MatchConditions {
		matchConditions = append(matchConditions, CelMatchConditionsInfo{
			name:       matchCondition.Name,
			expression: matchCondition.Expression,
		})
	}

	return &CelInformation{
		name, namespace, variables, validations, auditAnnotations, matchConditions, nil,
	}
}

func extractVAPV1CelInformation(policy *v1.ValidatingAdmissionPolicy) *CelInformation {
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

	auditAnnotations := []CelAuditAnnotationsInfo{}
	for _, auditAnnotation := range policy.Spec.AuditAnnotations {
		auditAnnotations = append(auditAnnotations, CelAuditAnnotationsInfo{
			key:        auditAnnotation.Key,
			expression: auditAnnotation.ValueExpression,
		})
	}

	matchConditions := []CelMatchConditionsInfo{}
	for _, matchCondition := range policy.Spec.MatchConditions {
		matchConditions = append(matchConditions, CelMatchConditionsInfo{
			name:       matchCondition.Name,
			expression: matchCondition.Expression,
		})
	}

	return &CelInformation{
		name, namespace, variables, validations, auditAnnotations, matchConditions, nil,
	}
}

func extractVWV1Beta1CelInformation(webhookConfig *v1beta1.ValidatingWebhookConfiguration) *CelInformation {
	webhookMatchConditions := [][]CelMatchConditionsInfo{}
	for _, webhook := range webhookConfig.Webhooks {
		matchConditions := []CelMatchConditionsInfo{}
		for _, matchCondition := range webhook.MatchConditions {
			matchConditions = append(matchConditions, CelMatchConditionsInfo{
				name:       matchCondition.Name,
				expression: matchCondition.Expression,
			})
		}
		webhookMatchConditions = append(webhookMatchConditions, matchConditions)
	}
	return &CelInformation{
		"", "", nil, nil, nil, nil, webhookMatchConditions,
	}
}

func extractVWV1CelInformation(webhookConfig *v1.ValidatingWebhookConfiguration) *CelInformation {
	webhookMatchConditions := [][]CelMatchConditionsInfo{}
	for _, webhook := range webhookConfig.Webhooks {
		matchConditions := []CelMatchConditionsInfo{}
		for _, matchCondition := range webhook.MatchConditions {
			matchConditions = append(matchConditions, CelMatchConditionsInfo{
				name:       matchCondition.Name,
				expression: matchCondition.Expression,
			})
		}
		webhookMatchConditions = append(webhookMatchConditions, matchConditions)
	}
	return &CelInformation{
		"", "", nil, nil, nil, nil, webhookMatchConditions,
	}
}

func extractMWV1Beta1CelInformation(webhookConfig *v1beta1.MutatingWebhookConfiguration) *CelInformation {
	webhookMatchConditions := [][]CelMatchConditionsInfo{}
	for _, webhook := range webhookConfig.Webhooks {
		matchConditions := []CelMatchConditionsInfo{}
		for _, matchCondition := range webhook.MatchConditions {
			matchConditions = append(matchConditions, CelMatchConditionsInfo{
				name:       matchCondition.Name,
				expression: matchCondition.Expression,
			})
		}
		webhookMatchConditions = append(webhookMatchConditions, matchConditions)
	}
	return &CelInformation{
		"", "", nil, nil, nil, nil, webhookMatchConditions,
	}
}

func extractMWV1CelInformation(webhookConfig *v1.MutatingWebhookConfiguration) *CelInformation {
	webhookMatchConditions := [][]CelMatchConditionsInfo{}
	for _, webhook := range webhookConfig.Webhooks {
		matchConditions := []CelMatchConditionsInfo{}
		for _, matchCondition := range webhook.MatchConditions {
			matchConditions = append(matchConditions, CelMatchConditionsInfo{
				name:       matchCondition.Name,
				expression: matchCondition.Expression,
			})
		}
		webhookMatchConditions = append(webhookMatchConditions, matchConditions)
	}
	return &CelInformation{
		"", "", nil, nil, nil, nil, webhookMatchConditions,
	}
}
