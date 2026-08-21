package recovery

import (
	"sort"
	"strings"
)

func normalizeLabels(labels map[string]string) {
	for key, value := range labels {
		normalizedKey := strings.ToLower(strings.TrimSpace(key))
		normalizedValue := strings.TrimSpace(value)
		if normalizedKey != key {
			delete(labels, key)
		}
		if normalizedKey != "" && normalizedValue != "" {
			labels[normalizedKey] = normalizedValue
		}
	}
}

func sortedLabels(labels map[string]string) [][2]string {
	keys := make([]string, 0, len(labels))
	for key := range labels {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([][2]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, [2]string{key, labels[key]})
	}
	return result
}
