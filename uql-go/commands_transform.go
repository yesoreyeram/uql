package uql

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// evalProject evaluates a project command
func evalProject(prev CommandResult, cmd Command) (CommandResult, error) {
	projections, ok := cmd.Value.([]interface{})
	if !ok {
		return prev, errors.New("invalid project arguments")
	}

	output := prev.Output
	if output == nil {
		return prev, nil
	}

	// Handle array of objects
	slice, err := toSlice(output)
	if err == nil {
		result := make([]interface{}, 0, len(slice))
		for _, item := range slice {
			projected := projectItem(item, projections)
			result = append(result, projected)
		}
		return CommandResult{Output: result, Context: prev.Context}, nil
	}

	// Handle single object
	projected := projectItem(output, projections)
	return CommandResult{Output: projected, Context: prev.Context}, nil
}

// projectItem projects a single item
func projectItem(item interface{}, projections []interface{}) interface{} {
	result := make(map[string]interface{})

	for _, proj := range projections {
		switch p := proj.(type) {
		case TypedValue:
			if p.Type == "ref" {
				field := p.Value.(string)
				key := p.Alias
				if key == "" {
					key = field
				}
				result[key] = getValue(item, field)
			}
		case FunctionCall:
			key := p.Alias
			if key == "" {
				key = string(p.Operator)
			}
			args := make([]interface{}, len(p.Args))
			for i, arg := range p.Args {
				args[i] = resolveArg(item, arg)
			}
			value := evaluateFunction(p.Operator, args)
			result[key] = value
		}
	}

	return result
}

// evalProjectAway evaluates a project-away command
func evalProjectAway(prev CommandResult, cmd Command) (CommandResult, error) {
	fields, ok := cmd.Value.([]TypedValue)
	if !ok {
		return prev, errors.New("invalid project-away arguments")
	}

	output := prev.Output
	if output == nil {
		return prev, nil
	}

	// Get field names to remove
	removeFields := make(map[string]bool)
	for _, f := range fields {
		if f.Type == "ref" {
			removeFields[f.Value.(string)] = true
		}
	}

	// Handle array of objects
	slice, err := toSlice(output)
	if err == nil {
		result := make([]interface{}, 0, len(slice))
		for _, item := range slice {
			projected := removeFieldsFromItem(item, removeFields)
			result = append(result, projected)
		}
		return CommandResult{Output: result, Context: prev.Context}, nil
	}

	// Handle single object
	projected := removeFieldsFromItem(output, removeFields)
	return CommandResult{Output: projected, Context: prev.Context}, nil
}

// removeFieldsFromItem removes specified fields from an item
func removeFieldsFromItem(item interface{}, removeFields map[string]bool) interface{} {
	m, ok := toMap(item)
	if !ok {
		return item
	}

	result := make(map[string]interface{})
	for k, v := range m {
		if !removeFields[k] {
			result[k] = v
		}
	}

	return result
}

// evalProjectReorder evaluates a project-reorder command
func evalProjectReorder(prev CommandResult, cmd Command) (CommandResult, error) {
	fields, ok := cmd.Value.([]TypedValue)
	if !ok {
		return prev, errors.New("invalid project-reorder arguments")
	}

	output := prev.Output
	if output == nil {
		return prev, nil
	}

	// Get ordered field names
	orderedFields := make([]string, 0, len(fields))
	for _, f := range fields {
		if f.Type == "ref" {
			orderedFields = append(orderedFields, f.Value.(string))
		}
	}

	// Handle array of objects
	slice, err := toSlice(output)
	if err == nil {
		result := make([]interface{}, 0, len(slice))
		for _, item := range slice {
			reordered := reorderFieldsInItem(item, orderedFields)
			result = append(result, reordered)
		}
		return CommandResult{Output: result, Context: prev.Context}, nil
	}

	// Handle single object
	reordered := reorderFieldsInItem(output, orderedFields)
	return CommandResult{Output: reordered, Context: prev.Context}, nil
}

// reorderFieldsInItem reorders fields in an item
func reorderFieldsInItem(item interface{}, orderedFields []string) interface{} {
	m, ok := toMap(item)
	if !ok {
		return item
	}

	result := make(map[string]interface{})

	// Add ordered fields first
	for _, field := range orderedFields {
		if val, exists := m[field]; exists {
			result[field] = val
		}
	}

	// Add remaining fields
	for k, v := range m {
		if _, exists := result[k]; !exists {
			result[k] = v
		}
	}

	return result
}

