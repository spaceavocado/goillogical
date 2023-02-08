package mock

import (
	"regexp"

	. "github.com/spaceavocado/goillogical/evaluable"
	collection "github.com/spaceavocado/goillogical/internal/operand/collection"
	reference "github.com/spaceavocado/goillogical/internal/operand/reference"
	value "github.com/spaceavocado/goillogical/internal/operand/value"
)

func Val(val any) Evaluable {
	e, _ := value.New(val)
	return e
}

func Ref(val string) Evaluable {
	serOpts := reference.DefaultSerializeOptions()
	simOpts := reference.SimplifyOptions{
		IgnoredPaths:   []string{"ignored"},
		IgnoredPathsRx: []regexp.Regexp{},
	}
	e, _ := reference.New(val, &serOpts, &simOpts)
	return e
}

func Col(items ...Evaluable) Evaluable {
	opts := collection.DefaultSerializeOptions()
	e, _ := collection.New(items, &opts)
	return e
}

func ExpUnary(op string, factory func(string, Evaluable) (Evaluable, error), eval Evaluable) Evaluable {
	e, _ := factory(op, eval)
	return e
}

func ExpBinary(op string, factory func(string, Evaluable, Evaluable) (Evaluable, error), left, right Evaluable) Evaluable {
	e, _ := factory(op, left, right)
	return e
}

func ExpMany(op string, factory func(string, []Evaluable, string, string) (Evaluable, error), operands ...Evaluable) Evaluable {
	e, _ := factory(op, operands, "", "")
	return e
}
