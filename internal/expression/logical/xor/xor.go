package xor

import (
	"errors"
	. "goillogical/internal"
	l "goillogical/internal/expression/logical"
)

func xor(a, b bool) bool {
	return (a || b) && !(a && b)
}

func handler(ctx Context, operands []Evaluable) (bool, error) {
	var out bool
	for i, o := range operands {
		res, err := l.Evaluate(ctx, o)
		if err != nil {
			return false, err
		}

		if i == 0 {
			out = res
		} else {
			out = xor(out, res)
		}
	}
	return out, nil
}

func New(operands []Evaluable) (Evaluable, error) {
	if len(operands) < 2 {
		return nil, errors.New("logical XOR expression must have at least 2 operands")
	}

	return l.New(Xor, "XOR", operands, handler)
}
