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

func TestUQLSummarize(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"patron": "a", "age": 48, "country": "foo"},
		map[string]interface{}{"patron": "b", "age": 34, "country": "foo"},
		map[string]interface{}{"patron": "c", "age": 12, "country": "bar"},
		map[string]interface{}{"patron": "d", "age": 40, "country": "bar"},
		map[string]interface{}{"patron": "e", "age": 36, "country": "baz"},
	}

	// Test summarize with single group by
	result, err := UQL(`summarize "user"=first("patron") by "country"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 3 {
		t.Errorf("expected 3 groups, got %d", len(slice))
	}

	// Test summarize with sum
	result2, err := UQL(`summarize "age"=sum("age") by "country"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice2, ok := result2.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result2)
	}

	if len(slice2) != 3 {
		t.Errorf("expected 3 groups, got %d", len(slice2))
	}
}

func TestUQLSummarizeMultiple(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"age": 1, "name": "foo1", "city": "chennai", "country": "india"},
		map[string]interface{}{"age": 2, "name": "foo1", "city": "chennai", "country": "india"},
		map[string]interface{}{"age": 3, "name": "foo1", "city": "mumbai", "country": "india"},
		map[string]interface{}{"age": 4, "name": "foo1", "city": "london", "country": "england"},
	}

	result, err := UQL(`summarize "age"=sum("age") by "country", "city"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 3 {
		t.Errorf("expected 3 groups, got %d", len(slice))
	}
}

func TestUQLPivot(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"fruit": "apple", "size": "sm", "qty": 1},
		map[string]interface{}{"fruit": "apple", "size": "md", "qty": 2},
		map[string]interface{}{"fruit": "apple", "size": "lg", "qty": 3},
		map[string]interface{}{"fruit": "banana", "size": "sm", "qty": 1},
		map[string]interface{}{"fruit": "banana", "size": "lg", "qty": 6},
		map[string]interface{}{"fruit": "banana", "size": "xl", "qty": 5},
	}

	// Test pivot with no fields
	result, err := UQL(`pivot count()`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 6 {
		t.Errorf("expected 6, got %v", result)
	}

	// Test pivot with sum
	result2, err := UQL(`pivot sum("qty")`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result2 != 18.0 {
		t.Errorf("expected 18, got %v", result2)
	}
}

func TestUQLPivotWithRow(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"fruit": "apple", "size": "sm", "qty": 1},
		map[string]interface{}{"fruit": "apple", "size": "md", "qty": 2},
		map[string]interface{}{"fruit": "apple", "size": "lg", "qty": 3},
		map[string]interface{}{"fruit": "banana", "size": "sm", "qty": 1},
		map[string]interface{}{"fruit": "banana", "size": "lg", "qty": 6},
		map[string]interface{}{"fruit": "banana", "size": "xl", "qty": 5},
	}

	result, err := UQL(`pivot sum("qty"), "fruit"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 2 {
		t.Errorf("expected 2 rows, got %d", len(slice))
	}
}

func TestUQLPivotWithRowAndCol(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"fruit": "apple", "size": "sm", "qty": 1},
		map[string]interface{}{"fruit": "apple", "size": "md", "qty": 2},
		map[string]interface{}{"fruit": "apple", "size": "lg", "qty": 3},
		map[string]interface{}{"fruit": "banana", "size": "sm", "qty": 1},
		map[string]interface{}{"fruit": "banana", "size": "lg", "qty": 6},
		map[string]interface{}{"fruit": "banana", "size": "xl", "qty": 5},
	}

	result, err := UQL(`pivot sum("qty"), "fruit", "size"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 2 {
		t.Errorf("expected 2 rows, got %d", len(slice))
	}

	// Check first row has correct structure
	row1, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", slice[0])
	}

	if _, hasFruit := row1["fruit"]; !hasFruit {
		t.Error("expected 'fruit' field in result")
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

func TestUQLRange(t *testing.T) {
	// Test numeric range
	result, err := UQL(`range from 1 to 5`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 5 {
		t.Errorf("expected 5 items, got %d", len(slice))
	}

	if slice[0] != 1.0 {
		t.Errorf("expected first item to be 1, got %v", slice[0])
	}
}

func TestUQLRangeWithStep(t *testing.T) {
	result, err := UQL(`range from 1 to 10 step 2`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 5 {
		t.Errorf("expected 5 items (1,3,5,7,9), got %d", len(slice))
	}
}

func TestUQLMvExpand(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"group": "A", "users": []interface{}{"user a1", "user a2"}},
		map[string]interface{}{"group": "B", "users": []interface{}{"user b1"}},
	}

	result, err := UQL(`mv-expand "users"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 3 {
		t.Errorf("expected 3 items, got %d", len(slice))
	}

	// Check first expanded item
	first, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", slice[0])
	}

	if first["group"] != "A" {
		t.Errorf("expected group A, got %v", first["group"])
	}

	if first["users"] != "user a1" {
		t.Errorf("expected users to be 'user a1', got %v", first["users"])
	}
}

func TestUQLMvExpandWithAlias(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"group": "A", "users": []interface{}{"user a1", "user a2"}},
		map[string]interface{}{"group": "B", "users": []interface{}{"user b1"}},
	}

	result, err := UQL(`mv-expand "user"="users"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 3 {
		t.Errorf("expected 3 items, got %d", len(slice))
	}

	// Check first expanded item
	first, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", slice[0])
	}

	if first["group"] != "A" {
		t.Errorf("expected group A, got %v", first["group"])
	}

	if first["user"] != "user a1" {
		t.Errorf("expected user to be 'user a1', got %v", first["user"])
	}

	if _, hasUsers := first["users"]; hasUsers {
		t.Errorf("users field should not be present after alias")
	}
}

func TestUQLMvExpandIgnoresNonArray(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"group": "A", "users": []interface{}{"user a1", "user a2"}},
		map[string]interface{}{"group": "B", "users": []interface{}{}},
		map[string]interface{}{"group": "C"},
		map[string]interface{}{"group": "D", "users": []interface{}{"user d1"}},
	}

	result, err := UQL(`mv-expand "user"="users"`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	// Should only include items from groups A and D (3 total: a1, a2, d1)
	if len(slice) != 3 {
		t.Errorf("expected 3 items, got %d", len(slice))
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
