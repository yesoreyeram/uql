package main

import (
	"encoding/json"
	"fmt"
	"log"

	uql "github.com/yesoreyeram/uql/uql-go"
)

func main() {
	fmt.Println("=== JSONata Examples ===\n")

	// Example 1: Basic aggregation
	fmt.Println("1. Basic Aggregation:")
	data1 := map[string]interface{}{
		"example": []interface{}{
			map[string]interface{}{"value": 4},
			map[string]interface{}{"value": 7},
			map[string]interface{}{"value": 13},
		},
	}

	result1, err := uql.UQL(`jsonata "$sum(example.value)"`, &uql.Options{Data: data1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query: jsonata \"$sum(example.value)\"\n")
	fmt.Printf("Result: %v\n\n", result1)

	// Example 2: Extract array values
	fmt.Println("2. Extract Array Values:")
	result2, err := uql.UQL(`jsonata "example.value"`, &uql.Options{Data: data1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query: jsonata \"example.value\"\n")
	jsonBytes, _ := json.MarshalIndent(result2, "", "  ")
	fmt.Printf("Result: %s\n\n", string(jsonBytes))

	// Example 3: Filter data
	fmt.Println("3. Filter Data:")
	data3 := map[string]interface{}{
		"Countries": []interface{}{
			map[string]interface{}{"Country": "India", "Population": 1380},
			map[string]interface{}{"Country": "America", "Population": 331},
			map[string]interface{}{"Country": "United Kingdom", "Population": 68},
		},
	}

	result3, err := uql.UQL(`scope "Countries" | jsonata "$[Country='India']"`, &uql.Options{Data: data3})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query: scope \"Countries\" | jsonata \"$[Country='India']\"\n")
	jsonBytes, _ = json.MarshalIndent(result3, "", "  ")
	fmt.Printf("Result: %s\n\n", string(jsonBytes))

	// Example 4: Complex library query
	fmt.Println("4. Complex Library Query:")
	libraryData := map[string]interface{}{
		"library": map[string]interface{}{
			"books": []interface{}{
				map[string]interface{}{
					"title": "Structure and Interpretation of Computer Programs",
					"price": 38.9,
				},
				map[string]interface{}{
					"title": "The C Programming Language",
					"price": 33.59,
				},
				map[string]interface{}{
					"title": "Compilers: Principles, Techniques, and Tools",
					"price": 23.38,
				},
			},
		},
	}

	result4, err := uql.UQL(`jsonata "$sum(library.books.price)"`, &uql.Options{Data: libraryData})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query: jsonata \"$sum(library.books.price)\"\n")
	fmt.Printf("Result: %v\n\n", result4)

	// Example 5: Combining JSONata with other commands
	fmt.Println("5. Combining JSONata with Other Commands:")
	result5, err := uql.UQL(`scope "library" | jsonata "books" | count`, &uql.Options{Data: libraryData})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query: scope \"library\" | jsonata \"books\" | count\n")
	fmt.Printf("Result: %v books\n\n", result5)

	// Example 6: Multiple filters
	fmt.Println("6. Multiple Filters with IN operator:")
	result6, err := uql.UQL(`scope "Countries" | jsonata "$[Country in ['India', 'United Kingdom']]"`, &uql.Options{Data: data3})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query: scope \"Countries\" | jsonata \"$[Country in ['India', 'United Kingdom']]\"\n")
	jsonBytes, _ = json.MarshalIndent(result6, "", "  ")
	fmt.Printf("Result: %s\n\n", string(jsonBytes))

	// Example 7: Object transformation
	fmt.Println("7. Object Transformation:")
	data7 := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"firstName": "John", "lastName": "Doe", "age": 30},
			map[string]interface{}{"firstName": "Jane", "lastName": "Smith", "age": 25},
		},
	}

	result7, err := uql.UQL(`jsonata "users.{'name': firstName & ' ' & lastName, 'age': age}"`, &uql.Options{Data: data7})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Query: jsonata \"users.{'name': firstName & ' ' & lastName, 'age': age}\"\n")
	jsonBytes, _ = json.MarshalIndent(result7, "", "  ")
	fmt.Printf("Result: %s\n", string(jsonBytes))
}
