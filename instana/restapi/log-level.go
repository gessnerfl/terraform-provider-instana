package restapi

// LogLevel custom type for log level
type LogLevel string

// LogLevels custom type for a slice of LogLevel
type LogLevels []LogLevel

// ToStringSlice Returns the corresponding string representations
func (levels LogLevels) ToStringSlice() []string {
	result := make([]string, len(levels))
	for i, v := range levels {
		result[i] = string(v)
	}
	return result
}

const (
	//LogLevelWarning constant value for the warning log level
	LogLevelWarning = LogLevel("WARN")
	//LogLevelError constant value for the error log level
	LogLevelError = LogLevel("ERROR")
	//LogLevelAny constant value for the any log level
	LogLevelAny = LogLevel("ANY")
)

// SupportedLogLevels list of all supported LogLevel
var SupportedLogLevels = LogLevels{LogLevelWarning, LogLevelError, LogLevelAny}
