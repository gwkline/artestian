package prompt_utils

import (
	"strings"
	"testing"
)

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello_world"},
		{"ABC", "a_b_c"},
		{"simpleText", "simple_text"},
		{"ID", "i_d"},
		{"", ""},
		{"alreadysnakecase", "alreadysnakecase"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toSnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("toSnakeCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStructToXMLString(t *testing.T) {
	type TestStruct struct {
		Name        string
		Age         int
		Scores      []int
		unexported  string
		StringSlice []string
	}

	test := TestStruct{
		Name:        "John Doe",
		Age:         30,
		Scores:      []int{95, 87, 92},
		unexported:  "should not appear",
		StringSlice: []string{"one", "two", "three"},
	}

	result, err := StructToXMLString(test)
	if err != nil {
		t.Fatalf("StructToXMLString failed: %v", err)
	}

	// Check for expected contents
	expectedParts := []string{
		"<name>\nJohn Doe\n</name>",
		"<age>\n30\n</age>",
		"<scores>\n95\n87\n92\n</scores>",
		"<string_slice>\none\ntwo\nthree\n</string_slice>",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected result to contain %q, but it didn't.\nGot: %s", part, result)
		}
	}

	// Check that unexported field is not included
	if strings.Contains(result, "unexported") {
		t.Error("Result should not contain unexported field")
	}

	// Test error case
	_, err = StructToXMLString("not a struct")
	if err == nil {
		t.Error("Expected error when passing non-struct value, got nil")
	}
}

type NamedInterface interface {
	GetName() string
}

type MockNamed struct{}

func (m MockNamed) GetName() string {
	return "MockObject"
}

type SimpleInterface interface {
	DoSomething()
}

type MockSimple struct{}

func (m MockSimple) DoSomething() {}

func TestStructToXMLString_Interface(t *testing.T) {
	type TestInterfaceStruct struct {
		Named  NamedInterface
		Simple SimpleInterface
	}

	test := TestInterfaceStruct{
		Named:  MockNamed{},
		Simple: MockSimple{},
	}

	result, err := StructToXMLString(test)
	if err != nil {
		t.Fatalf("StructToXMLString failed: %v", err)
	}

	// Check for expected contents
	expectedParts := []string{
		"<named>\nMockObject\n</named>",
		"<simple>\nprompt_utils.MockSimple\n</simple>",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected result to contain %q, but it didn't.\nGot: %s", part, result)
		}
	}
}

func TestStructToXMLString_EmbeddedStruct(t *testing.T) {
	type Address struct {
		Street string
		City   string
	}

	type PersonWithAddress struct {
		Name    string
		Age     int
		Address // embedded struct
	}

	test := PersonWithAddress{
		Name: "Jane Doe",
		Age:  25,
		Address: Address{
			Street: "123 Main St",
			City:   "Springfield",
		},
	}

	result, err := StructToXMLString(test)
	if err != nil {
		t.Fatalf("StructToXMLString failed: %v", err)
	}

	// Check for expected contents
	expectedParts := []string{
		"<name>\nJane Doe\n</name>",
		"<age>\n25\n</age>",
		"<street>\n123 Main St\n</street>",
		"<city>\nSpringfield\n</city>",
	}

	for _, part := range expectedParts {
		if !strings.Contains(result, part) {
			t.Errorf("Expected result to contain %q, but it didn't.\nGot: %s", part, result)
		}
	}
}
