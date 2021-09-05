package filterexpression

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewTagFilterMapper creates a new instance of the TagFilterMapper
func NewTagFilterMapper() TagFilterMapper {
	return &tagFilterMapperImpl{}
}

//TagFilterMapper interface of the tag filter expression mapper
type TagFilterMapper interface {
	FromAPIModel(input restapi.TagFilterExpressionElement) (*FilterExpression, error)
	ToAPIModel(input *FilterExpression) restapi.TagFilterExpressionElement
}

//struct for the filter expression mapper implementation for tag filter expressions
type tagFilterMapperImpl struct{}
