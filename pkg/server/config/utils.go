package config

import (
	"fmt"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func camelToUpperSnake(s string) string {
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}

func flattenConfig(prefix string, out map[string]string, v any) {
	switch t := v.(type) {
	case map[string]any:
		for key, value := range t {
			snakeKey := camelToUpperSnake(key)
			newPrefix := strings.TrimPrefix(prefix+"_"+snakeKey, "_")
			flattenConfig(newPrefix, out, value)
		}
	case []any:
		var strValues []string
		for _, val := range t {
			strValues = append(strValues, formatAny(val))
		}
		out[prefix] = strings.Join(strValues, ",")
	case string:
		out[prefix] = t
	case bool, int, float64:
		out[prefix] = formatAny(t)
	}
}

func formatAny(val any) string {
	switch v := val.(type) {
	case string:
		return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(v, "\n", ""), "\r", ""), "\t", ""))
	case bool:
		return fmt.Sprintf("%t", v)
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
