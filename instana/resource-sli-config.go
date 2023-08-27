package instana

import (
	"context"
	"errors"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaSliConfig the name of the terraform-provider-instana resource to manage SLI configurations
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
	//SliConfigFieldSliEntityApplication constant value for the schema field sli_entity.application
	SliConfigFieldSliEntityApplication = "application"
	//SliConfigFieldSliEntityAvailability constant value for the schema field sli_entity.availability
	SliConfigFieldSliEntityAvailability = "availability"
	//SliConfigFieldSliEntityWebsiteEventBased constant value for the schema field sli_entity.website_event_based
	SliConfigFieldSliEntityWebsiteEventBased = "website_event_based"
	//SliConfigFieldSliEntityWebsiteTimeBased constant value for the schema field sli_entity.website_time_based
	SliConfigFieldSliEntityWebsiteTimeBased = "website_time_based"
	//SliConfigFieldApplicationID constant value for the schema field sli_entity.*.application_id
	SliConfigFieldApplicationID = "application_id"
	//SliConfigFieldServiceID constant value for the schema field sli_entity.*.service_id
	SliConfigFieldServiceID = "service_id"
	//SliConfigFieldEndpointID constant value for the schema field sli_entity.*.endpoint_id
	SliConfigFieldEndpointID = "endpoint_id"
	//SliConfigFieldWebsiteID constant value for the schema field sli_entity.*.website_id
	SliConfigFieldWebsiteID = "website_id"
	//SliConfigFieldBeaconType constant value for the schema field sli_entity.*.beacon_Type
	SliConfigFieldBeaconType = "beacon_type"
	//SliConfigFieldBoundaryScope constant value for the schema field sli_entity.boundary_scope
	SliConfigFieldBoundaryScope = "boundary_scope"
	//SliConfigFieldBadEventFilterExpression constant value for the schema field sli_entity.*.bad_event_filter_expression
	SliConfigFieldBadEventFilterExpression = "bad_event_filter_expression"
	//SliConfigFieldFilterExpression constant value for the schema field sli_entity.*.filter_expression
	SliConfigFieldFilterExpression = "filter_expression"
	//SliConfigFieldGoodEventFilterExpression constant value for the schema field sli_entity.*.good_event_filter_expression
	SliConfigFieldGoodEventFilterExpression = "good_event_filter_expression"
	//SliConfigFieldIncludeInternal constant value for the schema field sli_entity.*.good_event_filter_expression
	SliConfigFieldIncludeInternal = "include_internal"
	//SliConfigFieldIncludeSynthetic constant value for the schema field sli_entity.*.good_event_filter_expression
	SliConfigFieldIncludeSynthetic = "include_synthetic"
)