// evalExtend evaluates an extend command
func evalExtend(prev CommandResult, cmd Command) (CommandResult, error) {
	extensions, ok := cmd.Value.([]interface{})
	if !ok {
		return prev, errors.New("invalid extend arguments")
	}

	output := prev.Output
	if output == nil {
		return prev, nil
	}

	// Handle array of objects
	slice, err := toSlice(output)
	if err == nil {
		result := make([]interface{}, 0, len(slice))
		for _, item := range slice {
			extended := extendItem(item, extensions)
			result = append(result, extended)
		}
		return CommandResult{Output: result, Context: prev.Context}, nil
	}

	// Handle single object
	extended := extendItem(output, extensions)
	return CommandResult{Output: extended, Context: prev.Context}, nil
}

// extendItem extends a single item with new fields
func extendItem(item interface{}, extensions []interface{}) interface{} {
	m, ok := toMap(item)
	if !ok {
		return item
	}

	result := make(map[string]interface{})
	// Copy existing fields
	for k, v := range m {
		result[k] = v
	}

	// Add new fields
	for _, ext := range extensions {
		switch e := ext.(type) {
		case TypedValue:
			if e.Type == "ref" {
				field := e.Value.(string)
				key := e.Alias
				if key == "" {
					key = field
				}
				result[key] = getValue(item, field)
			}
		case FunctionCall:
			key := e.Alias
			if key == "" {
				key = string(e.Operator)
			}
			args := make([]interface{}, len(e.Args))
			for i, arg := range e.Args {
				args[i] = resolveArg(item, arg)
			}
			value := evaluateFunction(e.Operator, args)
			result[key] = value
		}
	}

	return result
}

