package restapi_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldSuccessfullyUnmarshalCustomEventSpecifications(t *testing.T) {
	description := "event-description"
	query := "event-query"
	expirationTime := 60000
	systemRule := NewSystemRuleSpecification("system-rule-id", SeverityWarning.GetAPIRepresentation())
	customEventSpecification := CustomEventSpecification{
		ID:             "event-id",
		Name:           "event-name",
		EntityType:     "entity-type",
		Enabled:        true,
		Triggering:     true,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Query:          &query,
		Rules:          []RuleSpecification{systemRule},
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

	if !cmp.Equal(result, CustomEventSpecification{}) {
		t.Fatal("Expected empty custom event specification")
	}
}
