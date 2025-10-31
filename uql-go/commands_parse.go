package uql

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
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

	reader := csv.NewReader(strings.NewReader(csvStr))
	records, err := reader.ReadAll()
	if err != nil {
		return prev, err
	}

	if len(records) == 0 {
		return CommandResult{Output: []interface{}{}, Context: prev.Context}, nil
	}

	// First row as headers
	headers := records[0]
	result := make([]interface{}, 0, len(records)-1)

	for i := 1; i < len(records); i++ {
		row := make(map[string]interface{})
		for j, header := range headers {
			if j < len(records[i]) {
				row[header] = records[i][j]
			}
		}
		result = append(result, row)
	}

	return CommandResult{Output: result, Context: prev.Context}, nil
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
