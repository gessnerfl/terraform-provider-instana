package restapi_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTheProperRepresentationsForSeverityWarning(t *testing.T) {
	assert.Equal(t, 5, SeverityWarning.GetAPIRepresentation())
	assert.Equal(t, "warning", SeverityWarning.GetTerraformRepresentation())
}

func TestShouldReturnTheProperRepresentationsForSeverityCritical(t *testing.T) {
	assert.Equal(t, 10, SeverityCritical.GetAPIRepresentation())
	assert.Equal(t, "critical", SeverityCritical.GetTerraformRepresentation())
}

func TestShouldReturnSupportedSeveritiesAsStringSliceOfTerraformRepresentations(t *testing.T) {
	expected := []string{"warning", "critical"}
	require.Equal(t, expected, SupportedSeverities.TerraformRepresentations())
}

func TestShouldReturnSupportedSeveritiesAsStringSliceOfInstanaAPIRepresentations(t *testing.T) {
	expected := []int{5, 10}
	require.Equal(t, expected, SupportedSeverities.APIRepresentations())
}
