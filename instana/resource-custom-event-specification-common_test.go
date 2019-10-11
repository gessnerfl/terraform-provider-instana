package instana_test

import (
	"fmt"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/hashicorp/terraform/terraform"
)

func TestShouldFailToMigrateCustomEventStateFromVersionOtherThan0ForCreatedMigrationFunction(t *testing.T) {
	function := CreateMigrateCustomEventConfigStateFunction(make(map[int](func(inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error))))
	state := &terraform.InstanceState{}
	state.ID = "TEST_ID"

	for i := 1; i <= 10; i++ {
		t.Run(fmt.Sprintf("TestShouldFailToMigrateCustomEventStateFromVersionOtherThan0ForCreatedMigrationFunction_%d", i), func(t *testing.T) {
			_, err := function(i, state, &ProviderMeta{})

			if err == nil {
				t.Fatalf("Expected error when triggering migration of state from for unknown version %d", i)
			}
		})
	}
}

func TestShouldSkipMigrationWhenEmptyCustomEventStateIsProvided(t *testing.T) {
	state := &terraform.InstanceState{}

	function := CreateMigrateCustomEventConfigStateFunction(make(map[int](func(inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error))))

	for i := 0; i <= 10; i++ {
		t.Run(fmt.Sprintf("TestShouldSkipMigrationWhenEmptyCustomEventStateIsProvided_%d", i), func(t *testing.T) {
			result, err := function(i, state, &ProviderMeta{})

			if err != nil {
				t.Fatalf("No error expected during migration of state from version 0 to 1 but got %s", err)
			}

			if result.Attributes != nil {
				t.Fatal("No changes should be applied during migration of an empty state from version 0 to 1")
			}
		})
	}
}

func TestShouldMigrateCustomEventStateAndAddFullNameWithSameValueAsNameWhenMigratingFromVersion0To1(t *testing.T) {
	name := "Test Name"
	data := make(map[string]string)
	data[CustomEventSpecificationFieldName] = name
	state := &terraform.InstanceState{}
	state.ID = "TEST_ID"
	state.Attributes = data

	function := CreateMigrateCustomEventConfigStateFunction(make(map[int](func(inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error))))
	result, err := function(0, state, &ProviderMeta{})

	if err != nil {
		t.Fatalf("No error expected during migration of state from version 0 to 1 but got %s", err)
	}

	if result.Attributes[CustomEventSpecificationFieldFullName] != name {
		t.Fatal("Full name should be initialized with value of name when migrating from version 0 to 1")
	}
}

func TestShouldMigrateCustomEventStateAndApplyProvidedCustomMigrationFunction(t *testing.T) {
	testAttribute := "test_attribute"
	testString := "test_string"
	data := make(map[string]string)
	state := &terraform.InstanceState{}
	state.ID = "TEST_ID"
	state.Attributes = data

	specificFunctions := map[int](func(inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error)){
		0: func(inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
			inst.Attributes[testAttribute] = testString
			return inst, nil
		},
	}

	function := CreateMigrateCustomEventConfigStateFunction(specificFunctions)
	result, err := function(0, state, &ProviderMeta{})

	if err != nil {
		t.Fatalf("No error expected during migration of state from version 0 to 1 but got %s", err)
	}

	if result.Attributes[testAttribute] != testString {
		t.Fatal("Custom migration could should be executed")
	}
}
