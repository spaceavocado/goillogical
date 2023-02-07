package suffix

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
		// Truthy
		{Val("bogus"), Val("us"), true},
		// Falsy
		{Val("something"), Val("else"), false},
		// Diff types
		{Val(1), Val(1.1), false},
		{Val(1), Val("1"), false},
		{Val(1), Val(true), false},
		{Val(1.1), Val("1"), false},
		{Val(1.1), Val(true), false},
		{Val("1"), Val(true), false},
		// Slices
		{Col(Val(1)), Col(Val(1)), false},
		{Val(1), Col(Val(1)), false},
	}

	for _, test := range tests {
		c, _ := New("SUFFIX", test.left, test.right)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.left.String(), test.right.String(), test.expected, output, err)
		}
	}
}
