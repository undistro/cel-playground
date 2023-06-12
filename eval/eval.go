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
	"github.com/google/cel-go/ext"
	k8s "k8s.io/apiserver/pkg/cel/library"
)

var celEnvOptions = []cel.EnvOption{
	cel.HomogeneousAggregateLiterals(),
	cel.EagerlyValidateDeclarations(true),
	cel.DefaultUTCTimeZone(true),
	ext.Strings(ext.StringsVersion(2)),
	cel.CrossTypeNumericComparisons(true),
	cel.OptionalTypes(),
	k8s.URLs(),
	k8s.Regex(),
	k8s.Lists(),
}

var celProgramOptions = []cel.ProgramOption{
	cel.EvalOptions(cel.OptOptimize, cel.OptTrackCost),
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
		return "", fmt.Errorf("failed to intantiate CEL program: %w", err)
	}
	val, _, err := prog.Eval(input)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate: %w", err)
	}
	b, err := json.MarshalIndent(val.Value(), "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal the output: %w", err)
	}
	return string(b), nil
}
