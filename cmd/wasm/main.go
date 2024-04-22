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

//go:build js && wasm

package main

import (
	"errors"
	"fmt"
	"syscall/js"

	"github.com/undistro/cel-playground/eval"
	"github.com/undistro/cel-playground/k8s"
)

type execFunction func(mode string, argMap js.Value) (string, error)

func getArg(value js.Value, name string) []byte {
	arg := value.Get(name)
	if arg.Type() == js.TypeString {
		return []byte(arg.String())
	}
	return []byte{}
}

var modeExecFns = map[string]execFunction{
	"cel": func(mode string, argMap js.Value) (string, error) {
		return eval.CelEval(
			getArg(argMap, "cel"),
			getArg(argMap, "dataInput"),
		)
	},
	"vap": func(mode string, argMap js.Value) (string, error) {
		return k8s.EvalValidatingAdmissionPolicy(
			getArg(argMap, "vap"),
			getArg(argMap, "dataOriginal"),
			getArg(argMap, "dataUpdated"),
			getArg(argMap, "dataNamespace"),
			getArg(argMap, "dataRequest"),
			getArg(argMap, "dataAuthorizer"),
		)
	},
	"webhooks": func(mode string, argMap js.Value) (string, error) {
		return k8s.EvalWebhook(
			getArg(argMap, "webhooks"),
			getArg(argMap, "dataOriginal"),
			getArg(argMap, "dataUpdated"),
			getArg(argMap, "dataRequest"),
			getArg(argMap, "dataAuthorizer"),
		)
	},
}

func main() {
	defer addFunction("eval", dynamicEvalWrapper).Release()
	<-make(chan bool)
}

func addFunction(name string, fn func(js.Value, []js.Value) any) js.Func {
	function := js.FuncOf(fn)
	js.Global().Set(name, function)
	return function
}

func dynamicEvalWrapper(_ js.Value, args []js.Value) any {
	if len(args) < 2 {
		return response("", errors.New("invalid arguments"))
	}
	if args[0].Type() != js.TypeString || args[1].Type() != js.TypeObject {
		return response("", errors.New("invalid argument types, expecting string and object"))
	}
	mode := args[0].String()
	fn, ok := modeExecFns[mode]
	if !ok {
		return response("", fmt.Errorf("unknown mode %s", mode))
	}

	output, err := fn(mode, args[1])
	if err != nil {
		return response("", err)
	}
	return response(output, nil)
}

func response(out string, err error) any {
	if err != nil {
		out = err.Error()
	}
	return map[string]any{"output": out, "isError": err != nil}
}
