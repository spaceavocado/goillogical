package notin

import (
	"fmt"
	. "goillogical/internal"
	"testing"
)

type mock struct {
	val any
}

func (m mock) String() string {
	return fmt.Sprintf("%v", m.val)
}

func (m mock) Evaluate(ctx Context) (any, error) {
	return m.val, nil
}

func e(val any) Evaluable {
	return mock{val}
}

func TestHandler(t *testing.T) {
	var tests = []struct {
		left     Evaluable
		right    Evaluable
		expected bool
	}{
		// Truthy
		{e(0), e([]any{1}), true},
		{e([]any{1}), e(0), true},
		{e("0"), e([]any{"1"}), true},
		{e(false), e([]any{true}), true},
		{e(1.0), e([]any{1.1}), true},
		{e(1), e(1), true},
		{e([]any{1}), e([]any{1}), true},
		{e(1), e([]any{"1"}), true},
		// Falsy
		{e(1), e([]any{1}), false},
		{e([]any{1}), e(1), false},
	}

	for _, test := range tests {
		c, _ := New(test.left, test.right)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.left.String(), test.right.String(), test.expected, output, err)
		}
	}
}