// evalSummarize evaluates a summarize command
func evalSummarize(prev CommandResult, cmd Command) (CommandResult, error) {
	item, ok := cmd.Value.(SummarizeItem)
	if !ok {
		return prev, errors.New("invalid summarize arguments")
	}

	output := prev.Output
	if output == nil {
		return prev, nil
	}

	slice, err := toSlice(output)
	if err != nil {
		return prev, nil
	}

	var result interface{}

	// Group by fields
	if len(item.By) == 0 {
		// No grouping, summarize all data
		result = summarizeGroup(nil, item.Metrics, slice)
	} else if len(item.By) == 1 {
		// Single field grouping
		groupByKey := item.By[0].Value.(string)
		groups := groupByField(slice, groupByKey)

		resultArray := make([]interface{}, 0, len(groups))
		for key, group := range groups {
			summarized := summarizeGroup(map[string]interface{}{groupByKey: key}, item.Metrics, group)
			resultArray = append(resultArray, summarized)
		}
		result = resultArray
	} else {
		// Multiple field grouping
		groups := groupByFields(slice, item.By)

		resultArray := make([]interface{}, 0, len(groups))
		for _, group := range groups {
			if len(group) == 0 {
				continue
			}
			// Extract group keys from first item
			baseObj := make(map[string]interface{})
			for _, byField := range item.By {
				fieldName := byField.Value.(string)
				baseObj[fieldName] = getValue(group[0], fieldName)
			}
			summarized := summarizeGroup(baseObj, item.Metrics, group)
			resultArray = append(resultArray, summarized)
		}
		result = resultArray
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
}

// evalPivot evaluates a pivot command
func evalPivot(prev CommandResult, cmd Command) (CommandResult, error) {
	item, ok := cmd.Value.(PivotItem)
	if !ok {
		return prev, errors.New("invalid pivot arguments")
	}

	input := prev.Output
	if input == nil {
		return CommandResult{Output: nil, Context: prev.Context}, nil
	}

	slice, err := toSlice(input)
	if err != nil {
		return CommandResult{Output: nil, Context: prev.Context}, nil
	}

	// No fields - just aggregate all data
	if len(item.Fields) == 0 {
		result := summarizeGroup(nil, []SummarizeAssignment{item.Metric}, slice)
		// Try to extract the value from the result
		if m, ok := result.(map[string]interface{}); ok {
			// Try different possible key names
			metricName := string(item.Metric.Operator)
			if val, exists := m[metricName]; exists {
				return CommandResult{Output: val, Context: prev.Context}, nil
			}
			// If there's only one key, return its value
			if len(m) == 1 {
				for _, v := range m {
					return CommandResult{Output: v, Context: prev.Context}, nil
				}
			}
		}
		return CommandResult{Output: result, Context: prev.Context}, nil
	}

	// One field - pivot by rows
	if len(item.Fields) == 1 {
		rowField := item.Fields[0].Value.(string)
		rows := getUniqueValues(slice, rowField)

		resultArray := make([]interface{}, 0, len(rows))
		for _, row := range rows {
			filteredData := filterByField(slice, rowField, row)
			summarized := summarizeGroup(nil, []SummarizeAssignment{item.Metric}, filteredData)

			rowObj := map[string]interface{}{
				rowField: row,
				"value":  extractMetricValue(summarized, string(item.Metric.Operator)),
			}
			resultArray = append(resultArray, rowObj)
		}
		return CommandResult{Output: resultArray, Context: prev.Context}, nil
	}

	// Two fields - pivot by rows and columns
	if len(item.Fields) >= 2 {
		rowField := item.Fields[0].Value.(string)
		colField := item.Fields[1].Value.(string)

		rows := getUniqueValues(slice, rowField)
		cols := getUniqueValues(slice, colField)

		resultArray := make([]interface{}, 0, len(rows))
		for _, row := range rows {
			rowObj := map[string]interface{}{rowField: row}

			for _, col := range cols {
				filteredData := filterByTwoFields(slice, rowField, row, colField, col)

				var value interface{}
				if len(filteredData) == 0 {
					// Default values for empty groups
					op := string(item.Metric.Operator)
					if op == "count" || op == "dcount" || op == "sum" {
						value = 0
					} else {
						value = nil
					}
				} else {
					summarized := summarizeGroup(nil, []SummarizeAssignment{item.Metric}, filteredData)
					value = extractMetricValue(summarized, string(item.Metric.Operator))
				}

				colStr := fmt.Sprintf("%v", col)
				rowObj[colStr] = value
			}
			resultArray = append(resultArray, rowObj)
		}
		return CommandResult{Output: resultArray, Context: prev.Context}, nil
	}

	return CommandResult{Output: slice, Context: prev.Context}, nil
}

// resolveArg resolves an argument value
func resolveArg(item interface{}, arg TypedValue) interface{} {
	switch arg.Type {
	case "ref":
		return getValue(item, arg.Value.(string))
	case "string":
		return arg.Value
	case "number":
		return arg.Value
	default:
		return nil
	}
}

// toMap converts an interface{} to map[string]interface{}
func toMap(v interface{}) (map[string]interface{}, bool) {
	if m, ok := v.(map[string]interface{}); ok {
		return m, true
	}

	// Try reflection
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Map {
		result := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			result[key.String()] = val.MapIndex(key).Interface()
		}
		return result, true
	}

	return nil, false
}

