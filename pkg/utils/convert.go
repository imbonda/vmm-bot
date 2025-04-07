package utils

import (
	"fmt"
	"strconv"
)

// FormatIntToString converts an int to a string
func FormatIntToString(value int) string {
	return fmt.Sprint(value)
}

// FormatFloatToString converts a float64 to a string
func FormatFloatToString(value float64, precision int) string {
	if precision >= 0 {
		return strconv.FormatFloat(value, 'f', precision, 64)
	}
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func ParseFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}
