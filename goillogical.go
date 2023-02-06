package goillogical

import (
	. "goillogical/internal"
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
	return e.Evaluate(ctx)
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

func New(opts Options) Goillogical {
	return &illogical{opts, p.New(opts)}
}
