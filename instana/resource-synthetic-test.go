package instana

import (
	"encoding/json"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceInstanaSyntheticTest the name of the terraform-provider-instana resource to manage synthetic tests
const ResourceInstanaSyntheticTest = "instana_synthetic_test"

const (
	//SyntheticTestFieldLabel constant value for the schema field label
	SyntheticTestFieldLabel = "label"
	//SyntheticTestFieldDescription constant value for the computed schema field description
	SyntheticTestFieldDescription = "description"
	//SyntheticTestFieldActive constant value for the schema field active
	SyntheticTestFieldActive = "active"
	//SyntheticTestFieldApplicationID constant value for the schema field application_id
	SyntheticTestFieldApplicationID = "application_id"
	//SyntheticTestFieldConfiguration constant value for the schema field configuration
	SyntheticTestFieldConfiguration = "configuration"
	//SyntheticTestFieldCustomProperties constant value for the schema field custom_properties
	SyntheticTestFieldCustomProperties = "custom_properties"
	//SyntheticTestFieldLocations constant value for the schema field locations
	SyntheticTestFieldLocations = "locations"
	//SyntheticTestFieldPlaybackMode constant value for the schema field playback_mode
	SyntheticTestFieldPlaybackMode = "playback_mode"
	//SyntheticTestFieldTestFrequency constant value for the schema field test_frequency
	SyntheticTestFieldTestFrequency = "test_frequency"
)

// NewSyntheticTestResourceHandle creates the resource handle Synthetic Tests
func NewSyntheticTestResourceHandle() ResourceHandle {
	return &syntheticTestResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaSyntheticTest,
			Schema: map[string]*schema.Schema{
				SyntheticTestFieldLabel: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The label of the synthetic alert",
				},
				SyntheticTestFieldDescription: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The description of the synthetic alert",
				},
				SyntheticTestFieldActive: {
					Type:        schema.TypeBool,
					Required:    true,
					Description: "The status of the synthetic alert. Can be enabled or disabled",
				},
				SyntheticTestFieldApplicationID: {
					Type:        schema.TypeBool,
					Computed:    true,
					Description: "The application ID of the synthetic alert",
				},
				SyntheticTestFieldConfiguration: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The configuration of the synthetic alert",
					DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
						return NormalizeJSONString(old) == NormalizeJSONString(new)
					},
					StateFunc: func(val interface{}) string {
						return NormalizeJSONString(val.(string))
					},
				},
				SyntheticTestFieldCustomProperties: {
					Type:        schema.TypeMap,
					Optional:    true,
					Description: "",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				SyntheticTestFieldLocations: {
					Type:        schema.TypeList,
					Required:    true,
					Description: "",
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				SyntheticTestFieldPlaybackMode: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "",
				},
				SyntheticTestFieldTestFrequency: {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "",
				},
			},
			SchemaVersion: 0,
		},
	}
}

type syntheticTestResource struct {
	metaData ResourceMetaData
}

func (r *syntheticTestResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *syntheticTestResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *syntheticTestResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.SyntheticTestConfig()
}

func (r *syntheticTestResource) SetComputedFields(d *schema.ResourceData) {
	//No computed fields defined
}

func (r *syntheticTestResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject, formatter utils.ResourceNameFormatter) error {
	// mz: todo
	// dashboard := obj.(*restapi.CustomDashboard)

	// widgetsBytes, _ := dashboard.Widgets.MarshalJSON()
	// widgets := NormalizeJSONString(string(widgetsBytes))

	// d.Set(CustomDashboardFieldTitle, formatter.UndoFormat(dashboard.Title))
	// d.Set(CustomDashboardFieldFullTitle, dashboard.Title)
	// d.Set(CustomDashboardFieldWidgets, widgets)
	// d.Set(CustomDashboardFieldAccessRule, r.mapAccessRuleToState(dashboard))
	// d.SetId(dashboard.ID)
	return nil
}

func (r *syntheticTestResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	configuration := d.Get(SyntheticTestFieldConfiguration).(string)
	return &restapi.SyntheticTest{
		ID:               d.Id(),
		Label:            "",
		Description:      "",
		Active:           true,
		ApplicationID:    "",
		Configuration:    json.RawMessage(configuration),
		CustomProperties: map[string]string{},
		Locations:        ReadStringSetParameterFromResource(d, SyntheticTestFieldLocations),
		PlaybackMode:     "",
		TestFrequency:    120,
	}, nil
}
