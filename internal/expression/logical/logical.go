package logical

import (
	"fmt"
	. "goillogical/internal"
)

type logical struct {
	kind     Kind
	operator string
	operands []Evaluable
	handler  func(Context, []Evaluable) (bool, error)
}

func (l logical) Kind() Kind {
	return l.kind
}

func (l logical) Evaluate(ctx Context) (any, error) {
	return l.handler(ctx, l.operands)
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

func New(kind Kind, op string, operands []Evaluable, handler func(Context, []Evaluable) (bool, error)) (Evaluable, error) {
	return logical{kind: kind, operator: op, operands: operands, handler: handler}, nil
}
