package or

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
		operands []Evaluable
		expected bool
	}{
		// Truthy
		{[]Evaluable{e(true), e(false)}, true},
		{[]Evaluable{e(false), e(false), e(true)}, true},
		// Falsy
		{[]Evaluable{e(false), e(false)}, false},
	}

	for _, test := range tests {
		c, _ := New("OR", test.operands)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.operands, test.expected, output, err)
		}
	}

	var errs = []struct {
		operands []Evaluable
		expected error
	}{
		// Truthy
		{[]Evaluable{}, errors.New("logical OR expression must have at least 2 operands")},
		{[]Evaluable{e(true)}, errors.New("logical OR expression must have at least 2 operands")},
	}

	for _, test := range errs {

		if _, err := New("OR", test.operands); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.operands, test.expected, err)
		}
	}
}
