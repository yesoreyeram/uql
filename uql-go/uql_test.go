package uql

import (
	"encoding/json"
	"testing"
)

func TestUQLHello(t *testing.T) {
	result, err := UQL("hello", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello" {
		t.Errorf("expected 'hello', got %v", result)
	}
}

func TestUQLPing(t *testing.T) {
	result, err := UQL("ping", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "pong" {
		t.Errorf("expected 'pong', got %v", result)
	}
}

func TestUQLEcho(t *testing.T) {
	result, err := UQL(`echo "test message"`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "test message" {
		t.Errorf("expected 'test message', got %v", result)
	}
}

func TestUQLCount(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"name": "foo", "age": 2},
		map[string]interface{}{"name": "bar", "age": 3},
	}

	result, err := UQL("count", &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 2 {
		t.Errorf("expected 2, got %v", result)
	}
}

func TestUQLLimit(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"name": "foo", "age": 2},
		map[string]interface{}{"name": "bar", "age": 3},
		map[string]interface{}{"name": "baz", "age": 4},
	}

	result, err := UQL("limit 2", &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}
	if len(slice) != 2 {
		t.Errorf("expected length 2, got %d", len(slice))
	}
}

func TestUQLProject(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"name": "foo", "age": 2, "location": "uk"},
		map[string]interface{}{"name": "bar", "age": 3, "location": "usa"},
	}

	result, err := UQL(`project "name", "age"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 2 {
		t.Fatalf("expected 2 items, got %d", len(slice))
	}

	first, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", slice[0])
	}

	if first["name"] != "foo" || first["age"] != 2 {
		t.Errorf("unexpected values: %v", first)
	}

	if _, hasLocation := first["location"]; hasLocation {
		t.Errorf("location should not be present")
	}
}

func TestUQLProjectAway(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"name": "foo", "age": 2, "location": "uk"},
		map[string]interface{}{"name": "bar", "age": 3, "location": "usa"},
	}

	result, err := UQL(`project-away "location"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	first, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", slice[0])
	}

	if _, hasLocation := first["location"]; hasLocation {
		t.Errorf("location should not be present")
	}

	if first["name"] != "foo" || first["age"] != 2 {
		t.Errorf("unexpected values: %v", first)
	}
}

func TestUQLOrderBy(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"name": "bar", "age": 3},
		map[string]interface{}{"name": "foo", "age": 2},
	}

	result, err := UQL(`order by "name" asc`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	first, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", slice[0])
	}

	if first["name"] != "bar" {
		t.Errorf("expected first name to be 'bar', got %v", first["name"])
	}
}

func TestUQLParseJSON(t *testing.T) {
	jsonData := `{"name": "foo", "age": 25}`

	result, err := UQL("parse-json", &Options{Data: jsonData})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}

	if m["name"] != "foo" {
		t.Errorf("expected name 'foo', got %v", m["name"])
	}
}

func TestUQLScope(t *testing.T) {
	data := map[string]interface{}{
		"result": map[string]interface{}{
			"value": 42,
		},
	}

	result, err := UQL(`scope "result"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}

	if m["value"] != 42 {
		t.Errorf("expected value 42, got %v", m["value"])
	}
}

func TestUQLDistinct(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"name": "foo", "age": 2},
		map[string]interface{}{"name": "bar", "age": 3},
		map[string]interface{}{"name": "foo", "age": 2},
	}

	result, err := UQL("distinct", &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 2 {
		t.Errorf("expected 2 distinct items, got %d", len(slice))
	}
}

func TestUQLPipeline(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"name": "foo", "age": 2, "location": "uk"},
		map[string]interface{}{"name": "bar", "age": 3, "location": "usa"},
	}

	result, err := UQL(`project "name", "location" | order by "name" asc`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	first, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", slice[0])
	}

	if first["name"] != "bar" {
		t.Errorf("expected first name to be 'bar', got %v", first["name"])
	}

	if _, hasAge := first["age"]; hasAge {
		t.Errorf("age should not be present after project")
	}
}

func TestUQLComplexQuery(t *testing.T) {
	users := []interface{}{
		map[string]interface{}{"name": "foo", "age": 2, "location": "uk"},
		map[string]interface{}{"name": "bar", "age": 3, "location": "usa"},
	}

	query := `parse-json | order by "name" asc | project "name", "location"`

	jsonData, _ := json.Marshal(users)
	result, err := UQL(query, &Options{Data: string(jsonData)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 2 {
		t.Fatalf("expected 2 items, got %d", len(slice))
	}

	first, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", slice[0])
	}

	if first["name"] != "bar" {
		t.Errorf("expected first name to be 'bar', got %v", first["name"])
	}

	if first["location"] != "usa" {
		t.Errorf("expected first location to be 'usa', got %v", first["location"])
	}
}
