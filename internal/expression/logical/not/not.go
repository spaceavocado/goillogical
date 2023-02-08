package not

import (
	"errors"

	e "github.com/spaceavocado/goillogical/evaluable"
	l "github.com/spaceavocado/goillogical/internal/expression/logical"
)

func handler(ctx e.Context, operands []e.Evaluable) (bool, error) {
	res, err := operands[0].Evaluate(ctx)
	if err != nil {
		return false, err
	}
	switch b := res.(type) {
	case bool:
		return !b, nil
	default:
		return false, errors.New("logical NOT expression's operand must be evaluated to boolean value")
	}
}

func simplify(operator string, ctx e.Context, operands []e.Evaluable) (any, e.Evaluable) {
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

func New(operator string, operand e.Evaluable) (e.Evaluable, error) {
	return l.New(operator, "NOT", []e.Evaluable{operand}, handler, simplify)
}
