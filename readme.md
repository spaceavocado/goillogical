# (go)illogical

A micro conditional javascript engine used to parse the raw logical and comparison expressions, evaluate the expression in the given data context, and provide access to a text form of the given expressions.

> Revision: Feb 7, 2023.

<p align="center">
<a href="https://goreportcard.com/report/github.com/spaceavocado/goillogical"><img class="badge" tag="https://github.com/spaceavocado/goillogical" src="https://goreportcard.com/badge/github.com/spaceavocado/goillogical"></a>
<a href="https://codecov.io/gh/spaceavocado/goillogical" >
<img src="https://codecov.io/gh/spaceavocado/goillogical/branch/master/graph/badge.svg?token=16C3FQ0IA2"/>
</a>
<a href="https://godoc.org/github.com/spaceavocado/goillogical"><img alt="Go Doc" src="https://godoc.org/github.com/spaceavocado/goillogical?status.svg"></a>
</p>

## About

This project has been developed to provide Go(lang) implementation of [spaceavocado/illogical](https://github.com/spaceavocado/illogical).


## Getting Started

```sh
go get -u github.com/spaceavocado/goillogical@latest
```

**Table of Content**

---

- [(go)illogical](#goillogical)
  - [About](#about)
  - [Getting Started](#getting-started)
  - [Basic Usage](#basic-usage)
    - [Evaluate](#evaluate)
    - [Statement](#statement)
    - [Parse](#parse)
    - [Evaluable](#evaluable)
      - [Simplify](#simplify)
      - [Serialize](#serialize)
  - [Working with Expressions](#working-with-expressions)
    - [Evaluation Data Context](#evaluation-data-context)
      - [Accessing Array Element:](#accessing-array-element)
      - [Accessing Array Element via Reference:](#accessing-array-element-via-reference)
      - [Nested Referencing](#nested-referencing)
      - [Composite Reference Key](#composite-reference-key)
      - [Data Type Casting](#data-type-casting)
    - [Operand Types](#operand-types)
      - [Value](#value)
      - [Reference](#reference)
      - [Collection](#collection)
    - [Comparison Expressions](#comparison-expressions)
      - [Equal](#equal)
      - [Not Equal](#not-equal)
      - [Greater Than](#greater-than)
      - [Greater Than or Equal](#greater-than-or-equal)
      - [Less Than](#less-than)
      - [Less Than or Equal](#less-than-or-equal)
      - [In](#in)
      - [Not In](#not-in)
      - [Prefix](#prefix)
      - [Suffix](#suffix)
      - [Overlap](#overlap)
      - [Nil](#nil)
      - [Present](#present)
    - [Logical Expressions](#logical-expressions)
      - [And](#and)
      - [Or](#or)
      - [Nor](#nor)
      - [Xor](#xor)
      - [Not](#not)
  - [Engine Options](#engine-options)
    - [Reference Serialize Options](#reference-serialize-options)
      - [From](#from)
      - [To](#to)
    - [Collection Serialize Options](#collection-serialize-options)
      - [Escape Character](#escape-character)
    - [Simplify Options](#simplify-options)
      - [Ignored Paths](#ignored-paths)
      - [Ignored Paths RegEx](#ignored-paths-regex)
    - [Operator Mapping](#operator-mapping)
    - [Multiple Options](#multiple-options)
  - [Contributing](#contributing)
  - [License](#license)

---


## Basic Usage

```go
import (
	illogical "github.com/spaceavocado/goillogical
)

// Create a new instance of the engine
i := illogical.New()

// Evaluate the raw expression
res, err := i.Evaluate([]any{"==", 1, 1}, map[string]any{})
```

> For advanced usage, please [Engine Options](#engine-options).

### Evaluate

Evaluate comparison or logical expression:

`i.evaluate(`[Comparison Expression](#comparison-expressions) or [Logical Expression](#logical-expressions), [Evaluation Data Context](#evaluation-data-context)`)` => `boolean`

> Data context is optional.

**Example**

```go
ctx := map[string]any{
  "name": "peter",
}

// Comparison expression
i.Evaluate([]any{"==", 5, 5}, ctx)
i.Evaluate([]any{"==", "circle", "circle"}, ctx)
i.Evaluate([]any{"==", true, true}, ctx)
i.Evaluate([]any{"==", "$name", "peter"}, ctx)
i.Evaluate([]any{"NIL", "$RefA"}, ctx)

// Logical expression
i.Evaluate([]any{"AND", []any{"==", 5, 5}, []any{"==", 10, 10}}, ctx)
i.Evaluate([]any{"AND", []any{"==", "circle", "circle"}, []any{"==", 10, 10}}, ctx)
i.Evaluate([]any{"OR", []any{"==", "$name", "peter"}, []any{"==", 5, 10}}, ctx)
```

### Statement

Get expression string representation:

`i.Statement(`[Comparison Expression](#comparison-expressions) or [Logical Expression](#logical-expressions)`)` => `string`

**Example**

```go
// Comparison expression

i.Statement([]any{"==", 5, 5}) // (5 == 5)
i.Statement([]any{"==", "circle", "circle"}) // ("circle" == "circle")
i.Statement([]any{"==", true, true}) // (true == true)
i.Statement([]any{"==", "$name", "peter"}) // ({name} == "peter")
i.Statement([]any{"NIL", "$RefA"}) // ({RefA} <is nil>)

// Logical expression

i.Statement([]any{"AND", []any{"==", 5, 5}, []any{"==", 10, 10}}) // ((5 == 5) AND (10 == 10))
i.Statement([]any{"AND", []any{"==", "circle", "circle"}, []any{"==", 10, 10}}) // (("circle" == "circle") AND (10 == 10))
i.Statement([]any{"OR", []any{"==", "$name", "peter"}, []any{"==", 5, 10}}) // (({name} == "peter") OR (5 == 10))
```

### Parse

Parse the expression into a **Evaluable** object, i.e. it returns the parsed self-evaluable condition expression.

`i.Parse(`[Comparison Expression](#comparison-expressions) or [Logical Expression](#logical-expressions)`)` => `Evaluable`

### Evaluable

- `evaluable.Evaluate(context)` please see [Evaluation Data Context](#evaluation-data-context).
- `evaluable.String()` please see [Statement](#statement).
- `evaluable.Simplify(context)` please see [Simplify](#simplify).

**Example**

```go
e, err := i.Parse([]any{"==", "$name", "peter"})

e.Evaluate(map[string]any{"name": "peter"}) // true
e.String() // ({name} == "peter")
```

#### Simplify

Simplifies an expression with a given context. This is useful when you already have some of
the properties of context and wants to try to evaluate the expression.

**Example**

```go
e, err := i.Parse([]any{"AND", []any{"==", "$a", 10}, []any{"==", "$b", 20}})

e.Simplify(map[string]any{"a": 10}) // ({b} == 20)
e.simplify(map[string]any{"a": 20}) // false
```

Values not found in the context will cause the parent operand not to be evaluated and returned
as part of the simplified expression.

In some situations we might want to evaluate the expression even if referred value is not
present. You can provide a list of keys that will be strictly evaluated even if they are not
present in the context.

**Example**

```go
simplifyOptions := illogical.WithReferenceSimplifyOptions(illogical.SimplifyOptions{
  IgnoredPaths: []string{"ignored"},
  IgnoredPathsRx: []regexp.Regexp{*regexp.MustCompile("^ignored")},
})
i := illogical.New(simplifyOptions)

e, err := i.Parse([]any{"AND", []any{"==", "$a", 10}, []any{"==", "$ignored", 20}})
e.Simplify(map[string]any{"a": 10})

// false
// $ignored" will be evaluated to nil.
```

Alternatively we might want to do the opposite and strictly evaluate the expression for all referred
values not present in the context except for a specified list of optional keys.

**Example**

```go
simplifyOptions := illogical.WithReferenceSimplifyOptions(illogical.SimplifyOptions{
  IgnoredPaths: []string{"b"},
  IgnoredPathsRx: []regexp.Regexp{},
})
i := illogical.New(simplifyOptions)

e, err := i.Parse([]any{"OR", []any{"==", "$a", 10}, []any{"==", "$b", 20}, []any{"==", "$c", 20}})
e.Simplify(map[string]any{"c": 10})

// ({a} == 10)
// except for "$b" everything not in context will be evaluated to nil.
```

#### Serialize

Serializes an expression into the raw data form, reverse parse operation.

**Example**

```go
e, err := i.Parse([]any{"AND", []any{"==", "$a", 10}, []any{"==", 10, 20}})
e.Serialize() // [AND [== $a 10] [== 10 20]]
```

## Working with Expressions

### Evaluation Data Context

The evaluation data context is used to provide the expression with variable references, i.e. this allows for the dynamic expressions. The data context is object with properties used as the references keys, and its values as reference values.

> Valid reference values: object, string, number, [] boolean | string | number.

To reference the nested reference, please use "." delimiter, e.g.:
`$address.city`

#### Accessing Array Element:

`$options[1]`

#### Accessing Array Element via Reference:

`$options[{index}]`

- The **index** reference is resolved within the data context as an array index.

#### Nested Referencing

`$address.{segment}`

- The **segment** reference is resolved within the data context as a property key.

#### Composite Reference Key

`$shape{shapeType}`

- The **shapeType** reference is resolved within the data context, and inserted into the outer reference key.
- E.g. **shapeType** is resolved as "**B**" and would compose the **$shapeB** outer reference.
- This resolution could be n-nested.

#### Data Type Casting

`$payment.amount.(Type)`

Cast the given data context into the desired data type before being used as an operand in the evaluation.

> Note: If the conversion is invalid, then a warning message is being logged.

Supported data type conversions:

- .(String): cast a given reference to String.
- .(Number): cast a given reference to Number.
- .(Integer): cast a given reference to Integer.
- .(Float): cast a given reference to Float.
- .(Boolean): cast a given reference to Boolean.

**Example**

```go
// Data context
ctx := map[string]any{
  "name":    "peter",
  "country": "canada",
  "age":     21,
  "options": []int{1, 2, 3},
  "address": struct {
    city    string
    country string
  }{
    city:    "Toronto",
    country: "Canada",
  },
  "index":     2,
  "segment":   "city",
  "shapeA":    "box",
  "shapeB":    "circle",
  "shapeType": "B",
}

// Evaluate an expression in the given data context
i.Evaluate([]any{">", "$age", 20}, ctx) // true

// Evaluate an expression in the given data context
i.Evaluate([]any{"==", "$address.city", "Toronto"}, ctx) // true

// Accessing Array Element
i.Evaluate([]any{"==", "$options[1]", 2}, ctx) // true

// Accessing Array Element via Reference
i.Evaluate([]any{"==", "$options[{index}]", 3}, ctx) // true

// Nested Referencing
i.Evaluate([]any{"==", "$address.{segment}", "Toronto"}, ctx) // true

// Composite Reference Key
i.Evaluate([]any{"==", "$shape{shapeType}", "circle"}, ctx) // true

// Data Type Casting
i.Evaluate([]any{"==", "$age.(String)", "21"}, ctx) // true
```

### Operand Types

The [Comparison Expression](#comparison-expression) expect operands to be one of the below:

#### Value

Simple value types: string, number, boolean.

**Example**

```go
val1 := 5
var2 := "cirle"
var3 := true

i.Parse([]any{"AND", []any{"==", val1, var2}, []any{"==", var3, var3}})
```

#### Reference

The reference operand value is resolved from the [Evaluation Data Context](#evaluation-data-context), where the the operands name is used as key in the context.

The reference operand must be prefixed with `$` symbol, e.g.: `$name`. This might be customized via [Reference Predicate Parser Option](#reference-predicate).

**Example**

| Expression                    | Data Context      |
| ----------------------------- | ----------------- |
| `["==", "$age", 21]`          | `{age: 21}`       |
| `["==", "circle", "$shape"] ` | `{shape: "circle"}` |
| `["==", "$visible", true]`    | `{visible: true}` |

#### Collection

The operand could be an array mixed from [Value](#value) and [Reference](#reference).

**Example**

| Expression                               | Data Context                        |
| ---------------------------------------- | ----------------------------------- |
| `["IN", [1, 2], 1]`                      | `{}`                                |
| `["IN", "circle", ["$shapeA", "$shapeB"] ` | `{shapeA: "circle", shapeB: "box"}` |
| `["IN", ["$number", 5], 5]`                | `{number: 3}`                       |

### Comparison Expressions

#### Equal

Expression format: `["==", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types: string, number, boolean.

```json
["==", 5, 5]
```

```go
i.Evaluate([]any{"==", 5, 5}, ctx) // true
```

#### Not Equal

Expression format: `["!=", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types: string, number, boolean.

```json
["!=", "circle", "square"]
```

```go
i.Evaluate([]any{"!=", "circle", "square"}, ctx) // true
```

#### Greater Than

Expression format: `[">", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types: number.

```json
[">", 10, 5]
```

```go
i.Evaluate([]any{">", 10, 5}, ctx) // true
```

#### Greater Than or Equal

Expression format: `[">=", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types: number.

```json
[">=", 5, 5]
```

```go
i.Evaluate([]any{">=", 5, 5}, ctx) // true
```

#### Less Than

Expression format: `["<", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types: number.

```json
["<", 5, 10]
```

```go
i.Evaluate([]any{"<", 5, 10}, ctx) // true
```

#### Less Than or Equal

Expression format: `["<=", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types: number.

```json
["<=", 5, 5]
```

```go
i.Evaluate([]any{"<=", 5, 5}, ctx) // true
```

#### In

Expression format: `["IN", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types: number and number[] or string and string[].

```json
["IN", 5, [1, 2, 3, 4, 5]]
["IN", ["circle", "square", "triangle"], "square"]
```

```go
i.Evaluate([]any{"IN", 5, []int{1, 2, 3, 4, 5}}, ctx) // true
i.Evaluate([]any{"IN", []string{"circle", "square", "triangle"}, "square"}, ctx) // true
```

#### Not In

Expression format: `["NOT IN", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types: number and number[] or string and string[].

```json
["IN", 10, [1, 2, 3, 4, 5]]
["IN", ["circle", "square", "triangle"], "oval"]
```

```go
i.Evaluate([]any{"NOT IN", 10, []int{1, 2, 3, 4, 5}}, ctx) // true
i.Evaluate([]any{"NOT IN", []string{"circle", "square", "triangle"}, "oval"}, ctx) // true
```

#### Prefix

Expression format: `["PREFIX", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types: string.

- Left operand is the PREFIX term.
- Right operand is the tested word.

```json
["PREFIX", "hemi", "hemisphere"]
```

```go
i.Evaluate([]any{"PREFIX", "hemi", "hemisphere"}, ctx) // true
i.Evaluate([]any{"PREFIX", "hemi", "sphere"}, ctx) // false
```

#### Suffix

Expression format: `["SUFFIX", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types: string.

- Left operand is the tested word.
- Right operand is the SUFFIX term.

```json
["SUFFIX", "establishment", "ment"]
```

```go
i.Evaluate([]any{"SUFFIX", "establishment", "ment"}, ctx) // true
i.Evaluate([]any{"SUFFIX", "establish", "ment"}, ctx) // false
```

#### Overlap

Expression format: `["OVERLAP", `[Left Operand](#operand-types), [Right Operand](#operand-types)`]`.

> Valid operand types number[] or string[].

```json
["OVERLAP", [1, 2], [1, 2, 3, 4, 5]]
["OVERLAP", ["circle", "square", "triangle"], ["square"]]
```

```go
i.Evaluate([]any{"OVERLAP", []int{1, 2, 6}, []int{1, 2, 3, 4, 5}}, ctx) // true
i.Evaluate([]any{"OVERLAP", []string{"circle", "square", "triangle"}, []string{"square", "oval"}}, ctx) // true
```

#### Nil

Expression format: `["NIL", `[Reference Operand](#reference)`]`.

```json
["NIL", "$RefA"]
```

```go
i.Evaluate([]any{"NIL", "RefA"}, map[string]any{}) // true
i.Evaluate([]any{"NIL", "RefA"}, map[string]any{"RefA": 10}) // false
```

#### Present

Evaluates as FALSE when the operand is UNDEFINED or NULL.

Expression format: `["PRESENT", `[Reference Operand](#reference)`]`.

```json
["PRESENT", "$RefA"]
```

```go
i.Evaluate([]any{"PRESENT", "RefA"}, map[string]any{}) // false
i.Evaluate([]any{"PRESENT", "RefA"}, map[string]any{"RefA": 10}) // true
i.Evaluate([]any{"PRESENT", "RefA"}, map[string]any{"RefA": false}) // true
i.Evaluate([]any{"PRESENT", "RefA"}, map[string]any{"RefA": "val"}) // true
```

### Logical Expressions

#### And

The logical AND operator (&&) returns the boolean value TRUE if both operands are TRUE and returns FALSE otherwise.

Expression format: `["AND", Left Operand 1, Right Operand 2, ... , Right Operand N]`.

> Valid operand types: [Comparison Expression](#comparison-expressions) or [Nested Logical Expression](#logical-expressions).

```json
["AND", ["==", 5, 5], ["==", 10, 10]]
```

```go
i.Evaluate([]any{"AND", []any{"==", 5, 5}, []any{"==", 10, 10}}, ctx) // true
```

#### Or

The logical OR operator (||) returns the boolean value TRUE if either or both operands is TRUE and returns FALSE otherwise.

Expression format: `["OR", Left Operand 1, Right Operand 2, ... , Right Operand N]`.

> Valid operand types: [Comparison Expression](#comparison-expressions) or [Nested Logical Expression](#logical-expressions).

```json
["OR", ["==", 5, 5], ["==", 10, 5]]
```

```go
i.Evaluate([]any{"OR", []any{"==", 5, 5}, []any{"==", 10, 5}}, ctx) // true
```

#### Nor

The logical NOR operator returns the boolean value TRUE if both operands are FALSE and returns FALSE otherwise.

Expression format: `["NOR", Left Operand 1, Right Operand 2, ... , Right Operand N]`

> Valid operand types: [Comparison Expression](#comparison-expressions) or [Nested Logical Expression](#logical-expressions).

```json
["NOR", ["==", 5, 1], ["==", 10, 5]]
```

```go
i.Evaluate([]any{"NOR", []any{"==", 5, 1}, []any{"==", 10, 5}}, ctx) // true
```

#### Xor

The logical NOR operator returns the boolean value TRUE if both operands are FALSE and returns FALSE otherwise.

Expression format: `["XOR", Left Operand 1, Right Operand 2, ... , Right Operand N]`

> Valid operand types: [Comparison Expression](#comparison-expressions) or [Nested Logical Expression](#logical-expressions).

```json
["XOR", ["==", 5, 5], ["==", 10, 5]]
```

```go
i.Evaluate([]any{"XOR", []any{"==", 5, 5}, []any{"==", 10, 5}}, ctx) // true
```

```json
["XOR", ["==", 5, 5], ["==", 10, 10]]
```

```go
i.Evaluate([]any{"XOR", []any{"==", 5, 5}, []any{"==", 10, 10}}, ctx) // false
```

#### Not

The logical NOT operator returns the boolean value TRUE if the operand is FALSE, TRUE otherwise.

Expression format: `["NOT", Operand]`

> Valid operand types: [Comparison Expression](#comparison-expressions) or [Nested Logical Expression](#logical-expressions).

```json
["NOT", ["==", 5, 5]]
```

```go
i.Evaluate([]any{"NOT", []any{"==", 5, 5}}, ctx) // true
```

## Engine Options

### Reference Serialize Options

**Usage**

```go
referenceSerializeOptions := illogical.WithReferenceSerializeOptions(illogical.ReferenceSerializeOptions{
  From: func(string) (string, error)
	To:   func(string) string
})

i := illogical.New(referenceSerializeOptions)
```

#### From

A function used to determine if the operand is a reference type, otherwise evaluated as a static value.

```go
func(string) (string, error)
```

**Return value:**

- `true` = reference type
- `false` = value type

**Default reference predicate:**

> The `$` symbol at the begging of the operand is used to predicate the reference type., E.g. `$State`, `$Country`.

#### To

A function used to transform the operand into the reference annotation stripped form. I.e. remove any annotation used to detect the reference type. E.g. "$Reference" => "Reference".

```go
func(string) string
```

> **Default reference transform:**
> It removes the `$` symbol at the begging of the operand name.

### Collection Serialize Options

**Usage**

```go
collectionSerializeOptions := illogical.WithCollectionSerializeOptions(illogical.CollectionSerializeOptions{
	EscapeCharacter  string
})

i := illogical.New(collectionSerializeOptions)
```

#### Escape Character

Charter used to escape fist value within a collection, if the value contains operator value.

**Example**
- `["==", 1, 1]` // interpreted as EQ expression
- `["\==", 1, 1]` // interpreted as a collection

```go
EscapeCharacter  string
```

> **Default escape character:**
> `\`

### Simplify Options

Options applied while an expression is being simplified.

**Usage**

```go
referenceSimplifyOptions := illogical.WithReferenceSimplifyOptions(illogical.SimplifyOptions{
	IgnoredPaths   []string
	IgnoredPathsRx []regexp.Regexp
})

i := illogical.New(referenceSimplifyOptions)
```

#### Ignored Paths

Reference paths which should be ignored while simplification is applied. Must be an exact match.

```go
IgnoredPaths   []string
```

#### Ignored Paths RegEx

Reference paths which should be ignored while simplification is applied. Matching regular expression patterns.

```go
IgnoredPathsRx []regexp.Regexp
```

### Operator Mapping

Mapping of the operators. The key is unique operator key, and the value is the key used to represent the given operator in the raw expression.

**Usage**

```go
import (
	e "github.com/spaceavocado/goillogical/evaluable"
)

operatorMapping := illogical.WithOperatorMappingOptions(map[e.Kind]string{})

i := illogical.New(operatorMapping)
```

**Default operator mapping:**

```go
operatorMapping := e.OperatorMapping{
  // Comparison
  e.Eq: "==",
  e.Ne: "!=",
  e.Gt: ">",
  e.Ge: ">=",
  e.Lt: "<",
  e.Le: "<=",
  e.In: "IN",
  e.Nin: "NOT IN",
  e.Prefix: "PREFIX",
  e.Suffix: "SUFFIX",
  e.Overlap:, "OVERLAP",
  e.Nil: "NIL",
  e.Present: "PRESENT",
  // Logical
  e.And: "AND",
  e.Or: "OR",
  e.Nor: "NOR",
  e.Xor: "XOR",
  e.Not: "NOT",
}
```

### Multiple Options
All options could be used simultaneously.

```go
operatorMapping := illogical.WithOperatorMappingOptions(map[e.Kind]string{})

referenceSimplifyOptions := illogical.WithReferenceSimplifyOptions(illogical.SimplifyOptions{
  IgnoredPaths:   []string{},
  IgnoredPathsRx: []regexp.Regexp{},
})

i := illogical.New(operatorMapping, referenceSimplifyOptions)
```

---

## Contributing

See [contributing.md](contributing.md).

## License

Illogical is released under the MIT license. See [license.txt](license.txt).
