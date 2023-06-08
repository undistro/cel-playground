package main

import (
	"errors"
	"fmt"
	"syscall/js"

	"gopkg.in/yaml.v3"

	"github.com/undistro/cel-playground/eval"
)

func main() {
	evalFunc := js.FuncOf(evalWrapper)
	js.Global().Set("eval", evalFunc)
	defer evalFunc.Release()
	<-make(chan bool)
}

// evalWrapper wraps the eval function with `syscall/js` parameters
func evalWrapper(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return errors.New("invalid arguments")
	}
	exp := args[0].String()
	is := args[1].String()

	var input map[string]any
	if err := yaml.Unmarshal([]byte(is), &input); err != nil {
		return fmt.Errorf("failed to decode input: %w", err)
	}
	output, err := eval.Eval(exp, input)
	if err != nil {
		output = err.Error()
	}
	return map[string]any{"output": output, "isError": err != nil}
}
