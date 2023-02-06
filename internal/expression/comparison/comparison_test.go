package comparison

import (
	. "goillogical/internal"
	. "goillogical/internal/mock"
	"testing"
)

func TestEvaluate(t *testing.T) {
	e1 := E(true, "e1")
	e2 := E(false, "e2")

	var tests = []struct {
		op       string
		operands []Evaluable
		expected bool
	}{
		{"==", []Evaluable{e1, e1}, true},
		{"==", []Evaluable{e1, e2}, false},
	}

	for _, test := range tests {
		c, _ := New("Unknown", test.op, test.operands, func(evaluated []any) bool { return evaluated[0] == evaluated[1] })
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.op, test.operands, test.expected, output, err)
		}
	}
}

func TestSerialize(t *testing.T) {
	e1 := E("e1", "e1")
	e2 := E("e2", "e2")

	var tests = []struct {
		op       string
		operands []Evaluable
		expected any
	}{
		{"->", []Evaluable{e1, e2}, []any{"->", "e1", "e2"}},
		{"X", []Evaluable{e1}, []any{"X", "e1"}},
	}

	for _, test := range tests {
		c, _ := New(test.op, test.op, test.operands, func(evaluated []any) bool { return false })
		if output := c.Serialize(); Fprint(output) != Fprint(test.expected) {
			t.Errorf("input (%v, %v): expected %v, got %v", test.op, test.operands, test.expected, output)
		}
	}
}

func TestString(t *testing.T) {
	e1 := E(true, "e1")
	e2 := E(false, "e2")

	var tests = []struct {
		op       string
		operands []Evaluable
		expected string
	}{
		{"==", []Evaluable{e1, e2}, "(e1 == e2)"},
		{"<nil>", []Evaluable{e1}, "(e1 <nil>)"},
	}

	for _, test := range tests {
		c, _ := New("Unknown", test.op, test.operands, func(evaluated []any) bool { return false })
		if output := c.String(); output != test.expected {
			t.Errorf("input (%v, %v): expected %v, got %v", test.op, test.operands, test.expected, output)
		}
	}
}
