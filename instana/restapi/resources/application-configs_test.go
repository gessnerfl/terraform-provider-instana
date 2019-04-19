package resources_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/resources"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestSuccessfulGetOneApplicationConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := makeTestApplicationConfig()
	serializedJSON, _ := json.Marshal(applicationConfig)

	client.EXPECT().GetOne(gomock.Eq(applicationConfig.ID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(serializedJSON, nil)

	data, err := sut.GetOne(applicationConfig.ID)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(applicationConfig, data) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, applicationConfig, data, cmp.Diff(applicationConfig, data))
	}
}

func TestFailedGetOneApplicationConfigWhenRestClientReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfigID := "test-application-config-id"

	client.EXPECT().GetOne(gomock.Eq(applicationConfigID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(applicationConfigID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneApplicationConfigWhenResponseContainsInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfigID := "test-application-config-id"

	client.EXPECT().GetOne(gomock.Eq(applicationConfigID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(applicationConfigID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneApplicationConfigWhenResponseContainsInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfigID := "test-application-config-id"

	client.EXPECT().GetOne(gomock.Eq(applicationConfigID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(applicationConfigID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneApplicationConfigWhenResponseContainsNoJsonAsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfigID := "test-application-config-id"

	client.EXPECT().GetOne(gomock.Eq(applicationConfigID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(applicationConfigID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneApplicationConfigWhenExpressionTypeIsNotSupported(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	//config is invalid because there is no DType for the match specification.
	applicationConfig := restapi.ApplicationConfig{
		ID:    "id",
		Label: "label",
		MatchSpecification: restapi.TagMatcherExpression{
			Key:      "foo",
			Operator: restapi.NotEmptyOperator,
		},
		Scope: "scope",
	}
	serializedJSON, _ := json.Marshal(applicationConfig)

	client.EXPECT().GetOne(gomock.Eq(applicationConfig.ID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(serializedJSON, nil)

	_, err := sut.GetOne(applicationConfig.GetID())

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneApplicationConfigWhenLeftSideOfBinaryExpressionTypeIsNotValid(t *testing.T) {
	left := restapi.TagMatcherExpression{
		Key:      "foo",
		Operator: restapi.NotEmptyOperator,
	}
	right := restapi.NewUnaryOperationExpression("foo", restapi.IsEmptyOperator)
	testFailGetOneApplicationConfigWhenOneSideOfBinaryExpressionIsNotValue(left, right, t)
}

func TestFailedGetOneApplicationConfigWhenRightSideOfBinaryExpressionTypeIsNotValid(t *testing.T) {
	left := restapi.NewUnaryOperationExpression("foo", restapi.IsEmptyOperator)
	right := restapi.TagMatcherExpression{
		Key:      "foo",
		Operator: restapi.NotEmptyOperator,
	}
	testFailGetOneApplicationConfigWhenOneSideOfBinaryExpressionIsNotValue(left, right, t)
}

func testFailGetOneApplicationConfigWhenOneSideOfBinaryExpressionIsNotValue(left restapi.MatchExpression, right restapi.MatchExpression, t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := restapi.ApplicationConfig{
		ID:                 "id",
		Label:              "label",
		MatchSpecification: restapi.NewBinaryOperator(left, restapi.LogicalOr, right),
		Scope:              "scope",
	}
	serializedJSON, _ := json.Marshal(applicationConfig)

	client.EXPECT().GetOne(gomock.Eq(applicationConfig.ID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(serializedJSON, nil)

	_, err := sut.GetOne(applicationConfig.GetID())

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestSuccessfulUpsertOfApplicationConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := makeTestApplicationConfig()
	serializedJSON, _ := json.Marshal(applicationConfig)

	client.EXPECT().Put(gomock.Eq(applicationConfig), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(serializedJSON, nil)

	result, err := sut.Upsert(applicationConfig)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(applicationConfig, result) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, applicationConfig, result, cmp.Diff(applicationConfig, result))
	}
}

func TestSuccessfulUpsertOfComplexApplicationConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)

	applicationConfig := restapi.ApplicationConfig{
		ID:    "id",
		Label: "label",
		MatchSpecification: restapi.NewBinaryOperator(
			restapi.NewBinaryOperator(
				restapi.NewComparisionExpression("key1", restapi.EqualsOperator, "value1"),
				restapi.LogicalOr,
				restapi.NewUnaryOperationExpression("key2", restapi.NotEmptyOperator),
			),
			restapi.LogicalAnd,
			restapi.NewComparisionExpression("key3", restapi.NotEqualOperator, "value3"),
		),
		Scope: "scope",
	}
	serializedJSON, _ := json.Marshal(applicationConfig)

	client.EXPECT().Put(gomock.Eq(applicationConfig), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(serializedJSON, nil)

	result, err := sut.Upsert(applicationConfig)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(applicationConfig, result) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, applicationConfig, result, cmp.Diff(applicationConfig, result))
	}
}

func TestFailedUpsertOfApplicationConfigWhenApplicationConfigIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := restapi.ApplicationConfig{
		Label:              "Label",
		MatchSpecification: restapi.NewComparisionExpression("key", "EQUAL", "value"),
		Scope:              "scope",
	}

	client.EXPECT().Put(gomock.Eq(applicationConfig), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Times(0)

	_, err := sut.Upsert(applicationConfig)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfApplicationConfigWhenResponseMessageIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := makeTestApplicationConfig()

	client.EXPECT().Put(gomock.Eq(applicationConfig), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return([]byte("invalid response"), nil)

	_, err := sut.Upsert(applicationConfig)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfApplicationConfigWhenApplicationConfigInResponseIsInvalid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := makeTestApplicationConfig()

	client.EXPECT().Put(gomock.Eq(applicationConfig), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return([]byte("{ \"invalid\" : \"application config\" }"), nil)

	_, err := sut.Upsert(applicationConfig)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfApplicationConfigWhenClientReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := makeTestApplicationConfig()

	client.EXPECT().Put(gomock.Eq(applicationConfig), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(nil, errors.New("Error during test"))

	_, err := sut.Upsert(applicationConfig)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestSuccessfulDeleteOfApplicationConfigByApplicationConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := makeTestApplicationConfig()

	client.EXPECT().Delete(gomock.Eq("test-application-config-id-1"), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(nil)

	err := sut.Delete(applicationConfig)

	if err != nil {
		t.Fatalf("Expected no error got %s", err)
	}
}

func TestFailedDeleteOfApplicationConfigByApplicationConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := makeTestApplicationConfig()

	client.EXPECT().Delete(gomock.Eq("test-application-config-id-1"), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(applicationConfig)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func makeTestApplicationConfig() restapi.ApplicationConfig {
	return makeTestApplicationConfigWithCounter(1)
}

func makeTestApplicationConfigWithCounter(counter int) restapi.ApplicationConfig {
	id := fmt.Sprintf("test-application-config-id-%d", counter)
	label := fmt.Sprintf("Test Application Config Label %d", counter)
	return restapi.ApplicationConfig{
		ID:                 id,
		Label:              label,
		MatchSpecification: restapi.NewComparisionExpression("key", restapi.EqualsOperator, "value"),
		Scope:              "scope",
	}
}
