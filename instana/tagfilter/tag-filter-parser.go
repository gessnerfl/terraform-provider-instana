package tagfilter

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//ExpressionRenderer interface definition to render an expression in its normalized form
type ExpressionRenderer interface {
	Render() string
}

//EntityOrigin custom type for the origin (source or destination) of a entity spec
type EntityOrigin interface {
	//Key returns the key of the entity origin
	Key() string
	//TagFilterEntity returns the Instana API Ta Filter Entity
	TagFilterEntity() restapi.TagFilterEntity
}

func newEntityOrigin(key string, tagFilterEntity restapi.TagFilterEntity) EntityOrigin {
	return &baseEntityOrigin{key: key, tagFilterEntity: tagFilterEntity}
}

type baseEntityOrigin struct {
	key             string
	tagFilterEntity restapi.TagFilterEntity
}

//Key interface implementation of EntityOrigin
func (o *baseEntityOrigin) Key() string {
	return o.key
}

//TagFilterEntity interface implementation of EntityOrigin
func (o *baseEntityOrigin) TagFilterEntity() restapi.TagFilterEntity {
	return o.tagFilterEntity
}

var (
	//EntityOriginSource constant value for the EntityOrigin source
	EntityOriginSource = newEntityOrigin("src", restapi.TagFilterEntitySource)
	//EntityOriginDestination constant value for the EntityOrigin destination
	EntityOriginDestination = newEntityOrigin("dest", restapi.TagFilterEntityDestination)
	//EntityOriginNotApplicable constant value when no EntityOrigin is applicable
	EntityOriginNotApplicable = newEntityOrigin("na", restapi.TagFilterEntityNotApplicable)
)

//EntityOrigins custom type for a slice of entity origins
type EntityOrigins []EntityOrigin

//ForInstanaAPIEntity returns the EntityOrigin for its corresponding TagFilterEntity from the Instana API
func (origins EntityOrigins) ForInstanaAPIEntity(input restapi.TagFilterEntity) EntityOrigin {
	for _, o := range origins {
		if o.TagFilterEntity() == input {
			return o
		}
	}
	log.Printf("tag filter entity %s is not supported; fall back to default origin %s", input, EntityOriginDestination.Key())
	return EntityOriginDestination
}

//ForKey returns the EntityOrigin for its string representation
func (origins EntityOrigins) ForKey(input string) EntityOrigin {
	for _, o := range origins {
		if o.Key() == input {
			return o
		}
	}
	log.Printf("entity origin with key %s is not supported; fall back to default origin %s", input, EntityOriginDestination.Key())
	return EntityOriginDestination
}

//SupportedEntityOrigins slice of supported EntityOrigins
var SupportedEntityOrigins = EntityOrigins{EntityOriginSource, EntityOriginDestination, EntityOriginNotApplicable}

//EntitySpec custom type for an entity path specification            Ident (@TagKeySeparator @Ident)? (@EntityOriginOperator @EntityOrigin)? "
type EntitySpec struct {
	Identifier string  `parser:"@Ident"`
	TagKey     *string `parser:"( \":\" @Ident )?"`
	Origin     *string `parser:"( \"@\" @EntityOrigin )?"`
}

//Render implementation of the ExpressionRenderer interface
func (o *EntitySpec) Render() string {
	origin := EntityOriginDestination.Key()
	tagKey := ""
	if o.TagKey != nil {
		tagKey = ":" + *o.TagKey
	}
	if o.Origin != nil {
		origin = SupportedEntityOrigins.ForKey(*o.Origin).Key()
	}
	return o.Identifier + tagKey + "@" + origin
}

//Operator custom type for an operator
type Operator string

//Capture captures the string representation of an operator from the given string and converts it to upper case. Interface of participle
func (o *Operator) Capture(values []string) error {
	*o = Operator(strings.ToUpper(values[0]))
	return nil
}

//FilterExpression representation of a tag filter expression
type FilterExpression struct {
	Expression *LogicalOrExpression `parser:"@@"`
}

