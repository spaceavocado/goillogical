package le

import (
	. "goillogical/internal"
	c "goillogical/internal/expression/comparison"
)

func op[T c.Number](a T, b T) bool {
	return a <= b
}

func right[T c.Number](a T, b any) bool {
	switch b.(type) {
	case T:
		return op(a, b.(T))
	default:
		return false
	}
}

func left(a any, b any) bool {
	switch a.(type) {
	case int:
		return right(a.(int), b)
	case float32:
		return right(a.(float32), b)
	case float64:
		return right(a.(float64), b)
	default:
		return false
	}
}

func handler(evaluated []any) bool {
	return left(evaluated[0], evaluated[1])
}

func New(left Evaluable, right Evaluable) (Evaluable, error) {
	return c.New(Le, "<=", []Evaluable{left, right}, handler)
}
