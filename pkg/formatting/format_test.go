package formatting

import "testing"

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
