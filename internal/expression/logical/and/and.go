package and

import (
	"errors"

	e "github.com/spaceavocado/goillogical/evaluable"
	l "github.com/spaceavocado/goillogical/internal/expression/logical"
)

func handler(ctx e.Context, operands []e.Evaluable) (bool, error) {
	var flattenContext = e.FlattenContext(ctx)

	for _, o := range operands {
		res, err := l.Evaluate(flattenContext, o)
		if err != nil {
			return false, err
		}
		if !res {
			return false, nil
		}
	}
	return true, nil
}

func simplify(operator string, ctx e.Context, operands []e.Evaluable) (any, e.Evaluable) {
	var flattenContext = e.FlattenContext(ctx)

	simplified := []e.Evaluable{}
	for _, o := range operands {
		res, e := o.Simplify(flattenContext)
		if b, ok := res.(bool); ok {
			if !b {
				return false, nil
			}
			continue
		}

		simplified = append(simplified, e)
	}

	if len(simplified) == 0 {
		return true, nil
	}

	if len(simplified) == 1 {
		return nil, simplified[0]
	}

	e, _ := New(operator, simplified, "", "")
	return nil, e
}

func New(operator string, operands []e.Evaluable, notOp string, norOp string) (e.Evaluable, error) {
	if len(operands) < 2 {
		return nil, errors.New("logical AND expression must have at least 2 operands")
	}

	return l.New(operator, "AND", operands, handler, func(operator string, ctx map[string]any, operands []e.Evaluable) (any, e.Evaluable) {
		return simplify(operator, ctx, operands)
	})
}
