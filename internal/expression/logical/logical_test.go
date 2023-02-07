package logical

import (
	"testing"

	. "github.com/spaceavocado/goillogical/internal"
	. "github.com/spaceavocado/goillogical/internal/mock"
	. "github.com/spaceavocado/goillogical/internal/test"
)

func TestEvaluate(t *testing.T) {
	var tests = []struct {
		evaluable Evaluable
		expected  bool
	}{
		{Val(true), true},
		{Val(false), false},
		{Val("val"), false},
	}

	for _, test := range tests {
		if output, err := Evaluate(map[string]any{}, test.evaluable); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.evaluable, test.expected, output, err)
		}
	}
}

func TestSerialize(t *testing.T) {
	var tests = []struct {
		op       string
		operands []Evaluable
		expected any
	}{
		{"->", []Evaluable{Val("e1"), Val("e2")}, []any{"->", "e1", "e2"}},
		{"X", []Evaluable{Val("e1")}, []any{"X", "e1"}},
	}

	for _, test := range tests {
		c, _ := New(test.op, test.op, test.operands, func(Context, []Evaluable) (bool, error) { return false, nil }, func(string, Context, []Evaluable) (any, Evaluable) { return nil, nil })
		if output := c.Serialize(); Fprint(output) != Fprint(test.expected) {
			t.Errorf("input (%v, %v): expected %v, got %v", test.op, test.operands, test.expected, output)
		}
	}
}

func TestString(t *testing.T) {
	var tests = []struct {
		op       string
		operands []Evaluable
		expected string
	}{
		{"AND", []Evaluable{Val("e1"), Val("e2")}, "(\"e1\" AND \"e2\")"},
		{"AND", []Evaluable{Val("e1"), Val("e2"), Val("e1")}, "(\"e1\" AND \"e2\" AND \"e1\")"},
	}

	for _, test := range tests {
		c, _ := New("Unknown", test.op, test.operands, func(ctx Context, evaluated []Evaluable) (bool, error) { return false, nil }, func(string, Context, []Evaluable) (any, Evaluable) { return nil, nil })
		if output := c.String(); output != test.expected {
			t.Errorf("input (%v, %v): expected %v, got %v", test.op, test.operands, test.expected, output)
		}
	}
}
