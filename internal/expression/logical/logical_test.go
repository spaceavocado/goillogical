package logical

import (
	. "goillogical/internal"
	"testing"
)

func TestEvaluate(t *testing.T) {
	e1 := EvaluableMock(true, "e1")
	e2 := EvaluableMock(false, "e2")
	e3 := EvaluableMock("bogus", "e3")

	var tests = []struct {
		evaluable Evaluable
		expected  bool
	}{
		{e1, true},
		{e2, false},
		{e3, false},
	}

	for _, test := range tests {
		if output, err := Evaluate(map[string]any{}, test.evaluable); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.evaluable, test.expected, output, err)
		}
	}
}

func TestString(t *testing.T) {
	e1 := EvaluableMock(true, "e1")
	e2 := EvaluableMock(false, "e2")

	var tests = []struct {
		op       string
		operands []Evaluable
		expected string
	}{
		{"AND", []Evaluable{e1, e2}, "(e1 AND e2)"},
		{"AND", []Evaluable{e1, e2, e1}, "(e1 AND e2 AND e1)"},
	}

	for _, test := range tests {
		c, _ := New("Unknown", test.op, test.operands, func(ctx Context, evaluated []Evaluable) (bool, error) { return false, nil })
		if output := c.String(); output != test.expected {
			t.Errorf("input (%v, %v): expected %v, got %v", test.op, test.operands, test.expected, output)
		}
	}
}
