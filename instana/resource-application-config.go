package instana

import (
	"context"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaApplicationConfig the name of the terraform-provider-instana resource to manage application config
const ResourceInstanaApplicationConfig = "instana_application_config"

const (
	//ApplicationConfigFieldLabel const for the label field of the application config
	ApplicationConfigFieldLabel = "label"
	//ApplicationConfigFieldFullLabel const for the full label field of the application config. The field is computed and contains the label which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level
	ApplicationConfigFieldFullLabel = "full_label"
	//ApplicationConfigFieldScope const for the scope field of the application config
	ApplicationConfigFieldScope = "scope"
	//ApplicationConfigFieldBoundaryScope const for the boundary_scope field of the application config
	ApplicationConfigFieldBoundaryScope = "boundary_scope"
	//ApplicationConfigFieldTagFilter const for the tag_filter field of the application config
	ApplicationConfigFieldTagFilter = "tag_filter"
)

var (
	//ApplicationConfigLabel schema for the application config field label
	ApplicationConfigLabel = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The label of the application config",
	}
	//ApplicationConfigFullLabel schema for the application config field full_label
	ApplicationConfigFullLabel = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The the full label field of the application config. The field is computed and contains the label which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
	}
	//ApplicationConfigScope schema for the application config field scope
	ApplicationConfigScope = &schema.Schema{
		Type:         schema.TypeString,
		Required:     false,
		Optional:     true,
		Default:      string(restapi.ApplicationConfigScopeIncludeNoDownstream),
		ValidateFunc: validation.StringInSlice(restapi.SupportedApplicationConfigScopes.ToStringSlice(), false),
		Description:  "The scope of the application config",
	}
	//ApplicationConfigBoundaryScope schema for the application config field boundary_scope
	ApplicationConfigBoundaryScope = &schema.Schema{
		Type:         schema.TypeString,
		Required:     false,
		Optional:     true,
		Default:      string(restapi.BoundaryScopeDefault),
		ValidateFunc: validation.StringInSlice(restapi.SupportedApplicationConfigBoundaryScopes.ToStringSlice(), false),
		Description:  "The boundary scope of the application config",
	}
	//ApplicationConfigTagFilter schema for the application config field tag_filter
	ApplicationConfigTagFilter = &schema.Schema{
		Type:             schema.TypeString,
		Optional:         true,
		Description:      "The tag filter expression",
		DiffSuppressFunc: tagFilterDiffSuppressFunc,
		StateFunc:        tagFilterStateFunc,
		ValidateFunc:     tagFilterValidateFunc,
	}
)

// NewApplicationConfigResourceHandle creates a new instance of the ResourceHandle for application configs
func NewApplicationConfigResourceHandle() ResourceHandle[*restapi.ApplicationConfig] {
	return &applicationConfigResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaApplicationConfig,
			Schema: map[string]*schema.Schema{
				ApplicationConfigFieldLabel:         ApplicationConfigLabel,
				ApplicationConfigFieldScope:         ApplicationConfigScope,
				ApplicationConfigFieldBoundaryScope: ApplicationConfigBoundaryScope,
				ApplicationConfigFieldTagFilter:     ApplicationConfigTagFilter,
			},
			SchemaVersion: 4,
		},
	}
}

type applicationConfigResource struct {
	metaData ResourceMetaData
}

func (r *applicationConfigResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *applicationConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV3().CoreConfigSchema().ImpliedType(),
			Upgrade: r.stateUpgradeV3,
			Version: 3,
		},
	}
}

func (r *applicationConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.ApplicationConfig] {
	return api.ApplicationConfigs()
}

func (r *applicationConfigResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *applicationConfigResource) UpdateState(d *schema.ResourceData, applicationConfig *restapi.ApplicationConfig) error {
	data := make(map[string]interface{})
	if applicationConfig.TagFilterExpression != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(applicationConfig.TagFilterExpression)
		if err != nil {
			return err
		}
		data[ApplicationConfigFieldTagFilter] = normalizedTagFilterString
	}
	data[ApplicationConfigFieldLabel] = applicationConfig.Label
	data[ApplicationConfigFieldScope] = string(applicationConfig.Scope)
	data[ApplicationConfigFieldBoundaryScope] = string(applicationConfig.BoundaryScope)

	d.SetId(applicationConfig.ID)
	return tfutils.UpdateState(d, data)
}

func (r *applicationConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.ApplicationConfig, error) {
	var tagFilter *restapi.TagFilter
	var err error

	if tagFilterString, ok := d.GetOk(ApplicationConfigFieldTagFilter); ok {
		tagFilter, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
		if err != nil {
			return &restapi.ApplicationConfig{}, err
		}
	}

	return &restapi.ApplicationConfig{
		ID:                  d.Id(),
		Label:               d.Get(ApplicationConfigFieldLabel).(string),
		Scope:               restapi.ApplicationConfigScope(d.Get(ApplicationConfigFieldScope).(string)),
		BoundaryScope:       restapi.BoundaryScope(d.Get(ApplicationConfigFieldBoundaryScope).(string)),
		TagFilterExpression: tagFilter,
	}, nil
}

func (r *applicationConfigResource) mapTagFilterStringToAPIModel(input string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (r *applicationConfigResource) stateUpgradeV3(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[ApplicationConfigFieldFullLabel]; ok {
		state[ApplicationConfigFieldLabel] = state[ApplicationConfigFieldFullLabel]
		delete(state, ApplicationConfigFieldFullLabel)
	}
	delete(state, "match_specification")
	return state, nil
}

func (r *applicationConfigResource) schemaV3() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			ApplicationConfigFieldLabel:         ApplicationConfigLabel,
			ApplicationConfigFieldFullLabel:     ApplicationConfigFullLabel,
			ApplicationConfigFieldScope:         ApplicationConfigScope,
			ApplicationConfigFieldBoundaryScope: ApplicationConfigBoundaryScope,
			"match_specification": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The match specification of the application config",
				Deprecated:  fmt.Sprintf("match_specification is deprecated. Please migrate to %s", ApplicationConfigFieldTagFilter),
			},
			ApplicationConfigFieldTagFilter: ApplicationConfigTagFilter,
		},
	}
}
