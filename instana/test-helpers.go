package instana

import (
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
)

//NewTestHelper creates a new instance of TestHelper
func NewTestHelper(t *testing.T) TestHelper {
	return &testHelperImpl{t: t}
}

//TestHelper definition of the test helper utility
type TestHelper interface {
	CreateEmptyRuleResourceData() *schema.ResourceData
	CreateRuleResourceData(data map[string]interface{}) *schema.ResourceData
	CreateEmptyRuleBindingResourceData() *schema.ResourceData
	CreateRuleBindingResourceData(data map[string]interface{}) *schema.ResourceData
	CreateEmptyUserRoleResourceData() *schema.ResourceData
	CreateUserRoleResourceData(data map[string]interface{}) *schema.ResourceData
	CreateEmptyApplicationConfigResourceData() *schema.ResourceData
	CreateApplicationConfigResourceData(data map[string]interface{}) *schema.ResourceData
	CreateEmptyCustomSystemEventSpecificationResourceData() *schema.ResourceData
	CreateCustomSystemEventSpecificationResourceData(data map[string]interface{}) *schema.ResourceData
}

type testHelperImpl struct {
	t *testing.T
}

func (inst *testHelperImpl) CreateEmptyRuleResourceData() *schema.ResourceData {
	data := make(map[string]interface{})
	return inst.CreateRuleResourceData(data)
}

func (inst *testHelperImpl) CreateRuleResourceData(data map[string]interface{}) *schema.ResourceData {
	schemaMap := CreateResourceRule().Schema
	return schema.TestResourceDataRaw(inst.t, schemaMap, data)
}

func (inst *testHelperImpl) CreateEmptyRuleBindingResourceData() *schema.ResourceData {
	data := make(map[string]interface{})
	return inst.CreateRuleBindingResourceData(data)
}

func (inst *testHelperImpl) CreateRuleBindingResourceData(data map[string]interface{}) *schema.ResourceData {
	schemaMap := CreateResourceRuleBinding().Schema
	return schema.TestResourceDataRaw(inst.t, schemaMap, data)
}

func (inst *testHelperImpl) CreateEmptyUserRoleResourceData() *schema.ResourceData {
	data := make(map[string]interface{})
	return inst.CreateUserRoleResourceData(data)
}

func (inst *testHelperImpl) CreateUserRoleResourceData(data map[string]interface{}) *schema.ResourceData {
	schemaMap := CreateResourceUserRole().Schema
	return schema.TestResourceDataRaw(inst.t, schemaMap, data)
}

func (inst *testHelperImpl) CreateEmptyApplicationConfigResourceData() *schema.ResourceData {
	data := make(map[string]interface{})
	return inst.CreateApplicationConfigResourceData(data)
}

func (inst *testHelperImpl) CreateApplicationConfigResourceData(data map[string]interface{}) *schema.ResourceData {
	schemaMap := CreateResourceApplicationConfig().Schema
	return schema.TestResourceDataRaw(inst.t, schemaMap, data)
}

func (inst *testHelperImpl) CreateEmptyCustomSystemEventSpecificationResourceData() *schema.ResourceData {
	data := make(map[string]interface{})
	return inst.CreateCustomSystemEventSpecificationResourceData(data)
}

func (inst *testHelperImpl) CreateCustomSystemEventSpecificationResourceData(data map[string]interface{}) *schema.ResourceData {
	schemaMap := CreateResourceCustomSystemEventSpecification().Schema
	return schema.TestResourceDataRaw(inst.t, schemaMap, data)
}
