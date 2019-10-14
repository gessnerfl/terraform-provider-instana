package instana

import (
	"errors"
	"fmt"
	"log"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/hashicorp/terraform/terraform"
)

const (
	//CustomEventSpecificationFieldName constant value for the schema field name
	CustomEventSpecificationFieldName = "name"
	//CustomEventSpecificationFieldFullName constant value for the schema field full_name. The field is computed and contains the name which is sent to instana. The computation depends on the activation of add_terraform_managed_string at provider level
	CustomEventSpecificationFieldFullName = "full_name"
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
		CustomEventSpecificationFieldFullName: &schema.Schema{
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The computed full name of the custom event specification. The field contains the name which is sent to instana. The computation depends on the activation of add_terraform_managed_string at provider level",
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
			Required: false,
			Optional: true,
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

//CreateMigrateCustomEventConfigStateFunction creates the function for migrating schemas in terraform for the different implementations of custom events
func CreateMigrateCustomEventConfigStateFunction(specificFunctions map[int](func(inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error))) func(v int, inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	return func(v int, inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
		if inst.Empty() {
			log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
			return inst, nil
		}

		switch v {
		case 0:
			log.Println("[INFO] Found Custom Event Config State v0; migrating to v1")
			temp, err := applySpecificFunction(v, inst, meta, specificFunctions)
			if err != nil {
				return temp, err
			}
			return migrateCustomEventConfigNameInStateFromV0toV1(inst)
		default:
			return inst, fmt.Errorf("Unexpected schema version: %d", v)
		}
	}
}

func applySpecificFunction(v int, inst *terraform.InstanceState, meta interface{}, specificFunctions map[int](func(inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error))) (*terraform.InstanceState, error) {
	if specificFunctions != nil && specificFunctions[v] != nil {
		return specificFunctions[v](inst, meta)
	}
	return inst, nil
}

func migrateCustomEventConfigNameInStateFromV0toV1(inst *terraform.InstanceState) (*terraform.InstanceState, error) {
	inst.Attributes[CustomEventSpecificationFieldFullName] = inst.Attributes[CustomEventSpecificationFieldName]
	return inst, nil
}

func createCustomEventSpecificationReadFunc(ruleSpecificMappingFunc func(*schema.ResourceData, restapi.CustomEventSpecification) error) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		providerMeta := meta.(*ProviderMeta)
		instanaAPI := providerMeta.InstanaAPI
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
		return updateCustomEventSpecificationState(d, spec, providerMeta.ResourceNameFormatter, ruleSpecificMappingFunc)
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
		providerMeta := meta.(*ProviderMeta)
		instanaAPI := providerMeta.InstanaAPI
		spec, err := createCustomEventSpecificationFromResourceData(d, providerMeta.ResourceNameFormatter, ruleSpecificationMapFunc)
		if err != nil {
			return err
		}
		updatedCustomEventSpecification, err := instanaAPI.CustomEventSpecifications().Upsert(spec)
		if err != nil {
			return err
		}
		return updateCustomEventSpecificationState(d, updatedCustomEventSpecification, providerMeta.ResourceNameFormatter, ruleSpecificResourceMappingFunc)
	}
}

func createCustomEventSpecificationDeleteFunc(ruleSpecificationMapFunc func(*schema.ResourceData) (restapi.RuleSpecification, error)) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		providerMeta := meta.(*ProviderMeta)
		instanaAPI := providerMeta.InstanaAPI
		spec, err := createCustomEventSpecificationFromResourceData(d, providerMeta.ResourceNameFormatter, ruleSpecificationMapFunc)
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

func createCustomEventSpecificationFromResourceData(d *schema.ResourceData, formatter ResourceNameFormatter, ruleSpecificationMapFunc func(*schema.ResourceData) (restapi.RuleSpecification, error)) (restapi.CustomEventSpecification, error) {
	name := computeFullCustomEventNameString(d, formatter)

	apiModel := restapi.CustomEventSpecification{
		ID:             d.Id(),
		Name:           name,
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

func computeFullCustomEventNameString(d *schema.ResourceData, formatter ResourceNameFormatter) string {
	if d.HasChange(CustomEventSpecificationFieldName) {
		return formatter.Format(d.Get(CustomEventSpecificationFieldName).(string))
	}
	return d.Get(CustomEventSpecificationFieldFullName).(string)
}

func updateCustomEventSpecificationState(d *schema.ResourceData, spec restapi.CustomEventSpecification, formatter ResourceNameFormatter, ruleSpecificMappingFunc func(*schema.ResourceData, restapi.CustomEventSpecification) error) error {
	d.Set(CustomEventSpecificationFieldFullName, spec.Name)
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
