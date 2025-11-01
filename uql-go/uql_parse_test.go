package uql

import (
	"encoding/json"
	"testing"
)

// TestParseJSON tests parse-json command with various options
func TestParseJSON(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		data     string
		expected interface{}
	}{
		{
			name:  "simple object",
			query: "parse-json",
			data:  `{"name":"John","age":30}`,
			expected: map[string]interface{}{
				"name": "John",
				"age":  float64(30),
			},
		},
		{
			name:  "array",
			query: "parse-json",
			data:  `[1,2,3]`,
			expected: []interface{}{
				float64(1),
				float64(2),
				float64(3),
			},
		},
		{
			name:  "nested object",
			query: "parse-json",
			data:  `{"user":{"name":"John","age":30}}`,
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
					"age":  float64(30),
				},
			},
		},
		{
			name:  "array of objects",
			query: "parse-json",
			data:  `[{"name":"John"},{"name":"Jane"}]`,
			expected: []interface{}{
				map[string]interface{}{"name": "John"},
				map[string]interface{}{"name": "Jane"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UQL(tt.query, &Options{Data: tt.data})
			if err != nil {
				t.Fatalf("UQL failed: %v", err)
			}
			// Compare as JSON strings for easier comparison
			expectedJSON := toJSONString(t, tt.expected)
			resultJSON := toJSONString(t, result)
			if expectedJSON != resultJSON {
				t.Errorf("Expected %s, got %s", expectedJSON, resultJSON)
			}
		})
	}
}

// TestParseCSV tests parse-csv command with various options
func TestParseCSV(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		data     string
		expected []interface{}
	}{
		{
			name:  "csv with headers",
			query: "parse-csv",
			data:  "a,b\n1,2\n3,4",
			expected: []interface{}{
				map[string]interface{}{"a": "1", "b": "2"},
				map[string]interface{}{"a": "3", "b": "4"},
			},
		},
		{
			name:  "csv without headers (columns false)",
			query: `parse-csv --columns "false"`,
			data:  "1,2\n3,4",
			expected: []interface{}{
				map[string]interface{}{"col_0": "1", "col_1": "2"},
				map[string]interface{}{"col_0": "3", "col_1": "4"},
			},
		},
		{
			name:  "csv with custom headers",
			query: `parse-csv --columns "x,y"`,
			data:  "1,2\n3,4",
			expected: []interface{}{
				map[string]interface{}{"x": "1", "y": "2"},
				map[string]interface{}{"x": "3", "y": "4"},
			},
		},
		{
			name:  "csv with custom delimiter semicolon",
			query: `parse-csv --delimiter ";"`,
			data:  "a;b\n1;2\n3;4",
			expected: []interface{}{
				map[string]interface{}{"a": "1", "b": "2"},
				map[string]interface{}{"a": "3", "b": "4"},
			},
		},
		{
			name:  "tsv (tab delimiter)",
			query: `parse-csv --delimiter "\t"`,
			data:  "a\tb\n1\t2\n3\t4",
			expected: []interface{}{
				map[string]interface{}{"a": "1", "b": "2"},
				map[string]interface{}{"a": "3", "b": "4"},
			},
		},
		{
			name:  "csv with pipe delimiter",
			query: `parse-csv --delimiter "|"`,
			data:  "a|b\n1|2\n3|4",
			expected: []interface{}{
				map[string]interface{}{"a": "1", "b": "2"},
				map[string]interface{}{"a": "3", "b": "4"},
			},
		},
		{
			name:  "csv with skip empty lines",
			query: `parse-csv --skipEmptyLines "true"`,
			data:  "a,b\n1,2\n\n3,4",
			expected: []interface{}{
				map[string]interface{}{"a": "1", "b": "2"},
				map[string]interface{}{"a": "3", "b": "4"},
			},
		},
		{
			name:  "csv with comment lines",
			query: `parse-csv --comment "#"`,
			data:  "a,b\n# This is a comment\n1,2\n3,4",
			expected: []interface{}{
				map[string]interface{}{"a": "1", "b": "2"},
				map[string]interface{}{"a": "3", "b": "4"},
			},
		},
		{
			name:  "csv with trim",
			query: `parse-csv --trim "true"`,
			data:  "a,b\n  1  ,  2  \n  3  ,  4  ",
			expected: []interface{}{
				map[string]interface{}{"a": "1", "b": "2"},
				map[string]interface{}{"a": "3", "b": "4"},
			},
		},
		{
			name:  "csv with relaxColumnCount",
			query: `parse-csv --relaxColumnCount "true"`,
			data:  "a,b,c\n1,2\n3,4,5,6",
			expected: []interface{}{
				map[string]interface{}{"a": "1", "b": "2", "c": ""},
				map[string]interface{}{"a": "3", "b": "4", "c": "5"},
			},
		},
		{
			name:  "csv with multiple options",
			query: `parse-csv --delimiter ";" --trim "true" --skipEmptyLines "true"`,
			data:  "a;b\n  1  ;  2  \n\n  3  ;  4  ",
			expected: []interface{}{
				map[string]interface{}{"a": "1", "b": "2"},
				map[string]interface{}{"a": "3", "b": "4"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UQL(tt.query, &Options{Data: tt.data})
			if err != nil {
				t.Fatalf("UQL failed: %v", err)
			}
			expectedJSON := toJSONString(t, tt.expected)
			resultJSON := toJSONString(t, result)
			if expectedJSON != resultJSON {
				t.Errorf("Expected %s, got %s", expectedJSON, resultJSON)
			}
		})
	}
}

