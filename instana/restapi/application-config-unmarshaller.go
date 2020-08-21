package restapi

import (
	"encoding/json"
	"errors"
)

//NewApplicationConfigUnmarshaller creates a new Unmarshaller instance for application configs
func NewApplicationConfigUnmarshaller() Unmarshaller {
	return &applicationConfigUnmarshaller{}
}

type applicationConfigUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *applicationConfigUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	var matchExpression json.RawMessage
	temp := ApplicationConfig{
		MatchSpecification: &matchExpression,
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return ApplicationConfig{}, err
	}
	matchSpecification, err := u.unmarshalMatchSpecification(matchExpression)
	if err != nil {
		return ApplicationConfig{}, err
	}
	return ApplicationConfig{
		ID:                 temp.ID,
		Label:              temp.Label,
		MatchSpecification: matchSpecification,
		Scope:              temp.Scope,
		BoundaryScope:      temp.BoundaryScope,
	}, nil
}

func (u *applicationConfigUnmarshaller) unmarshalMatchSpecification(raw json.RawMessage) (MatchExpression, error) {
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
