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
	collection "goillogical/internal/operand/collection"
	reference "goillogical/internal/operand/reference"
	value "goillogical/internal/operand/value"
	. "goillogical/internal/options"
	"testing"
)

func val(val any) Evaluable {
	e, _ := value.New(val)
	return e
}

func addr(val string, opts Options) string {
	return opts.Serialize.Reference.To(val)
}

func ref(val string) Evaluable {
	e, _ := reference.New(val)
	return e
}

func col(items ...Evaluable) Evaluable {
	e, _ := collection.New(items)
	return e
}

func expUnary(factory func(Evaluable) (Evaluable, error), eval Evaluable) Evaluable {
	e, _ := factory(eval)
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
		if output, err := toReferenceAddr(test.input, opts); output != test.expected || err != nil {
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
		if _, err := toReferenceAddr(test.input, opts); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}

func TestValue(t *testing.T) {
	opts := DefaultOptions()
	parser := New(opts)

	var tests = []struct {
		input    any
		expected Evaluable
	}{
		{1, val(1)},
		{1.1, val(1.1)},
		{"val", val("val")},
		{true, val(true)},
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
	parser := New(opts)

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
	parser := New(opts)

	var tests = []struct {
		input    []any
		expected Evaluable
	}{
		{[]any{1}, col(val(1))},
		{[]any{"val"}, col(val("val"))},
		{[]any{true}, col(val(true))},
		{[]any{addr("ref", opts)}, col(ref("ref"))},
		{[]any{1, "val", true, addr("ref", opts)}, col(val(1), val("val"), val(true), ref("ref"))},
		// escaped
		{[]any{fmt.Sprintf("%s%s", opts.Serialize.Collection.EscapeCharacter, opts.OperatorMapping[Eq]), 1}, col(val(opts.OperatorMapping[Eq]), val(1))},
	}

	for _, test := range tests {
		if output, err := parser.Parse(test.input); output.String() != test.expected.String() || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}
}

func TestComparison(t *testing.T) {
	opts := DefaultOptions()
	parser := New(opts)

	var tests = []struct {
		input    []any
		expected Evaluable
	}{
		{[]any{opts.OperatorMapping[Eq], 1, 1}, expBinary(eq.New, val(1), val(1))},
		{[]any{opts.OperatorMapping[Ne], 1, 1}, expBinary(ne.New, val(1), val(1))},
		{[]any{opts.OperatorMapping[Gt], 1, 1}, expBinary(gt.New, val(1), val(1))},
		{[]any{opts.OperatorMapping[Ge], 1, 1}, expBinary(ge.New, val(1), val(1))},
		{[]any{opts.OperatorMapping[Lt], 1, 1}, expBinary(lt.New, val(1), val(1))},
		{[]any{opts.OperatorMapping[Le], 1, 1}, expBinary(le.New, val(1), val(1))},
		{[]any{opts.OperatorMapping[In], 1, 1}, expBinary(in.New, val(1), val(1))},
		{[]any{opts.OperatorMapping[Nin], 1, 1}, expBinary(nin.New, val(1), val(1))},
		{[]any{opts.OperatorMapping[Nil], 1, 1}, expUnary(null.New, val(1))},
		{[]any{opts.OperatorMapping[Present], 1, 1}, expUnary(present.New, val(1))},
		{[]any{opts.OperatorMapping[Suffix], 1, 1}, expBinary(suffix.New, val(1), val(1))},
		{[]any{opts.OperatorMapping[Prefix], 1, 1}, expBinary(prefix.New, val(1), val(1))},
	}

	for _, test := range tests {
		if output, err := parser.Parse(test.input); output.String() != test.expected.String() || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}
}

func TestLogical(t *testing.T) {
	opts := DefaultOptions()
	parser := New(opts)

	var tests = []struct {
		input    []any
		expected Evaluable
	}{
		{[]any{opts.OperatorMapping[And], true, true}, expMany(and.New, val(true), val(true))},
		{[]any{opts.OperatorMapping[Or], true, true}, expMany(or.New, val(true), val(true))},
		{[]any{opts.OperatorMapping[Nor], true, true}, expMany(nor.New, val(true), val(true))},
		{[]any{opts.OperatorMapping[Xor], true, true}, expMany(xor.New, val(true), val(true))},
		{[]any{opts.OperatorMapping[Not], true, true}, expUnary(not.New, val(true))},
	}

	for _, test := range tests {
		if output, err := parser.Parse(test.input); output.String() != test.expected.String() || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.input, test.expected, output, err)
		}
	}
}

func TestInvalid(t *testing.T) {
	opts := DefaultOptions()
	parser := New(opts)

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
