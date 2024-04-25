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

type evalResponse struct {
	name       string
	val        ref.Val
	details    *cel.EvalDetails
	messageVal ref.Val
	message    string
}

type evalResponses []*evalResponse

func newEvalResponseErr(operation, expression string, err error) *evalResponse {
	return &evalResponse{
		val: types.NewErr("Unexpected error %s expression %s: %v", operation, expression, err),
	}
}

func newEvalResponse(name string, exprEval ref.Val, details *cel.EvalDetails, message string, messageVal ref.Val) *evalResponse {
	return &evalResponse{
		name:       name,
		val:        exprEval,
		details:    details,
		messageVal: messageVal,
		message:    message,
	}
}

type lazyVariableEval struct {
	name string
	ast  *cel.Ast
	val  *evalResponse
}

func (lve *lazyVariableEval) eval(env *cel.Env, activation interpreter.Activation) ref.Val {
	val := lve.evalExpression(env, activation)
	lve.val = val
	return val.val
}

func (lve *lazyVariableEval) evalExpression(env *cel.Env, activation interpreter.Activation) *evalResponse {
	prog, err := env.Program(lve.ast, celProgramOptions...)
	if err != nil {
		return newEvalResponseErr("parsing", lve.name, err)
	}
	val, details, err := prog.Eval(activation)
	if err != nil {
		return newEvalResponseErr("evaluating", lve.name, err)
	}
	return newEvalResponse(lve.name, val, details, "", nil)
}

type lazyEvalMap map[string]*lazyVariableEval

type EvalVariable struct {
	Name  string  `json:"name"`
	Value any     `json:"value"`
	Cost  *uint64 `json:"cost,omitempty"`
	Error *string `json:"error,omitempty"`
}

type EvalResult struct {
	Name    *string `json:"name,omitempty"`
	Result  any     `json:"result,omitempty"`
	Cost    *uint64 `json:"cost,omitempty"`
	Error   *string `json:"error,omitempty"`
	Message any     `json:"message,omitempty"`
}

type EvalResponse struct {
	MatchConditionsVariables []*EvalVariable `json:"matchConditionVariables,omitempty"`
	MatchConditions          []*EvalResult   `json:"matchConditions,omitempty"`
	ValidationVariables      []*EvalVariable `json:"validationVariables,omitempty"`
	Validations              []*EvalResult   `json:"validations,omitempty"`
	AuditAnnotations         []*EvalResult   `json:"auditAnnotations,omitempty"`
	WebhookMatchConditions   [][]*EvalResult `json:"webhookMatchConditions,omitempty"`
	Cost                     *uint64         `json:"cost, omitempty"`
}

func getResults(val *ref.Val) (any, *string) {
	if val == nil || *val == nil {
		return nil, nil
	}
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

func getCost(details *cel.EvalDetails) *uint64 {
	if details == nil {
		return nil
	}
	return details.ActualCost()
}

func generateEvalVariables(names []string, lazyEvals lazyEvalMap) []*EvalVariable {
	variables := []*EvalVariable{}
	for _, name := range names {
		if varLazyEval, ok := lazyEvals[name]; ok && varLazyEval.val != nil {
			value, err := getResults(&varLazyEval.val.val)
			variables = append(variables, &EvalVariable{
				Name:  varLazyEval.name,
				Value: value,
				Cost:  getCost(varLazyEval.val.details),
				Error: err,
			})
		}
	}
	return variables
}

func generateEvalResults(responses evalResponses) []*EvalResult {
	evals := []*EvalResult{}
	for _, eval := range responses {
		value, err := getResults(&eval.val)
		var message any
		if eval.messageVal != nil {
			message, _ = getResults(&eval.messageVal)
		} else if eval.message != "" {
			message = eval.message
		}
		var name *string
		if eval.name != "" {
			name = &eval.name
		}
		evals = append(evals, &EvalResult{
			Name:    name,
			Result:  value,
			Cost:    getCost(eval.details),
			Error:   err,
			Message: message,
		})
	}
	return evals
}

func generateEvalArrayResults(responses []evalResponses) [][]*EvalResult {
	evalsArray := [][]*EvalResult{}
	for _, response := range responses {
		evals := generateEvalResults(response)
		evalsArray = append(evalsArray, evals)
	}
	return evalsArray
}

func calculateLazyEvalCost(lazyEvals lazyEvalMap) uint64 {
	var cost uint64
	for _, lazyEval := range lazyEvals {
		if lazyEval.val != nil && lazyEval.val.details != nil {
			cost += *lazyEval.val.details.ActualCost()
		}
	}
	return cost
}

func calculateEvalResponsesCost(evals evalResponses) uint64 {
	var cost uint64
	for _, eval := range evals {
		if eval.details != nil {
			cost += *eval.details.ActualCost()
		}
	}
	return cost
}

func calculateEvalResponsesArrayCost(evalsArray []evalResponses) uint64 {
	var cost uint64
	for _, evals := range evalsArray {
		cost += calculateEvalResponsesCost(evals)
	}
	return cost
}

func generateEvalResponse(matchConditionsVariableNames []string, matchConditionsVariableLazyEvals lazyEvalMap, matchConditionsEvals evalResponses,
	validationVariableNames []string, validationVariableLazyEvals lazyEvalMap, validationEvals evalResponses,
	auditAnnotationEvals evalResponses, webhookMatchConditionsEvals []evalResponses) *EvalResponse {

	cost := calculateLazyEvalCost(matchConditionsVariableLazyEvals)
	cost += calculateEvalResponsesCost(matchConditionsEvals)
	cost += calculateLazyEvalCost(validationVariableLazyEvals)
	cost += calculateEvalResponsesCost(validationEvals)
	cost += calculateEvalResponsesCost(auditAnnotationEvals)
	cost += calculateEvalResponsesArrayCost(webhookMatchConditionsEvals)

	return &EvalResponse{
		MatchConditionsVariables: generateEvalVariables(matchConditionsVariableNames, matchConditionsVariableLazyEvals),
		MatchConditions:          generateEvalResults(matchConditionsEvals),
		ValidationVariables:      generateEvalVariables(validationVariableNames, validationVariableLazyEvals),
		Validations:              generateEvalResults(validationEvals),
		AuditAnnotations:         generateEvalResults(auditAnnotationEvals),
		WebhookMatchConditions:   generateEvalArrayResults(webhookMatchConditionsEvals),
		Cost:                     &cost,
	}
}
