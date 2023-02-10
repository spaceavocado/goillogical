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
	res, eval := operands[0].Simplify(ctx)
	if typed, ok := res.(bool); ok {
		return !typed, nil
	}

	if eval != nil {
		e, _ := New(operator, eval)
		return nil, e
	}

	return nil, nil
}

func New(operator string, operand e.Evaluable) (e.Evaluable, error) {
	return l.New(operator, "NOT", []e.Evaluable{operand}, handler, simplify)
}
