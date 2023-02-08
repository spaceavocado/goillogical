package parser

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	e "github.com/spaceavocado/goillogical/evaluable"
	eq "github.com/spaceavocado/goillogical/internal/expression/comparison/eq"
	ge "github.com/spaceavocado/goillogical/internal/expression/comparison/ge"
	gt "github.com/spaceavocado/goillogical/internal/expression/comparison/gt"
	in "github.com/spaceavocado/goillogical/internal/expression/comparison/in"
	le "github.com/spaceavocado/goillogical/internal/expression/comparison/le"
	lt "github.com/spaceavocado/goillogical/internal/expression/comparison/lt"
	ne "github.com/spaceavocado/goillogical/internal/expression/comparison/ne"
	null "github.com/spaceavocado/goillogical/internal/expression/comparison/nil"
	nin "github.com/spaceavocado/goillogical/internal/expression/comparison/nin"
	overlap "github.com/spaceavocado/goillogical/internal/expression/comparison/overlap"
	prefix "github.com/spaceavocado/goillogical/internal/expression/comparison/prefix"
	present "github.com/spaceavocado/goillogical/internal/expression/comparison/present"
	suffix "github.com/spaceavocado/goillogical/internal/expression/comparison/suffix"
	and "github.com/spaceavocado/goillogical/internal/expression/logical/and"
	nor "github.com/spaceavocado/goillogical/internal/expression/logical/nor"
	not "github.com/spaceavocado/goillogical/internal/expression/logical/not"
	or "github.com/spaceavocado/goillogical/internal/expression/logical/or"
	xor "github.com/spaceavocado/goillogical/internal/expression/logical/xor"
	collection "github.com/spaceavocado/goillogical/internal/operand/collection"
	reference "github.com/spaceavocado/goillogical/internal/operand/reference"
	value "github.com/spaceavocado/goillogical/internal/operand/value"
	o "github.com/spaceavocado/goillogical/internal/options"
)

type options struct {
	OperatorHandlers map[string]func([]e.Evaluable) (e.Evaluable, error)
	Serialize        struct {
		Reference  reference.SerializeOptions
		Collection collection.SerializeOptions
	}
	Simplify struct {
		Reference reference.SimplifyOptions
	}
}

func expressionUnary(op string, factory func(string, e.Evaluable) (e.Evaluable, error)) func([]e.Evaluable) (e.Evaluable, error) {
	return func(operands []e.Evaluable) (e.Evaluable, error) {
		return factory(op, operands[0])
	}
}

func expressionBinary(op string, factory func(string, e.Evaluable, e.Evaluable) (e.Evaluable, error)) func([]e.Evaluable) (e.Evaluable, error) {
	return func(operands []e.Evaluable) (e.Evaluable, error) {
		return factory(op, operands[0], operands[1])
	}
}

func expressionMany(op string, factory func(string, []e.Evaluable, string, string) (e.Evaluable, error), notOp string, norOp string) func([]e.Evaluable) (e.Evaluable, error) {
	return func(operands []e.Evaluable) (e.Evaluable, error) {
		return factory(op, operands, notOp, norOp)
	}
}

func operatorHandlers(opts e.OperatorMapping) map[string]func([]e.Evaluable) (e.Evaluable, error) {
	return map[string]func([]e.Evaluable) (e.Evaluable, error){
		// Logical
		opts[e.And]: expressionMany(opts[e.And], and.New, opts[e.Not], opts[e.Nor]),
		opts[e.Or]:  expressionMany(opts[e.Or], or.New, opts[e.Not], opts[e.Nor]),
		opts[e.Nor]: expressionMany(opts[e.Nor], nor.New, opts[e.Not], opts[e.Nor]),
		opts[e.Xor]: expressionMany(opts[e.Xor], xor.New, opts[e.Not], opts[e.Nor]),
		opts[e.Not]: expressionUnary(opts[e.Not], not.New),
		// Comparison
		opts[e.Eq]:      expressionBinary(opts[e.Eq], eq.New),
		opts[e.Ne]:      expressionBinary(opts[e.Ne], ne.New),
		opts[e.Gt]:      expressionBinary(opts[e.Gt], gt.New),
		opts[e.Ge]:      expressionBinary(opts[e.Ge], ge.New),
		opts[e.Lt]:      expressionBinary(opts[e.Lt], lt.New),
		opts[e.Le]:      expressionBinary(opts[e.Le], le.New),
		opts[e.In]:      expressionBinary(opts[e.In], in.New),
		opts[e.Nin]:     expressionBinary(opts[e.Nin], nin.New),
		opts[e.Overlap]: expressionBinary(opts[e.Overlap], overlap.New),
		opts[e.Nil]:     expressionUnary(opts[e.Nil], null.New),
		opts[e.Present]: expressionUnary(opts[e.Present], present.New),
		opts[e.Suffix]:  expressionBinary(opts[e.Suffix], suffix.New),
		opts[e.Prefix]:  expressionBinary(opts[e.Prefix], prefix.New),
	}
}

type Parser interface {
	Parse(exp any) (e.Evaluable, error)
}

type parser struct {
	opts options
}

func (p parser) Parse(exp any) (e.Evaluable, error) {
	return parse(exp, &p.opts)
}

func isEscaped(value string, escapeCharacter string) bool {
	return escapeCharacter != "" && strings.HasPrefix(value, escapeCharacter)
}

func toReferenceAddr(input any, opts *reference.SerializeOptions) (string, error) {
	switch typed := input.(type) {
	case string:
		return opts.From(typed)
	default:
		return "", errors.New("invalid reference path")
	}
}

func createOperand(input any, opts *options) (e.Evaluable, error) {
	if input == nil {
		return nil, errors.New("invalid undefined operand")
	}

	t := reflect.TypeOf(input).Kind()
	if t == reflect.Slice {
		v := reflect.ValueOf(input)
		if v.Len() == 0 {
			return nil, errors.New("invalid undefined operand")
		}

		operands := make([]e.Evaluable, v.Len())
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

	if !e.IsEvaluatedPrimitive(input) {
		return nil, fmt.Errorf("invalid operand, %v", input)
	}

	return value.New(input)
}

func createExpression(expression []any, opts *options) (e.Evaluable, error) {
	operator := expression[0]
	operands := expression[1:]
	switch typed := operator.(type) {
	case string:
		handler, ok := opts.OperatorHandlers[typed]
		if !ok {
			return nil, errors.New("unexpected logical operator")
		}

		ops := make([]e.Evaluable, len(operands))
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

func parse(input any, opts *options) (e.Evaluable, error) {
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

func New(opts *o.Options) Parser {
	return &parser{opts: options{
		OperatorHandlers: operatorHandlers(opts.OperatorMapping),
		Serialize:        opts.Serialize,
		Simplify:         opts.Simplify,
	}}
}
