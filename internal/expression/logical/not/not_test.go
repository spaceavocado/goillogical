package not

import (
	"errors"
	"fmt"
	. "goillogical/internal"
	"testing"
)

func e(val any) Evaluable {
	return EvaluableMock(val, fmt.Sprintf("%v", val))
}
func TestHandler(t *testing.T) {
	var tests = []struct {
		operand  Evaluable
		expected bool
	}{
		{e(true), false},
		{e(false), true},
	}

	for _, test := range tests {
		c, _ := New(test.operand)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.operand, test.expected, output, err)
		}
	}

	var errs = []struct {
		operand  Evaluable
		expected error
	}{
		{e("bogus"), errors.New("logical NOT expression's operand must be evaluated to boolean value")},
	}

	for _, test := range errs {
		c, _ := New(test.operand)
		if _, err := c.Evaluate(map[string]any{}); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.operand, test.expected, err)
		}
	}
}
