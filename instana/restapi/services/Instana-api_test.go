package services_test

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/services"
)

func TestShouldReturnResourcesFromInstanaAPI(t *testing.T) {
	api := NewInstanaAPI("api-token", "endpoint")

	t.Run("Should return CustomEventSpecificationResource instance", func(t *testing.T) {
		customEventSpecificationResource := api.CustomEventSpecifications()
		if customEventSpecificationResource == nil {
			t.Fatal("Expected instance of CustomEventSpecificationResource to be returned")
		}
	})
	t.Run("Should return UserRoleResource instance", func(t *testing.T) {
		userRoleResource := api.UserRoles()
		if userRoleResource == nil {
			t.Fatal("Expected instance of UserRoleResource to be returned")
		}
	})
	t.Run("Should return ApplicationConfigResource instance", func(t *testing.T) {
		applicationConfigResource := api.ApplicationConfigs()
		if applicationConfigResource == nil {
			t.Fatal("Expected instance of ApplicationConfigResource to be returned")
		}
	})
	t.Run("Should return AlertingChannelResource instance", func(t *testing.T) {
		alertingChannelResource := api.AlertingChannels()
		if alertingChannelResource == nil {
			t.Fatal("Expected instance of AlertingChannelResource to be returned")
		}
	})
}

//Add tests for unmarshallers
func TestShouldSuccessfullyUnmarshalApplicationConfig(t *testing.T) {
	id := "test-application-config-id"
	label := "Test Application Config Label"
	applicationConfig := restapi.ApplicationConfig{
		ID:                 id,
		Label:              label,
		MatchSpecification: restapi.NewBinaryOperator(restapi.NewComparisionExpression("key", restapi.EqualsOperator, "value"), restapi.LogicalAnd, restapi.NewUnaryOperationExpression("key", restapi.NotBlankOperator)),
		Scope:              "scope",
	}

	serializedJSON, _ := json.Marshal(applicationConfig)

	result, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	if err != nil {
		t.Fatal("Expected application config to be successfully unmarshalled")
	}

	if !cmp.Equal(result, applicationConfig) {
		t.Fatalf("Expected application config to be properly unmarshalled, %s", cmp.Diff(result, applicationConfig))
	}
}

func TestShouldFailToUnmarashalApplicationConfigWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewApplicationConfigUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldReturnEmptyApplicationConfigWhenNoFieldOfResponseMatchesToModel(t *testing.T) {
	response := `{"foo" : "bar"}`
	_, err := NewApplicationConfigUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail when response is not matching application config")
	}
}

