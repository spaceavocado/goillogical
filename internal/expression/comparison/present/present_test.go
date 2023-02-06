package present

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
		eval     Evaluable
		expected bool
	}{
		// Truthy
		{e(1), true},
		{e(1.1), true},
		{e("1"), true},
		{e(true), true},
		{e(false), true},
		{e([]int{1}), true},
		// Falsy
		{e(nil), false},
	}

	for _, test := range tests {
		c, _ := New("PRESENT", test.eval)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.eval.String(), test.expected, output, err)
		}
	}
}
