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
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	"google.golang.org/protobuf/types/known/structpb"
)

type lazyVariableEval struct {
	name string
	ast  *cel.Ast
	val  *ref.Val
}

func (lve *lazyVariableEval) eval(env *cel.Env, activation interpreter.Activation) ref.Val {
	val := lve.evalExpression(env, activation)
	lve.val = &val
	return val
}

func (lve *lazyVariableEval) evalExpression(env *cel.Env, activation interpreter.Activation) ref.Val {
	prog, err := env.Program(lve.ast, celProgramOptions...)
	if err != nil {
		return types.NewErr("Unexpected error parsing expression %s: %v", lve.name, err)
	}
	val, _, err := prog.Eval(activation)
	if err != nil {
		return types.NewErr("Unexpected error parsing expression %s: %v", lve.name, err)
	}
	return val
}

type EvalVariable struct {
	Name  string  `json:"name"`
	Value any     `json:"value"`
	Error *string `json:"error,omitempty"`
}

type EvalValidation struct {
	Result any     `json:"result"`
	Error  *string `json:"error,omitempty"`
}

type EvalResponse struct {
	Variables   []EvalVariable   `json:"variables,omitempty"`
	Validations []EvalValidation `json:"validations,omitempty"`
}

func getResults(val *ref.Val) (any, *string) {
	value := (*val).Value()
	if err, ok := value.(error); ok {
		errResponse := err.Error()
		return nil, &errResponse
	}
	if value, err := (*val).ConvertToNative(reflect.TypeOf(&structpb.Value{})); err != nil {
		errResponse := err.Error()
		return nil, &errResponse
	} else {
		return value, nil
	}
}

func generateResponse(variableLazyEvals map[string]*lazyVariableEval, validationEvals []ref.Val) *EvalResponse {
	variables := []EvalVariable{}
	for _, varLazyEval := range variableLazyEvals {
		if varLazyEval.val != nil {
			value, err := getResults(varLazyEval.val)
			variables = append(variables, EvalVariable{
				Name:  varLazyEval.name,
				Value: value,
				Error: err,
			})
		}
	}
	validations := []EvalValidation{}
	for _, validationEval := range validationEvals {
		value, err := getResults(&validationEval)
		validations = append(validations, EvalValidation{
			Result: value,
			Error:  err,
		})
	}
	return &EvalResponse{
		Variables:   variables,
		Validations: validations,
	}
}