var sliConfigSliEntityTypeKeys = []string{
	"sli_entity.0.application",
	"sli_entity.0.availability",
	"sli_entity.0.website_event_based",
	"sli_entity.0.website_time_based",
}

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
					ValidateFunc: validation.StringInSlice([]string{"SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99", "P99_9", "P99_99", "DISTRIBUTION", "DISTINCT_COUNT", "SUM_POSITIVE", "PER_SECOND"}, true),
					Description:  "The aggregation type for the metric configuration (SUM, MEAN, MAX, MIN, P25, P50, P75, P90, P95, P98, P99, P99_9, P99_99, DISTRIBUTION, DISTINCT_COUNT, SUM_POSITIVE, PER_SECOND)",
				},
				//SliConfigFieldMetricThreshold nested schema field definition of instana_sli_config field metric_configuration.threshold
				SliConfigFieldMetricThreshold: {
					Type:        schema.TypeFloat,
					Required:    true,
					Description: "The threshold for the metric configuration",
					ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
						v := val.(float64)
						if v <= 0 {
							errs = append(errs, fmt.Errorf("metric threshold must be greater than 0"))
						}
						return
					},
				},
			},
		},
	}

	sliConfigSchemaSliEntityApplicationId = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The application ID of the entity",
	}
	sliConfigSchemaSliEntityServiceId = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The service ID of the entity",
	}
	sliConfigSchemaSliEntityEndpointId = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The endpoint ID of the entity",
	}
	sliConfigSchemaSliEntityWebsiteId = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The website ID of the entity",
	}
	sliConfigSchemaSliEntityBoundaryScope = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"ALL", "INBOUND"}, true),
		Description:  "The boundary scope for the entity configuration (ALL, INBOUND)",
	}
	sliConfigSchemaSliEntityBeaconType = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"pageLoad", "resourceLoad", "httpRequest", "error", "custom", "pageChange"}, true),
		Description:  "The beacon type for the entity configuration (pageLoad, resourceLoad, httpRequest, error, custom, pageChange)",
	}
	sliConfigSchemaIncludeInternal = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether also internal calls are included",
	}
	sliConfigSchemaIncludeSynthetic = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether also synthetic calls are included in the scope or not",
	}

	//SliConfigSliEntity schema field definition of instana_sli_config field sli_entity
	SliConfigSliEntity = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "The entity to use for the SLI config.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				SliConfigFieldSliEntityApplication: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "The SLI entity of type application to use for the SLI config",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SliConfigFieldApplicationID: sliConfigSchemaSliEntityApplicationId,
							SliConfigFieldServiceID:     sliConfigSchemaSliEntityServiceId,
							SliConfigFieldEndpointID:    sliConfigSchemaSliEntityEndpointId,
							SliConfigFieldBoundaryScope: sliConfigSchemaSliEntityBoundaryScope,
						},
					},
					ExactlyOneOf: sliConfigSliEntityTypeKeys,
				},
				SliConfigFieldSliEntityAvailability: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "The SLI entity of type availability to use for the SLI config",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SliConfigFieldApplicationID:             sliConfigSchemaSliEntityApplicationId,
							SliConfigFieldBoundaryScope:             sliConfigSchemaSliEntityBoundaryScope,
							SliConfigFieldBadEventFilterExpression:  RequiredTagFilterExpressionSchema,
							SliConfigFieldGoodEventFilterExpression: RequiredTagFilterExpressionSchema,
							SliConfigFieldIncludeInternal:           sliConfigSchemaIncludeInternal,
							SliConfigFieldIncludeSynthetic:          sliConfigSchemaIncludeSynthetic,
						},
					},
					ExactlyOneOf: sliConfigSliEntityTypeKeys,
				},
				SliConfigFieldSliEntityWebsiteEventBased: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "The SLI entity of type websiteEventBased to use for the SLI config",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SliConfigFieldWebsiteID:                 sliConfigSchemaSliEntityWebsiteId,
							SliConfigFieldBadEventFilterExpression:  RequiredTagFilterExpressionSchema,
							SliConfigFieldGoodEventFilterExpression: RequiredTagFilterExpressionSchema,
							SliConfigFieldBeaconType:                sliConfigSchemaSliEntityBeaconType,
						},
					},
					ExactlyOneOf: sliConfigSliEntityTypeKeys,
				},
				SliConfigFieldSliEntityWebsiteTimeBased: {
					Type:        schema.TypeList,
					Optional:    true,
					Description: "The SLI entity of type websiteTimeBased to use for the SLI config",
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							SliConfigFieldWebsiteID:        sliConfigSchemaSliEntityWebsiteId,
							SliConfigFieldFilterExpression: OptionalTagFilterExpressionSchema,
							SliConfigFieldBeaconType:       sliConfigSchemaSliEntityBeaconType,
						},
					},
					ExactlyOneOf: sliConfigSliEntityTypeKeys,
				},
			},
		},
	}
)

// NewSliConfigResourceHandle creates the resource handle for SLI configuration
func NewSliConfigResourceHandle() ResourceHandle[*restapi.SliConfig] {
	return &sliConfigResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaSliConfig,
			Schema: map[string]*schema.Schema{
				SliConfigFieldName:                       SliConfigName,
				SliConfigFieldInitialEvaluationTimestamp: SliConfigInitialEvaluationTimestamp,
				SliConfigFieldMetricConfiguration:        SliConfigMetricConfiguration,
				SliConfigFieldSliEntity:                  SliConfigSliEntity,
			},
			SchemaVersion: 1,
		},
	}
}

type sliConfigResource struct {
	metaData ResourceMetaData
}

func (r *sliConfigResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *sliConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.sliConfigSchemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: r.sliConfigStateUpgradeV0,
			Version: 0,
		},
	}
}

func (r *sliConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.SliConfig] {
	return api.SliConfigs()
}

