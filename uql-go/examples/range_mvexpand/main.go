package main

import (
	"fmt"

	uql "github.com/yesoreyeram/uql/uql-go"
)

func main() {
	fmt.Println("=== UQL Range and MV-Expand Examples ===\n")

	// Example 1: Basic range
	fmt.Println("Example 1: Range - Generate numbers from 1 to 5")
	result1, _ := uql.UQL(`range from 1 to 5`, nil)
	fmt.Printf("Result: %v\n\n", result1)

	// Example 2: Range with step
	fmt.Println("Example 2: Range - Generate numbers from 0 to 10 with step 2")
	result2, _ := uql.UQL(`range from 0 to 10 step 2`, nil)
	fmt.Printf("Result: %v\n\n", result2)

	// Example 3: Range with decimal step
	fmt.Println("Example 3: Range - Generate numbers from 0 to 2 with step 0.5")
	result3, _ := uql.UQL(`range from 0 to 2 step 0.5`, nil)
	fmt.Printf("Result: %v\n\n", result3)

	// Example 4: Range with pipeline
	fmt.Println("Example 4: Range - Generate and limit results")
	result4, _ := uql.UQL(`range from 1 to 20 step 3 | limit 3`, nil)
	fmt.Printf("Result: %v\n\n", result4)

	// Example 5: Basic mv-expand
	fmt.Println("Example 5: MV-Expand - Expand array field")
	data5 := []interface{}{
		map[string]interface{}{"group": "A", "users": []interface{}{"user a1", "user a2"}},
		map[string]interface{}{"group": "B", "users": []interface{}{"user b1"}},
	}
	result5, _ := uql.UQL(`mv-expand "users"`, &uql.Options{Data: data5})
	fmt.Printf("Result: %v\n\n", result5)

	// Example 6: MV-Expand with alias
	fmt.Println("Example 6: MV-Expand - Expand with alias")
	data6 := []interface{}{
		map[string]interface{}{"group": "A", "users": []interface{}{"user a1", "user a2"}},
		map[string]interface{}{"group": "B", "users": []interface{}{"user b1"}},
	}
	result6, _ := uql.UQL(`mv-expand "user"="users"`, &uql.Options{Data: data6})
	fmt.Printf("Result: %v\n\n", result6)

	// Example 7: MV-Expand ignoring empty/non-array values
	fmt.Println("Example 7: MV-Expand - Handles empty arrays and missing fields")
	data7 := []interface{}{
		map[string]interface{}{"group": "A", "users": []interface{}{"user a1", "user a2"}},
		map[string]interface{}{"group": "B", "users": []interface{}{}},
		map[string]interface{}{"group": "C"},
		map[string]interface{}{"group": "D", "users": []interface{}{"user d1"}},
	}
	result7, _ := uql.UQL(`mv-expand "user"="users"`, &uql.Options{Data: data7})
	fmt.Printf("Result: %v\n\n", result7)

	// Example 8: MV-Expand with objects
	fmt.Println("Example 8: MV-Expand - Expand array of objects")
	data8 := []interface{}{
		map[string]interface{}{
			"server": "server1",
			"disks": []interface{}{
				map[string]interface{}{"drive": "C", "size": 1024},
				map[string]interface{}{"drive": "D", "size": 2048},
			},
		},
		map[string]interface{}{
			"server": "server2",
			"disks": []interface{}{
				map[string]interface{}{"drive": "C", "size": 2048},
			},
		},
	}
	result8, _ := uql.UQL(`mv-expand "disk"="disks"`, &uql.Options{Data: data8})
	fmt.Printf("Result: %v\n\n", result8)

	fmt.Println("=== All Examples Completed ===")
}
