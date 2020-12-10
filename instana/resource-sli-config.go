package instana

import (
	"fmt"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

//ResourceInstanaSliConfig the name of the terraform-provider-instana resource to manage SLI configurations
const ResourceInstanaSliConfig = "instana_sli_config"

const (
	//SliConfigFieldName constant value for the schema field name
	SliConfigFieldName = "name"
	//SliConfigFieldFullName constant value for schema field full_name
	SliConfigFieldFullName = "full_name"
	//SliConfigFieldInitialEvaluationTimestamp constant value for the schema field initial_evaluation_timestamp
	SliConfigFieldInitialEvaluationTimestamp = "initial_evaluation_timestamp"
	//SliConfigFieldMetricConfiguration constant value for the schema field metric_configuration
	SliConfigFieldMetricConfiguration = "metric_configuration"
	//SliConfigFieldMetricName constant value for the schema field metric_configuration.metric_name
	SliConfigFieldMetricName = "metric_name"
	//SliConfigFieldMetricAggregation constant value for the schema field metric_configuration.aggregation
	SliConfigFieldMetricAggregation = "aggregation"
	//SliConfigFieldMetricThreshold constant value for the schema field metric_configuration.threshold
	SliConfigFieldMetricThreshold = "threshold"
	//SliConfigFieldSliEntity constant value for the schema field sli_entity
	SliConfigFieldSliEntity = "sli_entity"
	//SliConfigFieldSliType constant value for the schema field sli_entity.type
	SliConfigFieldSliType = "type"
	//SliConfigFieldApplicationID constant value for the schema field sli_entity.application_id
	SliConfigFieldApplicationID = "application_id"
	//SliConfigFieldServiceID constant value for the schema field sli_entity.service_id
	SliConfigFieldServiceID = "service_id"
	//SliConfigFieldEndpointID constant value for the schema field sli_entity.endpoint_id
	SliConfigFieldEndpointID = "endpoint_id"
	//SliConfigFieldBoundaryScope constant value for the schema field sli_entity.boundary_scope
	SliConfigFieldBoundaryScope = "boundary_scope"
)

var (
	//SliConfigName schema field definition of instana_sli_config field name
	SliConfigName = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringLenBetween(0, 256),
		Description:  "The name of the SLI config",
	}

	//SliConfigFullName schema field definition of instana_sli_config field full_name
	SliConfigFullName = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full name of the SLI config. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
	}

	//SliConfigInitialEvaluationTimestamp schema field definition of instana_sli_config field initial_evaluation_timestamp
	SliConfigInitialEvaluationTimestamp = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Initial evaluation timestamp for the SLI config",
	}

	//SliConfigMetricConfiguration schema field definition of instana_sli_config field metric_configuration
	SliConfigMetricConfiguration = &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "Metric configuration for the SLI config",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				//SliConfigFieldMetricName nested schema field definition of instana_sli_config field metric_configuration.metric_name
				SliConfigFieldMetricName: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The metric name for the metric configuration",
				},
				//SliConfigFieldMetricAggregation nested schema field definition of instana_sli_config field metric_configuration.aggregation
				SliConfigFieldMetricAggregation: {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99", "DISTINCT_COUNT"}, true),
					Description:  "The aggregation type for the metric configuration (SUM, MEAN, MAX, MIN, P25, P50, P75, P90, P95, P98, P99, DISTINCT_COUNT)",
				},
				//SliConfigFieldMetricThreshold nested schema field definition of instana_sli_config field metric_configuration.threshold
				SliConfigFieldMetricThreshold: {
					Type:        schema.TypeFloat,
					Required:    true,
					Description: "The threshold for the metric configuration",
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(float64)
						if v <= 0 {
							errs = append(errs, fmt.Errorf("Metric threshold must be higher than 0"))
						}
						return
					},
				},
			},
		},
	}

	//SliConfigSliEntity schema field definition of instana_sli_config field sli_entity
	SliConfigSliEntity = &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "The entity to use for the SLI config",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				//SliConfigFieldSliType nested schema field definition of instana_sli_config field sli_entity.sli_type
				SliConfigFieldSliType: {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"application", "custom", "availability"}, true),
					Description:  "The entity type (application, custom, availability)",
				},
				//SliConfigFieldApplicationId nested schema field definition of instana_sli_config field sli_entity.application_id
				SliConfigFieldApplicationID: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The application ID of the entity",
				},
				//SliConfigFieldServiceId nested schema field definition of instana_sli_config field sli_entity.service_id
				SliConfigFieldServiceID: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The service ID of the entity",
				},
				//SliConfigFieldEndpointId nested schema field definition of instana_sli_config field sli_entity.endpoint_id
				SliConfigFieldEndpointID: {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The endpoint ID of the entity",
				},
				//SliConfigFieldBoundaryScope nested schema field definition of instana_sli_config field sli_entity.boundary_scope
				SliConfigFieldBoundaryScope: {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"ALL", "INBOUND"}, true),
					Description:  "The boundary scope for the entity configuration (ALL, INBOUND)",
				},
			},
		},
	}
)

