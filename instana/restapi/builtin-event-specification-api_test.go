package restapi_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldReturnIdOfBuiltinEventSpecification(t *testing.T) {
	sut := BuiltinEventSpecification{
		ID: "1234",
	}

	require.Equal(t, "1234", sut.GetIDForResourcePath())
}

func TestShouldReturnNoErrorForEmptyBuiltInSpecification(t *testing.T) {
	sut := BuiltinEventSpecification{}

	require.NoError(t, sut.Validate())
}

func TestShouldReturnNoErrorForCompleteBuiltInSpecification(t *testing.T) {
	description := "description"
	sut := BuiltinEventSpecification{
		ID:            "id",
		Name:          "name",
		Description:   &description,
		ShortPluginID: "shortPluginId",
		Severity:      10,
		Enabled:       true,
		Triggering:    true,
	}

	require.NoError(t, sut.Validate())
}
