package instana

import (
	"errors"
	"fmt"
	"log"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/hashicorp/terraform/terraform"
)

const (
	//CustomEventSpecificationFieldName constant value for the schema field name
	CustomEventSpecificationFieldName = "name"
	//CustomEventSpecificationFieldFullName constant value for the schema field full_name. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level
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

var customEventSpecificationSchemaName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "Configures the name of the custom event specification",
}
var customEventSpecificationSchemaFullName = &schema.Schema{
	Type:        schema.TypeString,
	Computed:    true,
	Description: "The computed full name of the custom event specification. The field contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
}
var customEventSpecificationSchemaQuery = &schema.Schema{
	Type:        schema.TypeString,
	Required:    false,
	Optional:    true,
	Description: "Configures the dynamic focus query for the custom event specification",
}
var customEventSpecificationSchemaTriggering = &schema.Schema{
	Type:        schema.TypeBool,
	Default:     false,
	Optional:    true,
	Description: "Configures the custom event specification should trigger an incident",
}
var customEventSpecificationSchemaDescription = &schema.Schema{
	Type:        schema.TypeString,
	Required:    false,
	Optional:    true,
	Description: "Configures the description text of the custom event specification",
}
var customEventSpecificationSchemaExpirationTime = &schema.Schema{
	Type:        schema.TypeInt,
	Required:    false,
	Optional:    true,
	Description: "Configures the expiration time (grace period) to wait before the issue is closed",
}
var customEventSpecificationSchemaEnabled = &schema.Schema{
	Type:        schema.TypeBool,
	Default:     true,
	Optional:    true,
	Description: "Configures if the custom event specification is enabled or not",
}
var customEventSpecificationSchemaDownstreamIntegrationIds = &schema.Schema{
	Type:     schema.TypeList,
	Required: false,
	Optional: true,
	MinItems: 0,
	MaxItems: 16,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Description: "Configures the list of integration ids which should be used for downstream reporting",
}
var customEventSpecificationSchemaDownstreamBroadcastToAllAlertingConfigs = &schema.Schema{
	Type:        schema.TypeBool,
	Default:     true,
	Optional:    true,
	Description: "Configures the downstream reporting should be sent to all integrations",
}
var customEventSpecificationSchemaRuleSeverity = &schema.Schema{
	Type:         schema.TypeString,
	Required:     true,
	ValidateFunc: validation.StringInSlice([]string{restapi.SeverityWarning.GetTerraformRepresentation(), restapi.SeverityCritical.GetTerraformRepresentation()}, false),
	Description:  "Configures the severity of the rule of the custom event specification",
}

var defaultCustomEventSchemaFieldsV0 = map[string]*schema.Schema{
	CustomEventSpecificationFieldName:                               customEventSpecificationSchemaName,
	CustomEventSpecificationFieldQuery:                              customEventSpecificationSchemaQuery,
	CustomEventSpecificationFieldTriggering:                         customEventSpecificationSchemaTriggering,
	CustomEventSpecificationFieldDescription:                        customEventSpecificationSchemaDescription,
	CustomEventSpecificationFieldExpirationTime:                     customEventSpecificationSchemaExpirationTime,
	CustomEventSpecificationFieldEnabled:                            customEventSpecificationSchemaEnabled,
	CustomEventSpecificationDownstreamIntegrationIds:                customEventSpecificationSchemaDownstreamIntegrationIds,
	CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs: customEventSpecificationSchemaDownstreamBroadcastToAllAlertingConfigs,
	CustomEventSpecificationRuleSeverity:                            customEventSpecificationSchemaRuleSeverity,
}

var defaultCustomEventSchemaFields = map[string]*schema.Schema{
	CustomEventSpecificationFieldName:                               customEventSpecificationSchemaName,
	CustomEventSpecificationFieldFullName:                           customEventSpecificationSchemaFullName,
	CustomEventSpecificationFieldQuery:                              customEventSpecificationSchemaQuery,
	CustomEventSpecificationFieldTriggering:                         customEventSpecificationSchemaTriggering,
	CustomEventSpecificationFieldDescription:                        customEventSpecificationSchemaDescription,
	CustomEventSpecificationFieldExpirationTime:                     customEventSpecificationSchemaExpirationTime,
	CustomEventSpecificationFieldEnabled:                            customEventSpecificationSchemaEnabled,
	CustomEventSpecificationDownstreamIntegrationIds:                customEventSpecificationSchemaDownstreamIntegrationIds,
	CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs: customEventSpecificationSchemaDownstreamBroadcastToAllAlertingConfigs,
	CustomEventSpecificationRuleSeverity:                            customEventSpecificationSchemaRuleSeverity,
}

//Keep this
func mergeSchemaMap(mapA map[string]*schema.Schema, mapB map[string]*schema.Schema) map[string]*schema.Schema {
	mergedMap := make(map[string]*schema.Schema)

	for k, v := range mapA {
		mergedMap[k] = v
	}
	for k, v := range mapB {
		mergedMap[k] = v
	}

	return mergedMap
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
		response, err := instanaAPI.CustomEventSpecifications().GetOne(specID)
		if err != nil {
			if err == restapi.ErrEntityNotFound {
				d.SetId("")
				return nil
			}
			return err
		}
		return updateCustomEventSpecificationState(d, response.(restapi.CustomEventSpecification), providerMeta.ResourceNameFormatter, ruleSpecificMappingFunc)
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
		response, err := instanaAPI.CustomEventSpecifications().Upsert(spec)
		if err != nil {
			return err
		}
		return updateCustomEventSpecificationState(d, response.(restapi.CustomEventSpecification), providerMeta.ResourceNameFormatter, ruleSpecificResourceMappingFunc)
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

//Keep this
func createCustomEventSpecificationFromResourceData(d *schema.ResourceData, formatter utils.ResourceNameFormatter, ruleSpecificationMapFunc func(*schema.ResourceData) (restapi.RuleSpecification, error)) (restapi.CustomEventSpecification, error) {
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

//Keep this
func computeFullCustomEventNameString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(CustomEventSpecificationFieldName) {
		return formatter.Format(d.Get(CustomEventSpecificationFieldName).(string))
	}
	return d.Get(CustomEventSpecificationFieldFullName).(string)
}

func updateCustomEventSpecificationState(d *schema.ResourceData, spec restapi.CustomEventSpecification, formatter utils.ResourceNameFormatter, ruleSpecificMappingFunc func(*schema.ResourceData, restapi.CustomEventSpecification) error) error {
	d.Set(CustomEventSpecificationFieldFullName, spec.Name)
	d.Set(CustomEventSpecificationFieldQuery, spec.Query)
	d.Set(CustomEventSpecificationFieldEntityType, spec.EntityType)
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

//Keep this
func updateStateForBasicCustomEventSpecification(d *schema.ResourceData, spec restapi.CustomEventSpecification, ruleSpecificMappingFunc func(*schema.ResourceData, restapi.CustomEventSpecification) error) error {
	d.Set(CustomEventSpecificationFieldFullName, spec.Name)
	d.Set(CustomEventSpecificationFieldQuery, spec.Query)
	d.Set(CustomEventSpecificationFieldEntityType, spec.EntityType)
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

//Keep this
func migrateCustomEventConfigFullNameInStateFromV0toV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState[CustomEventSpecificationFieldFullName] = rawState[CustomEventSpecificationFieldName]
	return rawState, nil
}
