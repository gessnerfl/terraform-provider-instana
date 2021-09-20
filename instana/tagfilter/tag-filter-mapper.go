package tagfilter

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewMapper creates a new instance of the Mapper
func NewMapper() Mapper {
	return &tagFilterMapper{}
}

//Mapper interface of the tag filter expression mapper
type Mapper interface {
	FromAPIModel(input restapi.TagFilterExpressionElement) (*FilterExpression, error)
	ToAPIModel(input *FilterExpression) restapi.TagFilterExpressionElement
}

//struct for the filter expression mapper implementation for tag filter expressions
type tagFilterMapper struct{}
