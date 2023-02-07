package goillogical

import (
	"errors"
	"testing"

	. "github.com/spaceavocado/goillogical/internal"
	eq "github.com/spaceavocado/goillogical/internal/expression/comparison/eq"
	and "github.com/spaceavocado/goillogical/internal/expression/logical/and"
	. "github.com/spaceavocado/goillogical/internal/mock"
)

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
		{1, Val(1)},
		{true, Val(true)},
		{"val", Val("val")},
		{"$refA", Ref("refA")},
		{[]any{"==", 1, 1}, ExpBinary("OP", eq.New, Val(1), Val(1))},
		{[]any{"==", "$refA", "resolvedA"}, ExpBinary("OP", eq.New, Ref("refA"), Val("resolvedA"))},
		{[]any{"AND", []any{"==", 1, 1}, []any{"==", 2, 1}}, ExpMany("OP", and.New, ExpBinary("OP", eq.New, Val(1), Val(1)), ExpBinary("OP", eq.New, Val(2), Val(1)))},
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

func TestWithOperatorMappingOptions(t *testing.T) {
	illogical := New(WithOperatorMappingOptions(map[Kind]string{Eq: "IS"}))
	ctx := map[string]any{
		"refA": "resolvedA",
	}

	var tests1 = []struct {
		input    any
		expected any
	}{
		{[]any{"IS", 1, 1}, true},
		{[]any{"IS", "$refA", "resolvedA"}, true},
	}

	for _, test := range tests1 {
		if output, err := illogical.Evaluate(test.input, ctx); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}

	var tests2 = []struct {
		input    any
		expected string
	}{
		{[]any{"IS", 1, 1}, "(1 == 1)"},
		{[]any{"IS", "$refA", "resolvedA"}, "({refA} == \"resolvedA\")"},
	}

	for _, test := range tests2 {
		if output, err := illogical.Statement(test.input); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}
}
