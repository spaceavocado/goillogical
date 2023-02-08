package goillogical

import (
	intr "github.com/spaceavocado/goillogical/internal"
	c "github.com/spaceavocado/goillogical/internal/operand/collection"
	r "github.com/spaceavocado/goillogical/internal/operand/reference"
	o "github.com/spaceavocado/goillogical/internal/options"
	p "github.com/spaceavocado/goillogical/internal/parser"
)

type Kind = intr.Kind

const And = intr.And
const Or = intr.Or
const Nor = intr.Nor
const Xor = intr.Xor
const Not = intr.Not

const Eq = intr.Eq
const Ne = intr.Ne
const Gt = intr.Gt
const Ge = intr.Ge
const Lt = intr.Lt
const Le = intr.Le
const Nil = intr.Nil
const Present = intr.Present
const In = intr.In
const Nin = intr.Nin
const Overlap = intr.Overlap
const Prefix = intr.Prefix
const Suffix = intr.Suffix

type OperatorMapping = intr.OperatorMapping
type SimplifyOptions = r.SimplifyOptions
type ReferenceSerializeOptions = r.SerializeOptions
type CollectionSerializeOptions = struct {
	EscapeCharacter string
}
type Context = intr.Context
type Evaluable = intr.Evaluable

type Goillogical interface {
	Evaluate(any, Context) (any, error)
	Parse(any) (Evaluable, error)
	Statement(any) (string, error)
	Simplify(any, Context) (any, Evaluable, error)
}

type illogical struct {
	opts   o.Options
	parser p.Parser
}

func (i illogical) Evaluate(exp any, ctx Context) (any, error) {
	e, err := i.parser.Parse(exp)
	if err != nil {
		return nil, err
	}
	return e.Evaluate(intr.FlattenContext(ctx))
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

func WithOperatorMappingOptions(m OperatorMapping) Option {
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
