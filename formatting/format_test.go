package formatting

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"gotest.tools/v3/assert"
)

func TestStructToString(t *testing.T) {
	type TestStruct struct {
		Name    string
		Age     int
		Balance float64
		Active  bool
	}

	testStruct := TestStruct{
		Name:    "John Doe",
		Age:     30,
		Balance: 1234.56,
		Active:  true,
	}

	expectedOutput := "Name: John Doe\nAge: 30\nBalance: 1234.56\nActive: true\n"
	result := StructToString(testStruct)

	if result != expectedOutput {
		t.Errorf("Expected %q but got %q", expectedOutput, result)
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
func TestQueryToMap(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a response recorder to capture the response
	w := httptest.NewRecorder()

	// Create a new Gin context with a test request
	c, _ := gin.CreateTestContext(w)

	// Define a test request with query parameters
	req, _ := http.NewRequest("GET", "/your-endpoint?name=John&age=30&active=true", nil)
	c.Request = req

	// Call the function to test
	result := queryToMap(c)

	// Expected map based on the test query parameters
	expected := map[string]interface{}{
		"name":   "John",
		"age":    "30",
		"active": "true",
	}

	// Use reflect.DeepEqual to verify the result
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}
