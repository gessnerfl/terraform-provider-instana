package services_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/services"
)

func TestShouldReturnResourcesFromInstanaAPI(t *testing.T) {
	api := NewInstanaAPI("api-token", "endpoint")

	t.Run("Should return RuleResource instance", func(t *testing.T) {
		ruleResource := api.Rules()
		if ruleResource == nil {
			t.Fatal("Expected instance of RuleResource to be returned")
		}
	})
	t.Run("Should return RuleBindingResource instance", func(t *testing.T) {
		ruleBindingResource := api.RuleBindings()
		if ruleBindingResource == nil {
			t.Fatal("Expected instance of RuleBindingResource to be returned")
		}
	})
}
