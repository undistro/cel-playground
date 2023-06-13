package eval

import "testing"

var input = map[string]any{
	"object": map[string]any{
		"replicas": 2,
		"href":     "https://user:pass@example.com:80/path?query=val#fragment",
		"image":    "registry.com/image:v0.0.0",
		"items":    []int{1, 2, 3},
		"abc":      []string{"a", "b", "c"},
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
			name: "regex",
			exp:  "object.image.find('v[0-9]+.[0-9]+.[0-9]*$')",
			want: "\"v0.0.0\"",
		},
		{
			name: "list",
			exp:  "object.items.isSorted() && object.items.sum() == 6 && object.items.max() == 3 && object.items.indexOf(1) == 0",
			want: "true",
		},
		{
			name: "optional",
			exp:  "object[?'foo'].orValue('fallback')",
			want: "\"fallback\"",
		},
		{
			name: "strings",
			exp:  "object.abc.join(', ')",
			want: "\"a, b, c\"",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Eval(tt.exp, input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Eval() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Eval() got = %v, want %v", got, tt.want)
			}
		})
	}
}
