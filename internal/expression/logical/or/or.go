package or

import (
	"errors"

	. "github.com/spaceavocado/goillogical/internal"
	l "github.com/spaceavocado/goillogical/internal/expression/logical"
)

func handler(ctx Context, operands []Evaluable) (bool, error) {
	for _, o := range operands {
		res, err := l.Evaluate(ctx, o)
		if err != nil {
			return false, err
		}
		if res == true {
			return true, nil
		}
	}
	return false, nil
}

func simplify(operator string, ctx Context, operands []Evaluable) (any, Evaluable) {
	simplified := []Evaluable{}
	for _, o := range operands {
		res, e := o.Simplify(ctx)
		if b, ok := res.(bool); ok {
			if b == true {
				return true, nil
			}
			continue
		}

		simplified = append(simplified, e)
	}

	if len(simplified) == 0 {
		return false, nil
	}

	if len(simplified) == 1 {
		return nil, simplified[0]
	}

	e, _ := New(operator, simplified, "", "")
	return nil, e
}

func New(operator string, operands []Evaluable, notOp string, norOp string) (Evaluable, error) {
	if len(operands) < 2 {
		return nil, errors.New("logical OR expression must have at least 2 operands")
	}

	return l.New(operator, "OR", operands, handler, func(operator string, ctx map[string]any, operands []Evaluable) (any, Evaluable) {
		return simplify(operator, ctx, operands)
	})
}
