package reference

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	. "github.com/spaceavocado/goillogical/evaluable"
	. "github.com/spaceavocado/goillogical/internal/test"
)

func ref(val string) Evaluable {
	serOpts := DefaultSerializeOptions()
	simOpts := SimplifyOptions{
		IgnoredPaths:   []string{},
		IgnoredPathsRx: []regexp.Regexp{},
	}
	e, _ := New(val, &serOpts, &simOpts)
	return e
}

func TestGetDataType(t *testing.T) {
	var tests = []struct {
		input    string
		expected DataType
	}{
		{"ref", Undefined},
		{"ref.(X)", Undefined},
		{"ref.(String)", String},
		{"ref.(Number)", Number},
		{"ref.(Integer)", Integer},
		{"ref.(Float)", Float},
		{"ref.(Boolean)", Boolean},
	}

	for _, test := range tests {
		if output, _ := getDataType(test.input); output != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}

	input := "ref.(Struct)"
	expected := fmt.Sprintf("unsupported \"%s\" type casting", "Struct")
	if output, err := getDataType(input); err.Error() != expected {
		t.Errorf("input (%v): expected %v, got %v", input, expected, output)
	}
}

func TestTrimDataType(t *testing.T) {
	var tests = []struct {
		input    string
		expected string
	}{
		{"ref", "ref"},
		{"ref.(X)", "ref.(X)"},
		{"ref.(String)", "ref"},
		{"ref.(Nonsense)", "ref"},
	}

	for _, test := range tests {
		if output := trimDataType(test.input); output != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}
}

func TestToNumber(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{1, 1},
		{"1", 1},
		{"1.1", 1.1},
		{true, 1},
		{false, 0},
	}

	for _, test := range tests {
		if output, err := toNumber(test.input); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}

	errs := []struct {
		input    any
		expected error
	}{
		{"1,1", errors.New("invalid conversion from from \"1,1\" (string) to number")},
		{struct{ a string }{a: "b"}, errors.New("invalid conversion from from \"{b}\" to number")},
	}
	for _, test := range errs {
		if _, err := toNumber(test.input); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}

func TestToInteger(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{1, 1},
		{1.1, 1},
		{"1", 1},
		{"1.1", 1},
		{"1.9", 1},
		{true, 1},
		{false, 0},
	}

	for _, test := range tests {
		if output, err := toInteger(test.input); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}

	errs := []struct {
		input    any
		expected error
	}{
		{"1,1", errors.New("invalid conversion from from \"1,1\" (string) to integer")},
		{struct{ a string }{a: "b"}, errors.New("invalid conversion from from \"{b}\" to integer")},
	}
	for _, test := range errs {
		if _, err := toInteger(test.input); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		input    any
		expected any
	}{
		{1, 1.0},
		{1.1, 1.1},
		{"1", 1.0},
		{"1.1", 1.1},
		{"1.9", 1.9},
	}

	for _, test := range tests {
		if output, err := toFloat(test.input); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}

	errs := []struct {
		input    any
		expected error
	}{
		{"1,1", errors.New("invalid conversion from from \"1,1\" (string) to float")},
		{true, errors.New("invalid conversion from from \"true\" to float")},
		{struct{ a string }{a: "b"}, errors.New("invalid conversion from from \"{b}\" to float")},
	}
	for _, test := range errs {
		if _, err := toFloat(test.input); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}

func TestToBoolean(t *testing.T) {
	tests := []struct {
		input    any
		expected bool
	}{
		{true, true},
		{false, false},
		{"true", true},
		{"false", false},
		{" True ", true},
		{" False ", false},
		{"1", true},
		{"0", false},
		{1, true},
		{0, false},
	}

	for _, test := range tests {
		if output, err := toBoolean(test.input); output != test.expected || err != nil {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}

	errs := []struct {
		input    any
		expected error
	}{
		{"yes", errors.New("invalid conversion from from \"yes\" to boolean")},
		{"bogus", errors.New("invalid conversion from from \"bogus\" to boolean")},
		{2, errors.New("invalid conversion from from \"2\" to boolean")},
		{[]int{1}, errors.New("invalid conversion from from \"[1]\" to boolean")},
	}
	for _, test := range errs {
		if _, err := toBoolean(test.input); err.Error() != test.expected.Error() {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, err)
		}
	}
}

func TestToString(t *testing.T) {
	var tests = []struct {
		input    any
		expected string
	}{
		{1, "1"},
		{"1", "1"},
		{true, "true"},
	}

	for _, test := range tests {
		if output := toString(test.input); output != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}
}

