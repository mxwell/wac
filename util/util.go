package util

import (
	"os"
)

func ContainsString(arr *[]string, value string) bool {
	for _, e := range *arr {
		if e == value {
			return true
		}
	}
	return false
}

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
