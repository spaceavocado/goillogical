package internal

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
	Kind() Kind
	Evaluate(Context) (any, error)
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