func (r *sliConfigResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *sliConfigResource) UpdateState(d *schema.ResourceData, sliConfig *restapi.SliConfig) error {
	metricConfiguration := map[string]interface{}{
		SliConfigFieldMetricName:        sliConfig.MetricConfiguration.Name,
		SliConfigFieldMetricAggregation: sliConfig.MetricConfiguration.Aggregation,
		SliConfigFieldMetricThreshold:   sliConfig.MetricConfiguration.Threshold,
	}

	sliEntity, err := r.mapSliEntityToState(sliConfig)
	if err != nil {
		return err
	}

	data := map[string]interface{}{
		SliConfigFieldName:                       sliConfig.Name,
		SliConfigFieldInitialEvaluationTimestamp: sliConfig.InitialEvaluationTimestamp,
		SliConfigFieldMetricConfiguration:        []map[string]interface{}{metricConfiguration},
		SliConfigFieldSliEntity:                  []interface{}{sliEntity},
	}

	d.SetId(sliConfig.ID)
	return tfutils.UpdateState(d, data)
}

func (r *sliConfigResource) mapSliEntityToState(sliConfig *restapi.SliConfig) (map[string]interface{}, error) {
	if sliConfig.SliEntity.Type == "application" {
		return r.mapSliApplicationEntityToState(sliConfig.SliEntity)
	} else if sliConfig.SliEntity.Type == "availability" {
		return r.mapSliAvailabilityEntityToState(sliConfig.SliEntity)
	} else if sliConfig.SliEntity.Type == "websiteTimeBased" {
		return r.mapSliWebsiteTimeBasedEntityToState(sliConfig.SliEntity)
	} else if sliConfig.SliEntity.Type == "websiteEventBased" {
		return r.mapSliWebsiteEventBasedEntityToState(sliConfig.SliEntity)
	}
	return nil, fmt.Errorf("unsupported sli entity type %s", sliConfig.SliEntity.Type)
}

func (r *sliConfigResource) mapSliApplicationEntityToState(sliEntity restapi.SliEntity) (map[string]interface{}, error) {
	result := map[string]interface{}{
		SliConfigFieldSliEntityApplication: []interface{}{
			map[string]interface{}{
				SliConfigFieldApplicationID: sliEntity.ApplicationID,
				SliConfigFieldServiceID:     sliEntity.ServiceID,
				SliConfigFieldEndpointID:    sliEntity.EndpointID,
				SliConfigFieldBoundaryScope: sliEntity.BoundaryScope,
			},
		},
	}
	return result, nil
}

func (r *sliConfigResource) mapSliAvailabilityEntityToState(sliEntity restapi.SliEntity) (map[string]interface{}, error) {
	var goodEventFilterExpression *string
	var badEventFilterExpression *string
	var err error
	if sliEntity.GoodEventFilterExpression != nil {
		goodEventFilterExpression, err = tagfilter.MapTagFilterToNormalizedString(sliEntity.GoodEventFilterExpression.(restapi.TagFilterExpressionElement))
		if err != nil {
			return nil, err
		}
	}
	if sliEntity.BadEventFilterExpression != nil {
		badEventFilterExpression, err = tagfilter.MapTagFilterToNormalizedString(sliEntity.BadEventFilterExpression.(restapi.TagFilterExpressionElement))
		if err != nil {
			return nil, err
		}
	}

	result := map[string]interface{}{
		SliConfigFieldSliEntityAvailability: []interface{}{
			map[string]interface{}{
				SliConfigFieldApplicationID:             sliEntity.ApplicationID,
				SliConfigFieldBoundaryScope:             sliEntity.BoundaryScope,
				SliConfigFieldGoodEventFilterExpression: goodEventFilterExpression,
				SliConfigFieldBadEventFilterExpression:  badEventFilterExpression,
				SliConfigFieldIncludeInternal:           sliEntity.IncludeInternal,
				SliConfigFieldIncludeSynthetic:          sliEntity.IncludeSynthetic,
			},
		},
	}
	return result, nil
}

func (r *sliConfigResource) mapSliWebsiteEventBasedEntityToState(sliEntity restapi.SliEntity) (map[string]interface{}, error) {
	var goodEventFilterExpression *string
	var badEventFilterExpression *string
	var err error
	if sliEntity.GoodEventFilterExpression != nil {
		goodEventFilterExpression, err = tagfilter.MapTagFilterToNormalizedString(sliEntity.GoodEventFilterExpression.(restapi.TagFilterExpressionElement))
		if err != nil {
			return nil, err
		}
	}
	if sliEntity.BadEventFilterExpression != nil {
		badEventFilterExpression, err = tagfilter.MapTagFilterToNormalizedString(sliEntity.BadEventFilterExpression.(restapi.TagFilterExpressionElement))
		if err != nil {
			return nil, err
		}
	}

	result := map[string]interface{}{
		SliConfigFieldSliEntityWebsiteEventBased: []interface{}{
			map[string]interface{}{
				SliConfigFieldWebsiteID:                 sliEntity.WebsiteId,
				SliConfigFieldGoodEventFilterExpression: goodEventFilterExpression,
				SliConfigFieldBadEventFilterExpression:  badEventFilterExpression,
				SliConfigFieldBeaconType:                sliEntity.BeaconType,
			},
		},
	}
	return result, nil
}

