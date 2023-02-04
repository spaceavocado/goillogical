package value

import (
	. "goillogical/internal"

	"errors"
	"fmt"
)

type value struct {
	val any
}

// func (v value) Kind() Kind {
// 	return Value
// }

func (v value) String() string {
	switch v.val.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v.val)
	default:
		return fmt.Sprintf("%s", v.val)
	}
}

func (v value) Evaluate(ctx Context) (any, error) {
	return v.val, nil
}

func isPrimitive(v any) bool {
	switch v.(type) {
	case string, int, float32, float64, bool:
		return true
	default:
		return false
	}
}

func New(val any) (Evaluable, error) {
	if !isPrimitive(val) {
		return nil, errors.New("value could be only primitive type, string, int or bool")
	}
	return value{val}, nil
}
