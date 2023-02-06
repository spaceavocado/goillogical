package ne

import (
	"fmt"
	. "goillogical/internal"
	. "goillogical/internal/mock"
	"testing"
)

func e(val any) Evaluable {
	return E(val, fmt.Sprintf("%v", val))
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
		c, _ := New("!=", test.left, test.right)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.left.String(), test.right.String(), test.expected, output, err)
		}
	}
}
