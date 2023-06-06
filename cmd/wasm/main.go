package main

import (
	"errors"
	"fmt"
	"syscall/js"
)

func main() {
	js.Global().Set("eval", js.FuncOf(evalWrapper))
	<-make(chan bool)
}

// evalWrapper wraps the eval function with `syscall/js` parameters
func evalWrapper(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return errors.New("invalid arguments")
	}
	exp := args[0].String()
	inp := args[1].String()
	return eval(exp, inp)
}

// eval evaluates the cel expression against the given input
func eval(exp, input string) string {
	// TODO: evaluate CEL expression
	return fmt.Sprintf("%d\n%d", len(exp), len(input))
}
