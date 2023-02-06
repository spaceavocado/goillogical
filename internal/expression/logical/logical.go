package logical

import (
	"fmt"
	. "goillogical/internal"
)

type logical struct {
	kind     string
	operator string
	operands []Evaluable
	handler  func(Context, []Evaluable) (bool, error)
}

func (l logical) Evaluate(ctx Context) (any, error) {
	return l.handler(ctx, l.operands)
}

func (l logical) Serialize() any {
	res := []any{l.kind}
	for i := 1; i < len(l.operands); i++ {
		res = append(res, l.operands[i].Serialize())
	}
	return res
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

func New(kind string, op string, operands []Evaluable, handler func(Context, []Evaluable) (bool, error)) (Evaluable, error) {
	return logical{kind: kind, operator: op, operands: operands, handler: handler}, nil
}
