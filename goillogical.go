// A micro conditional go engine used to parse the raw logical and comparison expressions,
// evaluate the expression in the given data context, and provide access to a text form of
// the given expressions.
package goillogical

import (
	e "github.com/spaceavocado/goillogical/evaluable"
	c "github.com/spaceavocado/goillogical/internal/operand/collection"
	r "github.com/spaceavocado/goillogical/internal/operand/reference"
	o "github.com/spaceavocado/goillogical/internal/options"
	p "github.com/spaceavocado/goillogical/internal/parser"
)

// Simplify options for reference expression.
type SimplifyOptions = r.SimplifyOptions

// Serialization options for reference expression.
type ReferenceSerializeOptions = r.SerializeOptions

// Serialization options for collection expression.
type CollectionSerializeOptions = struct {
	// Charter used to escape fist value within a collection, if the value contains operator value.
	//
	// 	Example:
	//
	// - `["==", 1, 1]` // interpreted as EQ expression
	// - `["\==", 1, 1]` // interpreted as a collection
	EscapeCharacter string
}

// Goillogical engine engine providing access to parsing, evaluation and simplification
// of expressions.
type Goillogical interface {
	Evaluate(any, e.Context) (any, error)
	Parse(any) (e.Evaluable, error)
	Statement(any) (string, error)
	Simplify(any, e.Context) (any, e.Evaluable, error)
}

type illogical struct {
	opts   o.Options
	parser p.Parser
}

// Evaluate given raw expression in the given context.
//
// Example:
//
//	ctx := map[string]any{
//	  "name": "peter",
//	}
//
// i.Evaluate([]any{"==", 5, 5}, ctx)
// i.Evaluate([]any{"AND", []any{"==", 5, 5}, []any{"==", 10, 10}}, ctx)
func (i illogical) Evaluate(exp any, ctx e.Context) (any, error) {
	eval, err := i.parser.Parse(exp)
	if err != nil {
		return nil, err
	}
	return eval.Evaluate(e.FlattenContext(ctx))
}

// Parse given raw expression in an Evaluable object., i.e. it returns the parsed
// self-evaluable condition expression.
//
// Example:
//
// e, err := i.Parse([]any{"==", "$name", "peter"})
//
// e.Evaluate(map[string]any{"name": "peter"}) // true
// e.String() // ({name} == "peter")
func (i illogical) Parse(exp any) (e.Evaluable, error) {
	return i.parser.Parse(exp)
}

// Get expression string representation.
//
// Example:
//
// i.Statement([]any{"==", 5, 5}) // (5 == 5)
// i.Statement([]any{"AND", []any{"==", 5, 5}, []any{"==", 10, 10}}) // ((5 == 5) AND (10 == 10))
func (i illogical) Statement(exp any) (string, error) {
	e, err := i.parser.Parse(exp)
	if err != nil {
		return "", err
	}
	return e.String(), nil
}

// Simplify an expression with a given context. This is useful when you already have some of
// the properties of context and wants to try to evaluate the expression.
//
// Example:
//
// e, err := i.Parse([]any{"AND", []any{"==", "$a", 10}, []any{"==", "$b", 20}})
//
// e.Simplify(map[string]any{"a": 10}) // ({b} == 20)
// e.simplify(map[string]any{"a": 20}) // false
func (i illogical) Simplify(exp any, ctx e.Context) (any, e.Evaluable, error) {
	eval, err := i.parser.Parse(exp)
	if err != nil {
		return nil, nil, err
	}

	val, eval := eval.Simplify(e.FlattenContext(ctx))
	return val, eval, nil
}

// Option customizing the operator mapping, simplification and serialization of evaluables.
type Option func(*illogical)

// Illogical with reference custom serialization options
//
// Example:
//
//	referenceSerializeOptions := illogical.WithReferenceSerializeOptions(illogical.ReferenceSerializeOptions{
//	  From: func(string) (string, error)
//		To:   func(string) string
//	})
//
// i := illogical.New(referenceSerializeOptions)
func WithReferenceSerializeOptions(o ReferenceSerializeOptions) Option {
	return func(i *illogical) {
		i.opts.Serialize.Reference = o
	}
}

// Illogical with collection custom serialization options
//
// Example:
//
//	collectionSerializeOptions := illogical.WithCollectionSerializeOptions(illogical.CollectionSerializeOptions{
//		EscapeCharacter  string
//	})
//
// i := illogical.New(collectionSerializeOptions)
func WithCollectionSerializeOptions(o CollectionSerializeOptions) Option {
	return func(i *illogical) {
		i.opts.Serialize.Collection = c.SerializeOptions{
			EscapeCharacter:  o.EscapeCharacter,
			EscapedOperators: map[string]bool{},
		}
	}
}

// Illogical with reference custom simplification options
//
// Example:
//
//	referenceSimplifyOptions := illogical.WithReferenceSimplifyOptions(illogical.SimplifyOptions{
//		IgnoredPaths   []string
//		IgnoredPathsRx []regexp.Regexp
//	})
//
// i := illogical.New(referenceSimplifyOptions)
func WithReferenceSimplifyOptions(o SimplifyOptions) Option {
	return func(i *illogical) {
		i.opts.Simplify.Reference = o
	}
}

// Illogical with custom operator mapping.
// Mapping of the operators. The key is unique operator key, and the value is the key used to
// represent the given operator in the raw expression.
//
// Example:
//
// import (
//
//	e "github.com/spaceavocado/goillogical/evaluable"
//
// )
//
// operatorMapping := illogical.WithOperatorMappingOptions(map[e.Kind]string{})
//
// i := illogical.New(operatorMapping)
//
// Default mapping:
//
//	operatorMapping := e.OperatorMapping{
//		// Comparison
//		e.Eq: "==",
//		e.Ne: "!=",
//		e.Gt: ">",
//		e.Ge: ">=",
//		e.Lt: "<",
//		e.Le: "<=",
//		e.In: "IN",
//		e.Nin: "NOT IN",
//		e.Prefix: "PREFIX",
//		e.Suffix: "SUFFIX",
//		e.Overlap:, "OVERLAP",
//		e.Nil: "NIL",
//		e.Present: "PRESENT",
//		// Logical
//		e.And: "AND",
//		e.Or: "OR",
//		e.Nor: "NOR",
//		e.Xor: "XOR",
//		e.Not: "NOT",
//	  }
func WithOperatorMappingOptions(m e.OperatorMapping) Option {
	return func(i *illogical) {
		i.opts.OperatorMapping = m
	}
}

// Create new instance of the (go)illogical
func New(opts ...Option) Goillogical {
	i := &illogical{o.DefaultOptions(), nil}

	for _, opt := range opts {
		opt(i)
	}

	for _, op := range i.opts.OperatorMapping {
		i.opts.Serialize.Collection.EscapedOperators[op] = true
	}

	i.parser = p.New(&i.opts)
	return i
}
