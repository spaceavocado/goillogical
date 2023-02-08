package suffix

import (
	"strings"

	e "github.com/spaceavocado/goillogical/evaluable"
	c "github.com/spaceavocado/goillogical/internal/expression/comparison"
)

func op(a string, b string) bool {
	return strings.HasSuffix(a, b)
}

func right(a string, b any) bool {
	switch typed := b.(type) {
	case string:
		return op(a, typed)
	default:
		return false
	}
}

func left(a any, b any) bool {
	switch typed := a.(type) {
	case string:
		return right(typed, b)
	default:
		return false
	}
}

func handler(evaluated []any) bool {
	return left(evaluated[0], evaluated[1])
}

func New(operator string, left e.Evaluable, right e.Evaluable) (e.Evaluable, error) {
	return c.New(operator, "<with suffix>", []e.Evaluable{left, right}, handler)
}
