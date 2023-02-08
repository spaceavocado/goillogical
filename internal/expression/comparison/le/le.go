package le

import (
	e "github.com/spaceavocado/goillogical/evaluable"
	c "github.com/spaceavocado/goillogical/internal/expression/comparison"
)

func op[T c.Number](a T, b T) bool {
	return a <= b
}

func right[T c.Number](a T, b any) bool {
	switch typed := b.(type) {
	case T:
		return op(a, typed)
	default:
		return false
	}
}

func left(a any, b any) bool {
	switch typed := a.(type) {
	case int:
		return right(typed, b)
	case float32:
		return right(typed, b)
	case float64:
		return right(typed, b)
	default:
		return false
	}
}

func handler(evaluated []any) bool {
	return left(evaluated[0], evaluated[1])
}

func New(operator string, left e.Evaluable, right e.Evaluable) (e.Evaluable, error) {
	return c.New(operator, "<=", []e.Evaluable{left, right}, handler)
}
