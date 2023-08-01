package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.uber.org/mock/gomock"
)

// NewTestHelper creates a new instance of TestHelper
func NewTestHelper[T restapi.InstanaDataObject](t *testing.T) TestHelper[T] {
	return &testHelperImpl[T]{t: t}
}

// TestHelper definition of the test helper utility
type TestHelper[T restapi.InstanaDataObject] interface {
	WithMocking(t *testing.T, testFunction func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanApi *mocks.MockInstanaAPI, mockFormatter *mocks.MockResourceNameFormatter))
	CreateProviderMetaMock(ctrl *gomock.Controller) (*ProviderMeta, *mocks.MockInstanaAPI, *mocks.MockResourceNameFormatter)
	CreateEmptyResourceDataForResourceHandle(resourceHandle ResourceHandle[T]) *schema.ResourceData
	CreateResourceDataForResourceHandle(resourceHandle ResourceHandle[T], data map[string]interface{}) *schema.ResourceData
	ResourceFormatter() utils.ResourceNameFormatter
}

type testHelperImpl[T restapi.InstanaDataObject] struct {
	t *testing.T
}

func (inst *testHelperImpl[T]) WithMocking(t *testing.T, testFunction func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanApi *mocks.MockInstanaAPI, mockFormatter *mocks.MockResourceNameFormatter)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	providerMeta, mockInstanaAPI, mockResourceNameFormatter := inst.CreateProviderMetaMock(ctrl)
	testFunction(ctrl, providerMeta, mockInstanaAPI, mockResourceNameFormatter)
}

func (inst *testHelperImpl[T]) CreateProviderMetaMock(ctrl *gomock.Controller) (*ProviderMeta, *mocks.MockInstanaAPI, *mocks.MockResourceNameFormatter) {
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)
	mockResourceNameFormatter := mocks.NewMockResourceNameFormatter(ctrl)
	providerMeta := &ProviderMeta{
		InstanaAPI:            mockInstanaAPI,
		ResourceNameFormatter: mockResourceNameFormatter,
	}
	return providerMeta, mockInstanaAPI, mockResourceNameFormatter
}

func (inst *testHelperImpl[T]) CreateEmptyResourceDataForResourceHandle(resourceHandle ResourceHandle[T]) *schema.ResourceData {
	data := make(map[string]interface{})
	return inst.CreateResourceDataForResourceHandle(resourceHandle, data)
}

func (inst *testHelperImpl[T]) CreateResourceDataForResourceHandle(resourceHandle ResourceHandle[T], data map[string]interface{}) *schema.ResourceData {
	return schema.TestResourceDataRaw(inst.t, resourceHandle.MetaData().Schema, data)
}

func (inst *testHelperImpl[T]) ResourceFormatter() utils.ResourceNameFormatter {
	return utils.NewResourceNameFormatter("prefix", "suffix")
}
