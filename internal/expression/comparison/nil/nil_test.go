package nil

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
		{e(nil), true},
		// Falsy
		{e(1), false},
		{e(1.1), false},
		{e("1"), false},
		{e(true), false},
		{e(false), false},
		{e([]int{1}), false},
	}

	for _, test := range tests {
		c, _ := New("NIL", test.eval)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.eval.String(), test.expected, output, err)
		}
	}
}
