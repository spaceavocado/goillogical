package value

import (
	"errors"
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{1, 1},
		{1.1, 1.1},
		{"val", "val"},
		{true, true},
		{false, false},
	}

	for _, test := range tests {
		e, _ := New(test.input)
		if value, err := e.Evaluate(map[string]any{}); value != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v", test.input, value, err)
		}
	}

	errs := []struct {
		input    any
		expected error
	}{
		{nil, errors.New("value could be only primitive type, string, number or bool")},
	}
	for _, test := range errs {
		if _, err := New(test.input); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		input    any
		expected string
	}{
		{1, "1"},
		{1.1, "1.1"},
		{"val", "\"val\""},
		{true, "true"},
		{false, "false"},
	}

	for _, test := range tests {
		e, _ := New(test.input)
		if value := e.String(); value != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, value)
		}
	}
}
