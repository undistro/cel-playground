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
	"reflect"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
)

var input = map[string]any{
	"object": map[string]any{
		"replicas": 2,
		"href":     "https://user:pass@example.com:80/path?query=val#fragment",
		"image":    "registry.com/image:v0.0.0",
		"items":    []int{1, 2, 3},
		"abc":      []string{"a", "b", "c"},
		"memory":   "1.3G",
	},
}

func TestEval(t *testing.T) {
	tests := []struct {
		name    string
		exp     string
		want    any
		wantErr bool
	}{
		{
			name: "lte",
			exp:  "object.replicas <= 5",
			want: true,
		},
		{
			name:    "error",
			exp:     "object.",
			wantErr: true,
		},
		{
			name: "url",
			exp:  "isURL(object.href) && url(object.href).getScheme() == 'https' && url(object.href).getEscapedPath() == '/path'",
			want: true,
		},
		{
			name: "query",
			exp:  "url(object.href).getQuery()",
			want: map[string]any{
				"query": []any{"val"},
			},
		},
		{
			name: "regex",
			exp:  "object.image.find('v[0-9]+.[0-9]+.[0-9]*$')",
			want: "v0.0.0",
		},
		{
			name: "list",
			exp:  "object.items.isSorted() && object.items.sum() == 6 && object.items.max() == 3 && object.items.indexOf(1) == 0",
			want: true,
		},
		{
			name: "optional",
			exp:  `object.?foo.orValue("fallback")`,
			want: "fallback",
		},
		{
			name: "strings",
			exp:  "object.abc.join(', ')",
			want: "a, b, c",
		},
		{
			name: "cross type numeric comparisons",
			exp:  "object.replicas > 1.4",
			want: true,
		},
		{
			name: "split",
			exp:  "object.image.split(':').size() == 2",
			want: true,
		},
		{
			name: "quantity",
			exp:  `isQuantity(object.memory) && quantity(object.memory).add(quantity("700M")).sub(1).isLessThan(quantity("2G"))`,
			want: true,
		},
		{
			name: "sets.contains test 1",
			exp:  `sets.contains([], [])`,
			want: true,
		},
		{
			name: "sets.contains test 2",
			exp:  `sets.contains([], [1])`,
			want: false,
		},
		{
			name: "sets.contains test 3",
			exp:  `sets.contains([1, 2, 3, 4], [2, 3])`,
			want: true,
		},
		{
			name: "sets.contains test 4",
			exp:  `sets.contains([1, 2, 3], [3, 2, 1])`,
			want: true,
		},
		{
			name: "sets.equivalent test 1",
			exp:  `sets.equivalent([], [])`,
			want: true,
		},
		{
			name: "sets.equivalent test 2",
			exp:  `sets.equivalent([1], [1, 1])`,
			want: true,
		},
		{
			name: "sets.equivalent test 3",
			exp:  `sets.equivalent([1], [1, 1])`,
			want: true,
		},
		{
			name: "sets.equivalent test 4",
			exp:  `sets.equivalent([1, 2, 3], [3, 2, 1])`,
			want: true,
		},

		{
			name: "sets.intersects test 1",
			exp:  `sets.intersects([1], [])`,
			want: false,
		},
		{
			name: "sets.intersects test 2",
			exp:  `sets.intersects([1], [1, 2])`,
			want: true,
		},
		{
			name: "sets.intersects test 3",
			exp:  `sets.intersects([[1], [2, 3]], [[1, 2], [2, 3]])`,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Eval(tt.exp, input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				evalResponse := EvalResponse{}
				if err := json.Unmarshal([]byte(got), &evalResponse); err != nil {
					t.Errorf("Eval() error = %v", err)
				}

				if !reflect.DeepEqual(tt.want, evalResponse.Result) {
					t.Errorf("Expected %v\n, received %v", tt.want, evalResponse.Result)
				}
				if evalResponse.Cost == nil || *evalResponse.Cost <= 0 {
					t.Errorf("Expected Cost, returned %v", evalResponse.Cost)
				}
			}
		})
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		exp     string
		wantErr bool
	}{
		// Duration Literals
		{
			name:    "Duration Validation test 1",
			exp:     `duration('1')`,
			wantErr: true,
		},
		{
			name:    "Duration Validation test 2",
			exp:     `duration('1d')`,
			wantErr: true,
		},
		{
			name:    "Duration Validation test 3",
			exp:     `duration('1us') < duration('1nns')`,
			wantErr: true,
		},
		{
			name: "Duration Validation test 4",
			exp:  `duration('2h3m4s5us')`,
		},
		{
			name: "Duration Validation test 5",
			exp:  `duration(x)`,
		},

		// Timestamp Literals
		{
			name:    "Timestamp Validation test 1",
			exp:     `timestamp('1000-00-00T00:00:00Z')`,
			wantErr: true,
		},
		{
			name:    "Timestamp Validation test 2",
			exp:     `timestamp('1000-01-01T00:00:00ZZ')`,
			wantErr: true,
		},
		{
			name: "Timestamp Validation test 3",
			exp:  `timestamp('1000-01-01T00:00:00Z')`,
		},
		{
			name: "Timestamp Validation test 4",
			exp:  `timestamp(-6213559680)`, // min unix epoch time.
		},
		{
			name:    "Timestamp Validation test 5",
			exp:     `timestamp(-62135596801)`,
			wantErr: true,
		},
		{
			name: "Timestamp Validation test 6",
			exp:  `timestamp(x)`,
		},

		// Regex Literals
		{
			name: "Regex Validation test 1",
			exp:  `'hello'.matches('el*')`,
		},
		{
			name:    "Regex Validation test 2",
			exp:     `'hello'.matches('x++')`,
			wantErr: true,
		},
		{
			name:    "Regex Validation test 3",
			exp:     `'hello'.matches('(?<name%>el*)')`,
			wantErr: true,
		},
		{
			name:    "Regex Validation test 4",
			exp:     `'hello'.matches('??el*')`,
			wantErr: true,
		},
		{
			name: "Regex Validation test 5",
			exp:  `'hello'.matches(x)`,
		},

		// Homogeneous Aggregate Literals
		{
			name:    "Homogeneous Aggregate Validation test 1",
			exp:     `name in ['hello', 0]`,
			wantErr: true,
		},
		{
			name:    "Homogeneous Aggregate Validation test 2",
			exp:     `{'hello':'world', 1:'!'}`,
			wantErr: true,
		},
		{
			name:    "Homogeneous Aggregate Validation test 3",
			exp:     `name in {'hello':'world', 'goodbye':true}`,
			wantErr: true,
		},
		{
			name: "Homogeneous Aggregate Validation test 4",
			exp:  `name in ['hello', 'world']`,
		},
		{
			name: "Homogeneous Aggregate Validation test 5",
			exp:  `name in ['hello', ?optional.ofNonZeroValue('')]`,
		},
		{
			name: "Homogeneous Aggregate Validation test 6",
			exp:  `name in [?optional.ofNonZeroValue(''), 'hello', ?optional.of('')]`,
		},
		{
			name: "Homogeneous Aggregate Validation test 7",
			exp:  `name in {'hello': false, 'world': true}`,
		},
		{
			name: "Homogeneous Aggregate Validation test 8",
			exp:  `{'hello': false, ?'world': optional.ofNonZeroValue(true)}`,
		},
		{
			name: "Homogeneous Aggregate Validation test 9",
			exp:  `{?'hello': optional.ofNonZeroValue(false), 'world': true}`,
		},
	}
	env, err := cel.NewEnv(append(celEnvOptions,
		cel.Variable("x", types.StringType),
		cel.Variable("name", types.StringType),
	)...)
	if err != nil {
		t.Errorf("failed to create CEL env: %v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, issues := env.Compile(tt.exp)
			if tt.wantErr {
				if issues.Err() == nil {
					t.Fatalf("Compilation should have failed, expr: %v", tt.exp)
				}
			} else if issues.Err() != nil {
				t.Fatalf("Compilation failed, expr: %v, error: %v", tt.exp, issues.Err())
			}
		})
	}
}
