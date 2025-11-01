package uql

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// evalParseJSON evaluates a parse-json command
func evalParseJSON(prev CommandResult, cmd Command) (CommandResult, error) {
	output := prev.Output
	if output == nil {
		return prev, nil
	}

	// Convert to string if needed
	var jsonStr string
	switch v := output.(type) {
	case string:
		jsonStr = v
	case []byte:
		jsonStr = string(v)
	default:
		// Already parsed
		return prev, nil
	}

	var result interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return prev, err
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
}

// evalParseCSV evaluates a parse-csv command
func evalParseCSV(prev CommandResult, cmd Command) (CommandResult, error) {
	output := prev.Output
	if output == nil {
		return prev, nil
	}

	var csvStr string
	switch v := output.(type) {
	case string:
		csvStr = v
	case []byte:
		csvStr = string(v)
	default:
		return prev, errors.New("CSV data must be a string")
	}

	// Get options from command value
	options := getParseCSVOptions(cmd.Value)

	reader := csv.NewReader(strings.NewReader(csvStr))
	
	// Apply options
	if options.Delimiter != "" {
		reader.Comma = rune(options.Delimiter[0])
	}
	if options.Comment != "" {
		reader.Comment = rune(options.Comment[0])
	}
	reader.TrimLeadingSpace = options.Trim
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1 // Allow variable number of fields
	
	records, err := reader.ReadAll()
	if err != nil {
		return prev, err
	}

	if len(records) == 0 {
		return CommandResult{Output: []interface{}{}, Context: prev.Context}, nil
	}

	// Handle skipEmptyLines
	if options.SkipEmptyLines {
		filtered := make([][]string, 0)
		for _, record := range records {
			isEmpty := true
			for _, field := range record {
				if strings.TrimSpace(field) != "" {
					isEmpty = false
					break
				}
			}
			if !isEmpty {
				filtered = append(filtered, record)
			}
		}
		records = filtered
	}

	if len(records) == 0 {
		return CommandResult{Output: []interface{}{}, Context: prev.Context}, nil
	}

	// Determine headers
	var headers []string
	startRow := 0

	if options.Columns != nil && len(options.Columns) > 0 {
		// Custom headers provided
		headers = options.Columns
	} else if options.UseColumns {
		// First row as headers
		headers = records[0]
		startRow = 1
	} else {
		// Generate default column names
		if len(records) > 0 {
			headers = make([]string, len(records[0]))
			for i := range headers {
				headers[i] = fmt.Sprintf("col_%d", i)
			}
		}
	}

	result := make([]interface{}, 0, len(records)-startRow)

	for i := startRow; i < len(records); i++ {
		row := make(map[string]interface{})
		for j, header := range headers {
			value := ""
			if j < len(records[i]) {
				value = records[i][j]
				if options.Trim {
					value = strings.TrimSpace(value)
				}
			}
			row[header] = value
		}
		result = append(result, row)
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
}

// CSVOptions holds CSV parsing options
type CSVOptions struct {
	Delimiter       string
	Comment         string
	Columns         []string
	UseColumns      bool
	Trim            bool
	SkipEmptyLines  bool
	RelaxColumnCount bool
}

// getParseCSVOptions extracts CSV options from parse args
func getParseCSVOptions(value interface{}) CSVOptions {
	options := CSVOptions{
		Delimiter:  ",",
		UseColumns: true, // Default to using first row as headers
	}

	if value == nil {
		return options
	}

	// value should be [][]ParseArg
	argsSlice, ok := value.([][]ParseArg)
	if !ok || len(argsSlice) == 0 {
		return options
	}

	args := argsSlice[0]
	for _, arg := range args {
		switch arg.Identifier {
		case "delimiter":
			// Handle escape sequences
			delimiter := arg.Value
			delimiter = strings.ReplaceAll(delimiter, "\\t", "\t")
			delimiter = strings.ReplaceAll(delimiter, "\\n", "\n")
			delimiter = strings.ReplaceAll(delimiter, "\\r", "\r")
			options.Delimiter = delimiter
		case "comment":
			options.Comment = arg.Value
		case "columns":
			if strings.ToLower(arg.Value) == "false" {
				options.UseColumns = false
			} else if strings.ToLower(arg.Value) == "true" {
				options.UseColumns = true
			} else {
				// Custom column names
				options.Columns = strings.Split(arg.Value, ",")
				options.UseColumns = false
			}
		case "trim":
			options.Trim = strings.ToLower(arg.Value) == "true"
		case "skipEmptyLines":
			options.SkipEmptyLines = strings.ToLower(arg.Value) == "true"
		case "relaxColumnCount":
			options.RelaxColumnCount = strings.ToLower(arg.Value) == "true"
		}
	}

	return options
}

// evalParseXML evaluates a parse-xml command
func evalParseXML(prev CommandResult, cmd Command) (CommandResult, error) {
	output := prev.Output
	if output == nil {
		return prev, nil
	}

	var xmlStr string
	switch v := output.(type) {
	case string:
		xmlStr = v
	case []byte:
		xmlStr = string(v)
	default:
		return prev, errors.New("XML data must be a string")
	}

	var result interface{}
	err := xml.Unmarshal([]byte(xmlStr), &result)
	if err != nil {
		return prev, err
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
}

// evalParseYAML evaluates a parse-yaml command
func evalParseYAML(prev CommandResult, cmd Command) (CommandResult, error) {
	output := prev.Output
	if output == nil {
		return prev, nil
	}

	var yamlStr string
	switch v := output.(type) {
	case string:
		yamlStr = v
	case []byte:
		yamlStr = string(v)
	default:
		return prev, errors.New("YAML data must be a string")
	}

	var result interface{}
	err := yaml.Unmarshal([]byte(yamlStr), &result)
	if err != nil {
		return prev, err
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
}
