package mock

import (
	"errors"
	"regexp"

	e "github.com/spaceavocado/goillogical/evaluable"
	collection "github.com/spaceavocado/goillogical/internal/operand/collection"
	reference "github.com/spaceavocado/goillogical/internal/operand/reference"
	value "github.com/spaceavocado/goillogical/internal/operand/value"
)

type eval struct {
	evaluate  func(e.Context) (any, error)
	serialize func() any
	simplify  func(e.Context) (any, e.Evaluable)
	string    func() string
}

func (e eval) Evaluate(ctx e.Context) (any, error)       { return e.evaluate(ctx) }
func (e eval) Serialize() any                            { return e.serialize() }
func (e eval) Simplify(ctx e.Context) (any, e.Evaluable) { return e.simplify(ctx) }
func (e eval) String() string                            { return e.string() }

func Invalid() e.Evaluable {
	return eval{
		evaluate:  func(ctx e.Context) (any, error) { return nil, errors.New("invalid") },
		serialize: func() any { return nil },
		simplify:  func(ctx e.Context) (any, e.Evaluable) { return nil, nil },
		string:    func() string { return "invalid" },
	}
}

func Val(val any) e.Evaluable {
	e, _ := value.New(val)
	return e
}

func Ref(val string) e.Evaluable {
	serOpts := reference.DefaultSerializeOptions()
	simOpts := reference.SimplifyOptions{
		IgnoredPaths:   []string{"ignored"},
		IgnoredPathsRx: []regexp.Regexp{},
	}
	e, _ := reference.New(val, &serOpts, &simOpts)
	return e
}

func Col(items ...e.Evaluable) e.Evaluable {
	opts := collection.DefaultSerializeOptions()
	e, _ := collection.New(items, &opts)
	return e
}

func ExpUnary(op string, factory func(string, e.Evaluable) (e.Evaluable, error), eval e.Evaluable) e.Evaluable {
	e, _ := factory(op, eval)
	return e
}

func ExpBinary(op string, factory func(string, e.Evaluable, e.Evaluable) (e.Evaluable, error), left, right e.Evaluable) e.Evaluable {
	e, _ := factory(op, left, right)
	return e
}

func ExpMany(op string, factory func(string, []e.Evaluable, string, string) (e.Evaluable, error), operands ...e.Evaluable) e.Evaluable {
	e, _ := factory(op, operands, "", "")
	return e
}
