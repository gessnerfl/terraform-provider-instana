package restapi_test

import (
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestShouldSuccessfullyUnmarshalCustomDashboardIntoStruct(t *testing.T) {
	widgetJson := `[
    {
      "id": "6jK0w8KmdHtABCs3",
      "title": "Latency",
      "width": 4,
      "height": 13,
      "x": 4,
      "y": 26,
      "type": "chart",
      "config": {
        "y1": {
          "formatter": "millis.detailed",
          "renderer": "line",
          "metrics": [
            {
              "metric": "latency",
              "timeShift": 0,
              "tagFilters": [
                {
                  "stringValue": "my-app",
                  "name": "application.name",
                  "entity": "DESTINATION",
                  "operator": "EQUALS"
                },
                {
                  "name": "call.inbound_of_application",
                  "entity": "NOT_APPLICABLE",
                  "operator": "NOT_EMPTY"
                }
              ],
              "aggregation": "MEAN",
              "label": "Mean Latency",
              "source": "APPLICATION"
            },
            {
              "metric": "latency",
              "timeShift": 0,
              "tagFilters": [
                {
                  "stringValue": "my-app",
                  "name": "application.name",
                  "entity": "DESTINATION",
                  "operator": "EQUALS"
                },
                {
                  "name": "call.inbound_of_application",
                  "entity": "NOT_APPLICABLE",
                  "operator": "NOT_EMPTY"
                }
              ],
              "aggregation": "P99",
              "label": "99th latency",
              "source": "APPLICATION"
            }
          ]
        },
        "y2": {
          "formatter": "number.detailed",
          "renderer": "line",
          "metrics": []
        },
        "type": "TIME_SERIES"
      }
    }
  ]`
	dashboadJson := `
{
  "id": "dashboard-id-1234",
  "title": "My Dashboard",
  "accessRules": [
    {
      "accessType": "READ_WRITE",
      "relationType": "USER",
      "relatedId": "user-id-1"
    },
    {
      "accessType": "READ",
      "relationType": "GLOBAL",
      "relatedId": ""
    },
    {
      "accessType": "READ_WRITE",
      "relationType": "USER",
      "relatedId": "user-id-2"
    }
  ],
  "widgets": __WIDGETS__,
  "writable": false
}
`
	jsonMessage := strings.ReplaceAll(dashboadJson, "__WIDGETS__", widgetJson)

	result, err := NewCustomDashboardUnmarshaller().Unmarshal([]byte(jsonMessage))

	require.NoError(t, err)
	require.Equal(t, "dashboard-id-1234", result.(*CustomDashboard).ID)
	require.Equal(t, "My Dashboard", result.(*CustomDashboard).Title)
	userId1 := "user-id-1"
	userId2 := "user-id-2"
	require.Equal(t, []AccessRule{
		{AccessType: AccessTypeReadWrite, RelationType: RelationTypeUser, RelatedID: &userId1},
		{AccessType: AccessTypeRead, RelationType: RelationTypeGlobal},
		{AccessType: AccessTypeReadWrite, RelationType: RelationTypeUser, RelatedID: &userId2},
	}, result.(*CustomDashboard).AccessRules)
	require.Equal(t, widgetJson, result.(*CustomDashboard).Widgets)
}

func TestShouldFailToUnmarshalCustomDashboardWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewCustomDashboardUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldReturnEmptyCustomDashboardWhenNoFieldOfResponseMatchesToModel(t *testing.T) {
	response := `{"foo" : "bar"}`
	config, err := NewCustomDashboardUnmarshaller().Unmarshal([]byte(response))

	require.NoError(t, err)
	require.Equal(t, &CustomDashboard{}, config)
}

func TestShouldFailToUnmarshalCustomDashboardWhenResponseIsNotAValidJson(t *testing.T) {
	response := `Invalid Data`

	_, err := NewCustomDashboardUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}
