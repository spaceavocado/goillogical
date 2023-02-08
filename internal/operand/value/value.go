package value

import (
	e "github.com/spaceavocado/goillogical/evaluable"

	"errors"
	"fmt"
)

type value struct {
	val any
}

func (v value) Evaluate(ctx e.Context) (any, error) {
	return v.val, nil
}

func (v value) Serialize() any {
	return v.val
}

func (v value) Simplify(e.Context) (any, e.Evaluable) {
	return v.val, nil
}

func (v value) String() string {
	switch v.val.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v.val)
	default:
		return fmt.Sprintf("%v", v.val)
	}
}

func isPrimitive(v any) bool {
	switch v.(type) {
	case string, int, int8, int16, int32, int64, float32, float64, bool:
		return true
	default:
		return false
	}
}

func New(val any) (e.Evaluable, error) {
	if !isPrimitive(val) {
		return nil, errors.New("value could be only primitive type, string, number or bool")
	}
	return value{val}, nil
}
