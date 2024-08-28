package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Smartly parses a string as a float using only the numeric components.
func parseAsFloat(value string, fallback float64) (float64, error) {
	numericValue := regexp.MustCompile("[^0-9.]+").ReplaceAllString(value, "") // Strip any non-numeric characters (except for the decimal point)
	numericValue = regexp.MustCompile(`\.{2,}`).ReplaceAllString(numericValue, ".") // Collapse multiple decimal points into one
	numericValue = strings.Trim(numericValue, ".") // Trim any leading/trailing decimal points
	numericValue = strings.TrimSpace(numericValue) // Trim any whitespace

	if numericValue == "" { return fallback, nil }

	return strconv.ParseFloat(numericValue, 64)
}
