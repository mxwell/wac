package util

func ContainsString(arr *[]string, value string) bool {
	for _, e := range *arr {
		if e == value {
			return true
		}
	}
	return false
}
