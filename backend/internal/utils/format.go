package utils

import "strconv"

// FormatInt64 is a small convenience wrapper used in audit descriptions.
func FormatInt64(v int64) string {
	return strconv.FormatInt(v, 10)
}
