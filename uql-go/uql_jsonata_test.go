package uql

import (
	"encoding/json"
	"testing"
)

func TestUQLJSONata(t *testing.T) {
	data := map[string]interface{}{
		"example": []interface{}{
			map[string]interface{}{"value": 4},
			map[string]interface{}{"value": 7},
			map[string]interface{}{"value": 13},
		},
	}

	result, err := UQL(`jsonata "$sum(example.value)"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("UQL failed: %v", err)
	}

	// Check if result is a number (could be float64 or int)
	switch v := result.(type) {
	case float64:
		if v != 24.0 {
			t.Errorf("Expected 24, got %v", v)
		}
	case int:
		if v != 24 {
			t.Errorf("Expected 24, got %v", v)
		}
	default:
		t.Errorf("Expected number, got type %T: %v", result, result)
	}
}

func TestUQLJSONataArray(t *testing.T) {
	data := map[string]interface{}{
		"example": []interface{}{
			map[string]interface{}{"value": 4},
			map[string]interface{}{"value": 7},
			map[string]interface{}{"value": 13},
		},
	}

	result, err := UQL(`jsonata "example.value"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("UQL failed: %v", err)
	}

	// Convert result to JSON and compare
	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	expectedJSON := `[4,7,13]`
	if string(resultJSON) != expectedJSON {
		t.Errorf("Expected %s, got %s", expectedJSON, string(resultJSON))
	}
}

func TestUQLJSONataLibrary(t *testing.T) {
	libraryData := map[string]interface{}{
		"library": map[string]interface{}{
			"books": []interface{}{
				map[string]interface{}{
					"title":   "Structure and Interpretation of Computer Programs",
					"authors": []interface{}{"Abelson", "Sussman"},
					"isbn":    "9780262510875",
					"price":   38.9,
					"copies":  2,
				},
				map[string]interface{}{
					"title":   "The C Programming Language",
					"authors": []interface{}{"Kernighan", "Richie"},
					"isbn":    "9780131103627",
					"price":   33.59,
					"copies":  3,
				},
				map[string]interface{}{
					"title":   "The AWK Programming Language",
					"authors": []interface{}{"Aho", "Kernighan", "Weinberger"},
					"isbn":    "9780201079814",
					"copies":  1,
				},
				map[string]interface{}{
					"title":   "Compilers: Principles, Techniques, and Tools",
					"authors": []interface{}{"Aho", "Lam", "Sethi", "Ullman"},
					"isbn":    "9780201100884",
					"price":   23.38,
					"copies":  1,
				},
			},
		},
	}

	result, err := UQL(`jsonata "$sum(library.books.price)"`, &Options{Data: libraryData})
	if err != nil {
		t.Fatalf("UQL failed: %v", err)
	}

	// Check if result is approximately 95.87 (sum of prices)
	switch v := result.(type) {
	case float64:
		expected := 95.87
		if v < expected-0.01 || v > expected+0.01 {
			t.Errorf("Expected approximately %v, got %v", expected, v)
		}
	default:
		t.Errorf("Expected float64, got type %T: %v", result, result)
	}
}

func TestUQLJSONataCombinedWithOtherQueries(t *testing.T) {
	data := map[string]interface{}{
		"library": map[string]interface{}{
			"library": map[string]interface{}{
				"books": []interface{}{
					map[string]interface{}{
						"title": "Structure and Interpretation of Computer Programs",
						"isbn":  "9780262510875",
						"price": 38.9,
					},
					map[string]interface{}{
						"title": "The C Programming Language",
						"isbn":  "9780131103627",
						"price": 33.59,
					},
				},
				"loans": []interface{}{
					map[string]interface{}{
						"customer": "10001",
						"isbn":     "9780262510875",
						"return":   "2016-12-05",
					},
				},
				"customers": []interface{}{
					map[string]interface{}{
						"id":   "10001",
						"name": "Joe Doe",
					},
				},
			},
		},
	}

	// Combine scope and jsonata with count
	query := `scope "library" | jsonata "library.loans" | count`
	result, err := UQL(query, &Options{Data: data})
	if err != nil {
		t.Fatalf("UQL failed: %v", err)
	}

	// Should return count of loans (1)
	count, ok := result.(int)
	if !ok {
		t.Fatalf("Expected int, got type %T: %v", result, result)
	}

	if count != 1 {
		t.Errorf("Expected 1, got %d", count)
	}
}

func TestUQLJSONataFilter(t *testing.T) {
	data := map[string]interface{}{
		"Countries": []interface{}{
			map[string]interface{}{"Country": "India"},
			map[string]interface{}{"Country": "America"},
			map[string]interface{}{"Country": "United Kingdom"},
		},
	}

	// Filter for single element - use $ to reference the array after scope
	result, err := UQL(`scope "Countries" | jsonata "$[Country='India']"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("UQL failed: %v", err)
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal result: %v", err)
	}

	expectedJSON := `{"Country":"India"}`
	if string(resultJSON) != expectedJSON {
		t.Errorf("Expected %s, got %s", expectedJSON, string(resultJSON))
	}
}

func TestUQLJSONataFilterMultiple(t *testing.T) {
	data := map[string]interface{}{
		"Countries": []interface{}{
			map[string]interface{}{"Country": "India"},
			map[string]interface{}{"Country": "America"},
			map[string]interface{}{"Country": "United Kingdom"},
		},
	}

	// Filter for multiple elements - use $[predicate] to filter array
	result, err := UQL(`scope "Countries" | jsonata "$[Country in ['India','United Kingdom']]"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("UQL failed: %v", err)
	}

	resultArray, ok := result.([]interface{})
	if !ok {
		t.Fatalf("Expected array, got type %T: %v", result, result)
	}

	if len(resultArray) != 2 {
		t.Errorf("Expected 2 results, got %d", len(resultArray))
	}

	// Check that we have India and United Kingdom
	countries := make([]string, 0)
	for _, item := range resultArray {
		if m, ok := item.(map[string]interface{}); ok {
			if country, ok := m["Country"].(string); ok {
				countries = append(countries, country)
			}
		}
	}

	hasIndia := false
	hasUK := false
	for _, c := range countries {
		if c == "India" {
			hasIndia = true
		}
		if c == "United Kingdom" {
			hasUK = true
		}
	}

	if !hasIndia || !hasUK {
		t.Errorf("Expected India and United Kingdom, got %v", countries)
	}
}
