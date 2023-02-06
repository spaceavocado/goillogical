package nor

import (
	"errors"
	. "goillogical/internal"
	l "goillogical/internal/expression/logical"
)

func handler(ctx Context, operands []Evaluable) (bool, error) {
	for _, o := range operands {
		res, err := l.Evaluate(ctx, o)
		if err != nil {
			return false, err
		}
		if res == true {
			return false, nil
		}
	}
	return true, nil
}

func New(operands []Evaluable) (Evaluable, error) {
	if len(operands) < 2 {
		return nil, errors.New("logical NOR expression must have at least 2 operands")
	}

	return l.New(Nor, "NOR", operands, handler)
}
