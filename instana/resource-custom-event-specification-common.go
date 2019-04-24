package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	//CustomEventSpecificationFieldName constant value for the schema field name
	CustomEventSpecificationFieldName = "name"
	//CustomEventSpecificationFieldEntityType constant value for the schema field entity type
	CustomEventSpecificationFieldEntityType = "entity_type"
	//CustomEventSpecificationFieldQuery constant value for the schema field query
	CustomEventSpecificationFieldQuery = "query"
	//CustomEventSpecificationFieldTriggering constant value for the schema field triggering
	CustomEventSpecificationFieldTriggering = "triggering"
	//CustomEventSpecificationFieldDescription constant value for the schema field description
	CustomEventSpecificationFieldDescription = "description"
	//CustomEventSpecificationFieldExpirationTime constant value for the schema field expiration_time
	CustomEventSpecificationFieldExpirationTime = "expiration_time"
	//CustomEventSpecificationFieldEnabled constant value for the schema field enabled
	CustomEventSpecificationFieldEnabled = "enabled"

	ruleFieldPrefix = "rule_"

	//CustomEventSpecificationRuleSeverity constant value for the schema field severity of a rule specification
	CustomEventSpecificationRuleSeverity = ruleFieldPrefix + "severity"

	downstreamFieldPrefix = "downstream_"

	//CustomEventSpecificationDownstreamIntegrationIds constant value for the schema field integration_ids of a event specification downstream
	CustomEventSpecificationDownstreamIntegrationIds = downstreamFieldPrefix + "integration_ids"
	//CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs constant value for the schema field broadcast_to_all_alerting_configs of a event specification downstream
	CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs = downstreamFieldPrefix + "broadcast_to_all_alerting_configs"
)

func createCustomEventSpecificationSchema(ruleSpecificSchemaFields map[string]*schema.Schema) map[string]*schema.Schema {
	defaultMap := map[string]*schema.Schema{
		CustomEventSpecificationFieldName: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Configures the name of the custom event specification",
		},
		CustomEventSpecificationFieldEntityType: &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			Description: "Configures the entity type of the custom event specification",
		},
		CustomEventSpecificationFieldQuery: &schema.Schema{
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
			Description: "Configures the dynamic focus query for the custom event specification",
		},
		CustomEventSpecificationFieldTriggering: &schema.Schema{
			Type:        schema.TypeBool,
			Default:     false,
			Optional:    true,
			Description: "Configures the custom event specification should trigger an incident",
		},
		CustomEventSpecificationFieldDescription: &schema.Schema{
			Type:        schema.TypeString,
			Required:    false,
			Optional:    true,
			Description: "Configures the description text of the custom event specification",
		},
		CustomEventSpecificationFieldExpirationTime: &schema.Schema{
			Type:        schema.TypeInt,
			Required:    false,
			Optional:    true,
			Description: "Configures the expiration time (grace period) to wait before the issue is closed",
		},
		CustomEventSpecificationFieldEnabled: &schema.Schema{
			Type:        schema.TypeBool,
			Default:     true,
			Optional:    true,
			Description: "Configures if the custom event specification is enabled or not",
		},
		CustomEventSpecificationDownstreamIntegrationIds: &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			MinItems: 0,
			MaxItems: 16,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Description: "Configures the list of integration ids which should be used for downstream reporting",
		},
		CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs: &schema.Schema{
			Type:        schema.TypeBool,
			Default:     true,
			Optional:    true,
			Description: "Configures the downstream reporting should be sent to all integrations",
		},
		CustomEventSpecificationRuleSeverity: &schema.Schema{
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{restapi.SeverityWarning.GetTerraformRepresentation(), restapi.SeverityCritical.GetTerraformRepresentation()}, false),
			Description:  "Configures the severity of the rule of the custom event specification",
		},
	}

	for k, v := range ruleSpecificSchemaFields {
		defaultMap[k] = v
	}

	return defaultMap
}

func createCustomEventSpecificationReadFunc(ruleSpecificMappingFunc func(*schema.ResourceData, restapi.CustomEventSpecification) error) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		instanaAPI := meta.(restapi.InstanaAPI)
		specID := d.Id()
		if len(specID) == 0 {
			return errors.New("ID of custom event specification is missing")
		}
		spec, err := instanaAPI.CustomEventSpecifications().GetOne(specID)
		if err != nil {
			if err == restapi.ErrEntityNotFound {
				d.SetId("")
				return nil
			}
			return err
		}
		return updateCustomEventSpecificationState(d, spec, ruleSpecificMappingFunc)
	}
}

