package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"go.uber.org/mock/gomock"
)

// NewTestHelper creates a new instance of TestHelper
func NewTestHelper[T restapi.InstanaDataObject](t *testing.T) TestHelper[T] {
	return &testHelperImpl[T]{t: t}
}

// TestHelper definition of the test helper utility
type TestHelper[T restapi.InstanaDataObject] interface {
	WithMocking(t *testing.T, testFunction func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanApi *mocks.MockInstanaAPI))
	CreateProviderMetaMock(ctrl *gomock.Controller) (*ProviderMeta, *mocks.MockInstanaAPI)
	CreateEmptyResourceDataForResourceHandle(resourceHandle ResourceHandle[T]) *schema.ResourceData
	CreateResourceDataForResourceHandle(resourceHandle ResourceHandle[T], data map[string]interface{}) *schema.ResourceData
}

type testHelperImpl[T restapi.InstanaDataObject] struct {
	t *testing.T
}

func (inst *testHelperImpl[T]) WithMocking(t *testing.T, testFunction func(ctrl *gomock.Controller, meta *ProviderMeta, mockInstanApi *mocks.MockInstanaAPI)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	providerMeta, mockInstanaAPI := inst.CreateProviderMetaMock(ctrl)
	testFunction(ctrl, providerMeta, mockInstanaAPI)
}

func (inst *testHelperImpl[T]) CreateProviderMetaMock(ctrl *gomock.Controller) (*ProviderMeta, *mocks.MockInstanaAPI) {
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)
	providerMeta := &ProviderMeta{
		InstanaAPI: mockInstanaAPI,
	}
	return providerMeta, mockInstanaAPI
}

func (inst *testHelperImpl[T]) CreateEmptyResourceDataForResourceHandle(resourceHandle ResourceHandle[T]) *schema.ResourceData {
	data := make(map[string]interface{})
	return inst.CreateResourceDataForResourceHandle(resourceHandle, data)
}

func (inst *testHelperImpl[T]) CreateResourceDataForResourceHandle(resourceHandle ResourceHandle[T], data map[string]interface{}) *schema.ResourceData {
	return schema.TestResourceDataRaw(inst.t, resourceHandle.MetaData().Schema, data)
}
