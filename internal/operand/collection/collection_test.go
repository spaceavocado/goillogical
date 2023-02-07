package collection

import (
	"regexp"
	"testing"

	. "github.com/spaceavocado/goillogical/internal"
	eq "github.com/spaceavocado/goillogical/internal/expression/comparison/eq"
	reference "github.com/spaceavocado/goillogical/internal/operand/reference"
	value "github.com/spaceavocado/goillogical/internal/operand/value"
	. "github.com/spaceavocado/goillogical/internal/test"
)

func val(val any) Evaluable {
	e, _ := value.New(val)
	return e
}

func ref(val string) Evaluable {
	serOpts := reference.DefaultSerializeOptions()
	simOpts := reference.SimplifyOptions{
		IgnoredPaths:   []string{"ignored"},
		IgnoredPathsRx: []regexp.Regexp{},
	}
	e, _ := reference.New(val, &serOpts, &simOpts)
	return e
}

func expBinary(factory func(string, Evaluable, Evaluable) (Evaluable, error), left, right Evaluable) Evaluable {
	e, _ := factory("EXP", left, right)
	return e
}

func TestEvaluate(t *testing.T) {
	opts := DefaultSerializeOptions()
	ctx := map[string]any{
		"RefA": "A",
	}
	tests := []struct {
		input    []Evaluable
		expected any
	}{
		{[]Evaluable{val(1)}, []any{1}},
		{[]Evaluable{val("1")}, []any{"1"}},
		{[]Evaluable{val(true)}, []any{true}},
		{[]Evaluable{ref("RefA")}, []any{"A"}},
		{[]Evaluable{val(1), ref("RefA")}, []any{1, "A"}},
		{[]Evaluable{expBinary(eq.New, val(1), val(1)), ref("RefA")}, []any{true, "A"}},
	}

	for _, test := range tests {
		e, _ := New(test.input, &opts)
		if output, err := e.Evaluate(ctx); Fprint(output) != Fprint(test.expected) || err != nil {
			t.Errorf("input (%v): expected %v, got %v", test.input, output, err)
		}
	}
}

func TestSerialize(t *testing.T) {
	opts := SerializeOptions{
		EscapedOperators: map[string]bool{"==": true},
		EscapeCharacter:  "\\",
	}
	tests := []struct {
		input    []Evaluable
		expected any
	}{
		{[]Evaluable{val(1)}, []any{1}},
		{[]Evaluable{val("1")}, []any{"1"}},
		{[]Evaluable{val(true)}, []any{true}},
		{[]Evaluable{ref("RefA")}, []any{"$RefA"}},
		{[]Evaluable{val(1), ref("RefA")}, []any{1, "$RefA"}},
		{[]Evaluable{expBinary(eq.New, val(1), val(1)), ref("RefA")}, []any{[]any{"EXP", 1, 1}, "$RefA"}},
		{[]Evaluable{val("=="), val(1), val(1)}, []any{"\\==", 1, 1}},
	}

	for _, test := range tests {
		e, _ := New(test.input, &opts)
		if value := e.Serialize(); Fprint(value) != Fprint(test.expected) {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, value)
		}
	}
}

func TestSimplify(t *testing.T) {
	serOpts := DefaultSerializeOptions()

	ctx := map[string]any{
		"RefA": "A",
	}

	col := func(items ...Evaluable) Evaluable {
		e, _ := New(items, &serOpts)
		return e
	}

	tests := []struct {
		input []Evaluable
		value any
		e     any
	}{
		{[]Evaluable{ref("RefB")}, nil, col(ref("RefB"))},
		{[]Evaluable{ref("RefA")}, []any{"A"}, nil},
		{[]Evaluable{ref("RefA"), ref("RefB")}, nil, col(ref("RefA"), ref("RefB"))},
	}

	for _, test := range tests {
		e, _ := New(test.input, &serOpts)
		if value, self := e.Simplify(ctx); Fprint(value) != Fprint(test.value) || Fprint(self) != Fprint(test.e) {
			t.Errorf("input (%v): expected %v/%v, got %v/%v", test.input, test.value, test.e, value, self)
		}
	}
}

func TestString(t *testing.T) {
	opts := DefaultSerializeOptions()
	tests := []struct {
		input    []Evaluable
		expected string
	}{
		{[]Evaluable{val(1)}, "[1]"},
		{[]Evaluable{val("1")}, "[\"1\"]"},
		{[]Evaluable{val(true)}, "[true]"},
		{[]Evaluable{ref("RefA")}, "[{RefA}]"},
		{[]Evaluable{val(1), ref("RefA")}, "[1, {RefA}]"},
		{[]Evaluable{expBinary(eq.New, val(1), val(1)), ref("RefA")}, "[(1 == 1), {RefA}]"},
	}

	for _, test := range tests {
		e, _ := New(test.input, &opts)
		if value := e.String(); value != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, value)
		}
	}
}

func TestShouldBeEscaped(t *testing.T) {
	opts := SerializeOptions{
		EscapedOperators: map[string]bool{"==": true},
		EscapeCharacter:  "\\",
	}
	tests := []struct {
		input    any
		expected bool
	}{
		{"==", true},
		{"!=", false},
		{nil, false},
		{true, false},
	}

	for _, test := range tests {
		if value := shouldBeEscaped(test.input, &opts); value != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, value)
		}
	}
}

func TestEscapeOperator(t *testing.T) {
	opts := SerializeOptions{
		EscapedOperators: map[string]bool{"==": true},
		EscapeCharacter:  "\\",
	}
	tests := []struct {
		input    string
		expected string
	}{
		{"==", "\\=="},
	}

	for _, test := range tests {
		if value := escapeOperator(test.input, &opts); value != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, value)
		}
	}
}
