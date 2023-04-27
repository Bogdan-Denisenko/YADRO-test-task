package utils

import "strings"

func RemoveMatching(slice []string, str string) []string {
	result := []string{}

	for _, s := range slice {
		if !strings.EqualFold(s, str) {
			result = append(result, s)
		}
	}

	return result
}