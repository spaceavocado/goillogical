package internal

import (
	"fmt"
	"reflect"
)

type Context = map[string]any

type Kind byte

const (
	Unknown Kind = iota
	Value
	Reference
	Collection
	And
	Or
	Nor
	Xor
	Not
	Eq
	Ne
	Gt
	Ge
	Lt
	Le
	Nil
	Present
	In
	Nin
	Overlap
	Prefix
	Suffix
)

type OperatorMapping = map[Kind]string

type Evaluable interface {
	Evaluate(Context) (any, error)
	Serialize() any
	Simplify(Context) (any, Evaluable)
	String() string
}

func IsEvaluatedPrimitive(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return true
	case float32, float64:
		return true
	case bool:
		return true
	case string:
		return true
	default:
		return false
	}
}

func FlattenContext(ctx Context) map[string]any {
	res := make(map[string]any)
	var lookup func(p any, path string)

	joinPath := func(a string, b string) string {
		if len(a) == 0 {
			return b
		}
		return fmt.Sprintf("%s.%s", a, b)
	}

	lookup = func(val any, path string) {
		v := reflect.ValueOf(val)
		switch v.Kind() {
		case reflect.Bool:
			fallthrough
		case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
			fallthrough
		case reflect.Float32, reflect.Float64:
			fallthrough
		case reflect.String:
			res[path] = val
			break
		case reflect.Map:
			for prop, val := range val.(map[string]any) {
				lookup(val, joinPath(path, prop))
			}
			break
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				lookup(v.Index(i).Interface(), fmt.Sprintf("%s[%d]", path, i))
			}
			break
		default:
			return
		}
	}

	lookup(ctx, "")
	return res
}
