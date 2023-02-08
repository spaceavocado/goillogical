package overlap

import (
	"testing"

	. "github.com/spaceavocado/goillogical/evaluable"
	. "github.com/spaceavocado/goillogical/internal/mock"
)

func TestHandler(t *testing.T) {
	var tests = []struct {
		left     Evaluable
		right    Evaluable
		expected bool
	}{
		// Truthy
		{Col(Val(1)), Col(Val(1)), true},
		{Col(Val(1), Val(2)), Col(Val(1), Val(3)), true},
		{Col(Val(3), Val(2)), Col(Val(1), Val(2), Val(3)), true},
		{Col(Val("1")), Col(Val("1")), true},
		{Col(Val(true)), Col(Val(true)), true},
		{Col(Val(1.1)), Col(Val(1.1)), true},
		// Falsy
		{Val(1), Col(Val(1)), false},
		{Col(Val(1)), Val(1), false},
		{Val(1), Val(1), false},
		{Col(Val(1)), Col(Val(2)), false},
	}

	for _, test := range tests {
		c, _ := New("OVERLAP", test.left, test.right)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.left.String(), test.right.String(), test.expected, output, err)
		}
	}
}
