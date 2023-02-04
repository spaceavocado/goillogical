package xor

import (
	"errors"
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
		operands []Evaluable
		expected bool
	}{
		// Truthy
		{[]Evaluable{e(true), e(false)}, true},
		{[]Evaluable{e(false), e(true)}, true},
		{[]Evaluable{e(false), e(true), e(false)}, true},
		// Falsy
		{[]Evaluable{e(true), e(false), e(true)}, false},
		{[]Evaluable{e(false), e(false)}, false},
		{[]Evaluable{e(true), e(true)}, false},
	}

	for _, test := range tests {
		c, _ := New(test.operands)
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.operands, test.expected, output, err)
		}
	}

	var errs = []struct {
		operands []Evaluable
		expected error
	}{
		// Truthy
		{[]Evaluable{}, errors.New("logical XOR expression must have at least 2 operands")},
		{[]Evaluable{e(true)}, errors.New("logical XOR expression must have at least 2 operands")},
	}

	for _, test := range errs {

		if _, err := New(test.operands); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.operands, test.expected, err)
		}
	}
}
