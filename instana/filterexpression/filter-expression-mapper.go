package filterexpression

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewMapper creates a new instance of the Mapper
func NewMapper() Mapper {
	return new(mapperImpl)
}

//Mapper interface of the filter expression mapper
type Mapper interface {
	FromAPIModel(input restapi.MatchExpression) (*FilterExpression, error)
	ToAPIModel(input *FilterExpression) (*restapi.MatchExpression, error)
}

//struct for the filter expression mapper implementation
type mapperImpl struct{}

//ToAPIModel Implementation of the mapping form filter expression model to the Instana API model
func (m *mapperImpl) ToAPIModel(input *FilterExpression) (*restapi.MatchExpression, error) {
	return nil, nil
}
