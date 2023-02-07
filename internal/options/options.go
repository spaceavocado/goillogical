package options

import (
	. "goillogical/internal"
	c "goillogical/internal/operand/collection"
	r "goillogical/internal/operand/reference"
	"regexp"
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
		Nin:     "NON IT",
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
