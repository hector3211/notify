package utils

import (
	"os"
	"strings"
)

func IsProduction() bool {
	return os.Getenv("ENV") == "production"
}

func ContainsQuery(str, queryStr string) bool {
	strLower := strings.ToLower(str)
	queryStrLower := strings.ToLower(queryStr)
	for i := 0; i < len(queryStrLower); i++ {
		if strLower[i] != queryStrLower[i] {
			return false
		}
	}

	return true
}
