package test

import (
	"testing"

	. "github.com/spaceavocado/goillogical/internal/mock"
)

func TestFprint(t *testing.T) {
	var tests = []struct {
		input    any
		expected string
	}{
		{1, "1"},
		{"val", "\"val\""},
		{true, "true"},
		{[]any{1, "val"}, "[1,\"val\"]"},
		{struct {
			Age int
		}{35}, "{\"Age\":35}"},
		{Ref("Ref"), "{Ref}"},
	}

	for _, test := range tests {
		if output := Fprint(test.input); output != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}
}
