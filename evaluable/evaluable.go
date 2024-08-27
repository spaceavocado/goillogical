// Evaluable is defined in this package as the core data structure representing
// an evaluable expression. In addition this package contains additional utils
// to work with evaluation context.
package evaluable

import (
	"fmt"
	"reflect"
)

const FlattenContextKey string = "_flattenContext"

// Evaluation data context is used to provide the expression with variable references,
// i.e. this allows for the dynamic expressions. The data context is object with
// properties used as the references keys, and its values as reference values.
//
// Expected value types object, string, number, [] boolean | string | number and
// nested struct of type Context.
//
// Example:
//
//	ctx := Context{
//		"name":    "peter",
//		"country": "canada",
//		"age":     21,
//		"options": []int{1, 2, 3},
//		"address": struct {
//			city    string
//			country string
//		}{
//			city:    "Toronto",
//			country: "Canada",
//		},
//		"index":     2,
//		"segment":   "city",
//		"shapeA":    "box",
//		"shapeB":    "circle",
//		"shapeType": "B",
//	}
type Context = map[string]any

// Evaluation expression kind.
type Kind byte

// Evaluable Kind Identifier, i.e. an unique symbol representing an expression.
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

// Operator mapping represents a map between an expression kind (symbol) and the actual
// text literal denoting an expression.
//
// Example:
// ["==", 1, 1] to be mapped as EQ expression would be represented as:
//
// map[Kind]string{ Eq: "==" }
type OperatorMapping = map[Kind]string

// Evaluable is self-evaluable expression, allowing to serialize and simplify itself.
type Evaluable interface {
	// Evaluate given raw expression in the given context.
	Evaluate(Context) (any, error)
	// Serialized the evaluable into raw data form, i.e. parse-able expression.
	Serialize() any
	// Simplify the evaluable with a given context into an reduced Evaluable or evaluated value.
	Simplify(Context) (any, Evaluable)
	// Get string representation of the Evaluable.
	String() string
}

// Is evaluated primitive predicate.
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

// Flatten context into a map of map[property path]value.
//
// Example:
//
//		ctx := Context{
//			"name":    "peter",
//			"options": []int{1, 2, 3},
//			"address": struct {
//				city    string
//				country string
//			}{
//				city:    "Toronto",
//				country: "Canada",
//			},
//		}
//
//	 flattened = FlattenContext(ctx)
//
//		flattened := Context{
//			"name":    "peter",
//			"options[0]": 1,
//			"options[1]": 2,
//			"options[2]": 3,
//			"address.city": "Toronto",
//			"address.country": "Canada",
//		}
func FlattenContext(ctx Context) map[string]any {
	if ctx == nil {
		return nil
	}

	if value, ok := ctx[FlattenContextKey]; ok && value == FlattenContextKey {
		return ctx
	}

	res := make(map[string]any)
	res[FlattenContextKey] = FlattenContextKey

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
		case reflect.Map:
			for prop, val := range val.(map[string]any) {
				lookup(val, joinPath(path, prop))
			}
		case reflect.Slice:
			for i := 0; i < v.Len(); i++ {
				lookup(v.Index(i).Interface(), fmt.Sprintf("%s[%d]", path, i))
			}
		default:
			return
		}
	}

	lookup(ctx, "")
	return res
}
