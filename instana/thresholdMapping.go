package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	//ResourceFieldThreshold constant value for field threshold
	ResourceFieldThreshold = "threshold"
	//ResourceFieldThresholdLastUpdated constant value for field threshold.*.last_updated
	ResourceFieldThresholdLastUpdated = "last_updated"
	//ResourceFieldThresholdOperator constant value for field threshold.*.operator
	ResourceFieldThresholdOperator = "operator"
	//ResourceFieldThresholdHistoricBaseline constant value for field threshold.historic_baseline
	ResourceFieldThresholdHistoricBaseline = "historic_baseline"
	//ResourceFieldThresholdHistoricBaselineBaseline constant value for field threshold.historic_baseline.baseline
	ResourceFieldThresholdHistoricBaselineBaseline = "baseline"
	//ResourceFieldThresholdHistoricBaselineDeviationFactor constant value for field threshold.historic_baseline.deviation_factor
	ResourceFieldThresholdHistoricBaselineDeviationFactor = "deviation_factor"
	//ResourceFieldThresholdHistoricBaselineSeasonality constant value for field threshold.historic_baseline.seasonality
	ResourceFieldThresholdHistoricBaselineSeasonality = "seasonality"
	//ResourceFieldThresholdStatic constant value for field threshold.static
	ResourceFieldThresholdStatic = "static"
	//ResourceFieldThresholdStaticValue constant value for field threshold.static.value
	ResourceFieldThresholdStaticValue = "value"
)

var (
	resourceSchemaThresholdTypeKeys = []string{
		"threshold.0.historic_baseline",
		"threshold.0.static",
	}

	resourceSchemaRequiredThresholdOperator = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The operator which will be applied to evaluate the threshold",
		ValidateFunc: validation.StringInSlice(restapi.SupportedThresholdOperators.ToStringSlice(), true),
	}

	resourceSchemaOptionalThresholdLastUpdated = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(0),
		Description:  "The last updated value of the threshold",
	}
)

var thresholdSchema = &schema.Schema{
	Type:        schema.TypeList,
	MinItems:    1,
	MaxItems:    1,
	Required:    true,
	Description: "Indicates the type of threshold this alert rule is evaluated on.",
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			ResourceFieldThresholdHistoricBaseline: {
				Type:        schema.TypeList,
				MinItems:    0,
				MaxItems:    1,
				Optional:    true,
				Description: "Threshold based on a historic baseline.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ResourceFieldThresholdOperator:    resourceSchemaRequiredThresholdOperator,
						ResourceFieldThresholdLastUpdated: resourceSchemaOptionalThresholdLastUpdated,
						ResourceFieldThresholdHistoricBaselineBaseline: {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:     schema.TypeSet,
								Optional: false,
								Elem: &schema.Schema{
									Type: schema.TypeFloat,
								},
							},
							Description: "The baseline of the historic baseline threshold",
						},
						ResourceFieldThresholdHistoricBaselineDeviationFactor: {
							Type:         schema.TypeFloat,
							Optional:     true,
							ValidateFunc: validation.FloatBetween(0.5, 16),
							Description:  "The baseline of the historic baseline threshold",
						},
						ResourceFieldThresholdHistoricBaselineSeasonality: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(restapi.SupportedThresholdSeasonalities.ToStringSlice(), false),
							Description:  "The seasonality of the historic baseline threshold",
						},
					},
				},
				ExactlyOneOf: resourceSchemaThresholdTypeKeys,
			},
			ResourceFieldThresholdStatic: {
				Type:        schema.TypeList,
				MinItems:    0,
				MaxItems:    1,
				Optional:    true,
				Description: "Static threshold definition",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						ResourceFieldThresholdOperator:    resourceSchemaRequiredThresholdOperator,
						ResourceFieldThresholdLastUpdated: resourceSchemaOptionalThresholdLastUpdated,
						ResourceFieldThresholdStaticValue: {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "The value of the static threshold",
						},
					},
				},
				ExactlyOneOf: resourceSchemaThresholdTypeKeys,
			},
		},
	},
}

func newThresholdMapper() thresholdMapper {
	return &thresholdMapperImpl{}
}

type thresholdMapper interface {
	toState(threshold *restapi.Threshold) []map[string]interface{}
	fromState(d *schema.ResourceData) *restapi.Threshold
}

