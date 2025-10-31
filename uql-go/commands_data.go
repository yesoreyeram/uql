package uql

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
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
	conditions, ok := cmd.Value.([]TypedValue)
	if !ok {
		return prev, errors.New("invalid where arguments")
	}

	output := prev.Output
	if output == nil {
		return prev, nil
	}

	slice, err := toSlice(output)
	if err != nil {
		return prev, nil
	}

	if len(conditions) < 3 || conditions[1].Type != "operation" {
		return prev, errors.New("invalid where clause: expected field operator value")
	}

	result := make([]interface{}, 0)
	for _, item := range slice {
		if matchesCondition(item, conditions) {
			result = append(result, item)
		}
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
}

// matchesCondition checks if an item matches the where condition
func matchesCondition(item interface{}, conditions []TypedValue) bool {
	if len(conditions) < 3 {
		return false
	}

	// Get left-hand side value
	lhs := getWhereValue(item, conditions[0])

	// Get operator
	operator := conditions[1].Value.(string)

	// Get right-hand side value
	rhs := getWhereValue(item, conditions[2])

	// Handle array values for 'in', 'between', etc.
	var rhsArray []interface{}
	if conditions[2].Type == "value_array" {
		if arr, ok := conditions[2].Value.([]interface{}); ok {
			rhsArray = make([]interface{}, 0, len(arr))
			for _, v := range arr {
				if tv, ok := v.(TypedValue); ok {
					rhsArray = append(rhsArray, getWhereValue(item, tv))
				} else {
					rhsArray = append(rhsArray, v)
				}
			}
		}
	}

	// Apply operator
	switch operator {
	case "==":
		return compareEqual(lhs, rhs)
	case "!=":
		return !compareEqual(lhs, rhs)
	case ">":
		return compareGreater(lhs, rhs)
	case ">=":
		return compareGreaterEqual(lhs, rhs)
	case "<":
		return compareLess(lhs, rhs)
	case "<=":
		return compareLessEqual(lhs, rhs)
	case "=~":
		return caseInsensitiveContains(lhs, rhs)
	case "!~":
		return !caseInsensitiveContains(lhs, rhs)
	case "contains":
		return caseInsensitiveContains(lhs, rhs)
	case "!contains", "not contains":
		return !caseInsensitiveContains(lhs, rhs)
	case "contains_cs":
		return caseSensitiveContains(lhs, rhs)
	case "!contains_cs", "not contains_cs":
		return !caseSensitiveContains(lhs, rhs)
	case "startswith":
		return caseInsensitiveStartsWith(lhs, rhs)
	case "!startswith":
		return !caseInsensitiveStartsWith(lhs, rhs)
	case "startswith_cs":
		return caseSensitiveStartsWith(lhs, rhs)
	case "!startswith_cs":
		return !caseSensitiveStartsWith(lhs, rhs)
	case "endswith":
		return caseInsensitiveEndsWith(lhs, rhs)
	case "!endswith":
		return !caseInsensitiveEndsWith(lhs, rhs)
	case "endswith_cs":
		return caseSensitiveEndsWith(lhs, rhs)
	case "!endswith_cs":
		return !caseSensitiveEndsWith(lhs, rhs)
	case "in":
		return inArray(lhs, rhsArray)
	case "!in":
		return !inArray(lhs, rhsArray)
	case "in~":
		return inArrayCaseInsensitive(lhs, rhsArray)
	case "!in~":
		return !inArrayCaseInsensitive(lhs, rhsArray)
	case "between":
		if len(rhsArray) >= 2 {
			return betweenValues(lhs, rhsArray[0], rhsArray[1])
		}
	case "inside":
		if len(rhsArray) >= 2 {
			return insideValues(lhs, rhsArray[0], rhsArray[1])
		}
	case "outside":
		if len(rhsArray) >= 2 {
			return outsideValues(lhs, rhsArray[0], rhsArray[1])
		}
	case "matches regex":
		return matchesRegex(lhs, rhs)
	case "!matches regex":
		return !matchesRegex(lhs, rhs)
	}

	return false
}

// getWhereValue extracts the actual value from a TypedValue
func getWhereValue(item interface{}, tv TypedValue) interface{} {
	switch tv.Type {
	case "ref":
		// It's a field reference
		if fieldName, ok := tv.Value.(string); ok {
			return getValue(item, fieldName)
		}
	case "number", "string":
		return tv.Value
	case "value_array":
		return tv.Value
	}
	return tv.Value
}

