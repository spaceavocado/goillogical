package internal

type Context = map[string]any

// type Kind byte

// const (
// 	Value Kind = iota
// 	Reference
// 	Collection
// )

type Evaluable interface {
	// Kind() Kind
	Evaluate(Context) (any, error)
	String() string
}
