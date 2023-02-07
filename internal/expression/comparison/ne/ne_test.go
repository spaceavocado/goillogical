package ne

import (
	"testing"

	. "github.com/spaceavocado/goillogical/internal"
	. "github.com/spaceavocado/goillogical/internal/mock"
)

func TestHandler(t *testing.T) {
	var tests = []struct {
		left     Evaluable
		right    Evaluable
		expected bool
	}{
		// Same types
		{Val(1), Val(0), true},
		{Val(1), Val(1), false},
		{Val(1.1), Val(1.0), true},
		{Val(1.1), Val(1.1), false},
		{Val("1"), Val("2"), true},
		{Val("1"), Val("1"), false},
		{Val(true), Val(false), true},
		{Val(true), Val(true), false},
		// Diff types
		{Val(1), Val(1.1), true},
		{Val(1), Val("1"), true},
		{Val(1), Val(true), true},
		{Val(1.1), Val("1"), true},
		{Val(1.1), Val(true), true},
		{Val("1"), Val(true), true},
		// Slices
		{Col(Val(1)), Col(Val(1)), true},
		{Val(1), Col(Val(1)), true},
	}

	for _, test := range tests {
		c, _ := New("!=", test.left, test.right)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.left.String(), test.right.String(), test.expected, output, err)
		}
	}
}
