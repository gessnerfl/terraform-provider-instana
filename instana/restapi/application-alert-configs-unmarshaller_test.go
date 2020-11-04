package restapi

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_applicationAlertConfigsUnmarshaller_Unmarshal(t *testing.T) {

	var applicationAlertConfigs = ApplicationAlertConfigs{
		ID:              "1234",
		AlertName:       "Test Alert Name",
		ApplicationId:   "4567",
		Rule:            ApplicationAlertConfigsRule{},
		Description:     "Sample Description",
		Severity:        5,
		Threshold:       Threshold{},
		TagFilters: []ApplicationAlertConfigsTagFilter{
			{
				Type:         "TAG_FILTER",
				Name:         "endpoint.name",
				StringValue:  "foobar",
				NumberValue:  0,
				BooleanValue: false,
				Operator:     "EQUALS",
				Entity:       "NOT_APPLICABLE",
						},
		},
	}
	serializedJSON, _ := json.Marshal(applicationAlertConfigs)

	result, err := NewApplicationAlertConfigsUnmarshaller().Unmarshal(serializedJSON)
	assert.Nil(t, err)
	assert.Equal(t, applicationAlertConfigs, result)

}
