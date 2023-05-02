package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveFieldFilter(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]any
		fields []string
		expect map[string]any
	}{
		{
			name: "remove field",
			input: map[string]any{
				"a": "aa",
				"b": `{"c":"cc"}`,
			},
			fields: []string{"b"},
			expect: map[string]any{
				"a": "aa",
			},
		},
		{
			name: "remove field",
			input: map[string]any{
				"a": "aa",
				"b": `{"c":"cc"}`,
			},
			fields: []string{"c"},
			expect: map[string]any{
				"a": "aa",
				"b": `{"c":"cc"}`,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := RemoveFieldFilter(test.fields)(test.input)
			assert.EqualValues(t, test.expect, actual)
		})
	}
}
