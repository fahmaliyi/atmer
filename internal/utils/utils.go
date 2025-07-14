package utils

import "fmt"

func ToString(v any) string {
	switch val := v.(type) {
	case nil:
		return ""
	case string:
		return val
	case float64:
		return fmt.Sprintf("%.0f", val) // remove scientific notation
	case int, int64, int32:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
