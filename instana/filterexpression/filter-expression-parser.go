package filterexpression

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//ExpressionRenderer interface definition for all types of the Filter expression to render the corresponding value
type ExpressionRenderer interface {
	Render() string
}

//EntityOrigin custom type for the origin (source or destination) of a entity spec
type EntityOrigin interface {
	//Key returns the key of the entity origin
	Key() string
	//MatcherExpressionEntity returns the Instana API Matcher Expression Entity
	MatcherExpressionEntity() restapi.MatcherExpressionEntity
}

func newEntityOrigin(key string, entity restapi.MatcherExpressionEntity) EntityOrigin {
	return &baseEntityOrigin{key: key, instanaAPIEntity: entity}
}

type baseEntityOrigin struct {
	key              string
	instanaAPIEntity restapi.MatcherExpressionEntity
}

//Key interface implementation of EntityOrigin
func (o *baseEntityOrigin) Key() string {
	return o.key
}

//MatcherExpressionEntity interface implementation of EntityOrigin
func (o *baseEntityOrigin) MatcherExpressionEntity() restapi.MatcherExpressionEntity {
	return o.instanaAPIEntity
}

var (
	//EntityOriginSource constant value for the source EntityOrigin
	EntityOriginSource = newEntityOrigin("src", restapi.MatcherExpressionEntitySource)
	//EntityOriginDestination constant value for the destination EntityOrigin
	EntityOriginDestination = newEntityOrigin("dest", restapi.MatcherExpressionEntityDestination)
	//EntityOriginNotApplicable constant value for the not applicable EntityOrigin
	EntityOriginNotApplicable = newEntityOrigin("na", restapi.MatcherExpressionEntityNotApplicable)
)

//EntityOrigins custom type for a slice of entity origins
type EntityOrigins []EntityOrigin

//ForInstanaAPIEntity returns the EntityOrigin for its cooresponding MatchExpressionEntity from the Instana API
func (origins EntityOrigins) ForInstanaAPIEntity(input restapi.MatcherExpressionEntity) EntityOrigin {
	for _, o := range origins {
		if o.MatcherExpressionEntity() == input {
			return o
		}
	}
	log.Printf("match specification entity %s is not supported; fall back to default origin %s", input, EntityOriginDestination.Key())
	return EntityOriginDestination
}

//ForKey returns the EntityOrigin for its string representation
func (origins EntityOrigins) ForKey(input string) EntityOrigin {
	for _, o := range origins {
		if o.Key() == input {
			return o
		}
	}
	log.Printf("entity origin %s is not supported; fall back to default origin %s", input, EntityOriginDestination.Key())
	return EntityOriginDestination
}

//SupportedEntityOrigins slice of supported EntityOrigins
var SupportedEntityOrigins = EntityOrigins{EntityOriginSource, EntityOriginDestination, EntityOriginNotApplicable}

//EntitySpec custom type for any kind of entity path specification
type EntitySpec struct {
	Key           string
	Origin        EntityOrigin
	OriginDefined bool
}

//Capture captures the string representation of an entity path from the given string. Interface of participle
func (o *EntitySpec) Capture(values []string) error {
	val := values[0]
	if val == "@" {
		o.OriginDefined = true
	} else if o.OriginDefined {
		o.Origin = SupportedEntityOrigins.ForKey(val)
	} else {
		*o = EntitySpec{
			Key:    values[0],
			Origin: EntityOriginDestination,
		}
	}
	return nil
}

//Render implementation of the ExpressionRenderer interface of EntitySpec
func (o *EntitySpec) Render() string {
	return o.Key + "@" + string(o.Origin.Key())
}

//Operator custom type for any kind of operator
type Operator string

//Capture captures the string representation of an operator from the given string. Interface of participle
func (o *Operator) Capture(values []string) error {
	*o = Operator(strings.ToUpper(values[0]))
	return nil
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

//ComparisionExpression representation of a comparision expression.
type ComparisionExpression struct {
	Entity   *EntitySpec `parser:"@Ident (@EntityOriginOperator @EntityOrigin)? "`
	Operator Operator    `parser:"@( \"EQUALS\" | \"NOT_EQUAL\" | \"CONTAINS\" | \"NOT_CONTAIN\" | \"STARTS_WITH\" | \"ENDS_WITH\" | \"NOT_STARTS_WITH\" | \"NOT_ENDS_WITH\" | \"GREATER_OR_EQUAL_THAN\" | \"LESS_OR_EQUAL_THAN\" | \"LESS_THAN\" | \"GREATER_THAN\" )"`
	Value    string      `parser:"@String"`
}

//Render implementation of ExpressionRenderer.Render
func (e *ComparisionExpression) Render() string {
	return fmt.Sprintf("%s %s '%s'", e.Entity.Render(), e.Operator, e.Value)
}

//UnaryOperationExpression representation of a unary expression representing a unary operator
type UnaryOperationExpression struct {
	Entity   *EntitySpec `parser:"@Ident (@EntityOriginOperator @EntityOrigin)? "`
	Operator Operator    `parser:"@( \"IS_EMPTY\" | \"IS_BLANK\"  | \"NOT_EMPTY\" | \"NOT_BLANK\" )"`
}

//Render implementation of ExpressionRenderer.Render
func (e *UnaryOperationExpression) Render() string {
	return fmt.Sprintf("%s %s", e.Entity.Render(), e.Operator)
}

var (
	filterLexer = lexer.Must(lexer.Regexp(`(\s+)` +
		`|(?P<Keyword>(?i)OR|AND|TRUE|FALSE|IS_EMPTY|NOT_EMPTY|IS_BLANK|NOT_BLANK|EQUALS|NOT_EQUAL|CONTAINS|NOT_CONTAIN|STARTS_WITH|ENDS_WITH|NOT_STARTS_WITH|NOT_ENDS_WITH|GREATER_OR_EQUAL_THAN|LESS_OR_EQUAL_THAN|LESS_THAN|GREATER_THAN)` +
		`|(?P<EntityOrigin>(?i)src|dest|na)` +
		`|(?P<EntityOriginOperator>(?i)@)` +
		`|(?P<Ident>[a-zA-Z_][\.a-zA-Z0-9_\-/]*)` +
		`|(?P<Number>[-+]?\d+(\.\d+)?)` +
		`|(?P<String>'[^']*'|"[^"]*")`,
	))
	filterParser = participle.MustBuild(
		&FilterExpression{},
		participle.Lexer(filterLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
		participle.UseLookahead(3),
	)
)

//Normalize parses the input and returns the normalized representation of the input string
func Normalize(input string) (string, error) {
	parser := NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return input, err
	}
	return expr.Render(), nil
}

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
