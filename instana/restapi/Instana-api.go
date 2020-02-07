package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	//InstanaAPIBasePath path to Instana RESTful API
	InstanaAPIBasePath = "/api"
	//EventsBasePath path to Events resource of Instana RESTful API
	EventsBasePath = InstanaAPIBasePath + "/events"
	//settingsPathElement path element to settings
	settingsPathElement = "/settings"
	//EventSettingsBasePath path to Event Settings resource of Instana RESTful API
	EventSettingsBasePath = EventsBasePath + settingsPathElement
	//SettingsBasePath path to Event Settings resource of Instana RESTful API
	SettingsBasePath = InstanaAPIBasePath + settingsPathElement
)

//InstanaAPI is the interface to all resources of the Instana Rest API
type InstanaAPI interface {
	CustomEventSpecifications() RestResource
	UserRoles() RestResource
	ApplicationConfigs() RestResource
	AlertingChannels() RestResource
}

//NewInstanaAPI creates a new instance of the instana API
func NewInstanaAPI(apiToken string, endpoint string) InstanaAPI {
	client := NewClient(apiToken, endpoint)
	return &baseInstanaAPI{client: client}
}

type baseInstanaAPI struct {
	client RestClient
}

//CustomEventSpecifications implementation of InstanaAPI interface
func (api *baseInstanaAPI) CustomEventSpecifications() RestResource {
	return NewRestResource(CustomEventSpecificationResourcePath, NewCustomEventSpecificationUnmarshaller(), api.client)
}

//NewCustomEventSpecificationUnmarshaller creates a new instance of Unmarshaller for custom event specifications
func NewCustomEventSpecificationUnmarshaller() Unmarshaller {
	return &customEventSpecificationUnmarshaller{}
}

type customEventSpecificationUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *customEventSpecificationUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	customEventSpecification := CustomEventSpecification{}
	if err := json.Unmarshal(data, &customEventSpecification); err != nil {
		return customEventSpecification, fmt.Errorf("failed to parse json; %s", err)
	}
	return customEventSpecification, nil
}

//UserRoles implementation of InstanaAPI interface
func (api *baseInstanaAPI) UserRoles() RestResource {
	return NewRestResource(UserRolesResourcePath, NewUserRoleUnmarshaller(), api.client)
}

//NewUserRoleUnmarshaller creates a new Unmarshaller instance for user roles
func NewUserRoleUnmarshaller() Unmarshaller {
	return &userRoleUnmarshaller{}
}

type userRoleUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *userRoleUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	userRole := UserRole{}
	if err := json.Unmarshal(data, &userRole); err != nil {
		return userRole, fmt.Errorf("failed to parse json; %s", err)
	}
	return userRole, nil
}

//ApplicationConfigs implementation of InstanaAPI interface
func (api *baseInstanaAPI) ApplicationConfigs() RestResource {
	return NewRestResource(ApplicationConfigsResourcePath, NewApplicationConfigUnmarshaller(), api.client)
}

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

	if err := json.Unmarshal(raw, &temp); err != nil {
		return BinaryOperator{}, err
	}

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
	if err := json.Unmarshal(raw, &data); err != nil {
		return TagMatcherExpression{}, err
	}
	return data, nil
}

//AlertingChannels implementation of InstanaAPI interface
func (api *baseInstanaAPI) AlertingChannels() RestResource {
	return NewRestResource(AlertingChannelsResourcePath, NewAlertingChannelUnmarshaller(), api.client)
}

//NewAlertingChannelUnmarshaller creates a new Unmarshaller instance for AlertingChannels
func NewAlertingChannelUnmarshaller() Unmarshaller {
	return &alertingChannelUnmarshaller{}
}

type alertingChannelUnmarshaller struct{}

//Unmarshal Unmarshaller interface implementation
func (u *alertingChannelUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	alertingChannel := AlertingChannel{}
	if err := json.Unmarshal(data, &alertingChannel); err != nil {
		return alertingChannel, fmt.Errorf("failed to parse json; %s", err)
	}
	return alertingChannel, nil
}
