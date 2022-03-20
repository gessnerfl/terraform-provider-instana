package restapi

//ThresholdOperator custom type for the operator of a threshold data structure
type ThresholdOperator string

//ThresholdOperators custom type for a slice of LogLevel
type ThresholdOperators []ThresholdOperator

//IsSupported check if the provided ThresholdOperator is supported
func (operators ThresholdOperators) IsSupported(operator ThresholdOperator) bool {
	for _, g := range operators {
		if g == operator {
			return true
		}
	}
	return false
}

//ToStringSlice Returns the corresponding string representations
func (operators ThresholdOperators) ToStringSlice() []string {
	result := make([]string, len(operators))
	for i, v := range operators {
		result[i] = string(v)
	}
	return result
}

const (
	//ThresholdOperatorGreaterThan constant value for the threshold operator > (greater than)
	ThresholdOperatorGreaterThan = ThresholdOperator(">")
	//ThresholdOperatorGreaterThanOrEqual constant value for the threshold operator >= (greater than or equal)
	ThresholdOperatorGreaterThanOrEqual = ThresholdOperator(">=")
	//ThresholdOperatorLessThan constant value for the threshold operator < (less than)
	ThresholdOperatorLessThan = ThresholdOperator("<")
	//ThresholdOperatorLessThanOrEqual constant value for the threshold operator <= (less than or equal)
	ThresholdOperatorLessThanOrEqual = ThresholdOperator("<=")
)

//SupportedThresholdOperators list of all supported ThresholdOperator
var SupportedThresholdOperators = ThresholdOperators{ThresholdOperatorGreaterThan, ThresholdOperatorGreaterThanOrEqual, ThresholdOperatorLessThan, ThresholdOperatorLessThanOrEqual}

//ThresholdSeasonality custom type for the seasonality of a threshold data structure
type ThresholdSeasonality string

//ThresholdSeasonalities custom type for a slice of ThresholdSeasonality
type ThresholdSeasonalities []ThresholdSeasonality

//IsSupported check if the provided ThresholdSeasonality is supported
func (seasonalities ThresholdSeasonalities) IsSupported(seasonality ThresholdSeasonality) bool {
	for _, g := range seasonalities {
		if g == seasonality {
			return true
		}
	}
	return false
}

//ToStringSlice Returns the corresponding string representations
func (seasonalities ThresholdSeasonalities) ToStringSlice() []string {
	result := make([]string, len(seasonalities))
	for i, v := range seasonalities {
		result[i] = string(v)
	}
	return result
}

const (
	//ThresholdSeasonalityWeekly constant value for the threshold seasonality type weekly
	ThresholdSeasonalityWeekly = ThresholdSeasonality("WEEKLY")
	//ThresholdSeasonalityDaily constant value for the threshold seasonality type daily
	ThresholdSeasonalityDaily = ThresholdSeasonality("DAILY")
)

//SupportedThresholdSeasonalities list of all supported ThresholdSeasonality
var SupportedThresholdSeasonalities = ThresholdSeasonalities{ThresholdSeasonalityWeekly, ThresholdSeasonalityDaily}

//Threshold custom data structure representing the threshold type of the instana API
type Threshold struct {
	Type            string                `json:"type"`
	Operator        ThresholdOperator     `json:"operator"`
	Baseline        *[][]float64          `json:"baseline"`
	DeviationFactor *float32              `json:"deviationFactor"`
	LastUpdated     int64                 `json:"lastUpdated"`
	Seasonality     *ThresholdSeasonality `json:"seasonality"`
	Value           *float64              `json:"value"`
}

//TimeThreshold custom data structure representing the time threshold type of the instana API
type TimeThreshold struct {
	Type       string `json:"type"`
	TimeWindow int64  `json:"timeWindow"`
	Requests   *int32 `json:"requests"`
	Violations *int32 `json:"violations"`
}
