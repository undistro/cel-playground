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

package eval

import (
	"encoding/json"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	"github.com/google/cel-go/interpreter"
	"github.com/undistro/cel-playground/utils"
	"gopkg.in/yaml.v2"
	k8s "k8s.io/apiserver/pkg/cel/library"
)

type EvalResponse struct {
	Result any     `json:"result"`
	Cost   *uint64 `json:"cost, omitempty"`
}

var celEnvOptions = []cel.EnvOption{
	// 1.0 (1.23)
	cel.HomogeneousAggregateLiterals(),
	cel.EagerlyValidateDeclarations(true),
	cel.DefaultUTCTimeZone(true),
	k8s.URLs(),
	k8s.Regex(),
	k8s.Lists(),

	// 1.27
	// k8s.Authz(),

	// 1.28
	cel.CrossTypeNumericComparisons(true),
	cel.OptionalTypes(),
	k8s.Quantity(),

	// 1.29 (see also validator.ExtendedValidations())
	cel.ASTValidators(
		cel.ValidateDurationLiterals(),
		cel.ValidateTimestampLiterals(),
		cel.ValidateRegexLiterals(),
		cel.ValidateHomogeneousAggregateLiterals(),
	),

	// Strings (from 1.29 onwards)
	ext.Strings(ext.StringsVersion(2)),
	// Set library (1.29 onwards)
	ext.Sets(),

	// cel-go v0.17.7 introduced CostEstimatorOptions.
	// Previous the presence has a cost of 0 but cel fixed it to 1. We still set to 0 here to avoid breaking changes.
	cel.CostEstimatorOptions(checker.PresenceTestHasCost(false)),
}

var celProgramOptions = []cel.ProgramOption{
	cel.EvalOptions(cel.OptOptimize, cel.OptTrackCost),

	// cel-go v0.17.7 introduced CostTrackerOptions.
	// Previous the presence has a cost of 0 but cel fixed it to 1. We still set to 0 here to avoid breaking changes.
	cel.CostTrackerOptions(interpreter.PresenceTestHasCost(false)),
}

func CelEval(exp []byte, input []byte) (string, error) {
	var inputMap map[string]any
	if err := yaml.Unmarshal(input, &inputMap); err != nil {
		return "", fmt.Errorf("failed to decode input: %w", err)
	}
	return Eval(string(exp), inputMap)
}

// Eval evaluates the cel expression against the given input
func Eval(exp string, input map[string]any) (string, error) {
	inputVars := make([]cel.EnvOption, 0, len(input))
	for k := range input {
		inputVars = append(inputVars, cel.Variable(k, cel.DynType))
	}
	env, err := cel.NewEnv(append(celEnvOptions, inputVars...)...)
	if err != nil {
		return "", fmt.Errorf("failed to create CEL env: %w", err)
	}
	ast, issues := env.Compile(exp)
	if issues != nil {
		return "", fmt.Errorf("failed to compile the CEL expression: %s", issues.String())
	}
	prog, err := env.Program(ast, celProgramOptions...)
	if err != nil {
		return "", fmt.Errorf("failed to instantiate CEL program: %w", err)
	}
	val, costTracker, err := prog.Eval(input)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate: %w", err)
	}

	response, err := generateResponse(val, costTracker)
	if err != nil {
		return "", fmt.Errorf("failed to generate the response: %w", err)
	}

	out, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal the output: %w", err)
	}
	return string(out), nil
}

func getResults(val ref.Val) (any, error) {
	if value, err := utils.ConvertValToNative(val); err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

func generateResponse(val ref.Val, costTracker *cel.EvalDetails) (*EvalResponse, error) {
	result, evalError := getResults(val)
	if evalError != nil {
		return nil, evalError
	}
	cost := costTracker.ActualCost()
	return &EvalResponse{
		Result: result,
		Cost:   cost,
	}, nil
}
