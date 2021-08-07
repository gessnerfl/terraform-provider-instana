package filterexpression

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewMatchExpressionMapper creates a new instance of the MatchExpressionMapper
func NewMatchExpressionMapper() MatchExpressionMapper {
	return new(matchExpressionMapperImpl)
}

//MatchExpressionMapper interface of the filter expression mapper
type MatchExpressionMapper interface {
	FromAPIModel(input restapi.MatchExpression) (*FilterExpression, error)
	ToAPIModel(input *FilterExpression) restapi.MatchExpression
}

//struct for the filter expression mapper implementation for match expressions
type matchExpressionMapperImpl struct{}
