package formatting

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

// StructToString converts a struct to a string representation
func StructToString(v interface{}) string {
	val := reflect.ValueOf(v)
	typeOfVal := val.Type()

	var builder strings.Builder

	for i := 0; i < val.NumField(); i++ {
		field := typeOfVal.Field(i)
		value := val.Field(i).Interface()
		// Append each field and value to the string builder
		builder.WriteString(fmt.Sprintf("%s: %v\n", field.Name, value))
	}

	return builder.String()
}

// MapToString converts a map to a string representation
func MapToString(m map[string]interface{}) string {
	return mapToStringWithIndent(m, 0)
}

func mapToStringWithIndent(m map[string]interface{}, indentLevel int) string {
	var builder strings.Builder

	// Helper function to add indentation
	addIndent := func(level int) {
		for i := 0; i < level; i++ {
			builder.WriteString("  ")
		}
	}

	// Collect and sort the map keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Iterate over the sorted map keys
	for _, key := range keys {
		value := m[key]

		// Add key with indentation
		addIndent(indentLevel)
		builder.WriteString(fmt.Sprintf("%s:", key))

		// Process the value with appropriate formatting
		if subMap, ok := value.(map[string]interface{}); ok {
			builder.WriteString("\n")
			builder.WriteString(mapToStringWithIndent(subMap, indentLevel+1))
		} else if arr, ok := value.([]interface{}); ok {
			builder.WriteString(" [\n")
			for _, item := range arr {
				addIndent(indentLevel + 1)
				formatValue(item, &builder, indentLevel+1)
				builder.WriteString("\n")
			}
			addIndent(indentLevel)
			builder.WriteString("]")
		} else {
			builder.WriteString(" ")
			formatValue(value, &builder, indentLevel)
		}
		builder.WriteString("\n")
	}

	return builder.String()
}

func formatValue(value interface{}, builder *strings.Builder, indentLevel int) {
	switch v := value.(type) {
	case map[string]interface{}:
		builder.WriteString("\n")
		builder.WriteString(mapToStringWithIndent(v, indentLevel+1))
	case []interface{}:
		builder.WriteString("[\n")
		for _, item := range v {
			addIndent(indentLevel + 1)
			formatValue(item, builder, indentLevel+1)
			builder.WriteString("\n")
		}
		addIndent(indentLevel)
		builder.WriteString("]")
	default:
		builder.WriteString(fmt.Sprintf("%v", v))
	}
}
func addIndent(level int) string {
	var builder strings.Builder
	for i := 0; i < level; i++ {
		builder.WriteString("  ")
	}
	return builder.String()
}
func queryToMap(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{})

	// Get all query parameters from the request
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0] // Use the first value if there are multiple values for a key
		}
	}

	return params
}
