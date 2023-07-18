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

package main

import (
	"fmt"

	"gopkg.in/yaml.v3"

	evalcel "github.com/undistro/cel-playground/eval"
)

func main() {
}

// export eval
func eval(exp, is string) any {
	var input map[string]any
	if err := yaml.Unmarshal([]byte(is), &input); err != nil {
		return response("", fmt.Errorf("failed to decode input: %w", err))
	}
	output, err := evalcel.Eval(exp, input)
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
