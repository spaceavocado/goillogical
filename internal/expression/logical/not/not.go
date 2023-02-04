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

func New(operand Evaluable) (Evaluable, error) {
	return l.New("NOT", []Evaluable{operand}, handler)
}
