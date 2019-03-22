package resources_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/resources"
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestSuccessfulGetOneUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRole := makeTestUserRole()
	serializedJSON, _ := json.Marshal(userRole)

	client.EXPECT().GetOne(gomock.Eq(userRole.ID), gomock.Eq(restapi.UserRolesResourcePath)).Return(serializedJSON, nil)

	data, err := sut.GetOne(userRole.ID)

	if err != nil {
		t.Fatalf("Expected no error but got %s", err)
	}

	if !cmp.Equal(userRole, data) {
		t.Fatalf("Expected json to be unmarshalled to %v but got %v; diff %s", userRole, data, cmp.Diff(userRole, data))
	}
}

func TestFailedGetOneUserRoleBecauseOfErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRoleID := "test-user-role-id"

	client.EXPECT().GetOne(gomock.Eq(userRoleID), gomock.Eq(restapi.UserRolesResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(userRoleID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetOneUserRoleBecauseOfInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRoleID := "test-user-role-id"

	client.EXPECT().GetOne(gomock.Eq(userRoleID), gomock.Eq(restapi.UserRolesResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(userRoleID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetOneUserRoleBecauseOfInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRoleID := "test-user-role-id"

	client.EXPECT().GetOne(gomock.Eq(userRoleID), gomock.Eq(restapi.UserRolesResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(userRoleID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetOneUserRoleBecauseResponseIsNotAValidJsonDocument(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRoleID := "test-user-role-id"

	client.EXPECT().GetOne(gomock.Eq(userRoleID), gomock.Eq(restapi.UserRolesResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(userRoleID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestSuccessfulGetAllUserRoles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRole1 := makeTestUserRoleWithCounter(1)
	userRole2 := makeTestUserRoleWithCounter(2)
	userRoles := []restapi.UserRole{userRole1, userRole2}
	serializedJSON, _ := json.Marshal(userRoles)

	client.EXPECT().GetAll(gomock.Eq(restapi.UserRolesResourcePath)).Return(serializedJSON, nil)

	data, err := sut.GetAll()

	if err != nil {
		t.Fatalf("Expected no error but got %s", err)
	}

	if !cmp.Equal(userRoles, data) {
		t.Fatalf("Expected json to be unmarshalled to %v but got %v; diff %s", userRoles, data, cmp.Diff(userRoles, data))
	}
}

func TestFailedGetAllUserRolesWithErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)

	client.EXPECT().GetAll(gomock.Eq(restapi.UserRolesResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetAll()

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetAllUserRolesWithInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)

	client.EXPECT().GetAll(gomock.Eq(restapi.UserRolesResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetAllUserRolesWithInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)

	client.EXPECT().GetAll(gomock.Eq(restapi.UserRolesResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetAllUserRolesWithNoJsonAsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)

	client.EXPECT().GetAll(gomock.Eq(restapi.UserRolesResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestSuccessfulUpsertOfUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRole := makeTestUserRole()
	serializedJSON, _ := json.Marshal(userRole)

	client.EXPECT().Put(gomock.Eq(userRole), gomock.Eq(restapi.UserRolesResourcePath)).Return(serializedJSON, nil)

	result, err := sut.Upsert(userRole)

	if err != nil {
		t.Fatalf("Expected no error but got %s", err)
	}

	if !cmp.Equal(userRole, result) {
		t.Fatalf("Expected json to be unmarshalled to %v but got %v; diff %s", userRole, result, cmp.Diff(result, result))
	}
}

func TestFailedUpsertOfUserRoleBecauseOfClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRole := makeTestUserRole()

	client.EXPECT().Put(gomock.Eq(userRole), gomock.Eq(restapi.UserRolesResourcePath)).Return(nil, errors.New("Error during test"))

	_, err := sut.Upsert(userRole)

	if err == nil {
		t.Fatal("Expected to get error")
	}
}

func TestFailedUpsertOfUserRoleBecauseOfInvalidResponseMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRole := makeTestUserRole()

	client.EXPECT().Put(gomock.Eq(userRole), gomock.Eq(restapi.UserRolesResourcePath)).Return([]byte("invalid response"), nil)

	_, err := sut.Upsert(userRole)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedUpsertOfUserRoleBecauseOfInvalidUserRoleInResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRole := makeTestUserRole()

	client.EXPECT().Put(gomock.Eq(userRole), gomock.Eq(restapi.UserRolesResourcePath)).Return([]byte("{ \"invalid\" : \"userRole\" }"), nil)

	_, err := sut.Upsert(userRole)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedUpsertOfUserRoleBecauseOfInvalidUserRoleProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRole := restapi.UserRole{
		Name: "Test UserRole",
	}

	client.EXPECT().Put(gomock.Eq(userRole), gomock.Eq(restapi.UserRolesResourcePath)).Times(0)

	_, err := sut.Upsert(userRole)

	if err == nil {
		t.Fatal("Expected to get error")
	}
}

func TestSuccessfulDeleteOfUserRoleByUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRole := makeTestUserRole()

	client.EXPECT().Delete(gomock.Eq("test-user-role-id-1"), gomock.Eq(restapi.UserRolesResourcePath)).Return(nil)

	err := sut.Delete(userRole)

	if err != nil {
		t.Fatalf("Expected no error got %s", err)
	}
}

func TestFailedDeleteOfUserRoleByUserRole(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewUserRoleResource(client)
	userRole := makeTestUserRole()

	client.EXPECT().Delete(gomock.Eq("test-user-role-id-1"), gomock.Eq(restapi.UserRolesResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(userRole)

	if err == nil {
		t.Fatal("Expected to get error")
	}
}

func makeTestUserRole() restapi.UserRole {
	return makeTestUserRoleWithCounter(1)
}

func makeTestUserRoleWithCounter(counter int) restapi.UserRole {
	id := fmt.Sprintf("test-user-role-id-%d", counter)
	name := fmt.Sprintf("Test User Role %d", counter)
	return restapi.UserRole{
		ID:                                id,
		Name:                              name,
		ImplicitViewFilter:                "Test view filter",
		CanConfigureServiceMapping:        true,
		CanConfigureEumApplications:       true,
		CanConfigureUsers:                 true,
		CanInstallNewAgents:               true,
		CanSeeUsageInformation:            true,
		CanConfigureIntegrations:          true,
		CanSeeOnPremiseLicenseInformation: true,
		CanConfigureRoles:                 true,
		CanConfigureCustomAlerts:          true,
		CanConfigureAPITokens:             true,
		CanConfigureAgentRunMode:          true,
		CanViewAuditLog:                   true,
		CanConfigureObjectives:            true,
		CanConfigureAgents:                true,
		CanConfigureAuthenticationMethods: true,
		CanConfigureApplications:          true,
	}
}
