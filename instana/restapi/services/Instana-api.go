package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

//NewInstanaAPI creates a new instance of the instana API
func NewInstanaAPI(apiToken string, endpoint string) restapi.InstanaAPI {
	client := NewClient(apiToken, endpoint)
	return &baseInstanaAPI{client: client}
}

type baseInstanaAPI struct {
	client restapi.RestClient
}

//CustomEventSpecifications implementation of InstanaAPI interface
func (api *baseInstanaAPI) CustomEventSpecifications() restapi.RestResource {
	return NewRestResource(restapi.CustomEventSpecificationResourcePath, NewCustomEventSpecificationUnmarshaller(), api.client)
}

//NewCustomEventSpecificationUnmarshaller creates a new instance of Unmarshaller for custom event specifications
func NewCustomEventSpecificationUnmarshaller() Unmarshaller {
	return &customEventSpecificationUnmarshaller{}
}

type customEventSpecificationUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *customEventSpecificationUnmarshaller) Unmarshal(data []byte) (restapi.InstanaDataObject, error) {
	customEventSpecification := restapi.CustomEventSpecification{}
	if err := json.Unmarshal(data, &customEventSpecification); err != nil {
		return customEventSpecification, fmt.Errorf("failed to parse json; %s", err)
	}
	return customEventSpecification, nil
}

//UserRoles implementation of InstanaAPI interface
func (api *baseInstanaAPI) UserRoles() restapi.RestResource {
	return NewRestResource(restapi.UserRolesResourcePath, NewUserRoleUnmarshaller(), api.client)
}

//NewUserRoleUnmarshaller creates a new Unmarshaller instance for user roles
func NewUserRoleUnmarshaller() Unmarshaller {
	return &userRoleUnmarshaller{}
}

type userRoleUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *userRoleUnmarshaller) Unmarshal(data []byte) (restapi.InstanaDataObject, error) {
	userRole := restapi.UserRole{}
	if err := json.Unmarshal(data, &userRole); err != nil {
		return userRole, fmt.Errorf("failed to parse json; %s", err)
	}
	return userRole, nil
}

//ApplicationConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) ApplicationConfigs() restapi.RestResource {
	return NewRestResource(restapi.ApplicationConfigsResourcePath, NewApplicationConfigUnmarshaller(), api.client)
}

//NewApplicationConfigUnmarshaller creates a new Unmarshaller instance for application configs
func NewApplicationConfigUnmarshaller() Unmarshaller {
	return &applicationConfigUnmarshaller{}
}

type applicationConfigUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *applicationConfigUnmarshaller) Unmarshal(data []byte) (restapi.InstanaDataObject, error) {
	var matchExpression json.RawMessage
	temp := restapi.ApplicationConfig{
		MatchSpecification: &matchExpression,
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return restapi.ApplicationConfig{}, err
	}
	matchSpecification, err := u.unmarshalMatchSpecification(matchExpression)
	if err != nil {
		return restapi.ApplicationConfig{}, err
	}
	return restapi.ApplicationConfig{
		ID:                 temp.ID,
		Label:              temp.Label,
		MatchSpecification: matchSpecification,
		Scope:              temp.Scope,
	}, nil
}

func (u *applicationConfigUnmarshaller) unmarshalMatchSpecification(raw json.RawMessage) (restapi.MatchExpression, error) {
	temp := struct {
		Dtype restapi.MatchExpressionType `json:"type"`
	}{}

	if err := json.Unmarshal(raw, &temp); err != nil {
		return nil, err
	}

	if temp.Dtype == restapi.BinaryOperatorExpressionType {
		return u.unmarshalBinaryOperator(raw)
	} else if temp.Dtype == restapi.LeafExpressionType {
		return u.unmarshalTagMatcherExpression(raw)
	} else {
		return nil, errors.New("invalid expression type")
	}
}

func (u *applicationConfigUnmarshaller) unmarshalBinaryOperator(raw json.RawMessage) (restapi.BinaryOperator, error) {
	var leftRaw json.RawMessage
	var rightRaw json.RawMessage
	temp := restapi.BinaryOperator{
		Left:  &leftRaw,
		Right: &rightRaw,
	}

	if err := json.Unmarshal(raw, &temp); err != nil {
		return restapi.BinaryOperator{}, err
	}

	left, err := u.unmarshalMatchSpecification(leftRaw)
	if err != nil {
		return restapi.BinaryOperator{}, err
	}

	right, err := u.unmarshalMatchSpecification(rightRaw)
	if err != nil {
		return restapi.BinaryOperator{}, err
	}
	return restapi.BinaryOperator{
		Dtype:       temp.Dtype,
		Left:        left,
		Right:       right,
		Conjunction: temp.Conjunction,
	}, nil
}

func (u *applicationConfigUnmarshaller) unmarshalTagMatcherExpression(raw json.RawMessage) (restapi.TagMatcherExpression, error) {
	data := restapi.TagMatcherExpression{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return restapi.TagMatcherExpression{}, err
	}
	return data, nil
}

//AlertingChannels implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingChannels() restapi.RestResource {
	return NewRestResource(restapi.AlertingChannelsResourcePath, NewAlertingChannelUnmarshaller(), api.client)
}

//NewAlertingChannelUnmarshaller creates a new Unmarshaller instance for AlertingChannels
func NewAlertingChannelUnmarshaller() Unmarshaller {
	return &alertingChannelUnmarshaller{}
}

type alertingChannelUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *alertingChannelUnmarshaller) Unmarshal(data []byte) (restapi.InstanaDataObject, error) {
	alertingChannel := restapi.AlertingChannel{}
	if err := json.Unmarshal(data, &alertingChannel); err != nil {
		return alertingChannel, fmt.Errorf("failed to parse json; %s", err)
	}
	return alertingChannel, nil
}
