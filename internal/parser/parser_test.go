package parser

import (
	"errors"
	"fmt"
	. "goillogical/internal"
	eq "goillogical/internal/expression/comparison/eq"
	ge "goillogical/internal/expression/comparison/ge"
	gt "goillogical/internal/expression/comparison/gt"
	in "goillogical/internal/expression/comparison/in"
	le "goillogical/internal/expression/comparison/le"
	lt "goillogical/internal/expression/comparison/lt"
	ne "goillogical/internal/expression/comparison/ne"
	null "goillogical/internal/expression/comparison/nil"
	nin "goillogical/internal/expression/comparison/nin"
	prefix "goillogical/internal/expression/comparison/prefix"
	present "goillogical/internal/expression/comparison/present"
	suffix "goillogical/internal/expression/comparison/suffix"
	and "goillogical/internal/expression/logical/and"
	nor "goillogical/internal/expression/logical/nor"
	not "goillogical/internal/expression/logical/not"
	or "goillogical/internal/expression/logical/or"
	xor "goillogical/internal/expression/logical/xor"
	. "goillogical/internal/mock"
	reference "goillogical/internal/operand/reference"
	. "goillogical/internal/options"
	"regexp"
	"testing"
)

func addr(val string, opts Options) string {
	return opts.Serialize.Reference.To(val)
}

func ref(val string) Evaluable {
	serOpts := reference.DefaultSerializeOptions()
	simOpts := reference.SimplifyOptions{
		IgnoredPaths:   []string{},
		IgnoredPathsRx: []regexp.Regexp{},
	}
	e, _ := reference.New(val, &serOpts, &simOpts)
	return e
}

func expUnary(op string, factory func(string, Evaluable) (Evaluable, error), eval Evaluable) Evaluable {
	e, _ := factory(op, eval)
	return e
}
func expBinary(op string, factory func(string, Evaluable, Evaluable) (Evaluable, error), left, right Evaluable) Evaluable {
	e, _ := factory(op, left, right)
	return e
}
func expMany(op string, factory func(string, []Evaluable) (Evaluable, error), operands ...Evaluable) Evaluable {
	e, _ := factory(op, operands)
	return e
}

func TestIsEscaped(t *testing.T) {
	var tests = []struct {
		input           string
		escapeCharacter string
		expected        bool
	}{
		{"\\expected", "\\", true},
		{"unexpected", "\\", false},
		{"\\expected", "", false},
	}

	for _, test := range tests {
		if output := isEscaped(test.input, test.escapeCharacter); output != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}
}