//NewSliConfigResourceHandle creates the resource handle for SLI configuration
func NewSliConfigResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceInstanaSliConfig,
		Schema: map[string]*schema.Schema{
			SliConfigFieldName:                       SliConfigName,
			SliConfigFieldFullName:                   SliConfigFullName,
			SliConfigFieldInitialEvaluationTimestamp: SliConfigInitialEvaluationTimestamp,
			SliConfigFieldMetricConfiguration:        SliConfigMetricConfiguration,
			SliConfigFieldSliEntity:                  SliConfigSliEntity,
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    sliConfigSchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: sliConfigStateUpgradeV0,
				Version: 0,
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.SliConfigs() },
		UpdateState:          updateStateForSliConfig,
		MapStateToDataObject: mapStateToDataObjectForSliConfig,
	}
}

func updateStateForSliConfig(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	sliConfig := obj.(restapi.SliConfig)

	metricConfiguration := map[string]string{
		SliConfigFieldMetricName:        sliConfig.MetricConfiguration.Name,
		SliConfigFieldMetricAggregation: sliConfig.MetricConfiguration.Aggregation,
		SliConfigFieldMetricThreshold:   fmt.Sprintf("%f", sliConfig.MetricConfiguration.Threshold),
	}

	sliEntity := map[string]string{
		SliConfigFieldSliType:       sliConfig.SliEntity.Type,
		SliConfigFieldApplicationID: sliConfig.SliEntity.ApplicationID,
		SliConfigFieldServiceID:     sliConfig.SliEntity.ServiceID,
		SliConfigFieldEndpointID:    sliConfig.SliEntity.EndpointID,
		SliConfigFieldBoundaryScope: sliConfig.SliEntity.BoundaryScope,
	}

	d.Set(SliConfigFieldFullName, sliConfig.Name)
	d.Set(SliConfigFieldInitialEvaluationTimestamp, sliConfig.InitialEvaluationTimestamp)
	d.Set(SliConfigFieldMetricConfiguration, metricConfiguration)
	d.Set(SliConfigFieldSliEntity, sliEntity)

	d.SetId(sliConfig.ID)
	return nil
}

func mapStateToDataObjectForSliConfig(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	metricConfigurationsList := d.Get(SliConfigFieldMetricConfiguration).([]interface{})
	var metricConfiguration restapi.MetricConfiguration
	if len(metricConfigurationsList) > 0 {
		metricConfigurationState := metricConfigurationsList[0].(map[string]interface{})
		if len(metricConfigurationState) != 0 {
			metricConfiguration = restapi.MetricConfiguration{
				Name:        metricConfigurationState[SliConfigFieldMetricName].(string),
				Aggregation: metricConfigurationState[SliConfigFieldMetricAggregation].(string),
				Threshold:   metricConfigurationState[SliConfigFieldMetricThreshold].(float64),
			}
		}
	}

	sliEntitiesList := d.Get(SliConfigFieldSliEntity).([]interface{})
	var sliEntity restapi.SliEntity
	if len(sliEntitiesList) > 0 {
		sliEntitiesState := sliEntitiesList[0].(map[string]interface{})
		if len(sliEntitiesState) != 0 {
			sliEntity = restapi.SliEntity{
				Type:          sliEntitiesState[SliConfigFieldSliType].(string),
				ApplicationID: sliEntitiesState[SliConfigFieldApplicationID].(string),
				ServiceID:     sliEntitiesState[SliConfigFieldServiceID].(string),
				EndpointID:    sliEntitiesState[SliConfigFieldEndpointID].(string),
				BoundaryScope: sliEntitiesState[SliConfigFieldBoundaryScope].(string),
			}
		}
	}

	name := computeFullSliConfigNameString(d, formatter)
	return restapi.SliConfig{
		ID:                         d.Id(),
		Name:                       name,
		InitialEvaluationTimestamp: d.Get(SliConfigFieldInitialEvaluationTimestamp).(int),
		MetricConfiguration:        metricConfiguration,
		SliEntity:                  sliEntity,
	}, nil
}

func computeFullSliConfigNameString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(SliConfigFieldName) {
		return formatter.Format(d.Get(SliConfigFieldName).(string))
	}
	return d.Get(SliConfigFieldFullName).(string)
}

func sliConfigSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			SliConfigFieldName:                       SliConfigName,
			SliConfigFieldInitialEvaluationTimestamp: SliConfigInitialEvaluationTimestamp,
			SliConfigFieldMetricConfiguration:        SliConfigMetricConfiguration,
			SliConfigFieldSliEntity:                  SliConfigSliEntity,
		},
	}
}

func sliConfigStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	return rawState, nil
}
