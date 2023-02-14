package metricname

import (
	"strings"
)

func Join(a ...string) string {
	return strings.Join(a, ".")
}

func Group(s string) string {
	i := strings.LastIndex(s, ".")
	if i < 0 { // s = "name"
		return ""
	}
	return s[:i]
}
