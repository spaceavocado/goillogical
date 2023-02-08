package options

import (
	"regexp"

	e "github.com/spaceavocado/goillogical/evaluable"
	c "github.com/spaceavocado/goillogical/internal/operand/collection"
	r "github.com/spaceavocado/goillogical/internal/operand/reference"
)

type Options struct {
	Serialize struct {
		Reference  r.SerializeOptions
		Collection c.SerializeOptions
	}
	Simplify struct {
		Reference r.SimplifyOptions
	}
	OperatorMapping e.OperatorMapping
}

func DefaultOperatorMapping() e.OperatorMapping {
	return map[e.Kind]string{
		e.And:     "AND",
		e.Or:      "OR",
		e.Nor:     "NOR",
		e.Xor:     "XOR",
		e.Not:     "NOT",
		e.Eq:      "==",
		e.Ne:      "!=",
		e.Gt:      ">",
		e.Ge:      ">=",
		e.Lt:      "<",
		e.Le:      "<=",
		e.Nil:     "NIL",
		e.Present: "PRESENT",
		e.In:      "IN",
		e.Nin:     "NOT IN",
		e.Overlap: "OVERLAP",
		e.Prefix:  "PREFIX",
		e.Suffix:  "SUFFIX",
	}
}

func DefaultOptions() Options {
	return Options{
		Serialize: struct {
			Reference  r.SerializeOptions
			Collection c.SerializeOptions
		}{
			Reference:  r.DefaultSerializeOptions(),
			Collection: c.DefaultSerializeOptions(),
		},
		Simplify: struct {
			Reference r.SimplifyOptions
		}{
			Reference: r.SimplifyOptions{
				IgnoredPaths:   []string{},
				IgnoredPathsRx: []regexp.Regexp{},
			},
		},
		OperatorMapping: DefaultOperatorMapping(),
	}
}
