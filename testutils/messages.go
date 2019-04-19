package testutils

//ExpectedError message in case error is expected but no error received
const ExpectedError = "Expected to get error"

//ExpectedNoErrorButGotMessage message in case on error was expected
const ExpectedNoErrorButGotMessage = "Expected no error but got, %s"

//ExpectedUnmarshalledJSONWithStruct message when unmarshalled json does not match expected struct
const ExpectedUnmarshalledJSONWithStruct = "Expected json to be unmarshalled to %v but got %v; diff %s"