// String comparison functions
func caseInsensitiveContains(lhs, rhs interface{}) bool {
	lhsStr := strings.ToLower(fmt.Sprintf("%v", lhs))
	rhsStr := strings.ToLower(fmt.Sprintf("%v", rhs))
	return strings.Contains(lhsStr, rhsStr)
}

func caseSensitiveContains(lhs, rhs interface{}) bool {
	lhsStr := fmt.Sprintf("%v", lhs)
	rhsStr := fmt.Sprintf("%v", rhs)
	return strings.Contains(lhsStr, rhsStr)
}

func caseInsensitiveStartsWith(lhs, rhs interface{}) bool {
	lhsStr := strings.ToLower(fmt.Sprintf("%v", lhs))
	rhsStr := strings.ToLower(fmt.Sprintf("%v", rhs))
	return strings.HasPrefix(lhsStr, rhsStr)
}

func caseSensitiveStartsWith(lhs, rhs interface{}) bool {
	lhsStr := fmt.Sprintf("%v", lhs)
	rhsStr := fmt.Sprintf("%v", rhs)
	return strings.HasPrefix(lhsStr, rhsStr)
}

func caseInsensitiveEndsWith(lhs, rhs interface{}) bool {
	lhsStr := strings.ToLower(fmt.Sprintf("%v", lhs))
	rhsStr := strings.ToLower(fmt.Sprintf("%v", rhs))
	return strings.HasSuffix(lhsStr, rhsStr)
}

func caseSensitiveEndsWith(lhs, rhs interface{}) bool {
	lhsStr := fmt.Sprintf("%v", lhs)
	rhsStr := fmt.Sprintf("%v", rhs)
	return strings.HasSuffix(lhsStr, rhsStr)
}

func matchesRegex(lhs, rhs interface{}) bool {
	lhsStr := fmt.Sprintf("%v", lhs)
	rhsStr := fmt.Sprintf("%v", rhs)
	re, err := regexp.Compile(rhsStr)
	if err != nil {
		return false
	}
	return re.MatchString(lhsStr)
}

// Comparison functions
func compareEqual(lhs, rhs interface{}) bool {
	// Handle nil cases
	if lhs == nil && rhs == nil {
		return true
	}
	if lhs == nil || rhs == nil {
		return false
	}

	// Try numeric comparison
	lhsNum, lhsOk := toNumber(lhs)
	rhsNum, rhsOk := toNumber(rhs)
	if lhsOk && rhsOk {
		return lhsNum == rhsNum
	}

	// String comparison
	return fmt.Sprintf("%v", lhs) == fmt.Sprintf("%v", rhs)
}

func compareGreater(lhs, rhs interface{}) bool {
	lhsNum, lhsOk := toNumber(lhs)
	rhsNum, rhsOk := toNumber(rhs)
	if lhsOk && rhsOk {
		return lhsNum > rhsNum
	}
	return fmt.Sprintf("%v", lhs) > fmt.Sprintf("%v", rhs)
}

func compareGreaterEqual(lhs, rhs interface{}) bool {
	return compareGreater(lhs, rhs) || compareEqual(lhs, rhs)
}

func compareLess(lhs, rhs interface{}) bool {
	lhsNum, lhsOk := toNumber(lhs)
	rhsNum, rhsOk := toNumber(rhs)
	if lhsOk && rhsOk {
		return lhsNum < rhsNum
	}
	return fmt.Sprintf("%v", lhs) < fmt.Sprintf("%v", rhs)
}

func compareLessEqual(lhs, rhs interface{}) bool {
	return compareLess(lhs, rhs) || compareEqual(lhs, rhs)
}

// Array comparison functions
func inArray(lhs interface{}, arr []interface{}) bool {
	for _, item := range arr {
		if compareEqual(lhs, item) {
			return true
		}
	}
	return false
}

func inArrayCaseInsensitive(lhs interface{}, arr []interface{}) bool {
	lhsStr := strings.ToLower(fmt.Sprintf("%v", lhs))
	for _, item := range arr {
		itemStr := strings.ToLower(fmt.Sprintf("%v", item))
		if lhsStr == itemStr {
			return true
		}
	}
	return false
}

func betweenValues(lhs, lower, upper interface{}) bool {
	return (compareGreaterEqual(lhs, lower) && compareLessEqual(lhs, upper))
}

func insideValues(lhs, lower, upper interface{}) bool {
	return (compareGreater(lhs, lower) && compareLess(lhs, upper))
}

func outsideValues(lhs, lower, upper interface{}) bool {
	return (compareLess(lhs, lower) || compareGreater(lhs, upper))
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
