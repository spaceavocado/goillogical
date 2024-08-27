package logical

import (
	"errors"
	"fmt"

	e "github.com/spaceavocado/goillogical/evaluable"
)

type Handler func(e.Context, []e.Evaluable) (bool, error)
type Simplify func(string, e.Context, []e.Evaluable) (any, e.Evaluable)

type logical struct {
	kind     string
	operator string
	operands []e.Evaluable
	handler  Handler
	simplify Simplify
}

func (l logical) Evaluate(ctx e.Context) (any, error) {
	return l.handler(ctx, l.operands)
}

func (l logical) Serialize() any {
	res := []any{l.kind}
	for i := 0; i < len(l.operands); i++ {
		res = append(res, l.operands[i].Serialize())
	}
	return res
}

func (l logical) Simplify(ctx e.Context) (any, e.Evaluable) {
	return l.simplify(l.kind, ctx, l.operands)
}

func (l logical) String() string {
	res := "("
	for i := 0; i < len(l.operands); i++ {
		res += l.operands[i].String()
		if i < len(l.operands)-1 {
			res += fmt.Sprintf(" %s ", l.operator)
		}
	}
	return res + ")"
}

func Evaluate(ctx e.Context, o e.Evaluable) (bool, error) {
	res, err := o.Evaluate(ctx)
	if err != nil {
		return false, err
	}
	switch typed := res.(type) {
	case bool:
		return typed, nil
	default:
		return false, errors.New("invalid evaluated operand, must be boolean value")
	}
}

func New(kind string, op string, operands []e.Evaluable, handler Handler, simplify Simplify) (e.Evaluable, error) {
	return logical{kind, op, operands, handler, simplify}, nil
}