func createCustomEventSpecificationCreateFunc(ruleSpecificationMapFunc func(*schema.ResourceData) (restapi.RuleSpecification, error), ruleSpecificResourceMappingFunc func(*schema.ResourceData, restapi.CustomEventSpecification) error) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		updateFunc := createCustomEventSpecificationUpdateFunc(ruleSpecificationMapFunc, ruleSpecificResourceMappingFunc)

		d.SetId(RandomID())
		return updateFunc(d, meta)
	}
}

func createCustomEventSpecificationUpdateFunc(ruleSpecificationMapFunc func(*schema.ResourceData) (restapi.RuleSpecification, error), ruleSpecificResourceMappingFunc func(*schema.ResourceData, restapi.CustomEventSpecification) error) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		instanaAPI := meta.(restapi.InstanaAPI)
		spec, err := createCustomEventSpecificationFromResourceData(d, ruleSpecificationMapFunc)
		if err != nil {
			return err
		}
		updatedCustomEventSpecification, err := instanaAPI.CustomEventSpecifications().Upsert(spec)
		if err != nil {
			return err
		}
		return updateCustomEventSpecificationState(d, updatedCustomEventSpecification, ruleSpecificResourceMappingFunc)
	}
}

func createCustomEventSpecificationDeleteFunc(ruleSpecificationMapFunc func(*schema.ResourceData) (restapi.RuleSpecification, error)) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		instanaAPI := meta.(restapi.InstanaAPI)
		spec, err := createCustomEventSpecificationFromResourceData(d, ruleSpecificationMapFunc)
		if err != nil {
			return err
		}
		err = instanaAPI.CustomEventSpecifications().DeleteByID(spec.ID)
		if err != nil {
			return err
		}
		d.SetId("")
		return nil
	}
}

func createCustomEventSpecificationFromResourceData(d *schema.ResourceData, ruleSpecificationMapFunc func(*schema.ResourceData) (restapi.RuleSpecification, error)) (restapi.CustomEventSpecification, error) {
	apiModel := restapi.CustomEventSpecification{
		ID:             d.Id(),
		Name:           d.Get(CustomEventSpecificationFieldName).(string),
		EntityType:     d.Get(CustomEventSpecificationFieldEntityType).(string),
		Query:          GetStringPointerFromResourceData(d, CustomEventSpecificationFieldQuery),
		Triggering:     d.Get(CustomEventSpecificationFieldTriggering).(bool),
		Description:    GetStringPointerFromResourceData(d, CustomEventSpecificationFieldDescription),
		ExpirationTime: GetIntPointerFromResourceData(d, CustomEventSpecificationFieldExpirationTime),
		Enabled:        d.Get(CustomEventSpecificationFieldEnabled).(bool),
	}

	downstreamIntegrationIds := ReadStringArrayParameterFromResource(d, CustomEventSpecificationDownstreamIntegrationIds)
	if len(downstreamIntegrationIds) > 0 {
		apiModel.Downstream = &restapi.EventSpecificationDownstream{
			IntegrationIds:                downstreamIntegrationIds,
			BroadcastToAllAlertingConfigs: d.Get(CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs).(bool),
		}
	}

	rule, err := ruleSpecificationMapFunc(d)
	if err != nil {
		return apiModel, err
	}
	apiModel.Rules = []restapi.RuleSpecification{rule}
	return apiModel, nil
}

func updateCustomEventSpecificationState(d *schema.ResourceData, spec restapi.CustomEventSpecification, ruleSpecificMappingFunc func(*schema.ResourceData, restapi.CustomEventSpecification) error) error {
	d.Set(CustomEventSpecificationFieldName, spec.Name)
	d.Set(CustomEventSpecificationFieldEntityType, spec.EntityType)
	d.Set(CustomEventSpecificationFieldQuery, spec.Query)
	d.Set(CustomEventSpecificationFieldTriggering, spec.Triggering)
	d.Set(CustomEventSpecificationFieldDescription, spec.Description)
	d.Set(CustomEventSpecificationFieldExpirationTime, spec.ExpirationTime)
	d.Set(CustomEventSpecificationFieldEnabled, spec.Enabled)

	if spec.Downstream != nil {
		d.Set(CustomEventSpecificationDownstreamIntegrationIds, spec.Downstream.IntegrationIds)
		d.Set(CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs, spec.Downstream.BroadcastToAllAlertingConfigs)
	}

	err := ruleSpecificMappingFunc(d, spec)
	if err != nil {
		return err
	}

	d.SetId(spec.ID)
	return nil
}
