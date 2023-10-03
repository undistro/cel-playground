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
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	"github.com/google/cel-go/interpreter"
	"gopkg.in/yaml.v3"
	k8s "k8s.io/apiserver/pkg/cel/library"
)

var celEnvOptions = []cel.EnvOption{
	cel.EagerlyValidateDeclarations(true),
	cel.DefaultUTCTimeZone(true),
	ext.Strings(ext.StringsVersion(2)),
	cel.CrossTypeNumericComparisons(true),
	cel.OptionalTypes(),
	k8s.URLs(),
	k8s.Regex(),
	k8s.Lists(),
	k8s.Quantity(),
}

var celProgramOptions = []cel.ProgramOption{
	cel.EvalOptions(cel.OptOptimize, cel.OptTrackCost),
}

func EvalValidatingAdmissionPolicy(policyInput, origValueInput, updatedValueInput []byte) (string, error) {
	celInfo, err := extractCelInformation(policyInput)
	if err != nil {
		return "", err
	}

	var origValue map[string]any
	if err := yaml.Unmarshal(origValueInput, &origValue); err != nil {
		return "", fmt.Errorf("failed to decode input for the old resource value: %w", err)
	}

	var updatedValue map[string]any
	if err := yaml.Unmarshal(updatedValueInput, &updatedValue); err != nil {
		return "", fmt.Errorf("failed to decode input for the new resource value: %w", err)
	}

	celVars := []cel.EnvOption{}
	inputData := map[string]any{}

	if origValue != nil {
		celVars = append(celVars, cel.Variable("oldObject", cel.DynType))
		inputData["oldObject"] = origValue
	}

	if updatedValue != nil {
		celVars = append(celVars, cel.Variable("object", cel.DynType))
		inputData["object"] = updatedValue
	}

	env, err := cel.NewEnv(append(celEnvOptions, celVars...)...)
	if err != nil {
		return "", fmt.Errorf("failed to create CEL env: %w", err)
	}

	exprActivations, err := interpreter.NewActivation(inputData)
	if err != nil {
		return "", fmt.Errorf("failed to create CEL activations: %w", err)
	}

	variableLazyEvals := map[string]*lazyVariableEval{}

	if len(celInfo.variables) > 0 {
		for _, variable := range celInfo.variables {
			ast, issues := env.Parse(variable.expression)
			if issues.Err() != nil {
				return "", fmt.Errorf("failed to parse expression for variable %s: %w", variable.name, err)
			}
			env, err = env.Extend(cel.Variable(variable.name, ast.OutputType()))
			if err != nil {
				return "", fmt.Errorf("could not append variable %s to CEL env: %w", variable.name, err)
			}
			variableLazyEval := lazyVariableEval{
				name: variable.name,
				ast:  ast,
			}
			variableLazyEvals[variable.name] = &variableLazyEval
			name := "variables." + variable.name
			inputData[name] = func() ref.Val {
				return variableLazyEval.eval(env, exprActivations)
			}
			celVars = append(celVars, cel.Variable(name, ast.OutputType()))
		}
	}

	validationEvals := []ref.Val{}

	for _, validation := range celInfo.validations {
		ast, issues := env.Parse(validation.expression)
		if issues.Err() != nil {
			return "", fmt.Errorf("failed to parse expression %s: %w", validation.expression, err)
		}
		var val ref.Val
		if prog, err := env.Program(ast, celProgramOptions...); err != nil {
			val = types.NewErr("Unexpected error parsing expression %s: %v", validation.expression, err)
		} else if exprEval, _, err := prog.Eval(exprActivations); err != nil {
			val = types.NewErr("Unexpected error parsing expression %s: %v", validation.expression, err)
		} else {
			val = exprEval
		}
		validationEvals = append(validationEvals, val)
	}

	response := generateResponse(variableLazyEvals, validationEvals)

	out, err := json.Marshal(response)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
