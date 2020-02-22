package restapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

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

	assert.Nil(t, err)
	assert.Equal(t, customEventSpecification, result)
}

func TestShouldFailToUnmarshalCustomEventSpecificationWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewCustomEventSpecificationUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldFailToUnmarshalCustomEventSpecificationsWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewCustomEventSpecificationUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldReturnEmptyCustomEventSpecificationWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewCustomEventSpecificationUnmarshaller().Unmarshal([]byte(response))

	assert.Nil(t, err)
	assert.Equal(t, CustomEventSpecification{}, result)
}
