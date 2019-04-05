package resources_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/resources"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
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
		t.Fatalf("Expected no error but got %s", err)
	}

	if !cmp.Equal(applicationConfig, data) {
		t.Fatalf("Expected json to be unmarshalled to %v but got %v; diff %s", applicationConfig, data, cmp.Diff(applicationConfig, data))
	}
}

func TestFailedGetOneApplicationConfigBecauseOfErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfigID := "test-application-config-id"

	client.EXPECT().GetOne(gomock.Eq(applicationConfigID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(applicationConfigID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetOneApplicationConfigBecauseOfInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfigID := "test-application-config-id"

	client.EXPECT().GetOne(gomock.Eq(applicationConfigID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(applicationConfigID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetOneApplicationConfigBecauseOfInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfigID := "test-application-config-id"

	client.EXPECT().GetOne(gomock.Eq(applicationConfigID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(applicationConfigID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetOneApplicationConfigBecauseOfNoJsonAsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfigID := "test-application-config-id"

	client.EXPECT().GetOne(gomock.Eq(applicationConfigID), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(applicationConfigID)

	if err == nil {
		t.Fatalf("Expected to get error")
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
		t.Fatalf("Expected no error but got %s", err)
	}

	if !cmp.Equal(applicationConfig, result) {
		t.Fatalf("Expected json to be unmarshalled to %v but got %v; diff %s", applicationConfig, result, cmp.Diff(applicationConfig, result))
	}
}

func TestFailedUpsertOfApplicationConfigBecauseOfInvalidApplicationConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := restapi.ApplicationConfig{
		Label:              "Label",
		MatchSpecification: restapi.NewTagMatcherExpression("key", "EQUAL", "value"),
		Scope:              "scope",
	}

	client.EXPECT().Put(gomock.Eq(applicationConfig), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Times(0)

	_, err := sut.Upsert(applicationConfig)

	if err == nil {
		t.Fatal("Expected to get error")
	}
}

func TestFailedUpsertOfApplicationConfigBecauseOfInvalidResponseMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := makeTestApplicationConfig()

	client.EXPECT().Put(gomock.Eq(applicationConfig), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return([]byte("invalid response"), nil)

	_, err := sut.Upsert(applicationConfig)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedUpsertOfApplicationConfigBecauseOfInvalidApplicationConfigInResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := makeTestApplicationConfig()

	client.EXPECT().Put(gomock.Eq(applicationConfig), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return([]byte("{ \"invalid\" : \"application config\" }"), nil)

	_, err := sut.Upsert(applicationConfig)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedUpsertOfApplicationConfigBecauseOfClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewApplicationConfigResource(client)
	applicationConfig := makeTestApplicationConfig()

	client.EXPECT().Put(gomock.Eq(applicationConfig), gomock.Eq(restapi.ApplicationConfigsResourcePath)).Return(nil, errors.New("Error during test"))

	_, err := sut.Upsert(applicationConfig)

	if err == nil {
		t.Fatal("Expected to get error")
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
		t.Fatal("Expected to get error")
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
		MatchSpecification: restapi.NewTagMatcherExpression("key", "EQUAL", "value"),
		Scope:              "scope",
	}
}
