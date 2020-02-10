package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

//ResourceInstanaApplicationConfig the name of the terraform-provider-instana resource to manage application config
const ResourceInstanaApplicationConfig = "instana_application_config"

// matchSpecification := binaryOperation | tagMatcherExpression
// binaryOperation := matchSpecification conjunction matchSpecification
// conjunction := AND | OR
// tagMatcherExpreassion := key value operator
// operator := EQUALS | NOT_EQUAL | CONTAINS | NOT_CONTAIN | NOT_EMPTY

const (
	//ApplicationConfigScopeIncludeNoDownstream constant for the scope INCLUDE_NO_DOWNSTREAM
	ApplicationConfigScopeIncludeNoDownstream = "INCLUDE_NO_DOWNSTREAM"
	//ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging constant for the scope INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING
	ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging = "INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING"
	//ApplicationConfigScopeIncludeAllDownstream constant for the scope INCLUDE_ALL_DOWNSTREAM
	ApplicationConfigScopeIncludeAllDownstream = "INCLUDE_ALL_DOWNSTREAM"
)

const (
	//ApplicationConfigFieldLabel const for the label field of the application config
	ApplicationConfigFieldLabel = "label"
	//ApplicationConfigFieldFullLabel const for the full label field of the application config. The field is computed and contains the label which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level
	ApplicationConfigFieldFullLabel = "full_label"
	//ApplicationConfigFieldScope const for the scope field of the application config
	ApplicationConfigFieldScope = "scope"
	//ApplicationConfigFieldMatchSpecification const for the match_specification field of the application config
	ApplicationConfigFieldMatchSpecification = "match_specification"
)

//NewApplicationConfigResourceHandle creates a new instance of the ResourceHandle for application configs
func NewApplicationConfigResourceHandle() ResourceHandle {
	return &applicationConfigResourceHandle{}
}

type applicationConfigResourceHandle struct{}

func (h *applicationConfigResourceHandle) GetResourceFrom(api restapi.InstanaAPI) restapi.RestResource {
	return api.ApplicationConfigs()
}

func (h *applicationConfigResourceHandle) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		ApplicationConfigFieldLabel: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The label of the application config",
		},
		ApplicationConfigFieldFullLabel: &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The the full label field of the application config. The field is computed and contains the label which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
		},
		ApplicationConfigFieldScope: &schema.Schema{
			Type:         schema.TypeString,
			Required:     false,
			Optional:     true,
			Default:      ApplicationConfigScopeIncludeNoDownstream,
			ValidateFunc: validation.StringInSlice([]string{ApplicationConfigScopeIncludeNoDownstream, ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging, ApplicationConfigScopeIncludeAllDownstream}, false),
			Description:  "The scope of the application config",
		},
		ApplicationConfigFieldMatchSpecification: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "The match specification of the application config",
		},
	}
}

func (h *applicationConfigResourceHandle) SchemaVersion() int {
	return 1
}

func (h *applicationConfigResourceHandle) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    applicationConfigSchemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: applicationConfigStateUpgradeV0,
			Version: 0,
		},
	}
}

func (h *applicationConfigResourceHandle) ResourceName() string {
	return ResourceInstanaApplicationConfig
}

func (h *applicationConfigResourceHandle) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	applicationConfig := obj.(restapi.ApplicationConfig)
	normalizedExpressionString, err := h.convertAPIModelToNormalizedStringRepresentation(applicationConfig.MatchSpecification.(restapi.MatchExpression))
	if err != nil {
		return err
	}

	d.Set(ApplicationConfigFieldFullLabel, applicationConfig.Label)
	d.Set(ApplicationConfigFieldScope, applicationConfig.Scope)
	d.Set(ApplicationConfigFieldMatchSpecification, normalizedExpressionString)

	d.SetId(applicationConfig.ID)
	return nil
}

func (h *applicationConfigResourceHandle) convertAPIModelToNormalizedStringRepresentation(input restapi.MatchExpression) (string, error) {
	mapper := filterexpression.NewMapper()
	expr, err := mapper.FromAPIModel(input)
	if err != nil {
		return "", err
	}
	return expr.Render(), nil
}

func (h *applicationConfigResourceHandle) ConvertStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	matchSpecification, err := h.convertExpressionStringToAPIModel(d.Get(ApplicationConfigFieldMatchSpecification).(string))
	if err != nil {
		return restapi.ApplicationConfig{}, err
	}

	label := h.computeFullApplicationConfigLabelString(d, formatter)
	return restapi.ApplicationConfig{
		ID:                 d.Id(),
		Label:              label,
		Scope:              d.Get(ApplicationConfigFieldScope).(string),
		MatchSpecification: matchSpecification,
	}, nil
}

func (h *applicationConfigResourceHandle) convertExpressionStringToAPIModel(input string) (restapi.MatchExpression, error) {
	parser := filterexpression.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := filterexpression.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (h *applicationConfigResourceHandle) computeFullApplicationConfigLabelString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(ApplicationConfigFieldLabel) {
		return formatter.Format(d.Get(ApplicationConfigFieldLabel).(string))
	}
	return d.Get(ApplicationConfigFieldFullLabel).(string)
}

func applicationConfigSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			ApplicationConfigFieldLabel: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The label of the application config",
			},
			ApplicationConfigFieldScope: &schema.Schema{
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				Default:      ApplicationConfigScopeIncludeNoDownstream,
				ValidateFunc: validation.StringInSlice([]string{ApplicationConfigScopeIncludeNoDownstream, ApplicationConfigScopeIncludeImmediateDownstreamDatabaseAndMessaging, ApplicationConfigScopeIncludeAllDownstream}, false),
				Description:  "The scope of the application config",
			},
			ApplicationConfigFieldMatchSpecification: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The match specification of the application config",
			},
		},
	}
}

func applicationConfigStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState[ApplicationConfigFieldFullLabel] = rawState[ApplicationConfigFieldLabel]
	return rawState, nil
}
