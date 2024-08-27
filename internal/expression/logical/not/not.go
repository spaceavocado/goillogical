package not

import (
	e "github.com/spaceavocado/goillogical/evaluable"
	l "github.com/spaceavocado/goillogical/internal/expression/logical"
)

func handler(ctx e.Context, operands []e.Evaluable) (bool, error) {
	var flattenContext = e.FlattenContext(ctx)

	res, err := l.Evaluate(flattenContext, operands[0])
	if err != nil {
		return false, err
	}

	return !res, nil
}

func simplify(operator string, ctx e.Context, operands []e.Evaluable) (any, e.Evaluable) {
	var flattenContext = e.FlattenContext(ctx)

	res, eval := operands[0].Simplify(flattenContext)
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
