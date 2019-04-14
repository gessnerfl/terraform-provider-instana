package utils

import (
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func NewFilterExpressionMapper() FilterExpressionMapper {
	return new(filterExpressionMapperImpl)
}

type FilterExpressionMapper interface {
	FromApiModel(input restapi.MatchExpression) (*FilterExpression, error)
	ToApiModel(input *FilterExpression) (*restapi.MatchExpression, error)
}

type filterExpressionMapperImpl struct{}

func (i *filterExpressionMapperImpl) FromApiModel(input restapi.MatchExpression) (*FilterExpression, error) {
	if input.GetType() == restapi.BinaryOperatorExpressionType {

	} else if input.GetType() == restapi.LeafExpressionType {

	}
	return nil, fmt.Errorf("Unsupported match expression of type %s", input.GetType())
}

func (i *filterExpressionMapperImpl) ToApiModel(input *FilterExpression) (*restapi.MatchExpression, error) {
	return nil, nil
}