func (r *sliConfigResource) mapSliWebsiteTimeBasedEntityToState(sliEntity restapi.SliEntity) (map[string]interface{}, error) {
	var tagFilterExpression *string
	var err error
	if sliEntity.FilterExpression != nil {
		tagFilterExpression, err = tagfilter.MapTagFilterToNormalizedString(sliEntity.FilterExpression.(restapi.TagFilterExpressionElement))
		if err != nil {
			return nil, err
		}
	}

	result := map[string]interface{}{
		SliConfigFieldSliEntityWebsiteTimeBased: []interface{}{
			map[string]interface{}{
				SliConfigFieldWebsiteID:        sliEntity.WebsiteId,
				SliConfigFieldFilterExpression: tagFilterExpression,
				SliConfigFieldBeaconType:       sliEntity.BeaconType,
			},
		},
	}
	return result, nil
}

func (r *sliConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.SliConfig, error) {
	metricConfigurationsStateObject := d.Get(SliConfigFieldMetricConfiguration).([]interface{})
	var metricConfiguration restapi.MetricConfiguration
	if len(metricConfigurationsStateObject) > 0 {
		metricConfiguration = r.mapMetricConfigurationEntityFromState(metricConfigurationsStateObject)
	}

	sliEntitiesStateObject := d.Get(SliConfigFieldSliEntity).([]interface{})
	var sliEntity restapi.SliEntity
	var err error
	if len(sliEntitiesStateObject) == 1 {
		sliEntity, err = r.mapSliEntityListFromState(sliEntitiesStateObject[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("exactly one sli entity configuration is required")
	}

	return &restapi.SliConfig{
		ID:                         d.Id(),
		Name:                       d.Get(SliConfigFieldName).(string),
		InitialEvaluationTimestamp: d.Get(SliConfigFieldInitialEvaluationTimestamp).(int),
		MetricConfiguration:        metricConfiguration,
		SliEntity:                  sliEntity,
	}, nil
}

func (r *sliConfigResource) mapMetricConfigurationEntityFromState(stateObject []interface{}) restapi.MetricConfiguration {
	metricConfigurationState := stateObject[0].(map[string]interface{})
	if len(metricConfigurationState) != 0 {
		return restapi.MetricConfiguration{
			Name:        metricConfigurationState[SliConfigFieldMetricName].(string),
			Aggregation: metricConfigurationState[SliConfigFieldMetricAggregation].(string),
			Threshold:   metricConfigurationState[SliConfigFieldMetricThreshold].(float64),
		}
	}
	return restapi.MetricConfiguration{}
}

func (r *sliConfigResource) mapSliEntityListFromState(stateObject map[string]interface{}) (restapi.SliEntity, error) {
	if len(stateObject) > 0 {
		if details, ok := stateObject[SliConfigFieldSliEntityApplication]; ok && r.isSliEntitySet(details) {
			return r.mapSliEntityApplicationFromState(details.([]interface{})[0].(map[string]interface{}))
		}
		if details, ok := stateObject[SliConfigFieldSliEntityAvailability]; ok && r.isSliEntitySet(details) {
			return r.mapSliEntityAvailabilityFromState(details.([]interface{})[0].(map[string]interface{}))
		}
		if details, ok := stateObject[SliConfigFieldSliEntityWebsiteEventBased]; ok && r.isSliEntitySet(details) {
			return r.mapSliEntityWebsiteEventBasedFromState(details.([]interface{})[0].(map[string]interface{}))
		}
		if details, ok := stateObject[SliConfigFieldSliEntityWebsiteTimeBased]; ok && r.isSliEntitySet(details) {
			return r.mapSliEntityWebsiteTimeBasedFromState(details.([]interface{})[0].(map[string]interface{}))
		}
	}
	return restapi.SliEntity{}, fmt.Errorf("exactly one sli entity configuration of type %s, %s, %s or %s, is required", SliConfigFieldSliEntityApplication, SliConfigFieldSliEntityAvailability, SliConfigFieldSliEntityWebsiteEventBased, SliConfigFieldSliEntityWebsiteTimeBased)
}

func (r *sliConfigResource) isSliEntitySet(details interface{}) bool {
	list, ok := details.([]interface{})
	return ok && len(list) == 1
}

func (r *sliConfigResource) mapSliEntityApplicationFromState(data map[string]interface{}) (restapi.SliEntity, error) {
	return restapi.SliEntity{
		Type:          "application",
		ApplicationID: GetPointerFromMap[string](SliConfigFieldApplicationID, data),
		ServiceID:     GetPointerFromMap[string](SliConfigFieldServiceID, data),
		EndpointID:    GetPointerFromMap[string](SliConfigFieldEndpointID, data),
		BoundaryScope: GetPointerFromMap[string](SliConfigFieldBoundaryScope, data),
	}, nil
}

func (r *sliConfigResource) mapSliEntityAvailabilityFromState(data map[string]interface{}) (restapi.SliEntity, error) {
	var goodEventFilterExpression restapi.TagFilterExpressionElement
	var badEventFilterExpression restapi.TagFilterExpressionElement
	var err error
	if tagFilterString, ok := data[SliConfigFieldGoodEventFilterExpression]; ok {
		goodEventFilterExpression, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
		if err != nil {
			return restapi.SliEntity{}, err
		}
	}
	if tagFilterString, ok := data[SliConfigFieldBadEventFilterExpression]; ok {
		badEventFilterExpression, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
		if err != nil {
			return restapi.SliEntity{}, err
		}
	}

	return restapi.SliEntity{
		Type:                      "availability",
		ApplicationID:             GetPointerFromMap[string](SliConfigFieldApplicationID, data),
		BoundaryScope:             GetPointerFromMap[string](SliConfigFieldBoundaryScope, data),
		GoodEventFilterExpression: goodEventFilterExpression,
		BadEventFilterExpression:  badEventFilterExpression,
		IncludeInternal:           GetPointerFromMap[bool](SliConfigFieldIncludeInternal, data),
		IncludeSynthetic:          GetPointerFromMap[bool](SliConfigFieldIncludeSynthetic, data),
	}, nil
}

func (r *sliConfigResource) mapSliEntityWebsiteTimeBasedFromState(data map[string]interface{}) (restapi.SliEntity, error) {
	var tagFilter restapi.TagFilterExpressionElement
	var err error
	if tagFilterString, ok := data[SliConfigFieldFilterExpression]; ok {
		tagFilter, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
		if err != nil {
			return restapi.SliEntity{}, err
		}
	}

	return restapi.SliEntity{
		Type:             "websiteTimeBased",
		WebsiteId:        GetPointerFromMap[string](SliConfigFieldWebsiteID, data),
		FilterExpression: tagFilter,
		BeaconType:       GetPointerFromMap[string](SliConfigFieldBeaconType, data),
	}, nil
}

func (r *sliConfigResource) mapSliEntityWebsiteEventBasedFromState(data map[string]interface{}) (restapi.SliEntity, error) {
	var goodEventFilterExpression restapi.TagFilterExpressionElement
	var badEventFilterExpression restapi.TagFilterExpressionElement
	var err error
	if tagFilterString, ok := data[SliConfigFieldGoodEventFilterExpression]; ok {
		goodEventFilterExpression, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
		if err != nil {
			return restapi.SliEntity{}, err
		}
	}
	if tagFilterString, ok := data[SliConfigFieldBadEventFilterExpression]; ok {
		badEventFilterExpression, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
		if err != nil {
			return restapi.SliEntity{}, err
		}
	}

	return restapi.SliEntity{
		Type:                      "websiteEventBased",
		WebsiteId:                 GetPointerFromMap[string](SliConfigFieldWebsiteID, data),
		GoodEventFilterExpression: goodEventFilterExpression,
		BadEventFilterExpression:  badEventFilterExpression,
		BeaconType:                GetPointerFromMap[string](SliConfigFieldBeaconType, data),
	}, nil
}

func (r *sliConfigResource) mapTagFilterStringToAPIModel(input string) (restapi.TagFilterExpressionElement, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (r *sliConfigResource) sliConfigStateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[SliConfigFieldFullName]; ok {
		state[SliConfigFieldName] = state[SliConfigFieldFullName]
		delete(state, SliConfigFieldFullName)
	}
	return state, nil
}

func (r *sliConfigResource) sliConfigSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			SliConfigFieldName:                       SliConfigName,
			SliConfigFieldFullName:                   SliConfigFullName,
			SliConfigFieldInitialEvaluationTimestamp: SliConfigInitialEvaluationTimestamp,
			SliConfigFieldMetricConfiguration:        SliConfigMetricConfiguration,
			SliConfigFieldSliEntity:                  SliConfigSliEntity,
		},
	}
}
