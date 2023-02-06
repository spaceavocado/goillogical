package in

import (
	"fmt"
	. "goillogical/internal"
	"testing"
)

func e(val any) Evaluable {
	return EvaluableMock(val, fmt.Sprintf("%v", val))
}

func TestHandler(t *testing.T) {
	var tests = []struct {
		left     Evaluable
		right    Evaluable
		expected bool
	}{
		// Truthy
		{e(1), e([]any{1}), true},
		{e([]any{1}), e(1), true},
		{e("1"), e([]any{"1"}), true},
		{e(true), e([]any{true}), true},
		{e(1.1), e([]any{1.1}), true},
		// Falsy
		{e(1), e([]any{2}), false},
		{e([]any{2}), e(1), false},
		{e(1), e(1), false},
		{e([]any{1}), e([]any{1}), false},
		{e(1), e([]any{"1"}), false},
	}

	for _, test := range tests {
		c, _ := New(test.left, test.right)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.left.String(), test.right.String(), test.expected, output, err)
		}
	}
}
