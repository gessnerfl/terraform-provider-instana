package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
)

//NewApplicationConfigUnmarshaller creates a new Unmarshaller instance for application configs
func NewApplicationConfigUnmarshaller() JSONUnmarshaller {
	return &applicationConfigUnmarshaller{}
}

type applicationConfigUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *applicationConfigUnmarshaller) Unmarshal(data []byte) (interface{}, error) {
	var rawMatchSpecification json.RawMessage
	var rawTagFilterExpression json.RawMessage
	temp := &ApplicationConfig{
		MatchSpecification:  &rawMatchSpecification,
		TagFilterExpression: &rawTagFilterExpression,
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return &ApplicationConfig{}, err
	}
	matchSpecification, err := u.unmarshalMatchSpecification(rawMatchSpecification)
	if err != nil {
		return &ApplicationConfig{}, err
	}
	tagFilter, err := u.unmarshalTagFilterExpressionElement(rawTagFilterExpression)
	if err != nil {
		return &ApplicationConfig{}, err
	}
	return &ApplicationConfig{
		ID:                  temp.ID,
		Label:               temp.Label,
		MatchSpecification:  matchSpecification,
		TagFilterExpression: tagFilter,
		Scope:               temp.Scope,
		BoundaryScope:       temp.BoundaryScope,
	}, nil
}

func (u *applicationConfigUnmarshaller) unmarshalMatchSpecification(raw json.RawMessage) (MatchExpression, error) {
	if raw == nil {
		return nil, nil
	}
	temp := struct {
		Dtype MatchExpressionType `json:"type"`
	}{}

	if err := json.Unmarshal(raw, &temp); err != nil {
		return nil, err
	}

	if temp.Dtype == BinaryOperatorExpressionType {
		return u.unmarshalBinaryOperator(raw)
	} else if temp.Dtype == LeafExpressionType {
		return u.unmarshalTagMatcherExpression(raw)
	} else {
		return nil, errors.New("invalid expression type")
	}
}

func (u *applicationConfigUnmarshaller) unmarshalBinaryOperator(raw json.RawMessage) (BinaryOperator, error) {
	var leftRaw json.RawMessage
	var rightRaw json.RawMessage
	temp := BinaryOperator{
		Left:  &leftRaw,
		Right: &rightRaw,
	}

	json.Unmarshal(raw, &temp) //cannot fail as already successfully unmarshalled in unmarshalMatchSpecification
	left, err := u.unmarshalMatchSpecification(leftRaw)
	if err != nil {
		return BinaryOperator{}, err
	}

	right, err := u.unmarshalMatchSpecification(rightRaw)
	if err != nil {
		return BinaryOperator{}, err
	}
	return BinaryOperator{
		Dtype:       temp.Dtype,
		Left:        left,
		Right:       right,
		Conjunction: temp.Conjunction,
	}, nil
}

func (u *applicationConfigUnmarshaller) unmarshalTagMatcherExpression(raw json.RawMessage) (TagMatcherExpression, error) {
	data := TagMatcherExpression{}
	json.Unmarshal(raw, &data) //cannot fail as already successfully unmarshalled in unmarshalMatchSpecification
	return data, nil
}

func (u *applicationConfigUnmarshaller) unmarshalTagFilterExpressionElement(raw json.RawMessage) (TagFilterExpressionElement, error) {
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

func (u *applicationConfigUnmarshaller) unmarshalTagFilterExpression(raw json.RawMessage) (TagFilterExpressionElement, error) {
	temp := tempTagFilterExpression{}
	json.Unmarshal(raw, &temp) //cannot fail as already successfully unmarshalled in unmarshalTagFilterExpressionElement

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

func (u *applicationConfigUnmarshaller) unmarshalTagFilter(raw json.RawMessage) TagFilterExpressionElement {
	data := TagFilter{}
	json.Unmarshal(raw, &data) //cannot fail as already successfully unmarshalled in unmarshalTagFilterExpressionElement
	return &data
}

type tempTagFilterExpression struct {
	Elements        []json.RawMessage              `json:"elements"`
	LogicalOperator LogicalOperatorType            `json:"logicalOperator"`
	Type            TagFilterExpressionElementType `json:"type"`
}
