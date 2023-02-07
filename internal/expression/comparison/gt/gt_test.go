package gt

import (
	. "goillogical/internal"
	. "goillogical/internal/mock"
	"testing"
)

func TestHandler(t *testing.T) {
	var tests = []struct {
		left     Evaluable
		right    Evaluable
		expected bool
	}{
		// Truthy
		{Val(2), Val(1), true},
		{Val(1.2), Val(1.1), true},
		// Falsy
		{Val(1), Val(1), false},
		{Val(1.1), Val(1.1), false},
		{Val(0), Val(1), false},
		{Val(1.0), Val(1.1), false},
		// Non comparable
		{Val(1.1), Val(1), false},
	}

	for _, test := range tests {
		c, _ := New(">", test.left, test.right)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.left.String(), test.right.String(), test.expected, output, err)
		}
	}
}
