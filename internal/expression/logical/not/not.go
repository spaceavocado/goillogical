package not

import (
	"errors"
	. "goillogical/internal"
	l "goillogical/internal/expression/logical"
)

func handler(ctx Context, operands []Evaluable) (bool, error) {
	res, err := operands[0].Evaluate(ctx)
	if err != nil {
		return false, err
	}
	switch res.(type) {
	case bool:
		return !res.(bool), nil
	default:
		return false, errors.New("logical NOT expression's operand must be evaluated to boolean value")
	}
}

func simplify(operator string, ctx Context, operands []Evaluable) (any, Evaluable) {
	res, e := operands[0].Simplify(ctx)
	if b, ok := res.(bool); ok {
		return !b, nil
	}

	if e != nil {
		e, _ := New(operator, e)
		return nil, e
	}

	return nil, nil
}

func New(operator string, operand Evaluable) (Evaluable, error) {
	return l.New(operator, "NOT", []Evaluable{operand}, handler, simplify)
}
