package instana

import (
	"errors"
	"fmt"
	"log"

	"github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/hashicorp/terraform/terraform"
)

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

//CreateResourceApplicationConfig creates the resource definition for the resource instana_application_config
func CreateResourceApplicationConfig() *schema.Resource {
	return &schema.Resource{
		Create: CreateApplicationConfig,
		Read:   ReadApplicationConfig,
		Update: UpdateApplicationConfig,
		Delete: DeleteApplicationConfig,

		Schema: map[string]*schema.Schema{
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
		},
		SchemaVersion: 1,
		MigrateState:  MigrateApplicationConfigState,
	}
}

//MigrateApplicationConfigState migrates the terraform state from one version to the other
func MigrateApplicationConfigState(v int, inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Found Application Config State v0; migrating to v1")
		return migrateApplicationConfigStateV0toV1(inst)
	default:
		return inst, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func migrateApplicationConfigStateV0toV1(inst *terraform.InstanceState) (*terraform.InstanceState, error) {
	if inst.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return inst, nil
	}
	inst.Attributes[ApplicationConfigFieldFullLabel] = inst.Attributes[ApplicationConfigFieldLabel]
	return inst, nil
}

//CreateApplicationConfig defines the create operation for the resource instana_application_config
func CreateApplicationConfig(d *schema.ResourceData, meta interface{}) error {
	d.SetId(RandomID())
	return UpdateApplicationConfig(d, meta)
}

//ReadApplicationConfig defines the read operation for the resource instana_application_config
func ReadApplicationConfig(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	applicationConfigID := d.Id()
	if len(applicationConfigID) == 0 {
		return errors.New("ID of application config is missing")
	}
	applicationConfig, err := instanaAPI.ApplicationConfigs().GetOne(applicationConfigID)
	if err != nil {
		if err == restapi.ErrEntityNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	return updateApplicationConfigState(d, applicationConfig, providerMeta.ResourceNameFormatter)
}

//UpdateApplicationConfig defines the update operation for the resource instana_application_config
func UpdateApplicationConfig(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	applicationConfig, err := createApplicationConfigFromResourceData(d, providerMeta.ResourceNameFormatter)
	if err != nil {
		return err
	}
	updatedApplicationConfig, err := instanaAPI.ApplicationConfigs().Upsert(applicationConfig)
	if err != nil {
		return err
	}
	return updateApplicationConfigState(d, updatedApplicationConfig, providerMeta.ResourceNameFormatter)
}

//DeleteApplicationConfig defines the delete operation for the resource instana_application_config
func DeleteApplicationConfig(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(*ProviderMeta)
	instanaAPI := providerMeta.InstanaAPI
	applicationConfig, err := createApplicationConfigFromResourceData(d, providerMeta.ResourceNameFormatter)
	if err != nil {
		return err
	}
	err = instanaAPI.ApplicationConfigs().DeleteByID(applicationConfig.ID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createApplicationConfigFromResourceData(d *schema.ResourceData, formatter ResourceNameFormatter) (restapi.ApplicationConfig, error) {
	matchSpecification, err := convertExpressionStringToAPIModel(d.Get(ApplicationConfigFieldMatchSpecification).(string))
	if err != nil {
		return restapi.ApplicationConfig{}, err
	}

	label := computeFullApplicationConfigLabelString(d, formatter)
	return restapi.ApplicationConfig{
		ID:                 d.Id(),
		Label:              label,
		Scope:              d.Get(ApplicationConfigFieldScope).(string),
		MatchSpecification: matchSpecification,
	}, nil
}

func computeFullApplicationConfigLabelString(d *schema.ResourceData, formatter ResourceNameFormatter) string {
	if d.HasChange(ApplicationConfigFieldLabel) {
		return formatter.Format(d.Get(ApplicationConfigFieldLabel).(string))
	}
	return d.Get(ApplicationConfigFieldFullLabel).(string)
}

func convertExpressionStringToAPIModel(input string) (restapi.MatchExpression, error) {
	parser := filterexpression.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := filterexpression.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func updateApplicationConfigState(d *schema.ResourceData, applicationConfig restapi.ApplicationConfig, formatter ResourceNameFormatter) error {
	normalizedExpressionString, err := convertAPIModelToNormalizedStringRepresentation(applicationConfig.MatchSpecification.(restapi.MatchExpression))
	if err != nil {
		return err
	}

	d.Set(ApplicationConfigFieldFullLabel, applicationConfig.Label)
	d.Set(ApplicationConfigFieldScope, applicationConfig.Scope)
	d.Set(ApplicationConfigFieldMatchSpecification, normalizedExpressionString)

	d.SetId(applicationConfig.ID)
	return nil
}

func convertAPIModelToNormalizedStringRepresentation(input restapi.MatchExpression) (string, error) {
	mapper := filterexpression.NewMapper()
	expr, err := mapper.FromAPIModel(input)
	if err != nil {
		return "", err
	}
	return expr.Render(), nil
}
