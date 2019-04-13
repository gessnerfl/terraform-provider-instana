package utils

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

//Boolean custom type to represent a boolean value
type Boolean bool

//Capture captures a boolean value from the given string representation. Interface of participle
func (b *Boolean) Capture(values []string) error {
	*b = "TRUE" == strings.ToUpper(values[0])
	return nil
}

//ExpressionRenderer interface definition for all types of the Filter expression to render the corresponding value
type ExpressionRenderer interface {
	Render() string
}

//FilterExpression repressentation of a dynamic focus filter expression
type FilterExpression struct {
	Expression *LogicalOrExpression `parser:"@@"`
}

//Render implementation of ExpressionRenderer.Render
func (e *FilterExpression) Render() string {
	return e.Expression.Render()
}

//LogicalOrExpression representation of a logical OR or as a wrapper for a, LogicalAndExpression or a PrimaryExpression. The wrapping is required to handle precedence.
type LogicalOrExpression struct {
	Left     *LogicalAndExpression `parser:"  @@"`
	Operator string                `parser:"( @\"OR\""`
	Right    *LogicalOrExpression  `parser:"  @@ )?"`
}

//Render implementation of ExpressionRenderer.Render
func (e *LogicalOrExpression) Render() string {
	if "OR" == strings.ToUpper(e.Operator) {
		return fmt.Sprintf("%s OR %s", e.Left.Render(), e.Right.Render())
	}
	return e.Left.Render()
}

//LogicalAndExpression representation of a logical AND or as a wrapper for a PrimaryExpression only. The wrapping is required to handle precedence.
type LogicalAndExpression struct {
	Left     *PrimaryExpression    `parser:"  @@"`
	Operator string                `parser:"( @\"AND\""`
	Right    *LogicalAndExpression `parser:"  @@ )?"`
}

//Render implementation of ExpressionRenderer.Render
func (e *LogicalAndExpression) Render() string {
	if "AND" == strings.ToUpper(e.Operator) {
		return fmt.Sprintf("%s AND %s", e.Left.Render(), e.Right.Render())
	}
	return e.Left.Render()
}

//PrimaryExpression wrapper for either a comparision or a unary expression
type PrimaryExpression struct {
	Comparision     *ComparisionExpression `parser:"  @@"`
	UnaryExpression *UnaryExpression       `parser:"| @@"`
}

//Render implementation of ExpressionRenderer.Render
func (e *PrimaryExpression) Render() string {
	if e.Comparision != nil {
		return e.Comparision.Render()
	}
	return e.UnaryExpression.Render()
}

//ComparisionExpression representation of a comparision expression. Supported types: EQ (Equals), NE (Not Equal), CO (Contains), NC (Not Contain)
type ComparisionExpression struct {
	Key      string `parser:"@Ident"`
	Operator string `parser:"@( \"EQ\" | \"NE\" | \"CO\" | \"NC\" )"`
	Value    *Value `parser:"@@"`
}

//Render implementation of ExpressionRenderer.Render
func (e *ComparisionExpression) Render() string {
	return fmt.Sprintf("%s %s %s", e.Key, strings.ToUpper(e.Operator), e.Value.Render())
}

//UnaryExpression representation of a unary expression representing a function
type UnaryExpression struct {
	Key      string `parser:"@Ident"`
	Function string `parser:"@( \"IS\" \"EMPTY\" | \"NOT\" \"EMPTY\" )"`
}

//Render implementation of ExpressionRenderer.Render
func (e *UnaryExpression) Render() string {
	return fmt.Sprintf("%s %s", e.Key, e.formatFunctionName())
}

func (e *UnaryExpression) formatFunctionName() string {
	if strings.ToUpper(e.Function) == "NOTEMPTY" {
		return "NOT EMPTY"
	}
	if strings.ToUpper(e.Function) == "ISEMPTY" {
		return "IS EMPTY"
	}
	return "<unknown function>"
}

//Value representation of a term of an expression
type Value struct {
	String  *string  `parser:"  @String"`
	Number  *float64 `parser:"| @Number"`
	Boolean *Boolean `parser:"| @(\"TRUE\" | \"FALSE\")"`
}

//Render implementation of ExpressionRenderer.Render
func (e *Value) Render() string {
	if e.Boolean != nil {
		return fmt.Sprintf("%t", *e.Boolean)
	}
	if e.Number != nil {
		return fmt.Sprintf("%f", *e.Number)
	}
	return *e.String
}

var (
	filterLexer = lexer.Must(lexer.Regexp(`(\s+)` +
		`|(?P<Keyword>(?i)OR|AND|TRUE|FALSE|IS|NOT|EMPTY|EQ|NE|CO|NC)` +
		`|(?P<Ident>[a-zA-Z_][\.a-zA-Z0-9_]*)` +
		`|(?P<Number>[-+]?\d+(\.\d+)?)` +
		`|(?P<String>'[^']*'|"[^"]*")` +
		`|(?P<Operators>EQ|NE|CO|NC)`,
	))
	filterParser = participle.MustBuild(
		&FilterExpression{},
		participle.Lexer(filterLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword", "Operators"),
	)
)

//NewFilterExpressionParser creates a new instance of a FilterExpressionParser
func NewFilterExpressionParser() FilterExpressionParser {
	return new(filterExpressionParserImpl)
}

//FilterExpressionParser interface for working with Dynamic Focus filters of instana
type FilterExpressionParser interface {
	Parse(expression string) (*FilterExpression, error)
}

type filterExpressionParserImpl struct{}

//Parse implementation of the parsing of the FilterExpressionParser
func (f *filterExpressionParserImpl) Parse(expression string) (*FilterExpression, error) {
	parsedExpression := &FilterExpression{}
	err := filterParser.ParseString(expression, parsedExpression)
	if err != nil {
		return &FilterExpression{}, err
	}
	return parsedExpression, nil
}