//Render implementation of ExpressionRenderer.Render
func (e *FilterExpression) Render() string {
	return e.Expression.Render()
}

//LogicalOrExpression representation of a logical OR, or as a wrapper for a LogicalAndExpression or a PrimaryExpression. The wrapping is required to handle precedence.
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

//LogicalAndExpression representation of a logical AND, or as a wrapper for a PrimaryExpression. The wrapping is required to handle precedence.
type LogicalAndExpression struct {
	Left     *BracketExpression    `parser:"  @@"`
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

//BracketExpression representation of a bracket expression
type BracketExpression struct {
	Bracket *LogicalOrExpression `parser:"\"(\" @@ \")\""`
	Primary *PrimaryExpression   `parser:"| @@"`
}

//Render implementation of ExpressionRenderer.Render
func (e *BracketExpression) Render() string {
	if e.Bracket != nil {
		return "(" + e.Bracket.Render() + ")"
	}
	return e.Primary.Render()
}

//PrimaryExpression wrapper for either a comparison or a unary expression
type PrimaryExpression struct {
	Comparison     *ComparisonExpression     `parser:"  @@"`
	UnaryOperation *UnaryOperationExpression `parser:"| @@"`
}

//Render implementation of ExpressionRenderer.Render
func (e *PrimaryExpression) Render() string {
	if e.Comparison != nil {
		return e.Comparison.Render()
	}
	return e.UnaryOperation.Render()
}

//ComparisonExpression representation of a comparison expression.
type ComparisonExpression struct {
	Entity       *EntitySpec `parser:"@@"`
	Operator     Operator    `parser:"@( \"EQUALS\" | \"NOT_EQUAL\" | \"CONTAINS\" | \"NOT_CONTAIN\" | \"STARTS_WITH\" | \"ENDS_WITH\" | \"NOT_STARTS_WITH\" | \"NOT_ENDS_WITH\" | \"GREATER_OR_EQUAL_THAN\" | \"LESS_OR_EQUAL_THAN\" | \"LESS_THAN\" | \"GREATER_THAN\" )"`
	NumberValue  *int64      `parser:"( @Number"`
	BooleanValue *bool       `parser:"| @( \"FALSE\" | \"TRUE\" )"`
	StringValue  *string     `parser:"| @String )"`
}

//Render implementation of ExpressionRenderer.Render
func (e *ComparisonExpression) Render() string {
	if e.NumberValue != nil {
		return fmt.Sprintf("%s %s %d", e.Entity.Render(), e.Operator, *e.NumberValue)
	} else if e.BooleanValue != nil {
		return fmt.Sprintf("%s %s %t", e.Entity.Render(), e.Operator, *e.BooleanValue)
	}
	return fmt.Sprintf("%s %s '%s'", e.Entity.Render(), e.Operator, *e.StringValue)
}

//UnaryOperationExpression representation of a unary expression
type UnaryOperationExpression struct {
	Entity   *EntitySpec `parser:"@@"`
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
		`|(?P<Bracket>[\(\)])` +
		`|(?P<TagKeySeparator>(?i):)` +
		`|(?P<Ident>[a-zA-Z_][\.a-zA-Z0-9_\-/]*)` +
		`|(?P<Number>[-+]?\d+)` +
		`|(?P<String>'[^']*'|"[^"]*")`,
	))
	filterParser = participle.MustBuild(
		&FilterExpression{},
		participle.Lexer(filterLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
		participle.UseLookahead(5),
	)
)

//Normalize parses the input and returns the normalized representation of the input string
func Normalize(input string) (string, error) {
	parser := NewParser()
	mapper := NewMapper()

	parsed, err := parser.Parse(input)
	if err != nil {
		return input, err
	}

	apiModel := mapper.ToAPIModel(parsed)
	mapped, err := mapper.FromAPIModel(apiModel)
	if err != nil {
		return input, err
	}

	return mapped.Render(), nil
}

//NewParser creates a new instance of a Parser
func NewParser() Parser {
	return new(parserImpl)
}

//Parser interface for working with tag filter expressions of instana
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