func TestShouldFailToUnmarashalApplicationConfigWhenResponseIsNotAValidJson(t *testing.T) {
	response := `Invalid Data`

	_, err := NewApplicationConfigUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldFailToUnmarashalApplicationConfigWhenExpressionTypeIsNotSupported(t *testing.T) {
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

	_, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	if err == nil {
		t.Fatal("Expected unmarshalling to fail because of unsupported expression type")
	}
}

func TestShouldFailToUnmarashalApplicationConfigWhenLeftSideOfBinaryExpressionTypeIsNotValid(t *testing.T) {
	left := restapi.TagMatcherExpression{
		Key:      "foo",
		Operator: restapi.NotEmptyOperator,
	}
	right := restapi.NewUnaryOperationExpression("foo", restapi.IsEmptyOperator)
	testShouldFailToUnmarashalApplicationConfigWhenOneSideOfBinaryExpressionIsNotValid(left, right, t)
}

func TestShouldFailToUnmarashalApplicationConfigWhenRightSideOfBinaryExpressionTypeIsNotValid(t *testing.T) {
	left := restapi.NewUnaryOperationExpression("foo", restapi.IsEmptyOperator)
	right := restapi.TagMatcherExpression{
		Key:      "foo",
		Operator: restapi.NotEmptyOperator,
	}
	testShouldFailToUnmarashalApplicationConfigWhenOneSideOfBinaryExpressionIsNotValid(left, right, t)
}

func testShouldFailToUnmarashalApplicationConfigWhenOneSideOfBinaryExpressionIsNotValid(left restapi.MatchExpression, right restapi.MatchExpression, t *testing.T) {
	applicationConfig := restapi.ApplicationConfig{
		ID:                 "id",
		Label:              "label",
		MatchSpecification: restapi.NewBinaryOperator(left, restapi.LogicalOr, right),
		Scope:              "scope",
	}
	serializedJSON, _ := json.Marshal(applicationConfig)

	_, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	if err == nil {
		t.Fatal("Expected unmarshalling to fail because of invalid binary expression")
	}
}

func TestShouldSuccessfullyUnmarshalCustomEventSpecifications(t *testing.T) {
	description := "event-description"
	query := "event-query"
	expirationTime := 60000
	systemRule := restapi.NewSystemRuleSpecification("system-rule-id", restapi.SeverityWarning.GetAPIRepresentation())
	customEventSpecification := restapi.CustomEventSpecification{
		ID:             "event-id",
		Name:           "event-name",
		EntityType:     "entity-type",
		Enabled:        true,
		Triggering:     true,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Query:          &query,
		Rules:          []restapi.RuleSpecification{systemRule},
	}

	serializedJSON, _ := json.Marshal(customEventSpecification)

	result, err := NewCustomEventSpecificationUnmarshaller().Unmarshal(serializedJSON)

	if err != nil {
		t.Fatal("Expected custom event specification to be successfully unmarshalled")
	}

	if !cmp.Equal(result, customEventSpecification) {
		t.Fatalf("Expected custom event specification to be properly unmarshalled, %s", cmp.Diff(result, customEventSpecification))
	}
}

func TestShouldFailToUnmarshalCustomEventSpecificationWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewCustomEventSpecificationUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldFailToUnmarshalCustomEventSpecificationsWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewCustomEventSpecificationUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldReturnEmptyCustomEventSpecificationWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewCustomEventSpecificationUnmarshaller().Unmarshal([]byte(response))

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal custom event specification response, %s", err)
	}

	if !cmp.Equal(result, restapi.CustomEventSpecification{}) {
		t.Fatal("Expected empty custom event specification")
	}
}

func TestShouldSuccessfullyUnmarshalUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                                "role-id",
		Name:                              "role-name",
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

	serializedJSON, _ := json.Marshal(userRole)

	result, err := NewUserRoleUnmarshaller().Unmarshal(serializedJSON)

	if err != nil {
		t.Fatal("Expected user role to be successfully unmarshalled")
	}

	if !cmp.Equal(result, userRole) {
		t.Fatalf("Expected user role to be properly unmarshalled, %s", cmp.Diff(result, userRole))
	}
}

func TestShouldFailToUnmarshalUserRoleWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewUserRoleUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldFailToUnmarshalUserRoleWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewUserRoleUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldReturnEmptyUserRoleWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewUserRoleUnmarshaller().Unmarshal([]byte(response))

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal user role response, %s", err)
	}

	if !cmp.Equal(result, restapi.UserRole{}) {
		t.Fatal("Expected empty user role")
	}
}

func TestShouldSuccessfullyUnmarshalAlertingChannel(t *testing.T) {
	response := `{
		"id" : "test-id",
		"name" : "test-name",
		"kind" : "EMAIL",
		"emails" : ["test-email1","test-email2"]
	}`

	result, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal alerting channel response; %s", err)
	}

	alertingChannel, ok := result.(restapi.AlertingChannel)
	if !ok {
		t.Fatal("Expected result to be a alerting channel")
	}

	if alertingChannel.ID != "test-id" {
		t.Fatal("Expected ID to be properly mapped")
	}
	if alertingChannel.Name != "test-name" {
		t.Fatal("Expected name to be properly mapped")
	}
	if alertingChannel.Kind != restapi.EmailChannelType {
		t.Fatal("Expected kind to be properly mapped")
	}
	if !cmp.Equal(alertingChannel.Emails, []string{"test-email1", "test-email2"}) {
		t.Fatal("Expected emails to be properly mapped")
	}
}

func TestShouldFailToUnmarshalAlertingChannelWhenResponseIsAJsonArray(t *testing.T) {
	response := `["test-email1","test-email2"]`

	_, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldFailToUnmarshalAlertingChannelWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldReturnEmptyAlertingChannelWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewAlertingChannelUnmarshaller().Unmarshal([]byte(response))

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal alerting channel response, %s", err)
	}

	if !cmp.Equal(result, restapi.AlertingChannel{}) {
		t.Fatal("Expected empty alerting channel")
	}
}
