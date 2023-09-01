package restapi

// Aggregation custom type for an Aggregation
type Aggregation string

// Aggregations custom type for a slice of Aggregation
type Aggregations []Aggregation

// ToStringSlice Returns the corresponding string representations
func (aggregations Aggregations) ToStringSlice() []string {
	result := make([]string, len(aggregations))
	for i, v := range aggregations {
		result[i] = string(v)
	}
	return result
}

const (
	//SumAggregation constant value for the sum aggregation type
	SumAggregation = Aggregation("SUM")
	//MeanAggregation constant value for the mean aggregation type
	MeanAggregation = Aggregation("MEAN")
	//MaxAggregation constant value for the max aggregation type
	MaxAggregation = Aggregation("MAX")
	//MinAggregation constant value for the min aggregation type
	MinAggregation = Aggregation("MIN")
	//Percentile25Aggregation constant value for the 25th percentile aggregation type
	Percentile25Aggregation = Aggregation("P25")
	//Percentile50Aggregation constant value for the 50th percentile aggregation type
	Percentile50Aggregation = Aggregation("P50")
	//Percentile75Aggregation constant value for the 75th percentile aggregation type
	Percentile75Aggregation = Aggregation("P75")
	//Percentile90Aggregation constant value for the 90th percentile aggregation type
	Percentile90Aggregation = Aggregation("P90")
	//Percentile95Aggregation constant value for the 95th percentile aggregation type
	Percentile95Aggregation = Aggregation("P95")
	//Percentile98Aggregation constant value for the 98th percentile aggregation type
	Percentile98Aggregation = Aggregation("P98")
	//Percentile99Aggregation constant value for the 99th percentile aggregation type
	Percentile99Aggregation = Aggregation("P99")
	//Percentile99_9Aggregation constant value for the 99.9th percentile aggregation type
	Percentile99_9Aggregation = Aggregation("P99_9")
	//Percentile99_99Aggregation constant value for the 99.99th percentile aggregation type
	Percentile99_99Aggregation = Aggregation("P99_99")
	//DistributionAggregation constant value for the distribution aggregation type
	DistributionAggregation = Aggregation("DISTRIBUTION")
	//DistinctCountAggregation constant value for the distinct count aggregation type
	DistinctCountAggregation = Aggregation("DISTINCT_COUNT")
	//SumPositiveAggregation constant value for the sum positive aggregation type
	SumPositiveAggregation = Aggregation("SUM_POSITIVE")
)

// SupportedAggregations list of all supported Aggregation
var SupportedAggregations = Aggregations{SumAggregation, MeanAggregation, MaxAggregation, MinAggregation, Percentile25Aggregation, Percentile50Aggregation, Percentile75Aggregation, Percentile90Aggregation, Percentile95Aggregation, Percentile98Aggregation, Percentile99Aggregation, Percentile99_9Aggregation, Percentile99_99Aggregation, DistributionAggregation, DistinctCountAggregation, SumPositiveAggregation}