func TestToReferenceAddr(t *testing.T) {
	opts := reference.SerializeOptions{
		From: func(val string) (string, error) {
			if val == "expected" {
				return val, nil
			}
			return "", errors.New("unexpected")
		},
		To: func(val string) string { return "" },
	}

	var tests = []struct {
		input    any
		expected string
	}{
		{"expected", "expected"},
	}

	for _, test := range tests {
		if output, err := toReferenceAddr(test.input, &opts); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}

	var errs = []struct {
		input    any
		expected error
	}{
		{"unexpected", errors.New("unexpected")},
		{1, errors.New("invalid reference path")},
	}

	for _, test := range errs {
		if _, err := toReferenceAddr(test.input, &opts); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}

func TestValue(t *testing.T) {
	opts := DefaultOptions()
	parser := New(&opts)

	var tests = []struct {
		input    any
		expected Evaluable
	}{
		{1, Val(1)},
		{1.1, Val(1.1)},
		{"val", Val("val")},
		{true, Val(true)},
		// TODO: nil
	}

	for _, test := range tests {
		if output, err := parser.Parse(test.input); output.String() != test.expected.String() || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}
}

func TestReference(t *testing.T) {
	opts := DefaultOptions()
	parser := New(&opts)

	var tests = []struct {
		input    string
		expected Evaluable
	}{
		{addr("path", opts), ref("path")},
	}

	for _, test := range tests {
		if output, err := parser.Parse(test.input); output.String() != test.expected.String() || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}
}

func TestCollection(t *testing.T) {
	opts := DefaultOptions()
	parser := New(&opts)

	var tests = []struct {
		input    []any
		expected Evaluable
	}{
		{[]any{1}, Col(Val(1))},
		{[]any{"val"}, Col(Val("val"))},
		{[]any{true}, Col(Val(true))},
		{[]any{addr("ref", opts)}, Col(ref("ref"))},
		{[]any{1, "val", true, addr("ref", opts)}, Col(Val(1), Val("val"), Val(true), ref("ref"))},
		// escaped
		{[]any{fmt.Sprintf("%s%s", opts.Serialize.Collection.EscapeCharacter, opts.OperatorMapping[Eq]), 1}, Col(Val(opts.OperatorMapping[Eq]), Val(1))},
	}

	for _, test := range tests {
		if output, err := parser.Parse(test.input); output.String() != test.expected.String() || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}
}

func TestComparison(t *testing.T) {
	opts := DefaultOptions()
	parser := New(&opts)

	var tests = []struct {
		input    []any
		expected Evaluable
	}{
		{[]any{opts.OperatorMapping[Eq], 1, 1}, expBinary("OP", eq.New, Val(1), Val(1))},
		{[]any{opts.OperatorMapping[Ne], 1, 1}, expBinary("OP", ne.New, Val(1), Val(1))},
		{[]any{opts.OperatorMapping[Gt], 1, 1}, expBinary("OP", gt.New, Val(1), Val(1))},
		{[]any{opts.OperatorMapping[Ge], 1, 1}, expBinary("OP", ge.New, Val(1), Val(1))},
		{[]any{opts.OperatorMapping[Lt], 1, 1}, expBinary("OP", lt.New, Val(1), Val(1))},
		{[]any{opts.OperatorMapping[Le], 1, 1}, expBinary("OP", le.New, Val(1), Val(1))},
		{[]any{opts.OperatorMapping[In], 1, 1}, expBinary("OP", in.New, Val(1), Val(1))},
		{[]any{opts.OperatorMapping[Nin], 1, 1}, expBinary("OP", nin.New, Val(1), Val(1))},
		{[]any{opts.OperatorMapping[Nil], 1, 1}, expUnary("OP", null.New, Val(1))},
		{[]any{opts.OperatorMapping[Present], 1, 1}, expUnary("OP", present.New, Val(1))},
		{[]any{opts.OperatorMapping[Suffix], 1, 1}, expBinary("OP", suffix.New, Val(1), Val(1))},
		{[]any{opts.OperatorMapping[Prefix], 1, 1}, expBinary("OP", prefix.New, Val(1), Val(1))},
	}

	for _, test := range tests {
		if output, err := parser.Parse(test.input); output.String() != test.expected.String() || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}
}

func TestLogical(t *testing.T) {
	opts := DefaultOptions()
	parser := New(&opts)

	var tests = []struct {
		input    []any
		expected Evaluable
	}{
		{[]any{opts.OperatorMapping[And], true, true}, expMany("OP", and.New, Val(true), Val(true))},
		{[]any{opts.OperatorMapping[Or], true, true}, expMany("OP", or.New, Val(true), Val(true))},
		{[]any{opts.OperatorMapping[Nor], true, true}, expMany("OP", nor.New, Val(true), Val(true))},
		{[]any{opts.OperatorMapping[Xor], true, true}, expMany("OP", xor.New, Val(true), Val(true))},
		{[]any{opts.OperatorMapping[Not], true, true}, expUnary("OP", not.New, Val(true))},
	}

	for _, test := range tests {
		if output, err := parser.Parse(test.input); output.String() != test.expected.String() || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}
}

func TestInvalid(t *testing.T) {
	opts := DefaultOptions()
	parser := New(&opts)

	var tests = []struct {
		input    []any
		expected error
	}{
		{nil, errors.New("invalid undefined operand")},
		{[]any{}, errors.New("invalid undefined operand")},
		{[]any{struct{ int }{5}}, errors.New("invalid operand, {5}")},
	}

	for _, test := range tests {
		if _, err := parser.Parse(test.input); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}
