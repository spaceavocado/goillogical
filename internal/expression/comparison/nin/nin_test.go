package nin

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
		{Val(0), Col(Val(1)), true},
		{Col(Val(1)), Val(0), true},
		{Val("0"), Col(Val("1")), true},
		{Val(false), Col(Val(true)), true},
		{Val(1.0), Col(Val(1.1)), true},
		{Val(1), Val(1), true},
		{Col(Val(1)), Col(Val(1)), true},
		{Val(1), Col(Val("1")), true},
		// Falsy
		{Val(1), Col(Val(1)), false},
		{Col(Val(1)), Val(1), false},
	}

	for _, test := range tests {
		c, _ := New("NON IT", test.left, test.right)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.left.String(), test.right.String(), test.expected, output, err)
		}
	}
}
