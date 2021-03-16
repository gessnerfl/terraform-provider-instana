package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/assert"
)

const (
	sliConfigID                         = "sli-config-id"
	sliConfigName                       = "sli-config-name"
	sliConfigInitialEvaluationTimestamp = 0
	sliConfigMetricName                 = "sli-config-metric-name"
	sliConfigMetricAggregation          = "sli-config-metric-aggregation"
	sliConfigMetricThreshold            = 1.0
	sliConfigEntityType                 = "application"
	sliConfigEntityApplicationID        = "sli-config-entity-application-id"
	sliConfigEntityServiceID            = "sli-config-entity-service-id"
	sliConfigEntityEndpointID           = "sli-config-entity-endpoint-id"
	sliConfigEntityBoundaryScope        = "ALL"
)

func TestMinimalSliConfig(t *testing.T) {
	sliConfig := SliConfig{
		ID:   sliConfigID,
		Name: sliConfigName,
		MetricConfiguration: MetricConfiguration{
			Name:        sliConfigMetricName,
			Aggregation: sliConfigMetricAggregation,
			Threshold:   sliConfigMetricThreshold,
		},
		SliEntity: SliEntity{
			Type:          sliConfigEntityType,
			BoundaryScope: sliConfigEntityBoundaryScope,
		},
	}

	assert.Equal(t, sliConfigID, sliConfig.GetIDForResourcePath())

	err := sliConfig.Validate()
	assert.Nil(t, err)
}

func TestValidFullSliConfig(t *testing.T) {
	sliConfig := SliConfig{
		ID:                         sliConfigID,
		Name:                       sliConfigName,
		InitialEvaluationTimestamp: sliConfigInitialEvaluationTimestamp,
		MetricConfiguration: MetricConfiguration{
			Name:        sliConfigMetricName,
			Aggregation: sliConfigMetricAggregation,
			Threshold:   sliConfigMetricThreshold,
		},
		SliEntity: SliEntity{
			Type:          sliConfigEntityType,
			ApplicationID: sliConfigEntityApplicationID,
			ServiceID:     sliConfigEntityServiceID,
			EndpointID:    sliConfigEntityEndpointID,
			BoundaryScope: sliConfigEntityBoundaryScope,
		},
	}

	assert.Equal(t, sliConfigID, sliConfig.GetIDForResourcePath())

	err := sliConfig.Validate()
	assert.Nil(t, err)
}

func TestInvalidSliConfigBecauseOfMissingID(t *testing.T) {
	sliConfig := SliConfig{
		Name: sliConfigName,
	}

	err := sliConfig.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "id")
}

func TestInvalidSliConfigBecauseOfMissingName(t *testing.T) {
	sliConfig := SliConfig{
		ID: sliConfigID,
	}

	err := sliConfig.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestInvalidSliConfigBecauseOfMissingMetricName(t *testing.T) {
	sliConfig := SliConfig{
		ID:   sliConfigID,
		Name: sliConfigName,
	}

	err := sliConfig.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "metric name")
}

func TestInvalidSliConfigBecauseOfMissingMetricAggregation(t *testing.T) {
	sliConfig := SliConfig{
		ID:   sliConfigID,
		Name: sliConfigName,
		MetricConfiguration: MetricConfiguration{
			Name: sliConfigMetricName,
		},
	}

	err := sliConfig.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "aggregation")
}

func TestInvalidSliConfigBecauseOfMissingMetricThreshold(t *testing.T) {
	sliConfig := SliConfig{
		ID:   sliConfigID,
		Name: sliConfigName,
		MetricConfiguration: MetricConfiguration{
			Name:        sliConfigMetricName,
			Aggregation: sliConfigMetricAggregation,
		},
	}

	err := sliConfig.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "threshold")
}

func TestInvalidSliConfigBecauseOfInvalidMetricThreshold(t *testing.T) {
	sliConfig := SliConfig{
		ID:   sliConfigID,
		Name: sliConfigName,
		MetricConfiguration: MetricConfiguration{
			Name:        sliConfigMetricName,
			Aggregation: sliConfigMetricAggregation,
			Threshold:   -1.0,
		},
	}

	err := sliConfig.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "threshold")
}

func TestInvalidSliConfigBecauseOfMissingSliEntityType(t *testing.T) {
	sliConfig := SliConfig{
		ID:   sliConfigID,
		Name: sliConfigName,
		MetricConfiguration: MetricConfiguration{
			Name:        sliConfigMetricName,
			Aggregation: sliConfigMetricAggregation,
			Threshold:   sliConfigMetricThreshold,
		},
	}

	err := sliConfig.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "sli type")
}

func TestInvalidSliConfigBecauseOfMissingSliEntityBoundaryScope(t *testing.T) {
	sliConfig := SliConfig{
		ID:   sliConfigID,
		Name: sliConfigName,
		MetricConfiguration: MetricConfiguration{
			Name:        sliConfigMetricName,
			Aggregation: sliConfigMetricAggregation,
			Threshold:   sliConfigMetricThreshold,
		},
		SliEntity: SliEntity{
			Type: sliConfigEntityType,
		},
	}

	err := sliConfig.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "boundary scope")
}
