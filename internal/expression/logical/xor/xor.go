package xor

import (
	"errors"

	e "github.com/spaceavocado/goillogical/evaluable"
	l "github.com/spaceavocado/goillogical/internal/expression/logical"
	nor "github.com/spaceavocado/goillogical/internal/expression/logical/nor"
	not "github.com/spaceavocado/goillogical/internal/expression/logical/not"
)

func handler(ctx e.Context, operands []e.Evaluable) (bool, error) {
	var flattenContext = e.FlattenContext(ctx)
	var xor bool

	for i, o := range operands {
		res, err := l.Evaluate(flattenContext, o)
		if err != nil {
			return false, err
		}

		if i == 0 {
			xor = res
			continue
		}

		if xor && res {
			return false, nil
		}

		if res {
			xor = true
		}
	}
	return xor, nil
}

func simplify(operator string, ctx e.Context, operands []e.Evaluable, notOp string, norOp string) (any, e.Evaluable) {
	var flattenContext = e.FlattenContext(ctx)

	truthy := 0
	simplified := []e.Evaluable{}
	for _, o := range operands {
		res, e := o.Simplify(flattenContext)
		if b, ok := res.(bool); ok {
			if b {
				truthy++
			}
			if truthy > 1 {
				return false, nil
			}
			continue
		}

		simplified = append(simplified, e)
	}

	if len(simplified) == 0 {
		return truthy == 1, nil
	}

	if len(simplified) == 1 {
		if truthy == 1 {
			e, _ := not.New(notOp, simplified[0])
			return nil, e
		}
		return nil, simplified[0]
	}

	if truthy == 1 {
		e, _ := nor.New(norOp, simplified, notOp, norOp)
		return nil, e
	}

	e, _ := New(operator, simplified, notOp, norOp)
	return nil, e
}

func New(operator string, operands []e.Evaluable, notOp string, norOp string) (e.Evaluable, error) {
	if len(operands) < 2 {
		return nil, errors.New("logical XOR expression must have at least 2 operands")
	}

	return l.New(operator, "XOR", operands, handler, func(operator string, ctx map[string]any, operands []e.Evaluable) (any, e.Evaluable) {
		return simplify(operator, ctx, operands, notOp, norOp)
	})
}