// summarizeGroup applies summarize metrics to a group of data
func summarizeGroup(baseObj map[string]interface{}, metrics []SummarizeAssignment, data []interface{}) interface{} {
	result := make(map[string]interface{})

	// Copy base object fields
	if baseObj != nil {
		for k, v := range baseObj {
			result[k] = v
		}
	}

	for _, metric := range metrics {
		statName := metric.Alias
		if statName == "" {
			if len(metric.Args) > 0 {
				argValue := ""
				if metric.Args[0].Type == "ref" {
					argValue = metric.Args[0].Value.(string)
				} else if metric.Args[0].Type == "string" {
					argValue = metric.Args[0].Value.(string)
				}
				statName = fmt.Sprintf("%s (%s)", argValue, metric.Operator)
			} else {
				statName = string(metric.Operator)
			}
		}

		var value interface{}

		switch metric.Operator {
		case FnCount:
			value = len(data)
		case FnDCount:
			if len(metric.Args) > 0 {
				fieldName := metric.Args[0].Value.(string)
				uniqueVals := make(map[interface{}]bool)
				for _, item := range data {
					val := getValue(item, fieldName)
					uniqueVals[val] = true
				}
				value = len(uniqueVals)
			} else {
				value = len(data)
			}
		case FnSum:
			sum := 0.0
			if len(metric.Args) > 0 {
				fieldName := metric.Args[0].Value.(string)
				for _, item := range data {
					if num, ok := toNumber(getValue(item, fieldName)); ok {
						sum += num
					}
				}
			}
			value = sum
		case FnMean:
			if len(metric.Args) > 0 {
				fieldName := metric.Args[0].Value.(string)
				values := make([]float64, 0)
				for _, item := range data {
					if num, ok := toNumber(getValue(item, fieldName)); ok {
						values = append(values, num)
					}
				}
				if len(values) > 0 {
					sum := 0.0
					for _, v := range values {
						sum += v
					}
					value = sum / float64(len(values))
				}
			}
		case FnMin:
			if len(metric.Args) > 0 {
				fieldName := metric.Args[0].Value.(string)
				var minVal *float64
				for _, item := range data {
					if num, ok := toNumber(getValue(item, fieldName)); ok {
						if minVal == nil || num < *minVal {
							minVal = &num
						}
					}
				}
				if minVal != nil {
					value = *minVal
				}
			}
		case FnMax:
			if len(metric.Args) > 0 {
				fieldName := metric.Args[0].Value.(string)
				var maxVal *float64
				for _, item := range data {
					if num, ok := toNumber(getValue(item, fieldName)); ok {
						if maxVal == nil || num > *maxVal {
							maxVal = &num
						}
					}
				}
				if maxVal != nil {
					value = *maxVal
				}
			}
		case FnFirst:
			if len(data) > 0 {
				if len(metric.Args) > 0 {
					fieldName := metric.Args[0].Value.(string)
					value = getValue(data[0], fieldName)
				} else {
					value = data[0]
				}
			}
		case FnLast, FnLatest:
			if len(data) > 0 {
				if len(metric.Args) > 0 {
					fieldName := metric.Args[0].Value.(string)
					value = getValue(data[len(data)-1], fieldName)
				} else {
					value = data[len(data)-1]
				}
			}
		}

		result[statName] = value
	}

	return result
}

// groupByField groups data by a single field
func groupByField(data []interface{}, field string) map[interface{}][]interface{} {
	groups := make(map[interface{}][]interface{})
	for _, item := range data {
		key := getValue(item, field)
		groups[key] = append(groups[key], item)
	}
	return groups
}

// groupByFields groups data by multiple fields
func groupByFields(data []interface{}, fields []TypedValue) [][]interface{} {
	groupMap := make(map[string][]interface{})

	for _, item := range data {
		keyParts := make([]string, len(fields))
		for i, field := range fields {
			fieldName := field.Value.(string)
			val := getValue(item, fieldName)
			keyParts[i] = fmt.Sprintf("%v", val)
		}
		key := strings.Join(keyParts, "#___#")
		groupMap[key] = append(groupMap[key], item)
	}

	result := make([][]interface{}, 0, len(groupMap))
	for _, group := range groupMap {
		result = append(result, group)
	}
	return result
}

// getUniqueValues gets unique values for a field
func getUniqueValues(data []interface{}, field string) []interface{} {
	seen := make(map[interface{}]bool)
	result := make([]interface{}, 0)

	for _, item := range data {
		val := getValue(item, field)
		if val != nil && val != "" {
			if !seen[val] {
				seen[val] = true
				result = append(result, val)
			}
		}
	}

	return result
}

// filterByField filters data by a field value
func filterByField(data []interface{}, field string, value interface{}) []interface{} {
	result := make([]interface{}, 0)
	for _, item := range data {
		if getValue(item, field) == value {
			result = append(result, item)
		}
	}
	return result
}

// filterByTwoFields filters data by two field values
func filterByTwoFields(data []interface{}, field1 string, value1 interface{}, field2 string, value2 interface{}) []interface{} {
	result := make([]interface{}, 0)
	for _, item := range data {
		if getValue(item, field1) == value1 && getValue(item, field2) == value2 {
			result = append(result, item)
		}
	}
	return result
}

// extractMetricValue extracts a metric value from summarized result
func extractMetricValue(summarized interface{}, metricName string) interface{} {
	if m, ok := summarized.(map[string]interface{}); ok {
		// Try exact match first
		if val, exists := m[metricName]; exists {
			return val
		}
		// If only one key, return its value
		if len(m) == 1 {
			for _, v := range m {
				return v
			}
		}
	}
	return nil
}
