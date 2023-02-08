package or

import (
	"errors"

	e "github.com/spaceavocado/goillogical/evaluable"
	l "github.com/spaceavocado/goillogical/internal/expression/logical"
)

func handler(ctx e.Context, operands []e.Evaluable) (bool, error) {
	for _, o := range operands {
		res, err := l.Evaluate(ctx, o)
		if err != nil {
			return false, err
		}
		if res {
			return true, nil
		}
	}
	return false, nil
}

func simplify(operator string, ctx e.Context, operands []e.Evaluable) (any, e.Evaluable) {
	simplified := []e.Evaluable{}
	for _, o := range operands {
		res, e := o.Simplify(ctx)
		if b, ok := res.(bool); ok {
			if b {
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

func New(operator string, operands []e.Evaluable, notOp string, norOp string) (e.Evaluable, error) {
	if len(operands) < 2 {
		return nil, errors.New("logical OR expression must have at least 2 operands")
	}

	return l.New(operator, "OR", operands, handler, func(operator string, ctx map[string]any, operands []e.Evaluable) (any, e.Evaluable) {
		return simplify(operator, ctx, operands)
	})
}
