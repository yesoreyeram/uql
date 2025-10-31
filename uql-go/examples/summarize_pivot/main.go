package main

import (
	"fmt"

	uql "github.com/yesoreyeram/uql/uql-go"
)

func main() {
	fmt.Println("=== UQL Summarize and Pivot Examples ===\n")

	// Example 1: Summarize with single group by
	fmt.Println("Example 1: Summarize - first patron by country")
	users := []interface{}{
		map[string]interface{}{"patron": "a", "age": 48, "country": "foo"},
		map[string]interface{}{"patron": "b", "age": 34, "country": "foo"},
		map[string]interface{}{"patron": "c", "age": 12, "country": "bar"},
		map[string]interface{}{"patron": "d", "age": 40, "country": "bar"},
		map[string]interface{}{"patron": "e", "age": 36, "country": "baz"},
	}
	result1, _ := uql.UQL(`summarize "user"=first("patron") by "country"`, &uql.Options{Data: users})
	fmt.Printf("Result: %v\n\n", result1)

	// Example 2: Summarize with sum
	fmt.Println("Example 2: Summarize - sum of ages by country")
	result2, _ := uql.UQL(`summarize "total_age"=sum("age") by "country"`, &uql.Options{Data: users})
	fmt.Printf("Result: %v\n\n", result2)

	// Example 3: Summarize with multiple metrics
	fmt.Println("Example 3: Summarize - multiple metrics by country")
	result3, _ := uql.UQL(`summarize "total_age"=sum("age"),"avg_age"=mean("age") by "country"`, &uql.Options{Data: users})
	fmt.Printf("Result: %v\n\n", result3)

	// Example 4: Summarize with multiple group by
	fmt.Println("Example 4: Summarize - multiple group by")
	data := []interface{}{
		map[string]interface{}{"age": 1, "name": "foo1", "city": "chennai", "country": "india"},
		map[string]interface{}{"age": 2, "name": "foo1", "city": "chennai", "country": "india"},
		map[string]interface{}{"age": 3, "name": "foo1", "city": "mumbai", "country": "india"},
		map[string]interface{}{"age": 4, "name": "foo1", "city": "london", "country": "england"},
	}
	result4, _ := uql.UQL(`summarize "age"=sum("age") by "country", "city"`, &uql.Options{Data: data})
	fmt.Printf("Result: %v\n\n", result4)

	// Example 5: Pivot - count all
	fmt.Println("Example 5: Pivot - count all items")
	fruits := []interface{}{
		map[string]interface{}{"fruit": "apple", "size": "sm", "qty": 1},
		map[string]interface{}{"fruit": "apple", "size": "md", "qty": 2},
		map[string]interface{}{"fruit": "apple", "size": "lg", "qty": 3},
		map[string]interface{}{"fruit": "banana", "size": "sm", "qty": 1},
		map[string]interface{}{"fruit": "banana", "size": "lg", "qty": 6},
		map[string]interface{}{"fruit": "banana", "size": "xl", "qty": 5},
	}
	result5, _ := uql.UQL(`pivot count()`, &uql.Options{Data: fruits})
	fmt.Printf("Result: %v\n\n", result5)

	// Example 6: Pivot - sum with row field
	fmt.Println("Example 6: Pivot - sum by fruit")
	result6, _ := uql.UQL(`pivot sum("qty"), "fruit"`, &uql.Options{Data: fruits})
	fmt.Printf("Result: %v\n\n", result6)

	// Example 7: Pivot - sum with row and column fields
	fmt.Println("Example 7: Pivot - sum by fruit and size (cross-tab)")
	result7, _ := uql.UQL(`pivot sum("qty"), "fruit", "size"`, &uql.Options{Data: fruits})
	fmt.Printf("Result: %v\n\n", result7)

	// Example 8: Pivot - max with row and column fields
	fmt.Println("Example 8: Pivot - max quantity by fruit and size")
	result8, _ := uql.UQL(`pivot max("qty"), "fruit", "size"`, &uql.Options{Data: fruits})
	fmt.Printf("Result: %v\n\n", result8)

	fmt.Println("=== All Examples Completed ===")
}
