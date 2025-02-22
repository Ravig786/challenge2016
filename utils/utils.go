package utils

import "strings"

func NormalizeRegion(region string) string {
	return strings.ToUpper(region)
}

func SplitRegion(region string) []string {
	return strings.Split(strings.ToUpper(region), "-")
}
