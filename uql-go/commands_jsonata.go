package uql

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/xiatechs/jsonata-go"
)

// evaluateJSONata evaluates a JSONata expression against the given data
func evaluateJSONata(expression string, data interface{}) (interface{}, error) {
	// Convert data to JSON-compatible format if needed
	jsonData, err := convertToJSONData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to convert data: %w", err)
	}

	// Compile and evaluate the JSONata expression
	expr, err := jsonata.Compile(expression)
	if err != nil {
		return nil, fmt.Errorf("failed to compile expression: %w", err)
	}

	result, err := expr.Eval(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate expression: %w", err)
	}

	// Clean up the result (remove JSONata-specific metadata if present)
	cleanedResult := cleanJSONataResult(result)

	return cleanedResult, nil
}

// convertToJSONData converts Go data to JSON-compatible format
func convertToJSONData(data interface{}) (interface{}, error) {
	// If data is already nil, return as is
	if data == nil {
		return nil, nil
	}

	// Convert to JSON and back to ensure proper format
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// cleanJSONataResult removes JSONata-specific metadata from results
func cleanJSONataResult(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Map:
		// Create a new map without JSONata metadata fields
		result := make(map[string]interface{})
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key().Interface()
			keyStr, ok := key.(string)
			if !ok {
				continue
			}

			// Skip JSONata-specific metadata fields
			if keyStr == "sequence" || keyStr == "keepSingleton" {
				continue
			}

			value := iter.Value().Interface()
			result[keyStr] = cleanJSONataResult(value)
		}
		return result

	case reflect.Slice, reflect.Array:
		// Clean each element in the slice/array
		length := v.Len()
		result := make([]interface{}, length)
		for i := 0; i < length; i++ {
			result[i] = cleanJSONataResult(v.Index(i).Interface())
		}
		return result

	default:
		return data
	}
}
