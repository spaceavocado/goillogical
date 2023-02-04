package comparison

import (
	. "goillogical/internal"
	"testing"
)

type mock struct {
	val bool
	str string
}

func (m mock) String() string {
	return m.str
}

func (m mock) Evaluate(ctx Context) (any, error) {
	return m.val, nil
}

func TestEvaluate(t *testing.T) {
	e1 := mock{true, "e1"}
	e2 := mock{false, "e2"}

	var tests = []struct {
		op       string
		operands []Evaluable
		expected bool
	}{
		{"==", []Evaluable{e1, e1}, true},
		{"==", []Evaluable{e1, e2}, false},
	}

	for _, test := range tests {
		c, _ := New(test.op, test.operands, func(evaluated []any) bool { return evaluated[0] == evaluated[1] })
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.op, test.operands, test.expected, output, err)
		}
	}
}

func TestString(t *testing.T) {
	e1 := mock{true, "e1"}
	e2 := mock{false, "e2"}

	var tests = []struct {
		op       string
		operands []Evaluable
		expected string
	}{
		{"==", []Evaluable{e1, e2}, "e1 == e2"},
		{"==", []Evaluable{e1, e2, e1}, "e1 == e2, e1"},
	}

	for _, test := range tests {
		c, _ := New(test.op, test.operands, func(evaluated []any) bool { return false })
		if output := c.String(); output != test.expected {
			t.Errorf("input (%v, %v): expected %v, got %v", test.op, test.operands, test.expected, output)
		}
	}
}
