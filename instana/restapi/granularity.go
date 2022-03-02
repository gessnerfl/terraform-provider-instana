package restapi

//Granularity custom type for an Alter Granularity
type Granularity int32

//Granularities custom type for a slice of Granularity
type Granularities []Granularity

//IsSupported check if the provided Granularity is supported
func (granularities Granularities) IsSupported(granularity Granularity) bool {
	for _, g := range granularities {
		if g == granularity {
			return true
		}
	}
	return false
}

//ToIntSlice Returns the corresponding int representations
func (granularities Granularities) ToIntSlice() []int {
	result := make([]int, len(granularities))
	for i, v := range granularities {
		result[i] = int(v)
	}
	return result
}

const (
	//Granularity300000 constant value for granularity of 30sec
	Granularity300000 = Granularity(300000)
	//Granularity600000 constant value for granularity of 1min
	Granularity600000 = Granularity(600000)
	//Granularity900000 constant value for granularity of 1min 30sec
	Granularity900000 = Granularity(900000)
	//Granularity1200000 constant value for granularity of 2min
	Granularity1200000 = Granularity(1200000)
	//Granularity1800000 constant value for granularity of 3min
	Granularity1800000 = Granularity(1800000)
)

//SupportedGranularities list of all supported Granularities
var SupportedGranularities = Granularities{Granularity300000, Granularity600000, Granularity900000, Granularity1200000, Granularity1800000}
