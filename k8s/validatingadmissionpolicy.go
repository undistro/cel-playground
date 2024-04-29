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
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	"gopkg.in/yaml.v3"
)

const (
	metadata             = "metadata"
	metadataName         = "name"
	metadataGenerateName = "generateName"
)

var stringType = reflect.TypeOf("")

// From
//
//	pkg/apis/admissionregistration/types.go#Validation
//
// Expression represents the expression which will be evaluated by CEL.
// ref: https://github.com/google/cel-spec
// CEL expressions have access to the contents of the API request/response, organized into CEL variables as well as some other useful variables:
//
// 'object' - The object from the incoming request. The value is null for DELETE requests.
// 'oldObject' - The existing object. The value is null for CREATE requests.
// 'request' - Attributes of the API request([ref](/pkg/apis/admission/types.go#AdmissionRequest)).
// 'params' - Parameter resource referred to by the policy binding being evaluated. Only populated if the policy has a ParamKind.
// 'namespaceObject' - The namespace object that the incoming object belongs to. The value is null for cluster-scoped resources.
// 'variables' - Map of composited variables, from its name to its lazily evaluated value.
//
//		For example, a variable named 'foo' can be accessed as 'variables.foo'
//	  - 'authorizer' - A CEL Authorizer. May be used to perform authorization checks for the principal (user or service account) of the request.
//	    See https://pkg.go.dev/k8s.io/apiserver/pkg/cel/library#Authz
//	  - 'authorizer.requestResource' - A CEL ResourceCheck constructed from the 'authorizer' and configured with the
//	    request resource.
//
// KEV - check what metadata comment located below means
// The `apiVersion`, `kind`, `metadata.name` and `metadata.generateName` are always accessible from the root of the
// object. No other metadata properties are accessible.
//
// Only property names of the form `[a-zA-Z_.-/][a-zA-Z0-9_.-/]*` are accessible.
// Accessible property names are escaped according to the following rules when accessed in the expression:
//   - '__' escapes to '__underscores__'
//   - '.' escapes to '__dot__'
//   - '-' escapes to '__dash__'
//   - '/' escapes to '__slash__'
//   - Property names that exactly match a CEL RESERVED keyword escape to '__{keyword}__'. The keywords are:
//     "true", "false", "null", "in", "as", "break", "const", "continue", "else", "for", "function", "if",
//     "import", "let", "loop", "package", "namespace", "return".
//
// Examples:
//   - Expression accessing a property named "namespace": {"Expression": "object.__namespace__ > 0"}
//   - Expression accessing a property named "x-prop": {"Expression": "object.x__dash__prop > 0"}
//   - Expression accessing a property named "redact__d": {"Expression": "object.redact__underscores__d > 0"}
//
// Equality on arrays with list type of 'set' or 'map' ignores element order, i.e. [1, 2] == [2, 1].
// Concatenation on arrays with x-kubernetes-list-type use the semantics of the list type:
//   - 'set': `X + Y` performs a union where the array positions of all elements in `X` are preserved and
//     non-intersecting elements in `Y` are appended, retaining their partial order.
//   - 'map': `X + Y` performs a merge where the array positions of all keys in `X` are preserved but the values
//     are overwritten by values in `Y` when the key sets of `X` and `Y` intersect. Elements in `Y` with
//     non-intersecting keys are appended, retaining their partial order.
//
// TODO: Support parameters
func EvalValidatingAdmissionPolicy(policyInput, oldObjectInput, objectValueInput, namespaceInput, requestInput, authorizerInput []byte) (string, error) {
	celInfo, err := extractCelInformation(policyInput)
	if err != nil {
		return "", err
	}

	var oldObjectValue map[string]any
	if err := yaml.Unmarshal(oldObjectInput, &oldObjectValue); err != nil {
		return "", fmt.Errorf("failed to decode input for the old resource value: %w", err)
	}

	var objectValue map[string]any
	if err := yaml.Unmarshal(objectValueInput, &objectValue); err != nil {
		return "", fmt.Errorf("failed to decode input for the new resource value: %w", err)
	}

	namespaceObject, err := deserializeNamespace(namespaceInput)
	if err != nil {
		return "", err
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

	validationCelVars := []cel.EnvOption{}
	validationInputData := map[string]any{}
	matchConditionsCelVars := []cel.EnvOption{}
	matchConditionsInputData := map[string]any{}

	if objectValue != nil {
		cleanMetaData(objectValue)
		validationCelVars = updateVars("object", validationCelVars, validationInputData, objectValue)
		matchConditionsCelVars = updateVars("object", matchConditionsCelVars, matchConditionsInputData, objectValue)
	}

	if oldObjectValue != nil {
		cleanMetaData(oldObjectValue)
		validationCelVars = updateVars("oldObject", validationCelVars, validationInputData, oldObjectValue)
		matchConditionsCelVars = updateVars("oldObject", matchConditionsCelVars, matchConditionsInputData, oldObjectValue)
	}

	if request != nil {
		validationCelVars = updateVars("request", validationCelVars, validationInputData, request)
		matchConditionsCelVars = updateVars("request", matchConditionsCelVars, matchConditionsInputData, request)
	}

	if namespaceObject != nil {
		validationCelVars = updateVars("namespaceObject", validationCelVars, validationInputData, namespaceObject)
	}

	if authorizerRequestResource != nil {
		validationCelVars = updateVars("authorizer.requestResource", validationCelVars, validationInputData, authorizerRequestResource)
		matchConditionsCelVars = updateVars("authorizer.requestResource", matchConditionsCelVars, matchConditionsInputData, authorizerRequestResource)
	}

	validationCelVars = updateVars("authorizer", validationCelVars, validationInputData, &authorizer)
	matchConditionsCelVars = updateVars("authorizer", matchConditionsCelVars, matchConditionsInputData, &authorizer)

	// The exact matching logic is (in order):
	//   1. If ANY matchCondition evaluates to FALSE, the policy is skipped.
	//   2. If ALL matchConditions evaluate to TRUE, the policy is evaluated.
	//   3. If any matchCondition evaluates to an error (but none are FALSE):
	//      - If failurePolicy=Fail, reject the request
	//      - If failurePolicy=Ignore, the policy is skipped
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

	matchConditionsVariableLazyEvals := lazyEvalMap{}
	matchConditionsVariableNames := []string{}

	if len(celInfo.variables) > 0 {
		matchConditionsEnv, matchConditionsVariableNames, err = initVars(matchConditionsEnv, celInfo.variables /*matchConditionsCelVars,*/, matchConditionsVariableLazyEvals, matchConditionsExprActivations, matchConditionsInputData)
		if err != nil {
			return "", fmt.Errorf("failed to initialize variables: %w", err)
		}
	}

	matchConditions := true
	matchConditionsEvals := []*evalResponse{}

	for _, matchCondition := range celInfo.matchConditions {
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
			matchConditions = matchConditions && (exprEval.Value() == true)
			val = newEvalResponse(matchCondition.name, exprEval, details, "", nil)
		}
		matchConditionsEvals = append(matchConditionsEvals, val)
	}

	validationVariableLazyEvals := lazyEvalMap{}
	validationVariableNames := []string{}
	validationEvals := []*evalResponse{}
	auditAnnotationEvals := []*evalResponse{}

	// run validations only if matchConditions pass
	if matchConditions {
		validationEnvOptions := append([]cel.EnvOption(nil), celEnvOptions...)
		validationEnvOptions = append(validationEnvOptions, validationCelVars...)
		validationEnv, err := cel.NewEnv(validationEnvOptions...)
		if err != nil {
			return "", fmt.Errorf("failed to create CEL env: %w", err)
		}

		validationExprActivations, err := interpreter.NewActivation(validationInputData)
		if err != nil {
			return "", fmt.Errorf("failed to create CEL activations: %w", err)
		}
		if len(celInfo.variables) > 0 {
			validationEnv, validationVariableNames, err = initVars(validationEnv, celInfo.variables, validationVariableLazyEvals, validationExprActivations, validationInputData)
			if err != nil {
				return "", fmt.Errorf("failed to initialize variables: %w", err)
			}
		}

		validationResult := true
		for _, validation := range celInfo.validations {
			ast, issues := validationEnv.Parse(validation.expression)
			if issues.Err() != nil {
				return "", fmt.Errorf("failed to parse expression %s: %w", validation.expression, issues.Err())
			}
			var val *evalResponse
			if prog, err := validationEnv.Program(ast, celProgramOptions...); err != nil {
				val = newEvalResponseErr("parsing", validation.expression, err)
			} else if exprEval, details, err := prog.Eval(validationExprActivations); err != nil {
				val = newEvalResponseErr("evaluating", validation.expression, err)
			} else if exprEval.Value() == true {
				val = newEvalResponse("", exprEval, details, "", nil)
			} else {
				validationResult = false
				if validation.message != "" {
					val = newEvalResponse("", exprEval, details, validation.message, nil)
				} else if validation.messageExpression != "" {
					msgAst, issues := validationEnv.Parse(validation.messageExpression)
					if issues.Err() != nil {
						return "", fmt.Errorf("failed to parse expression %s: %w", validation.messageExpression, issues.Err())
					}
					if msgProg, err := validationEnv.Program(msgAst, celProgramOptions...); err != nil {
						val = newEvalResponseErr("parsing", validation.messageExpression, err)
					} else if msgExprEval, details, err := msgProg.Eval(validationExprActivations); err != nil {
						val = newEvalResponseErr("evaluating", validation.messageExpression, err)
					} else {
						val = newEvalResponse("", exprEval, details, "", msgExprEval)
					}
				} else {
					val = newEvalResponse("", exprEval, details, validation.message, nil)
				}
			}
			validationEvals = append(validationEvals, val)
		}

		if validationResult {
			for _, auditAnnotation := range celInfo.auditAnnotations {
				ast, issues := validationEnv.Parse(auditAnnotation.expression)
				if issues.Err() != nil {
					return "", fmt.Errorf("failed to parse expression %s: %w", auditAnnotation.expression, issues.Err())
				}
				var val *evalResponse
				if prog, err := validationEnv.Program(ast, celProgramOptions...); err != nil {
					val = newEvalResponseErr("parsing", auditAnnotation.expression, err)
				} else if exprEval, details, err := prog.Eval(validationExprActivations); err != nil {
					val = newEvalResponseErr("evaluating", auditAnnotation.expression, err)
				} else {
					val = newEvalResponse(auditAnnotation.key, nil, details, "", exprEval)
				}
				auditAnnotationEvals = append(auditAnnotationEvals, val)
			}
		}
	}

	response := generateEvalResponse(matchConditionsVariableNames, matchConditionsVariableLazyEvals, matchConditionsEvals,
		validationVariableNames, validationVariableLazyEvals, validationEvals,
		auditAnnotationEvals, nil)

	out, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func updateVars(name string, celVars []cel.EnvOption, inputData map[string]any, value any) []cel.EnvOption {
	celVars = append(celVars, cel.Variable(name, cel.DynType))
	inputData[name] = value
	return celVars
}

func cleanMetaData(obj map[string]any) {
	// KEV the comment says only a few, however examples and code suggest otherwise
	// KEV check to see what is really going on
	// if metadataVals, ok := obj[metadata]; ok {
	// 	switch mapVal := metadataVals.(type) {
	// 	case map[string]any:
	// 		for k, _ := range mapVal {
	// 			if k != metadataName && k != metadataGenerateName {
	// 				delete(mapVal, k)
	// 			}
	// 		}
	// 	}
	// }
}

func initVars(env *cel.Env, variableInfos []CelVariableInfo /*celVars []cel.EnvOption,*/, lazyEvals lazyEvalMap, activation interpreter.Activation, inputData map[string]any) (*cel.Env, []string, error) {
	names := []string{}
	for _, variable := range variableInfos {
		ast, issues := env.Parse(variable.expression)
		if issues.Err() != nil {
			return nil, nil, fmt.Errorf("failed to parse expression for variable %s: %w", variable.name, issues.Err())
		}
		env, err := env.Extend(cel.Variable(variable.name, ast.OutputType()))
		if err != nil {
			return nil, nil, fmt.Errorf("could not append variable %s to CEL env: %w", variable.name, err)
		}
		variableLazyEval := lazyVariableEval{
			name: variable.name,
			ast:  ast,
		}
		names = append(names, variable.name)
		lazyEvals[variable.name] = &variableLazyEval
		name := "variables." + variable.name
		inputData[name] = func() ref.Val {
			return variableLazyEval.eval(env, activation)
		}
	}
	return env, names, nil
}
