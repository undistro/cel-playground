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
	"strings"
	"testing"
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
		want    string
		wantErr bool
	}{
		{
			name: "lte",
			exp:  "object.replicas <= 5",
			want: "true",
		},
		{
			name:    "error",
			exp:     "object.",
			wantErr: true,
		},
		{
			name: "url",
			exp:  "isURL(object.href) && url(object.href).getScheme() == 'https' && url(object.href).getEscapedPath() == '/path'",
			want: "true",
		},
		{
			name: "query",
			exp:  "url(object.href).getQuery()",
			want: `{"query": ["val"]}`,
		},
		{
			name: "regex",
			exp:  "object.image.find('v[0-9]+.[0-9]+.[0-9]*$')",
			want: `"v0.0.0"`,
		},
		{
			name: "list",
			exp:  "object.items.isSorted() && object.items.sum() == 6 && object.items.max() == 3 && object.items.indexOf(1) == 0",
			want: "true",
		},
		{
			name: "optional",
			exp:  `object.?foo.orValue("fallback")`,
			want: `"fallback"`,
		},
		{
			name: "strings",
			exp:  "object.abc.join(', ')",
			want: `"a, b, c"`,
		},
		{
			name: "cross type numeric comparisons",
			exp:  "object.replicas > 1.4",
			want: "true",
		},
		{
			name: "split",
			exp:  "object.image.split(':').size() == 2",
			want: "true",
		},
		{
			name: "quantity",
			exp:  `isQuantity(object.memory) && quantity(object.memory).add(quantity("700M")).sub(1).isLessThan(quantity("2G"))`,
			want: "true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Eval(tt.exp, input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if stripWhitespace(got) != stripWhitespace(tt.want) {
				t.Errorf("Eval() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func stripWhitespace(a string) string {
	a = strings.ReplaceAll(a, " ", "")
	a = strings.ReplaceAll(a, "\n", "")
	a = strings.ReplaceAll(a, "\t", "")
	return strings.ReplaceAll(a, "\r", "")
}
