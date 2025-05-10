package helper

import (
	"fmt"
	"sort"
	"strings"
)

func FormatarLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	var parts []string
	for k, v := range labels {
		parts = append(parts, fmt.Sprintf("%s=\"%s\"", k, v))
	}
	sort.Strings(parts)
	return "{" + strings.Join(parts, ",") + "}"
}
