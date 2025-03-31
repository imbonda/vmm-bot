package utils

import "strconv"

// FormatFloatToString converts a float64 to a string without specifying decimals
func FormatFloatToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
