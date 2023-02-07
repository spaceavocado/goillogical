package parser

import (
	"errors"
	"fmt"
	. "goillogical/internal"
	eq "goillogical/internal/expression/comparison/eq"
	ge "goillogical/internal/expression/comparison/ge"
	gt "goillogical/internal/expression/comparison/gt"
	in "goillogical/internal/expression/comparison/in"
	le "goillogical/internal/expression/comparison/le"
	lt "goillogical/internal/expression/comparison/lt"
	ne "goillogical/internal/expression/comparison/ne"
	null "goillogical/internal/expression/comparison/nil"
	nin "goillogical/internal/expression/comparison/nin"
	overlap "goillogical/internal/expression/comparison/overlap"
	prefix "goillogical/internal/expression/comparison/prefix"
	present "goillogical/internal/expression/comparison/present"
	suffix "goillogical/internal/expression/comparison/suffix"
	and "goillogical/internal/expression/logical/and"
	nor "goillogical/internal/expression/logical/nor"
	not "goillogical/internal/expression/logical/not"
	or "goillogical/internal/expression/logical/or"
	xor "goillogical/internal/expression/logical/xor"
	collection "goillogical/internal/operand/collection"
	reference "goillogical/internal/operand/reference"
	value "goillogical/internal/operand/value"
	. "goillogical/internal/options"
	"reflect"
	"strings"
)

type options struct {
	OperatorHandlers map[string]func([]Evaluable) (Evaluable, error)
	Serialize        struct {
		Reference  reference.SerializeOptions
		Collection collection.SerializeOptions
	}
	Simplify struct {
		Reference reference.SimplifyOptions
	}
}

func expressionUnary(op string, factory func(string, Evaluable) (Evaluable, error)) func([]Evaluable) (Evaluable, error) {
	return func(operands []Evaluable) (Evaluable, error) {
		return factory(op, operands[0])
	}
}

func expressionBinary(op string, factory func(string, Evaluable, Evaluable) (Evaluable, error)) func([]Evaluable) (Evaluable, error) {
	return func(operands []Evaluable) (Evaluable, error) {
		return factory(op, operands[0], operands[1])
	}
}

func expressionMany(op string, factory func(string, []Evaluable) (Evaluable, error)) func([]Evaluable) (Evaluable, error) {
	return func(operands []Evaluable) (Evaluable, error) {
		return factory(op, operands)
	}
}

func operatorHandlers(opts OperatorMapping) map[string]func([]Evaluable) (Evaluable, error) {
	return map[string]func([]Evaluable) (Evaluable, error){
		// Logical
		opts[And]: expressionMany(opts[And], and.New),
		opts[Or]:  expressionMany(opts[Or], or.New),
		opts[Nor]: expressionMany(opts[Nor], nor.New),
		opts[Xor]: expressionMany(opts[Xor], xor.New),
		opts[Not]: expressionUnary(opts[Not], not.New),
		// Comparison
		opts[Eq]:      expressionBinary(opts[Eq], eq.New),
		opts[Ne]:      expressionBinary(opts[Ne], ne.New),
		opts[Gt]:      expressionBinary(opts[Gt], gt.New),
		opts[Ge]:      expressionBinary(opts[Ge], ge.New),
		opts[Lt]:      expressionBinary(opts[Lt], lt.New),
		opts[Le]:      expressionBinary(opts[Le], le.New),
		opts[In]:      expressionBinary(opts[In], in.New),
		opts[Nin]:     expressionBinary(opts[Nin], nin.New),
		opts[Overlap]: expressionBinary(opts[Overlap], overlap.New),
		opts[Nil]:     expressionUnary(opts[Nil], null.New),
		opts[Present]: expressionUnary(opts[Present], present.New),
		opts[Suffix]:  expressionBinary(opts[Suffix], suffix.New),
		opts[Prefix]:  expressionBinary(opts[Prefix], prefix.New),
	}
}

type Parser interface {
	Parse(exp any) (Evaluable, error)
}

type parser struct {
	opts options
}

func (p parser) Parse(exp any) (Evaluable, error) {
	return parse(exp, &p.opts)
}

func isEscaped(value string, escapeCharacter string) bool {
	return escapeCharacter != "" && strings.HasPrefix(value, escapeCharacter)
}

func toReferenceAddr(input any, opts *reference.SerializeOptions) (string, error) {
	switch input.(type) {
	case string:
		return opts.From(input.(string))
	default:
		return "", errors.New("invalid reference path")
	}
}

func createOperand(input any, opts *options) (Evaluable, error) {
	if input == nil {
		return nil, errors.New("invalid undefined operand")
	}

	t := reflect.TypeOf(input).Kind()
	if t == reflect.Slice {
		v := reflect.ValueOf(input)
		if v.Len() == 0 {
			return nil, errors.New("invalid undefined operand")
		}

		operands := make([]Evaluable, v.Len())
		for i := 0; i < v.Len(); i++ {
			e, err := parse(v.Index(i).Interface(), opts)
			if err != nil {
				return nil, err
			}
			operands[i] = e
		}
		return collection.New(operands, &opts.Serialize.Collection)
	}

	addr, err := toReferenceAddr(input, &opts.Serialize.Reference)
	if err == nil {
		return reference.New(addr, &opts.Serialize.Reference, &opts.Simplify.Reference)
	}

	if !IsEvaluatedPrimitive(input) {
		return nil, errors.New(fmt.Sprintf("invalid operand, %v", input))
	}

	return value.New(input)
}

func createExpression(expression []any, opts *options) (Evaluable, error) {
	operator := expression[0]
	operands := expression[1:]
	switch operator.(type) {
	case string:
		handler, ok := opts.OperatorHandlers[operator.(string)]
		if !ok {
			return nil, errors.New("unexpected logical operator")
		}

		ops := make([]Evaluable, len(operands))
		for i := 0; i < len(operands); i++ {
			e, err := parse(operands[i], opts)
			if err != nil {
				return nil, err
			}
			ops[i] = e
		}

		return handler(ops)
	default:
		return nil, errors.New("unexpected logical expression")
	}
}

func parse(input any, opts *options) (Evaluable, error) {
	if input == nil {
		return nil, errors.New("unexpected input")
	}

	t := reflect.TypeOf(input).Kind()
	if t != reflect.Slice {
		return createOperand(input, opts)
	}

	v := reflect.ValueOf(input)

	if v.Len() < 2 {
		return createOperand(input, opts)
	}

	operator := v.Index(0).Interface()
	if isEscaped(fmt.Sprintf("%v", operator), opts.Serialize.Collection.EscapeCharacter) {
		return createOperand(append([]any{operator.(string)[1:]}, v.Slice(1, v.Len()).Interface().([]any)...), opts)
	}

	e, err := createExpression(input.([]any), opts)
	if err != nil {
		return createOperand(input, opts)
	}

	return e, nil
}

func New(opts *Options) Parser {
	return &parser{opts: options{
		OperatorHandlers: operatorHandlers(opts.OperatorMapping),
		Serialize:        opts.Serialize,
		Simplify:         opts.Simplify,
	}}
}
