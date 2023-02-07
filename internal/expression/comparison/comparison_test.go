package comparison

import (
	"testing"

	. "github.com/spaceavocado/goillogical/internal"
	. "github.com/spaceavocado/goillogical/internal/mock"
	. "github.com/spaceavocado/goillogical/internal/test"
)

func TestEvaluate(t *testing.T) {
	var tests = []struct {
		op       string
		operands []Evaluable
		expected bool
	}{
		{"==", []Evaluable{Val(1), Val(1)}, true},
		{"==", []Evaluable{Val(1), Val(2)}, false},
	}

	for _, test := range tests {
		c, _ := New("Unknown", test.op, test.operands, func(evaluated []any) bool { return evaluated[0] == evaluated[1] })
		if output, err := c.Evaluate(map[string]any{}); output != test.expected || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v/%v", test.op, test.operands, test.expected, output, err)
		}
	}
}

func TestSerialize(t *testing.T) {
	var tests = []struct {
		op       string
		operands []Evaluable
		expected any
	}{
		{"->", []Evaluable{Val(1), Val(2)}, []any{"->", 1, 2}},
		{"X", []Evaluable{Val(1)}, []any{"X", 1}},
	}

	for _, test := range tests {
		c, _ := New(test.op, test.op, test.operands, func(evaluated []any) bool { return false })
		if output := c.Serialize(); Fprint(output) != Fprint(test.expected) {
			t.Errorf("input (%v, %v): expected %v, got %v", test.op, test.operands, test.expected, output)
		}
	}
}

func TestSimplify(t *testing.T) {
	ctx := map[string]any{
		"RefA": "A",
	}

	eq := func(operands ...Evaluable) Evaluable {
		e, _ := New("Unknown", "==", operands, func(evaluated []any) bool { return evaluated[0] == evaluated[1] })
		return e
	}

	tests := []struct {
		input []Evaluable
		value any
		e     any
	}{
		{[]Evaluable{Val(0), Ref("Missing")}, nil, eq(Val(0), Ref("Missing"))},
		{[]Evaluable{Ref("Missing"), Val(0)}, nil, eq(Ref("Missing"), Val(0))},
		{[]Evaluable{Ref("Missing"), Ref("Missing")}, nil, eq(Ref("Missing"), Ref("Missing"))},
		{[]Evaluable{Val(0), Val(0)}, true, nil},
		{[]Evaluable{Val(0), Val(1)}, false, nil},
		{[]Evaluable{Val("A"), Ref("RefA")}, true, nil},
	}

	for _, test := range tests {
		e := eq(test.input...)
		if value, self := e.Simplify(ctx); Fprint(value) != Fprint(test.value) || Fprint(self) != Fprint(test.e) {
			t.Errorf("input (%v): expected %v/%v, got %v/%v", test.input, test.value, test.e, value, self)
		}
	}
}

func TestString(t *testing.T) {
	var tests = []struct {
		op       string
		operands []Evaluable
		expected string
	}{
		{"==", []Evaluable{Val(1), Val(2)}, "(1 == 2)"},
		{"<nil>", []Evaluable{Val(1)}, "(1 <nil>)"},
	}

	for _, test := range tests {
		c, _ := New("Unknown", test.op, test.operands, func(evaluated []any) bool { return false })
		if output := c.String(); output != test.expected {
			t.Errorf("input (%v, %v): expected %v, got %v", test.op, test.operands, test.expected, output)
		}
	}
}
