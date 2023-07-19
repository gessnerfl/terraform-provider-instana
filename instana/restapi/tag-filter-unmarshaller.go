package restapi

import (
	"encoding/json"
	"fmt"
)

// TagFilterUnmarshaller interface for the unmarshaller for TagFilterExpressions
type TagFilterUnmarshaller interface {
	Unmarshal(raw json.RawMessage) (TagFilterExpressionElement, error)
}

// NewTagFilterUnmarshaller creates a new instance of TagFilterUnmarshaller
func NewTagFilterUnmarshaller() TagFilterUnmarshaller {
	return &tagFilterUnmarshaller{}
}

type tagFilterUnmarshaller struct{}

func (u *tagFilterUnmarshaller) Unmarshal(raw json.RawMessage) (TagFilterExpressionElement, error) {
	return u.unmarshalTagFilterExpressionElement(raw)
}

func (u *tagFilterUnmarshaller) unmarshalTagFilterExpressionElement(raw json.RawMessage) (TagFilterExpressionElement, error) {
	if raw == nil {
		return nil, nil
	}
	temp := struct {
		Type TagFilterExpressionElementType `json:"type"`
	}{}

	if err := json.Unmarshal(raw, &temp); err != nil {
		return nil, err
	}

	if temp.Type == TagFilterExpressionType {
		return u.unmarshalTagFilterExpression(raw)
	} else if temp.Type == TagFilterType {
		return u.unmarshalTagFilter(raw), nil
	} else {
		return nil, fmt.Errorf("invalid tag filter element type %s", temp.Type)
	}
}

func (u *tagFilterUnmarshaller) unmarshalTagFilterExpression(raw json.RawMessage) (TagFilterExpressionElement, error) {
	temp := tempTagFilterExpression{}
	json.Unmarshal(raw, &temp) // NOSONAR: cannot fail as already successfully unmarshalled in unmarshalTagFilterExpressionElement

	elements := make([]TagFilterExpressionElement, len(temp.Elements))
	for i, e := range temp.Elements {
		element, err := u.unmarshalTagFilterExpressionElement(e)
		if err != nil {
			return nil, err
		}
		elements[i] = element
	}
	return &TagFilterExpression{
		Type:            temp.Type,
		LogicalOperator: temp.LogicalOperator,
		Elements:        elements,
	}, nil
}

func (u *tagFilterUnmarshaller) unmarshalTagFilter(raw json.RawMessage) TagFilterExpressionElement {
	data := TagFilter{}
	json.Unmarshal(raw, &data) // NOSONAR: cannot fail as already successfully unmarshalled in unmarshalTagFilterExpressionElement
	return &data
}

type tempTagFilterExpression struct {
	Elements        []json.RawMessage              `json:"elements"`
	LogicalOperator LogicalOperatorType            `json:"logicalOperator"`
	Type            TagFilterExpressionElementType `json:"type"`
}
