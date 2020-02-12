package utils

//StringSliceElementsAreUnique checks if the given string slice contains unique elements only
func StringSliceElementsAreUnique(slice []string) bool {
	checkMap := make(map[string]bool)

	for _, v := range slice {
		if checkMap[v] {
			return false
		}
		checkMap[v] = true
	}
	return true
}
