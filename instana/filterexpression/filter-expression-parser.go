package filterexpression

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

//Operator Custom type for and kind of operator
type Operator string

//Capture captures the string representation of an operator from the given string. Interface of participle
func (c *Operator) Capture(values []string) error {
	*c = Operator(strings.ToUpper(values[0]))
	return nil
}

//UnaryOperator Custom type for a unary operations
type UnaryOperator string

//Capture captures the string representation of a unary operation from the given slice of strings. Interface of participle
func (c *UnaryOperator) Capture(values []string) error {
	*c = UnaryOperator(strings.ToUpper(strings.Join(values, " ")))
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

//Conjunction represents a logical and or a logical or conjunction
type Conjunction interface {
	GetLeft() Conjunction
	GetOperator() Operator
	GetRight() Conjunction
}

//LogicalOrExpression representation of a logical OR or as a wrapper for a, LogicalAndExpression or a PrimaryExpression. The wrapping is required to handle precedence.
type LogicalOrExpression struct {
	Left     *LogicalAndExpression `parser:"  @@"`
	Operator *Operator             `parser:"( @\"OR\""`
	Right    *LogicalOrExpression  `parser:"  @@ )?"`
}

//Render implementation of ExpressionRenderer.Render
func (e *LogicalOrExpression) Render() string {
	if e.Operator != nil {
		return fmt.Sprintf("%s OR %s", e.Left.Render(), e.Right.Render())
	}
	return e.Left.Render()
}

//LogicalAndExpression representation of a logical AND or as a wrapper for a PrimaryExpression only. The wrapping is required to handle precedence.
type LogicalAndExpression struct {
	Left     *PrimaryExpression    `parser:"  @@"`
	Operator *Operator             `parser:"( @\"AND\""`
	Right    *LogicalAndExpression `parser:"  @@ )?"`
}

//Render implementation of ExpressionRenderer.Render
func (e *LogicalAndExpression) Render() string {
	if e.Operator != nil {
		return fmt.Sprintf("%s AND %s", e.Left.Render(), e.Right.Render())
	}
	return e.Left.Render()
}

//PrimaryExpression wrapper for either a comparision or a unary expression
type PrimaryExpression struct {
	Comparision    *ComparisionExpression    `parser:"  @@"`
	UnaryOperation *UnaryOperationExpression `parser:"| @@"`
}

//Render implementation of ExpressionRenderer.Render
func (e *PrimaryExpression) Render() string {
	if e.Comparision != nil {
		return e.Comparision.Render()
	}
	return e.UnaryOperation.Render()
}

//ComparisionExpression representation of a comparision expression. Supported types: EQ (Equals), NE (Not Equal), CO (Contains), NC (Not Contain)
type ComparisionExpression struct {
	Key      string   `parser:"@Ident"`
	Operator Operator `parser:"@( \"EQ\" | \"NE\" | \"CO\" | \"NC\" )"`
	Value    string   `parser:"@String"`
}

//Render implementation of ExpressionRenderer.Render
func (e *ComparisionExpression) Render() string {
	return fmt.Sprintf("%s %s '%s'", e.Key, e.Operator, e.Value)
}

//UnaryOperationExpression representation of a unary expression representing a unary operator
type UnaryOperationExpression struct {
	Key      string        `parser:"@Ident"`
	Operator UnaryOperator `parser:"@( \"IS\" (\"EMPTY\" | \"BLANK\")  | \"NOT\" (\"EMPTY\" | \"BLANK\") )"`
}

//Render implementation of ExpressionRenderer.Render
func (e *UnaryOperationExpression) Render() string {
	return fmt.Sprintf("%s %s", e.Key, e.Operator)
}

var (
	filterLexer = lexer.Must(lexer.Regexp(`(\s+)` +
		`|(?P<Keyword>(?i)OR|AND|TRUE|FALSE|IS|NOT|EMPTY|BLANK|EQ|NE|CO|NC)` +
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

//NewParser creates a new instance of a Parser
func NewParser() Parser {
	return new(parserImpl)
}

//Parser interface for working with Dynamic Focus filters of instana
type Parser interface {
	Parse(expression string) (*FilterExpression, error)
}

type parserImpl struct{}

//Parse implementation of the parsing of the Parser
func (f *parserImpl) Parse(expression string) (*FilterExpression, error) {
	parsedExpression := &FilterExpression{}
	err := filterParser.ParseString(expression, parsedExpression)
	if err != nil {
		return &FilterExpression{}, err
	}
	return parsedExpression, nil
}
