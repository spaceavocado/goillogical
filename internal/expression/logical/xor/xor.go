package xor

import (
	"errors"

	. "github.com/spaceavocado/goillogical/evaluable"
	l "github.com/spaceavocado/goillogical/internal/expression/logical"
	nor "github.com/spaceavocado/goillogical/internal/expression/logical/nor"
	not "github.com/spaceavocado/goillogical/internal/expression/logical/not"
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

func simplify(operator string, ctx Context, operands []Evaluable, notOp string, norOp string) (any, Evaluable) {
	truthy := 0
	simplified := []Evaluable{}
	for _, o := range operands {
		res, e := o.Simplify(ctx)
		if b, ok := res.(bool); ok {
			if b == true {
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

// not, nor reference needed
func New(operator string, operands []Evaluable, notOp string, norOp string) (Evaluable, error) {
	if len(operands) < 2 {
		return nil, errors.New("logical XOR expression must have at least 2 operands")
	}

	return l.New(operator, "XOR", operands, handler, func(operator string, ctx map[string]any, operands []Evaluable) (any, Evaluable) {
		return simplify(operator, ctx, operands, notOp, norOp)
	})
}
