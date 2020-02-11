package instana

import (
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform/helper/schema"
)

//NewTestHelper creates a new instance of TestHelper
func NewTestHelper(t *testing.T) TestHelper {
	return &testHelperImpl{t: t}
}

//TestHelper definition of the test helper utility
type TestHelper interface {
	WithMocking(t *testing.T, testFunction func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanApi *mocks.MockInstanaAPI, mockFormatter *mocks.MockResourceNameFormatter))
	CreateProviderMetaMock(ctrl *gomock.Controller) (*ProviderMeta, *mocks.MockInstanaAPI, *mocks.MockResourceNameFormatter)
	CreateEmptyResourceDataForResourceHandle(resourceHandle *ResourceHandle) *schema.ResourceData
	CreateResourceDataForResourceHandle(resourceHandle *ResourceHandle, data map[string]interface{}) *schema.ResourceData
}

type testHelperImpl struct {
	t *testing.T
}

func (inst *testHelperImpl) WithMocking(t *testing.T, testFunction func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanApi *mocks.MockInstanaAPI, mockFormatter *mocks.MockResourceNameFormatter)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	providerMeta, mockInstanaAPI, mockResourceNameFormatter := inst.CreateProviderMetaMock(ctrl)
	testFunction(ctrl, providerMeta, mockInstanaAPI, mockResourceNameFormatter)
}

func (inst *testHelperImpl) CreateProviderMetaMock(ctrl *gomock.Controller) (*ProviderMeta, *mocks.MockInstanaAPI, *mocks.MockResourceNameFormatter) {
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)
	mockResourceNameFormatter := mocks.NewMockResourceNameFormatter(ctrl)
	providerMeta := &ProviderMeta{
		InstanaAPI:            mockInstanaAPI,
		ResourceNameFormatter: mockResourceNameFormatter,
	}
	return providerMeta, mockInstanaAPI, mockResourceNameFormatter
}

func (inst *testHelperImpl) CreateEmptyResourceDataForResourceHandle(resourceHandle *ResourceHandle) *schema.ResourceData {
	data := make(map[string]interface{})
	return inst.CreateResourceDataForResourceHandle(resourceHandle, data)
}

func (inst *testHelperImpl) CreateResourceDataForResourceHandle(resourceHandle *ResourceHandle, data map[string]interface{}) *schema.ResourceData {
	return schema.TestResourceDataRaw(inst.t, resourceHandle.Schema, data)
}
