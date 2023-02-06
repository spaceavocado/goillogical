package parser

import (
	"errors"
	"fmt"
	. "goillogical"
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
	"reflect"
	"strings"
)

type options struct {
	OperatorHandlers          map[string]func([]Evaluable) (Evaluable, error)
	EscapeCharacter           string
	ReferenceSerializeOptions reference.SerializeOptions
}

func expressionUnary(factory func(Evaluable) (Evaluable, error)) func([]Evaluable) (Evaluable, error) {
	return func(operands []Evaluable) (Evaluable, error) {
		return factory(operands[0])
	}
}

func expressionBinary(factory func(Evaluable, Evaluable) (Evaluable, error)) func([]Evaluable) (Evaluable, error) {
	return func(operands []Evaluable) (Evaluable, error) {
		return factory(operands[0], operands[1])
	}
}

func expressionMany(factory func([]Evaluable) (Evaluable, error)) func([]Evaluable) (Evaluable, error) {
	return func(operands []Evaluable) (Evaluable, error) {
		return factory(operands)
	}
}

func operatorHandlers(opts OperatorMapping) map[string]func([]Evaluable) (Evaluable, error) {
	return map[string]func([]Evaluable) (Evaluable, error){
		// Logical
		opts[And]: expressionMany(and.New),
		opts[Or]:  expressionMany(or.New),
		opts[Nor]: expressionMany(nor.New),
		opts[Xor]: expressionMany(xor.New),
		opts[Not]: expressionUnary(not.New),
		// Comparison
		opts[Eq]:      expressionBinary(eq.New),
		opts[Ne]:      expressionBinary(ne.New),
		opts[Gt]:      expressionBinary(gt.New),
		opts[Ge]:      expressionBinary(ge.New),
		opts[Lt]:      expressionBinary(lt.New),
		opts[Le]:      expressionBinary(le.New),
		opts[In]:      expressionBinary(in.New),
		opts[Nin]:     expressionBinary(nin.New),
		opts[Overlap]: expressionBinary(overlap.New),
		opts[Nil]:     expressionUnary(null.New),
		opts[Present]: expressionUnary(present.New),
		opts[Suffix]:  expressionBinary(suffix.New),
		opts[Prefix]:  expressionBinary(prefix.New),
	}
}

type Parser interface {
	Parse(exp any) (Evaluable, error)
}

type parser struct {
	opts options
}

func (p parser) Parse(exp any) (Evaluable, error) {
	return parse(exp, p.opts)
}

func isEscaped(value string, escapeCharacter string) bool {
	return escapeCharacter != "" && strings.HasPrefix(value, escapeCharacter)
}

func toReferenceAddr(input any, opts reference.SerializeOptions) (string, error) {
	switch input.(type) {
	case string:
		return opts.From(input.(string))
	default:
		return "", errors.New("invalid reference path")
	}
}

func createOperand(input any, opts options) (Evaluable, error) {
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
		return collection.New(operands)
	}

	addr, err := toReferenceAddr(input, opts.ReferenceSerializeOptions)
	if err == nil {
		return reference.New(addr)
	}

	if !IsEvaluatedPrimitive(input) {
		return nil, errors.New(fmt.Sprintf("invalid operand, %v", input))
	}

	return value.New(input)
}

func createExpression(expression []any, opts options) (Evaluable, error) {
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

func parse(input any, opts options) (Evaluable, error) {
	t := reflect.TypeOf(input).Kind()
	if t != reflect.Slice {
		return createOperand(input, opts)
	}

	v := reflect.ValueOf(input)

	if v.Len() < 2 {
		return createOperand(input, opts)
	}

	operator := v.Index(0).Interface()
	if isEscaped(fmt.Sprintf("%v", operator), opts.EscapeCharacter) {
		return createOperand(append([]any{operator.(string)[1:]}, v.Slice(1, v.Len()).Interface().([]any)...), opts)
	}

	e, err := createExpression(input.([]any), opts)
	if err != nil {
		return createOperand(input, opts)
	}

	return e, nil
}

func New(opts Options) Parser {
	return &parser{opts: options{
		OperatorHandlers:          operatorHandlers(opts.OperatorMapping),
		ReferenceSerializeOptions: opts.Serialize.Reference,
		EscapeCharacter:           opts.Serialize.Collection.EscapeCharacter,
	}}
}
