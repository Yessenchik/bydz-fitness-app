package validator

import "strings"

func Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

func PositiveFloat(value float64) bool {
	return value >= 0
}

func ValidDurationMonths(value int) bool {
	return value == 1 || value == 6 || value == 12
}
