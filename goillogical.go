package goillogical

import (
	e "github.com/spaceavocado/goillogical/evaluable"
	c "github.com/spaceavocado/goillogical/internal/operand/collection"
	r "github.com/spaceavocado/goillogical/internal/operand/reference"
	o "github.com/spaceavocado/goillogical/internal/options"
	p "github.com/spaceavocado/goillogical/internal/parser"
)

type SimplifyOptions = r.SimplifyOptions
type ReferenceSerializeOptions = r.SerializeOptions
type CollectionSerializeOptions = struct {
	EscapeCharacter string
}

type Goillogical interface {
	Evaluate(any, e.Context) (any, error)
	Parse(any) (e.Evaluable, error)
	Statement(any) (string, error)
	Simplify(any, e.Context) (any, e.Evaluable, error)
}

type illogical struct {
	opts   o.Options
	parser p.Parser
}

func (i illogical) Evaluate(exp any, ctx e.Context) (any, error) {
	eval, err := i.parser.Parse(exp)
	if err != nil {
		return nil, err
	}
	return eval.Evaluate(e.FlattenContext(ctx))
}

func (i illogical) Parse(exp any) (e.Evaluable, error) {
	return i.parser.Parse(exp)
}

func (i illogical) Statement(exp any) (string, error) {
	e, err := i.parser.Parse(exp)
	if err != nil {
		return "", err
	}
	return e.String(), nil
}

func (i illogical) Simplify(exp any, ctx e.Context) (any, e.Evaluable, error) {
	e, err := i.parser.Parse(exp)
	if err != nil {
		return nil, nil, err
	}

	val, eval := e.Simplify(ctx)
	return val, eval, nil
}

type Option func(*illogical)

func WithReferenceSerializeOptions(o ReferenceSerializeOptions) Option {
	return func(i *illogical) {
		i.opts.Serialize.Reference = o
	}
}

func WithCollectionSerializeOptions(o CollectionSerializeOptions) Option {
	return func(i *illogical) {
		i.opts.Serialize.Collection = c.SerializeOptions{
			EscapeCharacter:  o.EscapeCharacter,
			EscapedOperators: map[string]bool{},
		}
	}
}

func WithReferenceSimplifyOptions(o SimplifyOptions) Option {
	return func(i *illogical) {
		i.opts.Simplify.Reference = o
	}
}

func WithOperatorMappingOptions(m e.OperatorMapping) Option {
	return func(i *illogical) {
		i.opts.OperatorMapping = m
	}
}

func New(opts ...Option) Goillogical {
	i := &illogical{o.DefaultOptions(), nil}

	for _, opt := range opts {
		opt(i)
	}

	i.parser = p.New(&i.opts)
	return i
}
