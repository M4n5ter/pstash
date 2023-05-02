package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferFilter(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]any
		field  string
		target string
		expect map[string]any
	}{
		{
			name: "with target",
			input: map[string]any{
				"a": "aa",
				"b": `{"c":"cc"}`,
			},
			field:  "b",
			target: "data",
			expect: map[string]any{
				"a": "aa",
				"data": map[string]any{
					"c": "cc",
				},
			},
		},
		{
			name: "without target",
			input: map[string]any{
				"a": "aa",
				"b": `{"c":"cc"}`,
			},
			field: "b",
			expect: map[string]any{
				"a": "aa",
				"c": "cc",
			},
		},
		{
			name: "without field",
			input: map[string]any{
				"a": "aa",
				"b": `{"c":"cc"}`,
			},
			field: "c",
			expect: map[string]any{
				"a": "aa",
				"b": `{"c":"cc"}`,
			},
		},
		{
			name: "with not json",
			input: map[string]any{
				"a": "aa",
				"b": `{"c":"cc"`,
			},
			field: "b",
			expect: map[string]any{
				"a": "aa",
				"b": `{"c":"cc"`,
			},
		},
		{
			name: "with not string",
			input: map[string]any{
				"a": "aa",
				"b": map[string]any{"c": "cc"},
			},
			field: "b",
			expect: map[string]any{
				"a": "aa",
				"b": map[string]any{"c": "cc"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := TransferFilter(test.field, test.target)(test.input)
			assert.EqualValues(t, test.expect, actual)
		})
	}
}
