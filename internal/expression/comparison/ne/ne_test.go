package ne

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
		// Same types
		{e(1), e(0), true},
		{e(1), e(1), false},
		{e(1.1), e(1.0), true},
		{e(1.1), e(1.1), false},
		{e("1"), e("2"), true},
		{e("1"), e("1"), false},
		{e(true), e(false), true},
		{e(true), e(true), false},
		// Diff types
		{e(1), e(1.1), true},
		{e(1), e("1"), true},
		{e(1), e(true), true},
		{e(1.1), e("1"), true},
		{e(1.1), e(true), true},
		{e("1"), e(true), true},
		// Slices
		{e([]int{1}), e([]int{1}), true},
		{e(1), e([]int{1}), true},
	}

	for _, test := range tests {
		c, _ := New(test.left, test.right)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.left.String(), test.right.String(), test.expected, output, err)
		}
	}
}