func TestContextLookup(t *testing.T) {
	ctx := FlattenContext(map[string]any{
		"refA": 1,
		"refB": map[string]any{
			"refB1": 2,
			"refB2": "refB1",
			"refB3": true,
		},
		"refC": "refB1",
		"refD": "refB2",
		"refE": []any{1, []any{2, 3, 4}},
		"refF": "A",
		"refG": "1",
		"refH": "1.1",
	})

	var tests = []struct {
		input string
		path  string
		value any
	}{
		{"UNDEFINED", "UNDEFINED", nil},
		{"refA", "refA", 1},
		{"refB.refB1", "refB.refB1", 2},
		{"refB.{refC}", "refB.refB1", 2},
		{"refB.{UNDEFINED}", "refB.{UNDEFINED}", nil},
		{"refB.{refB.refB2}", "refB.refB1", 2},
		{"refB.{refB.{refD}}", "refB.refB1", 2},
		{"refE[0]", "refE[0]", 1},
		{"refE[2]", "refE[2]", nil},
		{"refE[1][0]", "refE[1][0]", 2},
		{"refE[1][3]", "refE[1][3]", nil},
		{"refE[{refA}][0]", "refE[1][0]", 2},
		{"refE[{refA}][{refB.refB1}]", "refE[1][2]", 4},
		{"ref{refF}", "refA", 1},
		{"ref{UNDEFINED}", "ref{UNDEFINED}", nil},
	}

	for _, test := range tests {
		if path, value := contextLookup(ctx, test.input); path != test.path || value != test.value {
			t.Errorf("input (%v): expected %v/%v, got %v/%v", test.input, test.path, test.value, path, value)
		}
	}
}

func TestEvaluate(t *testing.T) {
	ctx := FlattenContext(map[string]any{
		"refA": 1,
		"refB": map[string]any{
			"refB1": 2,
			"refB2": "refB1",
			"refB3": true,
		},
		"refC": "refB1",
		"refD": "refB2",
		"refE": []any{1, []any{2, 3, 4}},
		"refF": func() {},
		"refG": "1",
		"refH": "1.1",
	})

	tests := []struct {
		path  string
		dt    DataType
		value any
	}{
		{"refA", Undefined, 1},
		{"refA", String, "1"},
		{"refG", Number, 1},
		{"refH", Number, 1.1},
		{"refB.refB3", String, "true"},
		{"refB.refB3", Number, 1},
		{"refJ", Undefined, nil},
	}

	for _, test := range tests {
		if _, value, err := evaluate(ctx, test.path, test.dt); value != test.value || err != nil {
			t.Errorf("input (%v, %v): expected %v, got %v", test.path, test.dt, test.value, value)
		}
	}
}

func TestIsIgnoredPath(t *testing.T) {
	simOpts := SimplifyOptions{
		IgnoredPaths:   []string{"ignored"},
		IgnoredPathsRx: []regexp.Regexp{*regexp.MustCompile("^refC")},
	}

	tests := []struct {
		input    string
		expected bool
	}{
		{"ignored", true},
		{"not", false},
		{"refC", true},
		{"refC.(Number)", true},
	}

	for _, test := range tests {
		if output := isIgnoredPath(test.input, &simOpts); output != test.expected {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.expected, output)
		}
	}
}

func TestSerialize(t *testing.T) {
	serOpts := DefaultSerializeOptions()
	simOpts := SimplifyOptions{
		IgnoredPaths:   []string{},
		IgnoredPathsRx: []regexp.Regexp{},
	}

	tests := []struct {
		input string
		value any
	}{
		{"refA", "$refA"},
		{"refA.(Number)", "$refA.(Number)"},
	}

	for _, test := range tests {
		e, _ := New(test.input, &serOpts, &simOpts)
		if value := e.Serialize(); value != test.value {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.value, value)
		}
	}
}

func TestSimplify(t *testing.T) {
	opts := DefaultSerializeOptions()
	simOpts := SimplifyOptions{
		IgnoredPaths:   []string{"ignored"},
		IgnoredPathsRx: []regexp.Regexp{*regexp.MustCompile("^refC")},
	}
	ctx := FlattenContext(map[string]any{
		"refA": 1,
		"refB": map[string]any{
			"refB1": 2,
			"refB2": "refB1",
			"refB3": true,
		},
		"refC": "refB1",
		"refD": "refB2",
		"refE": []any{1, []any{2, 3, 4}},
		"refF": func() {},
		"refG": "1",
		"refH": "1.1",
	})

	tests := []struct {
		input string
		value any
		e     any
	}{
		{"refJ", nil, ref("refJ")},
		{"ignored", nil, nil},
		{"refA", 1, nil},
		{"refB.{refJ}", nil, ref("refB.{refJ}")},
		{"refC.{refJ}", nil, nil},
	}

	for _, test := range tests {
		e, _ := New(test.input, &opts, &simOpts)
		if value, self := e.Simplify(ctx); value != test.value || Fprint(self) != Fprint(test.e) {
			t.Errorf("input (%v): expected %v/%v, got %v/%v", test.input, test.value, test.e, value, self)
		}
	}
}

func TestString(t *testing.T) {
	opts := DefaultSerializeOptions()
	simOpts := SimplifyOptions{
		IgnoredPaths:   []string{},
		IgnoredPathsRx: []regexp.Regexp{},
	}

	tests := []struct {
		input string
		value string
	}{
		{"refA", "{refA}"},
		{"refA.(Number)", "{refA.(Number)}"},
	}

	for _, test := range tests {
		e, _ := New(test.input, &opts, &simOpts)
		if value := e.String(); value != test.value {
			t.Errorf("input (%v): expected %v, got %v", test.input, test.value, value)
		}
	}
}