type thresholdMapperImpl struct{}

func (m *thresholdMapperImpl) toState(input *restapi.Threshold) []map[string]interface{} {
	thresholdConfig := make(map[string]interface{})
	thresholdConfig[ResourceFieldThresholdOperator] = input.Operator
	thresholdConfig[ResourceFieldThresholdLastUpdated] = input.LastUpdated

	if input.Value != nil {
		thresholdConfig[ResourceFieldThresholdStaticValue] = *input.Value
	}
	if input.Baseline != nil {
		thresholdConfig[ResourceFieldThresholdHistoricBaselineBaseline] = *input.Baseline
	}
	if input.DeviationFactor != nil {
		thresholdConfig[ResourceFieldThresholdHistoricBaselineDeviationFactor] = float64(*input.DeviationFactor)
	}
	if input.Seasonality != nil {
		thresholdConfig[ResourceFieldThresholdHistoricBaselineSeasonality] = *input.Seasonality
	}

	thresholdType := m.mapThresholdTypeToSchema(input.Type)
	threshold := make(map[string]interface{})
	threshold[thresholdType] = []interface{}{thresholdConfig}
	result := make([]map[string]interface{}, 1)
	result[0] = threshold
	return result
}

func (m *thresholdMapperImpl) mapThresholdTypeToSchema(input string) string {
	if input == "historicBaseline" {
		return ResourceFieldThresholdHistoricBaseline
	} else if input == "staticThreshold" {
		return ResourceFieldThresholdStatic
	}
	return input
}

func (m *thresholdMapperImpl) fromState(d *schema.ResourceData) *restapi.Threshold {
	thresholdSlice := d.Get(ResourceFieldThreshold).([]interface{})
	threshold := thresholdSlice[0].(map[string]interface{})
	for thresholdType, v := range threshold {
		configSlice := v.([]interface{})
		if len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})
			return m.mapThresholdConfigFromSchema(config, thresholdType)
		}
	}
	return &restapi.Threshold{}
}

func (m *thresholdMapperImpl) mapThresholdConfigFromSchema(config map[string]interface{}, thresholdType string) *restapi.Threshold {
	var seasonalityPtr *restapi.ThresholdSeasonality
	if v, ok := config[ResourceFieldThresholdHistoricBaselineSeasonality]; ok {
		seasonality := restapi.ThresholdSeasonality(v.(string))
		seasonalityPtr = &seasonality
	}
	var lastUpdatePtr *int64
	if v, ok := config[ResourceFieldThresholdLastUpdated]; ok {
		lastUpdate := int64(v.(int))
		lastUpdatePtr = &lastUpdate
	}
	var valuePtr *float64
	if v, ok := config[ResourceFieldThresholdStaticValue]; ok {
		value := v.(float64)
		valuePtr = &value
	}
	var deviationFactorPtr *float32
	if v, ok := config[ResourceFieldThresholdHistoricBaselineDeviationFactor]; ok {
		deviationFactor := float32(v.(float64))
		deviationFactorPtr = &deviationFactor
	}
	var baselinePtr *[][]float64
	if v, ok := config[ResourceFieldThresholdHistoricBaselineBaseline]; ok {
		baselineSet := v.(*schema.Set)
		if baselineSet.Len() > 0 {
			baseline := make([][]float64, baselineSet.Len())
			for i, val := range baselineSet.List() {
				baseline[i] = ConvertInterfaceSlice[float64](val.(*schema.Set).List())
			}
			baselinePtr = &baseline
		}
	}
	return &restapi.Threshold{
		Type:            m.mapThresholdTypeFromSchema(thresholdType),
		Operator:        restapi.ThresholdOperator(config[ResourceFieldThresholdOperator].(string)),
		LastUpdated:     lastUpdatePtr,
		Value:           valuePtr,
		DeviationFactor: deviationFactorPtr,
		Baseline:        baselinePtr,
		Seasonality:     seasonalityPtr,
	}
}

func (m *thresholdMapperImpl) mapThresholdTypeFromSchema(input string) string {
	if input == ResourceFieldThresholdHistoricBaseline {
		return "historicBaseline"
	} else if input == ResourceFieldThresholdStatic {
		return "staticThreshold"
	}
	return input
}
