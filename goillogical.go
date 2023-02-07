package goillogical

import (
	. "goillogical/internal"
	c "goillogical/internal/operand/collection"
	r "goillogical/internal/operand/reference"
	. "goillogical/internal/options"
	p "goillogical/internal/parser"
)

type Goillogical interface {
	Evaluate(any, Context) (any, error)
	Parse(any) (Evaluable, error)
	Statement(any) (string, error)
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

func New(opts ...Option) Goillogical {
	i := &illogical{DefaultOptions(), nil}

	for _, opt := range opts {
		opt(i)
	}

	i.parser = p.New(&i.opts)
	return i
}
