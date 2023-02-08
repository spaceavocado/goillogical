package nor

import (
	"errors"

	. "github.com/spaceavocado/goillogical/evaluable"
	l "github.com/spaceavocado/goillogical/internal/expression/logical"
	not "github.com/spaceavocado/goillogical/internal/expression/logical/not"
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

func simplify(operator string, ctx Context, operands []Evaluable, notOp string) (any, Evaluable) {
	simplified := []Evaluable{}
	for _, o := range operands {
		res, e := o.Simplify(ctx)
		if b, ok := res.(bool); ok {
			if b == true {
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
		e, _ := not.New("NOT", simplified[0])
		return nil, e
	}

	e, _ := New(operator, simplified, notOp, "")
	return nil, e
}

// not reference needed
func New(operator string, operands []Evaluable, notOp string, norOp string) (Evaluable, error) {
	if len(operands) < 2 {
		return nil, errors.New("logical NOR expression must have at least 2 operands")
	}

	return l.New(operator, "NOR", operands, handler, func(operator string, ctx map[string]any, operands []Evaluable) (any, Evaluable) {
		return simplify(operator, ctx, operands, notOp)
	})
}
