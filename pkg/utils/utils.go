package utils

import "strconv"

func FormatFloat(num float64) string {
	return strconv.FormatFloat(num, 'f', 2, 64) // 'f' for fixed-point, 2 for 2 decimal places
}
