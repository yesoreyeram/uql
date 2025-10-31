package uql

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
)

// evalCount evaluates a count command
func evalCount(prev CommandResult, cmd Command) (CommandResult, error) {
	output := prev.Output
	if output == nil {
		return CommandResult{Output: 0, Context: prev.Context}, nil
	}

	val := reflect.ValueOf(output)
	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		return CommandResult{Output: val.Len(), Context: prev.Context}, nil
	}

	return CommandResult{Output: 1, Context: prev.Context}, nil
}

// evalLimit evaluates a limit command
func evalLimit(prev CommandResult, cmd Command) (CommandResult, error) {
	limit, ok := cmd.Value.(float64)
	if !ok {
		return prev, errors.New("limit value must be a number")
	}

	output := prev.Output
	if output == nil {
		return prev, nil
	}

	val := reflect.ValueOf(output)
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return prev, nil
	}

	limitInt := int(limit)
	if limitInt >= val.Len() {
		return prev, nil
	}

	// Create a new slice with limited elements
	result := make([]interface{}, limitInt)
	for i := 0; i < limitInt; i++ {
		result[i] = val.Index(i).Interface()
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
}

// evalOrderBy evaluates an order by command
func evalOrderBy(prev CommandResult, cmd Command) (CommandResult, error) {
	args, ok := cmd.Value.([]OrderByArg)
	if !ok || len(args) == 0 {
		return prev, errors.New("invalid order by arguments")
	}

	output := prev.Output
	if output == nil {
		return prev, nil
	}

	// Convert to slice of interfaces
	slice, err := toSlice(output)
	if err != nil {
		return prev, nil
	}

	// Sort the slice
	sort.Slice(slice, func(i, j int) bool {
		for _, arg := range args {
			vi := getValue(slice[i], arg.Field)
			vj := getValue(slice[j], arg.Field)

			cmp := compareValues(vi, vj)
			if cmp != 0 {
				if arg.Direction == "desc" {
					return cmp > 0
				}
				return cmp < 0
			}
		}
		return false
	})

	return CommandResult{Output: slice, Context: prev.Context}, nil
}

// evalScope evaluates a scope command
func evalScope(prev CommandResult, cmd Command) (CommandResult, error) {
	ref, ok := cmd.Value.(TypedValue)
	if !ok {
		return prev, errors.New("invalid scope argument")
	}

	output := prev.Output
	if output == nil {
		return prev, nil
	}

	value := getValue(output, ref.Value.(string))
	return CommandResult{Output: value, Context: prev.Context}, nil
}

// evalDistinct evaluates a distinct command
func evalDistinct(prev CommandResult, cmd Command) (CommandResult, error) {
	output := prev.Output
	if output == nil {
		return prev, nil
	}

	slice, err := toSlice(output)
	if err != nil {
		return prev, nil
	}

	if cmd.Value == nil {
		// Distinct on entire objects
		seen := make(map[string]bool)
		result := make([]interface{}, 0)

		for _, item := range slice {
			key, err := json.Marshal(item)
			if err != nil {
				continue
			}
			if !seen[string(key)] {
				seen[string(key)] = true
				result = append(result, item)
			}
		}
		return CommandResult{Output: result, Context: prev.Context}, nil
	}

	// Distinct on specific field
	ref, ok := cmd.Value.(TypedValue)
	if !ok {
		return prev, errors.New("invalid distinct argument")
	}

	field := ref.Value.(string)
	seen := make(map[interface{}]bool)
	result := make([]interface{}, 0)

	for _, item := range slice {
		val := getValue(item, field)
		key := val
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
}

// Placeholder implementations for other commands
func evalWhere(prev CommandResult, cmd Command) (CommandResult, error) {
	// TODO: Implement where clause evaluation
	return prev, nil
}

func evalMvExpand(prev CommandResult, cmd Command) (CommandResult, error) {
	mvExpandVal, ok := cmd.Value.(MvExpandValue)
	if !ok {
		return prev, errors.New("invalid mv-expand arguments")
	}

	output := prev.Output
	if output == nil {
		return prev, nil
	}

	slice, err := toSlice(output)
	if err != nil {
		return prev, nil
	}

	result := make([]interface{}, 0)

	for _, item := range slice {
		// Get the array to expand
		expandingItem := getValue(item, mvExpandVal.Field)

		// Check if it's an array
		if expandingItem == nil {
			continue
		}

		expandSlice, err := toSlice(expandingItem)
		if err != nil || len(expandSlice) == 0 {
			continue
		}

		// Expand each element
		for _, element := range expandSlice {
			// Create a copy of the item
			newItem := make(map[string]interface{})
			if m, ok := toMap(item); ok {
				for k, v := range m {
					newItem[k] = v
				}
			}

			// Set the expanded value
			fieldName := mvExpandVal.Field
			if mvExpandVal.Alias != "" {
				fieldName = mvExpandVal.Alias
				// Remove the original field if alias is used
				delete(newItem, mvExpandVal.Field)
			}
			newItem[fieldName] = element

			result = append(result, newItem)
		}
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
}

func evalCommandFunc(prev CommandResult, cmd Command) (CommandResult, error) {
	// TODO: Implement command function
	return prev, nil
}

func evalJSONata(prev CommandResult, cmd Command) (CommandResult, error) {
	// TODO: Implement JSONata
	return prev, nil
}

func evalRange(prev CommandResult, cmd Command) (CommandResult, error) {
	rangeVal, ok := cmd.Value.(RangeValue)
	if !ok {
		return prev, errors.New("invalid range arguments")
	}

	var result []interface{}

	// Check if it's numeric range
	if startNum, ok := rangeVal.Start.(float64); ok {
		endNum, ok := rangeVal.End.(float64)
		if !ok {
			return prev, errors.New("start and end must be both numbers or both strings")
		}

		stepNum := 1.0
		if s, ok := rangeVal.Step.(float64); ok {
			stepNum = s
		}

		if stepNum <= 0 {
			return prev, errors.New("step must be positive")
		}

		// Generate numeric range
		for i := startNum; i <= endNum; i += stepNum {
			result = append(result, i)
		}
	} else {
		// String range - for now, just return the start and end as strings
		// Full implementation would parse date/time strings and generate series
		// This is a simplified version
		result = append(result, rangeVal.Start)
		if rangeVal.Start != rangeVal.End {
			result = append(result, rangeVal.End)
		}
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
}