// TestParseXML tests parse-xml command
func TestParseXML(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		data     string
		wantErr  bool
	}{
		{
			name:  "simple xml",
			query: "parse-xml",
			data:  `<root><name>John</name><age>30</age></root>`,
			wantErr: false,
		},
		{
			name:  "xml with attributes",
			query: "parse-xml",
			data:  `<root id="1"><name>John</name></root>`,
			wantErr: false,
		},
		{
			name:  "xml array",
			query: "parse-xml",
			data:  `<root><item>1</item><item>2</item><item>3</item></root>`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UQL(tt.query, &Options{Data: tt.data})
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error=%v, got error=%v", tt.wantErr, err)
			}
		})
	}
}

// TestParseYAML tests parse-yaml command
func TestParseYAML(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		data     string
		expected interface{}
	}{
		{
			name:  "simple yaml object",
			query: "parse-yaml",
			data:  "name: John\nage: 30",
			expected: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
		},
		{
			name:  "yaml array",
			query: "parse-yaml",
			data:  "- item1\n- item2\n- item3",
			expected: []interface{}{
				"item1",
				"item2",
				"item3",
			},
		},
		{
			name:  "nested yaml",
			query: "parse-yaml",
			data:  "user:\n  name: John\n  age: 30",
			expected: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "John",
					"age":  30,
				},
			},
		},
		{
			name:  "yaml with multiple documents",
			query: "parse-yaml",
			data:  "---\nname: John\n---\nname: Jane",
			expected: map[string]interface{}{
				"name": "John",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UQL(tt.query, &Options{Data: tt.data})
			if err != nil {
				t.Fatalf("UQL failed: %v", err)
			}
			expectedJSON := toJSONString(t, tt.expected)
			resultJSON := toJSONString(t, result)
			if expectedJSON != resultJSON {
				t.Errorf("Expected %s, got %s", expectedJSON, resultJSON)
			}
		})
	}
}

// TestParseCombinations tests combinations of parse commands
func TestParseCombinations(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		data    interface{}
		wantErr bool
	}{
		{
			name:    "json then project",
			query:   `parse-json | project "name"`,
			data:    `{"name":"John","age":30}`,
			wantErr: false,
		},
		{
			name:    "csv then count",
			query:   "parse-csv | count",
			data:    "a,b\n1,2\n3,4",
			wantErr: false,
		},
		{
			name:    "csv with delimiter then filter",
			query:   `parse-csv --delimiter ";" | where "a" == '1'`,
			data:    "a;b\n1;2\n3;4",
			wantErr: false,
		},
		{
			name:    "yaml then extend",
			query:   `parse-yaml | extend "full"=strcat("name", " ", "age")`,
			data:    "name: John\nage: 30",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UQL(tt.query, &Options{Data: tt.data})
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error=%v, got error=%v", tt.wantErr, err)
			}
		})
	}
}

// Helper function to convert to JSON string for comparison
func toJSONString(t *testing.T, v interface{}) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}
	return string(b)
}
