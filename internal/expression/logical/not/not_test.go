package not

import (
	"errors"
	"testing"

	. "github.com/spaceavocado/goillogical/evaluable"
	. "github.com/spaceavocado/goillogical/internal/mock"
	. "github.com/spaceavocado/goillogical/internal/test"
)

func TestHandler(t *testing.T) {
	var tests = []struct {
		operand  Evaluable
		expected bool
	}{
		{Val(true), false},
		{Val(false), true},
	}

	for _, test := range tests {
		l, _ := New("NOT", test.operand)
		if output, err := l.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.operand, test.expected, output, err)
		}
	}

	errs := []struct {
		operand  Evaluable
		expected error
	}{
		{Val("bogus"), errors.New("invalid evaluated operand, must be boolean value")},
	}

	for _, test := range errs {
		l, _ := New("NOT", test.operand)
		if _, err := l.Evaluate(map[string]any{}); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.operand, test.expected, err)
		}
	}

	errs = []struct {
		operand  Evaluable
		expected error
	}{
		{Invalid(), errors.New("invalid")},
		{Val(1), errors.New("invalid evaluated operand, must be boolean value")},
	}

	for _, test := range errs {
		l, _ := New("NOT", test.operand)
		if _, err := l.Evaluate(map[string]any{}); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.operand, test.expected, err)
		}
	}
}

func TestSimplify(t *testing.T) {
	ctx := map[string]any{
		"RefA": true,
	}

	exp := func(operand Evaluable) Evaluable {
		e, _ := New("NOT", operand)
		return e
	}

	var tests = []struct {
		input Evaluable
		value any
		e     any
	}{
		{Val(true), false, nil},
		{Val(false), true, nil},
		{Ref("RefA"), false, nil},
		{Ref("Missing"), nil, exp(Ref("Missing"))},
		{Invalid(), nil, nil},
	}

	for _, test := range tests {
		if value, self := exp(test.input).Simplify(ctx); Fprint(value) != Fprint(test.value) || Fprint(self) != Fprint(test.e) {
			t.Errorf("input (%v): expected %v/%v, got %v/%v", test.input, test.value, test.e, value, self)
		}
	}
}
