package comparison

import (
	"fmt"
	. "goillogical/internal"
	"reflect"
)

type Number interface {
	~int | ~float32 | ~float64
}

type comparison struct {
	kind     Kind
	operator string
	operands []Evaluable
	handler  func([]any) bool
}

func (c comparison) Kind() Kind {
	return c.kind
}

func (c comparison) Evaluate(ctx Context) (any, error) {
	evaluated := make([]any, len(c.operands))
	for i, e := range c.operands {
		val, err := e.Evaluate(ctx)
		if err != nil {
			return nil, err
		}
		evaluated[i] = val
	}
	return c.handler(evaluated), nil
}

func (c comparison) String() string {
	res := fmt.Sprintf("(%s %s", c.operands[0].String(), c.operator)
	if len(c.operands) > 1 {
		res += fmt.Sprintf(" %s", c.operands[1].String())
	}
	return res + ")"
}

func IsComparable(left any, right any) bool {
	t1 := reflect.TypeOf(left).Kind()
	t2 := reflect.TypeOf(right).Kind()
	if t1 != t2 {
		return false
	}
	if t1 == reflect.Slice || t2 == reflect.Slice {
		return false
	}
	return true
}

func New(kind Kind, op string, operands []Evaluable, handler func([]any) bool) (Evaluable, error) {
	return comparison{operator: op, operands: operands, handler: handler}, nil
}
