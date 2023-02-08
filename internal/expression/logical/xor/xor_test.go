package xor

import (
	"errors"
	"testing"

	. "github.com/spaceavocado/goillogical/evaluable"
	nor "github.com/spaceavocado/goillogical/internal/expression/logical/nor"
	not "github.com/spaceavocado/goillogical/internal/expression/logical/not"
	. "github.com/spaceavocado/goillogical/internal/mock"
	. "github.com/spaceavocado/goillogical/internal/test"
)

func TestHandler(t *testing.T) {
	var tests = []struct {
		operands []Evaluable
		expected bool
	}{
		// Truthy
		{[]Evaluable{Val(true), Val(false)}, true},
		{[]Evaluable{Val(false), Val(true)}, true},
		{[]Evaluable{Val(false), Val(true), Val(false)}, true},
		// Falsy
		{[]Evaluable{Val(true), Val(false), Val(true)}, false},
		{[]Evaluable{Val(false), Val(false)}, false},
		{[]Evaluable{Val(true), Val(true)}, false},
	}

	for _, test := range tests {
		c, _ := New("XOR", test.operands, "NOT", "NOR")
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
		{[]Evaluable{Val(true)}, errors.New("logical XOR expression must have at least 2 operands")},
	}

	for _, test := range errs {

		if _, err := New("XOR", test.operands, "NOT", "NOR"); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.operands, test.expected, err)
		}
	}
}

func TestSimplify(t *testing.T) {
	ctx := map[string]any{
		"RefA": true,
	}

	exp := func(operands ...Evaluable) Evaluable {
		e, _ := New("XOR", operands, "NOT", "NOR")
		return e
	}

	neg := func(operand Evaluable) Evaluable {
		e, _ := not.New("NOT", operand)
		return e
	}

	flip := func(operands ...Evaluable) Evaluable {
		e, _ := nor.New("NOR", operands, "NOT", "NOR")
		return e
	}

	var tests = []struct {
		input []Evaluable
		value any
		e     any
	}{
		{[]Evaluable{Val(false), Val(false)}, false, nil},
		{[]Evaluable{Ref("RefA"), Val(true)}, false, nil},
		{[]Evaluable{Ref("Missing"), Val(true), Ref("Missing")}, nil, flip(Ref("Missing"), Ref("Missing"))},
		{[]Evaluable{Ref("Missing"), Val(true), Val(false)}, nil, neg(Ref("Missing"))},
		{[]Evaluable{Ref("RefA"), Ref("RefA"), Val(true)}, false, nil},
	}

	for _, test := range tests {
		if value, self := exp(test.input...).Simplify(ctx); Fprint(value) != Fprint(test.value) || Fprint(self) != Fprint(test.e) {
			t.Errorf("input (%v): expected %v/%v, got %v/%v", test.input, test.value, test.e, value, self)
		}
	}
}
