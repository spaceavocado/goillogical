package goillogical

import (
	. "goillogical/internal/options"
	"testing"
)

func TestEvaluate(t *testing.T) {
	opts := DefaultOptions()
	illogical := New(opts)
	ctx := map[string]any{
		"refA": "resolvedA",
	}

	var tests = []struct {
		input    any
		expected any
	}{
		{1, 1},
		{true, true},
		{"val", "val"},
		{"$refA", "resolvedA"},
		// {nil, nil},
		{[]any{"==", 1, 1}, true},
		{[]any{"==", "$refA", "resolvedA"}, true},
		{[]any{"AND", []any{"==", 1, 1}, []any{"!=", 2, 1}}, true},
	}

	for _, test := range tests {
		if output, err := illogical.Evaluate(test.input, ctx); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}
}
