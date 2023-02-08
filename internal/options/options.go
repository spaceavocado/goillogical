package options

import (
	"regexp"

	. "github.com/spaceavocado/goillogical/internal"
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
	OperatorMapping OperatorMapping
}

func DefaultOperatorMapping() OperatorMapping {
	return map[Kind]string{
		And:     "AND",
		Or:      "OR",
		Nor:     "NOR",
		Xor:     "XOR",
		Not:     "NOT",
		Eq:      "==",
		Ne:      "!=",
		Gt:      ">",
		Ge:      ">=",
		Lt:      "<",
		Le:      "<=",
		Nil:     "NIL",
		Present: "PRESENT",
		In:      "IN",
		Nin:     "NOT IN",
		Overlap: "OVERLAP",
		Prefix:  "PREFIX",
		Suffix:  "SUFFIX",
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
