package nil

import (
	"testing"

	. "github.com/spaceavocado/goillogical/internal"
	. "github.com/spaceavocado/goillogical/internal/mock"
)

func TestHandler(t *testing.T) {
	var tests = []struct {
		eval     Evaluable
		expected bool
	}{
		// Truthy
		{Ref("Missing"), true},
		// Falsy
		{Val(1), false},
		{Val(1.1), false},
		{Val("1"), false},
		{Val(true), false},
		{Val(false), false},
	}

	for _, test := range tests {
		c, _ := New("NIL", test.eval)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.eval.String(), test.expected, output, err)
		}
	}
}
