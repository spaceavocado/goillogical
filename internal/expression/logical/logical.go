package logical

import (
	"fmt"

	. "github.com/spaceavocado/goillogical/evaluable"
)

type Handler func(Context, []Evaluable) (bool, error)
type Simplify func(string, Context, []Evaluable) (any, Evaluable)

type logical struct {
	kind     string
	operator string
	operands []Evaluable
	handler  Handler
	simplify Simplify
}

func (l logical) Evaluate(ctx Context) (any, error) {
	return l.handler(ctx, l.operands)
}

func (l logical) Serialize() any {
	res := []any{l.kind}
	for i := 0; i < len(l.operands); i++ {
		res = append(res, l.operands[i].Serialize())
	}
	return res
}

func (l logical) Simplify(ctx Context) (any, Evaluable) {
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

func Evaluate(ctx Context, o Evaluable) (bool, error) {
	res, err := o.Evaluate(ctx)
	if err != nil {
		return false, err
	}
	switch res.(type) {
	case bool:
		return res.(bool), nil
	default:
		return false, nil
	}
}

func New(kind string, op string, operands []Evaluable, handler Handler, simplify Simplify) (Evaluable, error) {
	return logical{kind, op, operands, handler, simplify}, nil
}
