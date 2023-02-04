package present

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
		c, _ := New(test.eval)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.eval.String(), test.expected, output, err)
		}
	}
}
