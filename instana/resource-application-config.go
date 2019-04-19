package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/filterexpression"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
	//ApplicationConfigFieldLabel const for the lable field of the application config
	ApplicationConfigFieldLabel = "label"
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

//CreateApplicationConfig defines the create operation for the resource instana_application_config
func CreateApplicationConfig(d *schema.ResourceData, meta interface{}) error {
	d.SetId(RandomID())
	return UpdateApplicationConfig(d, meta)
}

//ReadApplicationConfig defines the read operation for the resource instana_application_config
func ReadApplicationConfig(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
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
	return updateApplicationConfigState(d, applicationConfig)
}

//UpdateApplicationConfig defines the update operation for the resource instana_application_config
func UpdateApplicationConfig(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	applicationConfig, err := createApplicationConfigFromResourceData(d)
	if err != nil {
		return err
	}
	updatedApplicationConfig, err := instanaAPI.ApplicationConfigs().Upsert(applicationConfig)
	if err != nil {
		return err
	}
	return updateApplicationConfigState(d, updatedApplicationConfig)
}

//DeleteApplicationConfig defines the delete operation for the resource instana_application_config
func DeleteApplicationConfig(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	applicationConfig, err := createApplicationConfigFromResourceData(d)
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

func createApplicationConfigFromResourceData(d *schema.ResourceData) (restapi.ApplicationConfig, error) {
	matchSpecification, err := convertExpressionStringToAPIModel(d.Get(ApplicationConfigFieldMatchSpecification).(string))
	if err != nil {
		return restapi.ApplicationConfig{}, err
	}
	return restapi.ApplicationConfig{
		ID:                 d.Id(),
		Label:              d.Get(ApplicationConfigFieldLabel).(string),
		Scope:              d.Get(ApplicationConfigFieldScope).(string),
		MatchSpecification: matchSpecification,
	}, nil
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

func updateApplicationConfigState(d *schema.ResourceData, applicationConfig restapi.ApplicationConfig) error {
	normalizedExpressionString, err := convertAPIModelToNormalizedStringRepresentation(applicationConfig.MatchSpecification.(restapi.MatchExpression))
	if err != nil {
		return err
	}

	d.Set(ApplicationConfigFieldLabel, applicationConfig.Label)
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
