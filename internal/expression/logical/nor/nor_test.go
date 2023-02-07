package nor

import (
	"errors"
	. "goillogical/internal"
	not "goillogical/internal/expression/logical/not"
	. "goillogical/internal/mock"
	. "goillogical/internal/test"
	"testing"
)

func TestHandler(t *testing.T) {
	var tests = []struct {
		operands []Evaluable
		expected bool
	}{
		// Truthy
		{[]Evaluable{Val(false), Val(false)}, true},
		{[]Evaluable{Val(false), Val(false), Val(false)}, true},
		// Falsy
		{[]Evaluable{Val(true), Val(false)}, false},
		{[]Evaluable{Val(false), Val(true)}, false},
		{[]Evaluable{Val(true), Val(true)}, false},
	}

	for _, test := range tests {
		c, _ := New("NOR", test.operands, "NOT", "NOR")
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.operands, test.expected, output, err)
		}
	}

	var errs = []struct {
		operands []Evaluable
		expected error
	}{
		// Truthy
		{[]Evaluable{}, errors.New("logical NOR expression must have at least 2 operands")},
		{[]Evaluable{Val(true)}, errors.New("logical NOR expression must have at least 2 operands")},
	}

	for _, test := range errs {

		if _, err := New("NOR", test.operands, "NOT", "NOR"); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.operands, test.expected, err)
		}
	}
}

func TestSimplify(t *testing.T) {
	ctx := map[string]any{
		"RefA": true,
	}

	exp := func(operands ...Evaluable) Evaluable {
		e, _ := New("NOR", operands, "NOT", "NOR")
		return e
	}

	neg := func(operand Evaluable) Evaluable {
		e, _ := not.New("NOT", operand)
		return e
	}

	var tests = []struct {
		input []Evaluable
		value any
		e     any
	}{
		{[]Evaluable{Val(false), Val(false)}, true, nil},
		{[]Evaluable{Val(true), Val(true)}, false, nil},
		{[]Evaluable{Val(true), Val(false)}, false, nil},
		{[]Evaluable{Ref("RefA"), Val(false)}, false, nil},
		{[]Evaluable{Ref("Missing"), Val(true)}, false, nil},
		{[]Evaluable{Ref("Missing"), Val(false)}, nil, neg(Ref("Missing"))},
		{[]Evaluable{Ref("Missing"), Ref("Missing")}, nil, exp(Ref("Missing"), Ref("Missing"))},
	}

	for _, test := range tests {
		if value, self := exp(test.input...).Simplify(ctx); Fprint(value) != Fprint(test.value) || Fprint(self) != Fprint(test.e) {
			t.Errorf("input (%v): expected %v/%v, got %v/%v", test.input, test.value, test.e, value, self)
		}
	}
}
