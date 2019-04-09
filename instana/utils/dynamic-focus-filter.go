package utils

import (
	"github.com/alecthomas/participle"
)

//Expression repressentation of a dynamic focus filter expression
type Expression struct {
	//Comparision, field is set if the expression is a comparision of a path with a comparator and a value
	Comparision *ComparisionExpression `  @@`
	//NotEmpty, field is set if the expression is a not empty function on a path
	NotEmpty *NotEmptyExpression `|  @@`
	//Conjunction, field is set if the expression is a conjunction
	Conjunction *Conjunction `| @@`
	//SubExpression, field is set if the expression is a sub expression in parenthes
	SubExpression *Expression `| "(" @@ ")"`
}

//ComparisionExpression representation of a compatision expression which compares a tag with a given value.
type ComparisionExpression struct {
	//Key, the key of the comparision. The key represents the entity/span/... path expression e.g. entity.os.type or in case of tags entity.agent.tag.env where env is the tag name
	Key string `@Ident @{ "." Ident }`
	//Operator, the comparision operator
	Operator string `@( "EQ" | "NE" | "CO" | "NC" )`
	//Value the value used for comparision
	Value string `@String`
}

//BinaryFunctionExpression a binary function on a path expression
type NotEmptyExpression struct {
	Key string `"NOT_EMPTY("@Ident @{ "." Ident }")"`
}

//Conjunction representation of a conjunction expression
type Conjunction struct {
	//Left the left expression of the conjunction
	Left *Expression `@@`
	//Operator the conjunction operator which is either a logical AND or a logical OR
	Operator string `@( "AND" | "OR" )`
	//Right the right expression of the conjunction
	Right *Expression `@@ `
}

//NewDynamicFocusFilter creates a new instance of a DynamicFocusFilter
func NewDynamicFocusFilter() DynamicFocusFilter {
	return new(dynamicFocusFilterImpl)
}

//DynamicFocusFilter interface for working with Dynamic Focus filters of instana
type DynamicFocusFilter interface {
	Parse(expression string) (*Expression, error)
	SPrint(expression Expression) (string, error)
}

type dynamicFocusFilterImpl struct{}

//Parse implementation of the parsing of the DynamicFocusFilter
func (f *dynamicFocusFilterImpl) Parse(expression string) (*Expression, error) {
	participle.UseLookahead(3)
	parser, err := participle.Build(&Expression{})
	if err != nil {
		return &Expression{}, err
	}

	parsedExpression := &Expression{}
	err = parser.ParseString(expression, parsedExpression)
	if err != nil {
		return &Expression{}, err
	}
	return &Expression{}, nil
}

//SPrint implementation of the printing of the DynamicFocusFilter. SPrint renders the expression as string
func (f *dynamicFocusFilterImpl) SPrint(expression Expression) (string, error) {
	return "", nil
}
