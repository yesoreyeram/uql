package uql

import (
	"errors"
	"reflect"
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
	// TODO: Implement summarize
	return prev, nil
}

// evalPivot evaluates a pivot command
func evalPivot(prev CommandResult, cmd Command) (CommandResult, error) {
	// TODO: Implement pivot
	return prev, nil
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
