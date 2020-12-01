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
	SliConfigFieldSliName                    = "name"
	SliConfigFieldInitialEvaluationTimestamp = "initial_evaluation_timestamp"
	SliConfigFieldMetricConfiguration        = "metric_configuration"
	SliConfigFieldMetricName                 = "name"
	SliConfigFieldMetricAggregation          = "aggregation"
	SliConfigFieldMetricThreshold            = "threshold"
	SliConfigFieldSliEntity                  = "sli_entity"
	SliConfigFieldSliType                    = "type"
	SliConfigFieldApplicationId              = "application_id"
	SliConfigFieldServiceId                  = "service_id"
	SliConfigFieldEndpointId                 = "endpoint_id"
	SliConfigFieldBoundaryScope              = "boundary_scope"
)

var (
	SliConfigName = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringLenBetween(0, 256),
		Description:  "The name of the SLI config",
	}

	SliConfigInitialEvaluationTimestamp = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     0,
		Description: "Initial evaluation timestamp for the SLI config",
	}

	SliConfigMetricConfiguration = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				SliConfigFieldMetricName: &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				SliConfigFieldMetricAggregation: &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99", "DISTINCT_COUNT"}, true),
				},
				SliConfigFieldMetricThreshold: &schema.Schema{
					Type:     schema.TypeFloat,
					Required: true,
				},
			},
		},
	}

	SliConfigSliEntity = &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				SliConfigFieldSliType: &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"application", "custom", "availability"}, true),
				},
				SliConfigFieldApplicationId: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				SliConfigFieldServiceId: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				SliConfigFieldEndpointId: &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
				},
				SliConfigFieldBoundaryScope: &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"ALL", "INBOUND"}, true),
				},
			},
		},
	}
)

func NewSliConfigResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceInstanaSliConfig,
		Schema: map[string]*schema.Schema{
			SliConfigFieldSliName:                    SliConfigName,
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
		SliConfigFieldApplicationId: sliConfig.SliEntity.ApplicationID,
		SliConfigFieldServiceId:     sliConfig.SliEntity.ServiceID,
		SliConfigFieldEndpointId:    sliConfig.SliEntity.EndpointID,
		SliConfigFieldBoundaryScope: sliConfig.SliEntity.BoundaryScope,
	}

	d.Set(SliConfigFieldSliName, sliConfig.Name)
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
		fmt.Println(sliEntitiesState)
		if len(sliEntitiesState) != 0 {
			sliEntity = restapi.SliEntity{
				Type:          sliEntitiesState[SliConfigFieldSliType].(string),
				ApplicationID: sliEntitiesState[SliConfigFieldApplicationId].(string),
				ServiceID:     sliEntitiesState[SliConfigFieldServiceId].(string),
				EndpointID:    sliEntitiesState[SliConfigFieldEndpointId].(string),
				BoundaryScope: sliEntitiesState[SliConfigFieldBoundaryScope].(string),
			}
		}
	}

	return restapi.SliConfig{
		ID:                         d.Id(),
		Name:                       d.Get(SliConfigFieldSliName).(string),
		InitialEvaluationTimestamp: d.Get(SliConfigFieldInitialEvaluationTimestamp).(int),
		MetricConfiguration:        metricConfiguration,
		SliEntity:                  sliEntity,
	}, nil
}

func sliConfigSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			SliConfigFieldSliName:                    SliConfigName,
			SliConfigFieldInitialEvaluationTimestamp: SliConfigInitialEvaluationTimestamp,
			SliConfigFieldMetricConfiguration:        SliConfigMetricConfiguration,
			SliConfigFieldSliEntity:                  SliConfigSliEntity,
		},
	}
}

func sliConfigStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	return rawState, nil
}
