package logical

import (
	"testing"

	"errors"

	. "github.com/spaceavocado/goillogical/evaluable"
	. "github.com/spaceavocado/goillogical/internal/mock"
	. "github.com/spaceavocado/goillogical/internal/test"
)

func TestEvaluate(t *testing.T) {
	var tests = []struct {
		op       string
		operands []Evaluable
		expected bool
	}{
		{"AND", []Evaluable{Val(true), Val(true)}, true},
		{"AND", []Evaluable{Val(false), Val(false)}, false},
	}

	handler := func(ctx Context, operands []Evaluable) (bool, error) {
		return Evaluate(ctx, operands[0])
	}
	simplify := func(string, Context, []Evaluable) (any, Evaluable) { return nil, nil }

	for _, test := range tests {
		l, _ := New("Unknown", test.op, test.operands, handler, simplify)
		if output, err := l.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.operands, test.expected, output, err)
		}
	}

	errs := []struct {
		evaluable Evaluable
		expected  error
	}{
		{Invalid(), errors.New("invalid")},
	}

	for _, test := range errs {
		if _, err := Evaluate(map[string]any{}, test.evaluable); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.evaluable, test.expected, err)
		}
	}
}

func TestSimplify(t *testing.T) {
	var tests = []struct {
		op       string
		operands []Evaluable
		expected bool
	}{
		{"AND", []Evaluable{Val(true), Val(true)}, true},
	}

	handler := func(ctx Context, operands []Evaluable) (bool, error) {
		return true, nil
	}
	simplify := func(string, Context, []Evaluable) (any, Evaluable) { return true, nil }

	for _, test := range tests {
		l, _ := New("Unknown", test.op, test.operands, handler, simplify)
		if output, err := l.Simplify(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v/%v", test.operands, test.expected, output, err)
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
		l, _ := New(test.op, test.op, test.operands, func(Context, []Evaluable) (bool, error) { return false, nil }, func(string, Context, []Evaluable) (any, Evaluable) { return nil, nil })
		if output := l.Serialize(); Fprint(output) != Fprint(test.expected) {
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
