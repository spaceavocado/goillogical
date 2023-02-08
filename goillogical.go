package goillogical

import (
	. "github.com/spaceavocado/goillogical/internal"
	c "github.com/spaceavocado/goillogical/internal/operand/collection"
	r "github.com/spaceavocado/goillogical/internal/operand/reference"
	. "github.com/spaceavocado/goillogical/internal/options"
	p "github.com/spaceavocado/goillogical/internal/parser"
)

type Goillogical interface {
	Evaluate(any, Context) (any, error)
	Parse(any) (Evaluable, error)
	Statement(any) (string, error)
	Simplify(any, Context) (any, Evaluable, error)
}

type illogical struct {
	opts   Options
	parser p.Parser
}

func (i illogical) Evaluate(exp any, ctx Context) (any, error) {
	e, err := i.parser.Parse(exp)
	if err != nil {
		return nil, err
	}
	return e.Evaluate(FlattenContext(ctx))
}

func (i illogical) Parse(exp any) (Evaluable, error) {
	return i.parser.Parse(exp)
}

func (i illogical) Statement(exp any) (string, error) {
	e, err := i.parser.Parse(exp)
	if err != nil {
		return "", err
	}
	return e.String(), nil
}

func (i illogical) Simplify(exp any, ctx Context) (any, Evaluable, error) {
	e, err := i.parser.Parse(exp)
	if err != nil {
		return nil, nil, err
	}

	val, eval := e.Simplify(ctx)
	return val, eval, nil
}

type Option func(*illogical)

func WithReferenceSerializeOptions(o r.SerializeOptions) Option {
	return func(i *illogical) {
		i.opts.Serialize.Reference = o
	}
}

func WithCollectionSerializeOptions(o c.SerializeOptions) Option {
	return func(i *illogical) {
		i.opts.Serialize.Collection = o
	}
}

func WithReferenceSimplifyOptions(o r.SimplifyOptions) Option {
	return func(i *illogical) {
		i.opts.Simplify.Reference = o
	}
}

func WithOperatorMappingOptions(m OperatorMapping) Option {
	return func(i *illogical) {
		i.opts.OperatorMapping = m
	}
}

type SimplifyOptions = r.SimplifyOptions
type ReferenceSerializeOptions = r.SerializeOptions
type CollectionSerializeOptions = c.SerializeOptions

func New(opts ...Option) Goillogical {
	i := &illogical{DefaultOptions(), nil}

	for _, opt := range opts {
		opt(i)
	}

	i.parser = p.New(&i.opts)
	return i
}
