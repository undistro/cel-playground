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
	"encoding/json"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/interpreter"
	"gopkg.in/yaml.v3"
)

func EvalWebhook(webhookInput, oldObjectInput, objectValueInput, requestInput, authorizerInput []byte) (string, error) {
	celInfo, err := extractCelInformation(webhookInput)
	if err != nil {
		return "", err
	}

	var oldObjectValue map[string]any
	if err := yaml.Unmarshal(oldObjectInput, &oldObjectValue); err != nil {
		return "", fmt.Errorf("failed to decode input for the old object resource value: %w", err)
	}

	var objectValue map[string]any
	if err := yaml.Unmarshal(objectValueInput, &objectValue); err != nil {
		return "", fmt.Errorf("failed to decode input for the object resource value: %w", err)
	}

	request, err := deserializeRequest(requestInput)
	if err != nil {
		return "", err
	}

	var authorizer Authorizer
	if err := yaml.Unmarshal(authorizerInput, &authorizer); err != nil {
		return "", fmt.Errorf("failed to decode input for the authorizer: %w", err)
	}
	initReceiver(&authorizer.receiverOnlyObjectVal, AuthorizerType)

	authorizerRequestResource, err := getAuthorizerRequestResource(&authorizer, request)
	if err != nil {
		return "", err
	}

	matchConditionsCelVars := []cel.EnvOption{}
	matchConditionsInputData := map[string]any{}

	if objectValue != nil {
		cleanMetaData(objectValue)
		matchConditionsCelVars = updateVars("object", matchConditionsCelVars, matchConditionsInputData, objectValue)
	}

	if oldObjectValue != nil {
		cleanMetaData(oldObjectValue)
		matchConditionsCelVars = updateVars("oldObject", matchConditionsCelVars, matchConditionsInputData, oldObjectValue)
	}

	if request != nil {
		matchConditionsCelVars = updateVars("request", matchConditionsCelVars, matchConditionsInputData, request)
	}

	if authorizerRequestResource != nil {
		matchConditionsCelVars = updateVars("authorizer.requestResource", matchConditionsCelVars, matchConditionsInputData, authorizerRequestResource)
	}

	matchConditionsCelVars = updateVars("authorizer", matchConditionsCelVars, matchConditionsInputData, &authorizer)

	// 'object' - The object from the incoming request. The value is null for DELETE requests.
	// 'oldObject' - The existing object. The value is null for CREATE requests.
	// 'request' - Attributes of the admission request(/pkg/apis/admission/types.go#AdmissionRequest).
	// 'authorizer' - A CEL Authorizer. May be used to perform authorization checks for the principal (user or service account) of the request.
	// 'authorizer.requestResource' - A CEL ResourceCheck constructed from the 'authorizer' and configured with the request resource.

	matchConditionsEnvOptions := append([]cel.EnvOption(nil), celEnvOptions...)
	matchConditionsEnvOptions = append(matchConditionsEnvOptions, matchConditionsCelVars...)
	matchConditionsEnv, err := cel.NewEnv(matchConditionsEnvOptions...)
	if err != nil {
		return "", fmt.Errorf("failed to create CEL env: %w", err)
	}

	matchConditionsExprActivations, err := interpreter.NewActivation(matchConditionsInputData)
	if err != nil {
		return "", fmt.Errorf("failed to create CEL activations: %w", err)
	}

	matchConditionsEvals := []evalResponses{}
	for _, webhookMatchConditions := range celInfo.webhookMatchConditions {
		matchConditionsEval := []*evalResponse{}
		for _, matchCondition := range webhookMatchConditions {
			ast, issues := matchConditionsEnv.Parse(matchCondition.expression)
			if issues.Err() != nil {
				return "", fmt.Errorf("failed to parse expression %s: %w", matchCondition.expression, issues.Err())
			}
			var val *evalResponse
			if prog, err := matchConditionsEnv.Program(ast, celProgramOptions...); err != nil {
				val = newEvalResponseErr("parsing", matchCondition.expression, err)
			} else if exprEval, details, err := prog.Eval(matchConditionsExprActivations); err != nil {
				val = newEvalResponseErr("evaluating", matchCondition.expression, err)
			} else {
				val = newEvalResponse(matchCondition.name, exprEval, details, "", nil)
			}
			matchConditionsEval = append(matchConditionsEval, val)
		}
		matchConditionsEvals = append(matchConditionsEvals, matchConditionsEval)
	}

	response := generateEvalResponse(nil, nil, nil, nil, nil, nil, nil, matchConditionsEvals)

	out, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
