package prefix

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
		{e("bo"), e("bogus"), true},
		// Falsy
		{e("bo"), e("something"), false},
		// Diff types
		{e(1), e(1.1), false},
		{e(1), e("1"), false},
		{e(1), e(true), false},
		{e(1.1), e("1"), false},
		{e(1.1), e(true), false},
		{e("1"), e(true), false},
		// Slices
		{e([]int{1}), e([]int{1}), false},
		{e(1), e([]int{1}), false},
	}

	for _, test := range tests {
		c, _ := New("PREFIX", test.left, test.right)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.left.String(), test.right.String(), test.expected, output, err)
		}
	}
}
