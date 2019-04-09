package utils

import (
	"github.com/alecthomas/participle"
)

type Expression struct {
	Tag           *Tag         `  @@`
	Conjunction   *Conjunction `| @@`
	SubExpression *Expression  `| "(" @@ ")"`
}

type Tag struct {
	Key      *string `@Ident`
	Operator string  `@( "EQUALS" | "NOT_EQUAL" | "CONTAINS" | "NOT_CONTAIN" | "NOT_EMPTY" )`
	Value    *string `@Ident?`
}

type Conjunction struct {
	Left     *Expression `@@`
	Operator string      `@( "AND" | "OR" )`
	Right    *Expression `@@ `
}

func NewDynamicFocusFiler() DynamicFocusFiler {
	return new(DynamicFocusFilerImpl)
}

type DynamicFocusFiler interface {
	Parse(expression string) (*Expression, error)
	SPrint(expression Expression) (string, error)
}

type DynamicFocusFilerImpl struct{}

func (f *DynamicFocusFilerImpl) Parse(expression string) (*Expression, error) {
	parser, err := participle.Build(&Expression{})
	if err != nil {
		return &Expression{}, err
	}

	parsedExpression := &Expression{}
	err = parser.ParseString(string, parsedExpression)
	if err != nil {
		return &Expression{}, err
	}
	return &Expression{}, nil
}

func (s *DynamicFocusFilerImpl) SPrint(expression Expression) (string, error) {
	return "", nil
}
