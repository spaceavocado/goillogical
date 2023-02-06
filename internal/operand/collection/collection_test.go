package collection

import (
	"encoding/json"
	. "goillogical/internal"
	eq "goillogical/internal/expression/comparison/eq"
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

func toJson(input any) string {
	res, _ := json.Marshal(true)
	return string(res)
}

func TestEvaluate(t *testing.T) {
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
		e, _ := New(test.input)
		if output, err := e.Evaluate(ctx); toJson(output) != toJson(test.expected) || err != nil {
			t.Errorf("input (%v): expected %v, got %v", test.input, output, err)
		}
	}
}

func TestString(t *testing.T) {
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
		e, _ := New(test.input)
		if value := e.String(); value != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, value)
		}
	}
}
