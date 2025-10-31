package main

import (
	"encoding/json"
	"fmt"
	"log"

	uql "github.com/yesoreyeram/uql/uql-go"
)

func main() {
	fmt.Println("=== UQL Go Examples ===\n")

	// Example 1: Basic Hello
	fmt.Println("Example 1: Basic Hello")
	result1, _ := uql.UQL("hello", nil)
	fmt.Printf("Result: %v\n\n", result1)

	// Example 2: Ping
	fmt.Println("Example 2: Ping")
	result2, _ := uql.UQL("ping", nil)
	fmt.Printf("Result: %v\n\n", result2)

	// Example 3: Echo
	fmt.Println("Example 3: Echo")
	result3, _ := uql.UQL(`echo "Hello from UQL!"`, nil)
	fmt.Printf("Result: %v\n\n", result3)

	// Example 4: Count
	fmt.Println("Example 4: Count")
	users := []interface{}{
		map[string]interface{}{"name": "Alice", "age": 30},
		map[string]interface{}{"name": "Bob", "age": 25},
		map[string]interface{}{"name": "Charlie", "age": 35},
	}
	result4, _ := uql.UQL("count", &uql.Options{Data: users})
	fmt.Printf("Result: %v\n\n", result4)

	// Example 5: Project
	fmt.Println("Example 5: Project specific fields")
	data := []interface{}{
		map[string]interface{}{"name": "Alice", "age": 30, "city": "NYC", "country": "USA"},
		map[string]interface{}{"name": "Bob", "age": 25, "city": "LA", "country": "USA"},
	}
	result5, _ := uql.UQL(`project "name", "city"`, &uql.Options{Data: data})
	fmt.Printf("Result: %v\n\n", result5)

	// Example 6: Order By
	fmt.Println("Example 6: Order by age descending")
	scores := []interface{}{
		map[string]interface{}{"name": "Charlie", "score": 85},
		map[string]interface{}{"name": "Alice", "score": 92},
		map[string]interface{}{"name": "Bob", "score": 78},
	}
	result6, _ := uql.UQL(`order by "score" desc`, &uql.Options{Data: scores})
	fmt.Printf("Result: %v\n\n", result6)

	// Example 7: Pipeline
	fmt.Println("Example 7: Pipeline - order, limit, and project")
	result7, _ := uql.UQL(`order by "score" desc | limit 2 | project "name", "score"`, &uql.Options{Data: scores})
	fmt.Printf("Result: %v\n\n", result7)

	// Example 8: Parse JSON
	fmt.Println("Example 8: Parse JSON")
	jsonData := `{"users": [{"name": "foo", "age": 25}, {"name": "bar", "age": 30}]}`
	result8, _ := uql.UQL(`parse-json | scope "users"`, &uql.Options{Data: jsonData})
	fmt.Printf("Result: %v\n\n", result8)

	// Example 9: Complex query with JSON parsing
	fmt.Println("Example 9: Complex query - parse JSON, order, and project")
	complexJSON, _ := json.Marshal([]interface{}{
		map[string]interface{}{"name": "foo", "age": 2, "location": "uk"},
		map[string]interface{}{"name": "bar", "age": 3, "location": "usa"},
	})
	query := `parse-json | order by "name" asc | project "name", "location"`
	result9, err := uql.UQL(query, &uql.Options{Data: string(complexJSON)})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Result: %v\n\n", result9)
	}

	// Example 10: Distinct
	fmt.Println("Example 10: Distinct values")
	dataWithDuplicates := []interface{}{
		map[string]interface{}{"name": "foo", "age": 2},
		map[string]interface{}{"name": "bar", "age": 3},
		map[string]interface{}{"name": "foo", "age": 2},
	}
	result10, _ := uql.UQL("distinct", &uql.Options{Data: dataWithDuplicates})
	fmt.Printf("Result: %v\n\n", result10)

	// Example 11: Project-away
	fmt.Println("Example 11: Project-away to remove fields")
	fullData := []interface{}{
		map[string]interface{}{"name": "Alice", "age": 30, "ssn": "123-45-6789", "city": "NYC"},
		map[string]interface{}{"name": "Bob", "age": 25, "ssn": "987-65-4321", "city": "LA"},
	}
	result11, _ := uql.UQL(`project-away "ssn"`, &uql.Options{Data: fullData})
	fmt.Printf("Result: %v\n\n", result11)

	// Example 12: Scope to navigate nested data
	fmt.Println("Example 12: Scope to access nested data")
	nestedData := map[string]interface{}{
		"metadata": map[string]interface{}{
			"version": "1.0",
			"timestamp": "2024-01-01",
		},
		"result": map[string]interface{}{
			"value": 42,
			"status": "success",
		},
	}
	result12, _ := uql.UQL(`scope "result"`, &uql.Options{Data: nestedData})
	fmt.Printf("Result: %v\n\n", result12)

	fmt.Println("=== All Examples Completed ===")
}
