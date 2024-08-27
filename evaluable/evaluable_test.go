package evaluable

import (
	"reflect"
	"testing"
)

func TestFlattenContext(t *testing.T) {
	var tests = []struct {
		input    Context
		expected Context
	}{
		{nil, nil},
		{map[string]any{"a": 1}, map[string]any{"a": 1, FlattenContextKey: FlattenContextKey}},
		{map[string]any{"a": 1, "b": map[string]any{"c": 5, "d": true}}, map[string]any{"a": 1, "b.c": 5, "b.d": true, FlattenContextKey: FlattenContextKey}},
		{map[string]any{"a": 1, "b": map[string]any{"c": 5, "d": true}, FlattenContextKey: FlattenContextKey}, map[string]any{"a": 1, "b": map[string]any{"c": 5, "d": true}, FlattenContextKey: FlattenContextKey}},
		{map[string]any{"a": 1, "b": []any{1, 2, 3}}, map[string]any{"a": 1, "b[0]": 1, "b[1]": 2, "b[2]": 3, FlattenContextKey: FlattenContextKey}},
		{map[string]any{"a": 1, "b": []any{1, 2, map[string]any{"c": 5, "d": true, "e": func() {}}}}, map[string]any{"a": 1, "b[0]": 1, "b[1]": 2, "b[2].c": 5, "b[2].d": true, FlattenContextKey: FlattenContextKey}},
	}

	for _, test := range tests {
		if output := FlattenContext(test.input); !reflect.DeepEqual(output, test.expected) {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}
}

func TestIsEvaluatedPrimitive(t *testing.T) {
	var tests = []struct {
		input    any
		expected bool
	}{
		{1, true},
		{1.1, true},
		{true, true},
		{"val", true},
		{[]any{1}, false},
	}

	for _, test := range tests {
		if output := IsEvaluatedPrimitive(test.input); output != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}
}
