package goillogical

import (
	"errors"
	. "goillogical/internal"
	eq "goillogical/internal/expression/comparison/eq"
	and "goillogical/internal/expression/logical/and"
	reference "goillogical/internal/operand/reference"
	value "goillogical/internal/operand/value"
	"testing"
)

func val(val any) Evaluable {
	e, _ := value.New(val)
	return e
}

func ref(val string) Evaluable {
	e, _ := reference.New(val)
	return e
}

func expBinary(factory func(Evaluable, Evaluable) (Evaluable, error), left, right Evaluable) Evaluable {
	e, _ := factory(left, right)
	return e
}

func expMany(factory func([]Evaluable) (Evaluable, error), operands ...Evaluable) Evaluable {
	e, _ := factory(operands)
	return e
}

func TestEvaluate(t *testing.T) {
	illogical := New()
	ctx := map[string]any{
		"refA": "resolvedA",
	}

	var tests = []struct {
		input    any
		expected any
	}{
		{1, 1},
		{true, true},
		{"val", "val"},
		{"$refA", "resolvedA"},
		{[]any{"==", 1, 1}, true},
		{[]any{"==", "$refA", "resolvedA"}, true},
		{[]any{"AND", []any{"==", 1, 1}, []any{"!=", 2, 1}}, true},
		{[]any{"NIL", "$refB"}, true},
		{[]any{"PRESENT", "$refB"}, false},
	}

	for _, test := range tests {
		if output, err := illogical.Evaluate(test.input, ctx); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}

	var errs = []struct {
		input    any
		expected error
	}{
		{nil, errors.New("unexpected input")},
		{struct{ int }{4}, errors.New("invalid operand, {4}")},
		{[]any{"==", struct{ int }{4}, 1}, errors.New("invalid operand, {4}")},
	}

	for _, test := range errs {
		if _, err := illogical.Evaluate(test.input, ctx); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}

func TestParse(t *testing.T) {
	illogical := New()

	var tests = []struct {
		input    any
		expected Evaluable
	}{
		{1, val(1)},
		{true, val(true)},
		{"val", val("val")},
		{"$refA", ref("refA")},
		{[]any{"==", 1, 1}, expBinary(eq.New, val(1), val(1))},
		{[]any{"==", "$refA", "resolvedA"}, expBinary(eq.New, ref("refA"), val("resolvedA"))},
		{[]any{"AND", []any{"==", 1, 1}, []any{"==", 2, 1}}, expMany(and.New, expBinary(eq.New, val(1), val(1)), expBinary(eq.New, val(2), val(1)))},
	}

	for _, test := range tests {
		if output, err := illogical.Parse(test.input); output.String() != test.expected.String() || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}

	var errs = []struct {
		input    any
		expected error
	}{
		{nil, errors.New("unexpected input")},
		{struct{ int }{4}, errors.New("invalid operand, {4}")},
		{[]any{"==", struct{ int }{4}, 1}, errors.New("invalid operand, {4}")},
	}

	for _, test := range errs {
		if _, err := illogical.Parse(test.input); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}

func TestStatement(t *testing.T) {
	illogical := New()

	var tests = []struct {
		input    any
		expected string
	}{
		{1, "1"},
		{true, "true"},
		{"val", "\"val\""},
		{"$refA", "{refA}"},
		{[]any{"==", 1, 1}, "(1 == 1)"},
		{[]any{"==", "$refA", "resolvedA"}, "({refA} == \"resolvedA\")"},
		{[]any{"AND", []any{"==", 1, 1}, []any{"!=", 2, 1}}, "((1 == 1) AND (2 != 1))"},
	}

	for _, test := range tests {
		if output, err := illogical.Statement(test.input); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}

	var errs = []struct {
		input    any
		expected error
	}{
		{nil, errors.New("unexpected input")},
		{struct{ int }{4}, errors.New("invalid operand, {4}")},
		{[]any{"==", struct{ int }{4}, 1}, errors.New("invalid operand, {4}")},
	}

	for _, test := range errs {
		if _, err := illogical.Statement(test.input); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}
