package formatting

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
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
	var builder strings.Builder

	// Collect and sort the map keys
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Iterate over the sorted map keys
	for _, key := range keys {
		value := m[key]

		// Add key and process value with appropriate formatting
		builder.WriteString(fmt.Sprintf("%s: %v\n", key, formatValue(value)))
	}

	return builder.String()
}

// formatValue formats the value based on its type
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case map[string]interface{}:
		// If the value is a nested map, format it with indentation
		var builder strings.Builder
		builder.WriteString("\n")
		for key, nestedValue := range v {
			builder.WriteString(fmt.Sprintf("  %s: %v\n", key, formatValue(nestedValue)))
		}
		return builder.String()
	case []interface{}:
		// Handle arrays/lists
		var builder strings.Builder
		builder.WriteString("[\n")
		for _, item := range v {
			builder.WriteString(fmt.Sprintf("  %v\n", formatValue(item)))
		}
		builder.WriteString("]")
		return builder.String()
	default:
		// For basic types, return the string representation
		return fmt.Sprintf("%v", v)
	}
}
func TestMapToString(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{
			name:     "Simple key-value pairs",
			input:    map[string]interface{}{"a": 1, "b": "string", "c": true},
			expected: "a: 1\nb: string\nc: true\n",
		},
		{
			name: "Nested map",
			input: map[string]interface{}{
				"parent": map[string]interface{}{
					"child1": "value1",
					"child2": 2,
				},
				"anotherKey": "anotherValue",
			},
			expected: "anotherKey: anotherValue\nparent:\n  child1: value1\n  child2: 2\n",
		},
		{
			name: "Array in map",
			input: map[string]interface{}{
				"list": []interface{}{"item1", 2, true},
			},
			expected: "list: [\n  item1\n  2\n  true\n]\n",
		},
		{
			name:     "Empty map",
			input:    map[string]interface{}{},
			expected: "",
		},
		{
			name: "Mixed nested structures",
			input: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": []interface{}{
						map[string]interface{}{
							"key1": "value1",
							"key2": 3,
						},
					},
				},
			},
			expected: "level1:\n  level2: [\n    key1: value1\n    key2: 3\n  ]\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := MapToString(test.input)
			assert.Equal(t, test.expected, actual)
		})
	}
}
