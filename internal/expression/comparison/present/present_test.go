package present

import (
	"testing"

	. "github.com/spaceavocado/goillogical/evaluable"
	. "github.com/spaceavocado/goillogical/internal/mock"
)

func TestHandler(t *testing.T) {
	var tests = []struct {
		eval     Evaluable
		expected bool
	}{
		// Truthy
		{Val(1), true},
		{Val(1.1), true},
		{Val("1"), true},
		{Val(true), true},
		{Val(false), true},
		// Falsy
		{Ref("Missing"), false},
	}

	for _, test := range tests {
		c, _ := New("PRESENT", test.eval)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.eval.String(), test.expected, output, err)
		}
	}
}
