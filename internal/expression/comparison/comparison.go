package comparison

import (
	"fmt"
	. "goillogical/internal"
	"reflect"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

type comparison struct {
	kind     string
	operator string
	operands []Evaluable
	handler  func([]any) bool
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

func (c comparison) Serialize() any {
	res := []any{c.kind}
	for i := 1; i < len(c.operands); i++ {
		res = append(res, c.operands[i].Serialize())
	}
	return res
}

func (c comparison) Simplify(ctx Context) (any, Evaluable) {
	res := []any{}
	for _, o := range c.operands {
		val, e := o.Simplify(ctx)
		if e != nil {
			return nil, &c
		}
		res = append(res, val)
	}

	return c.handler(res), nil
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

func New(kind string, op string, operands []Evaluable, handler func([]any) bool) (Evaluable, error) {
	return comparison{kind: kind, operator: op, operands: operands, handler: handler}, nil
}
