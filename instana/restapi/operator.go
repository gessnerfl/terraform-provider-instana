package restapi

//LogicalOperatorType custom type for logical operators
type LogicalOperatorType string

//LogicalOperatorTypes custom type for slice of logical operators
type LogicalOperatorTypes []LogicalOperatorType

//IsSupported check if the provided logical operator is supported
func (operators LogicalOperatorTypes) IsSupported(o LogicalOperatorType) bool {
	for _, v := range operators {
		if v == o {
			return true
		}
	}
	return false
}

const (
	//LogicalAnd constant for logical AND conjunction
	LogicalAnd = LogicalOperatorType("AND")
	//LogicalOr constant for logical OR conjunction
	LogicalOr = LogicalOperatorType("OR")
)

//SupportedLogicalOperatorTypes list of supported logical operators of Instana API
var SupportedLogicalOperatorTypes = LogicalOperatorTypes{LogicalAnd, LogicalOr}

//ExpressionOperator custom type for tag matcher operators
type ExpressionOperator string

//ExpressionOperators custom type representing a slice of ExpressionOperator
type ExpressionOperators []ExpressionOperator

//IsSupported check if the provided tag filter operator is supported
func (operators ExpressionOperators) IsSupported(o ExpressionOperator) bool {
	for _, v := range operators {
		if v == o {
			return true
		}
	}
	return false
}

const (
	//EqualsOperator constant for the EQUALS operator
	EqualsOperator = ExpressionOperator("EQUALS")
	//NotEqualOperator constant for the NOT_EQUAL operator
	NotEqualOperator = ExpressionOperator("NOT_EQUAL")
	//ContainsOperator constant for the CONTAINS operator
	ContainsOperator = ExpressionOperator("CONTAINS")
	//NotContainOperator constant for the NOT_CONTAIN operator
	NotContainOperator = ExpressionOperator("NOT_CONTAIN")

	//IsEmptyOperator constant for the IS_EMPTY operator
	IsEmptyOperator = ExpressionOperator("IS_EMPTY")
	//NotEmptyOperator constant for the NOT_EMPTY operator
	NotEmptyOperator = ExpressionOperator("NOT_EMPTY")
	//IsBlankOperator constant for the IS_BLANK operator
	IsBlankOperator = ExpressionOperator("IS_BLANK")
	//NotBlankOperator constant for the NOT_BLANK operator
	NotBlankOperator = ExpressionOperator("NOT_BLANK")

	//StartsWithOperator constant for the STARTS_WITH operator
	StartsWithOperator = ExpressionOperator("STARTS_WITH")
	//EndsWithOperator constant for the ENDS_WITH operator
	EndsWithOperator = ExpressionOperator("ENDS_WITH")
	//NotStartsWithOperator constant for the NOT_STARTS_WITH operator
	NotStartsWithOperator = ExpressionOperator("NOT_STARTS_WITH")
	//NotEndsWithOperator constant for the NOT_ENDS_WITH operator
	NotEndsWithOperator = ExpressionOperator("NOT_ENDS_WITH")
	//GreaterOrEqualThanOperator constant for the GREATER_OR_EQUAL_THAN operator
	GreaterOrEqualThanOperator = ExpressionOperator("GREATER_OR_EQUAL_THAN")
	//LessOrEqualThanOperator constant for the LESS_OR_EQUAL_THAN operator
	LessOrEqualThanOperator = ExpressionOperator("LESS_OR_EQUAL_THAN")
	//GreaterThanOperator constant for the GREATER_THAN operator
	GreaterThanOperator = ExpressionOperator("GREATER_THAN")
	//LessThanOperator constant for the LESS_THAN operator
	LessThanOperator = ExpressionOperator("LESS_THAN")
)

//SupportedComparisonOperators list of supported comparison operators of Instana API
var SupportedComparisonOperators = ExpressionOperators{
	EqualsOperator,
	NotEqualOperator,
	ContainsOperator,
	NotContainOperator,
	StartsWithOperator,
	EndsWithOperator,
	NotStartsWithOperator,
	NotEndsWithOperator,
	GreaterOrEqualThanOperator,
	LessOrEqualThanOperator,
	GreaterThanOperator,
	LessThanOperator,
}

//SupportedUnaryExpressionOperators list of supported unary expression operators of Instana API
var SupportedUnaryExpressionOperators = ExpressionOperators{
	IsEmptyOperator,
	NotEmptyOperator,
	IsBlankOperator,
	NotBlankOperator,
}
