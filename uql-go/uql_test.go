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

func TestUQLWhere(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"a": 1},
		map[string]interface{}{"a": 10},
		map[string]interface{}{"a": 20},
		map[string]interface{}{"a": 30},
	}

	result, err := UQL(`where "a" == 10`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 1 {
		t.Errorf("expected 1 item, got %d", len(slice))
	}

	first, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", slice[0])
	}

	if first["a"] != 10 {
		t.Errorf("expected a=10, got %v", first["a"])
	}
}

func TestUQLWhereGreaterThan(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"a": 1},
		map[string]interface{}{"a": 10},
		map[string]interface{}{"a": 20},
		map[string]interface{}{"a": 30},
	}

	result, err := UQL(`where "a" > 10`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 2 {
		t.Errorf("expected 2 items, got %d", len(slice))
	}
}

func TestUQLWhereStringContains(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"a": "FabriKam"},
		map[string]interface{}{"a": "banana"},
	}

	result, err := UQL(`where "a" contains 'BRik'`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 1 {
		t.Errorf("expected 1 item, got %d", len(slice))
	}

	first, ok := slice[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", slice[0])
	}

	if first["a"] != "FabriKam" {
		t.Errorf("expected 'FabriKam', got %v", first["a"])
	}
}

func TestUQLWhereIn(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"a": 1},
		map[string]interface{}{"a": 10},
		map[string]interface{}{"a": 20},
		map[string]interface{}{"a": 30},
	}

	result, err := UQL(`where "a" in (10,20)`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 2 {
		t.Errorf("expected 2 items, got %d", len(slice))
	}
}

func TestUQLWhereBetween(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"a": 1},
		map[string]interface{}{"a": 10},
		map[string]interface{}{"a": 20},
		map[string]interface{}{"a": 30},
	}

	result, err := UQL(`where "a" between (10,20)`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	slice, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected slice, got %T", result)
	}

	if len(slice) != 2 {
		t.Errorf("expected 2 items (10 and 20), got %d", len(slice))
	}
}

func TestUQLSplit(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"text": "a,b,c"},
	}

	result, err := UQL(`extend "parts"=split("text", ',')`, &Options{Data: data})
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

	parts, ok := first["parts"].([]interface{})
	if !ok {
		t.Fatalf("expected parts to be array, got %T", first["parts"])
	}

	if len(parts) != 3 {
		t.Errorf("expected 3 parts, got %d", len(parts))
	}

	if parts[0] != "a" || parts[1] != "b" || parts[2] != "c" {
		t.Errorf("unexpected split result: %v", parts)
	}
}

func TestUQLExtract(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"text": "The price is 123 dollars"},
	}

	result, err := UQL(`extend "price"=extract('([0-9]+)', 1, "text")`, &Options{Data: data})
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

	if first["price"] != "123" {
		t.Errorf("expected '123', got %v", first["price"])
	}
}

func TestUQLExtractWithNumberConversion(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"text": "The price is 123 dollars"},
	}

	result, err := UQL(`extend "price"=extract('([0-9]+)', 1, "text", 'number')`, &Options{Data: data})
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

	if price, ok := first["price"].(float64); !ok || price != 123.0 {
		t.Errorf("expected 123.0 as float64, got %v (type %T)", first["price"], first["price"])
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

// Test new math functions
func TestUQLMathFunctions(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"value": 10.0},
	}

	// Test diff
	result, err := UQL(`extend "diff"=diff(10, 3)`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice := result.([]interface{})
	m := slice[0].(map[string]interface{})
	if m["diff"] != 7.0 {
		t.Errorf("expected diff to be 7, got %v", m["diff"])
	}

	// Test mul
	result, err = UQL(`extend "product"=mul(3, 4, 2)`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice = result.([]interface{})
	m = slice[0].(map[string]interface{})
	if m["product"] != 24.0 {
		t.Errorf("expected product to be 24, got %v", m["product"])
	}

	// Test div
	result, err = UQL(`extend "quotient"=div(20, 4)`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice = result.([]interface{})
	m = slice[0].(map[string]interface{})
	if m["quotient"] != 5.0 {
		t.Errorf("expected quotient to be 5, got %v", m["quotient"])
	}

	// Test percentage
	result, err = UQL(`extend "pct"=percentage(25, 200)`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice = result.([]interface{})
	m = slice[0].(map[string]interface{})
	if m["pct"] != 12.5 {
		t.Errorf("expected percentage to be 12.5, got %v", m["pct"])
	}
}

// Test hyperbolic trig functions
func TestUQLHyperbolicTrigFunctions(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"value": 1.0},
	}

	// Test sinh
	result, err := UQL(`extend "result"=sinh(0)`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice := result.([]interface{})
	m := slice[0].(map[string]interface{})
	if m["result"] != 0.0 {
		t.Errorf("expected sinh(0) to be 0, got %v", m["result"])
	}

	// Test cosh
	result, err = UQL(`extend "result"=cosh(0)`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice = result.([]interface{})
	m = slice[0].(map[string]interface{})
	if m["result"] != 1.0 {
		t.Errorf("expected cosh(0) to be 1, got %v", m["result"])
	}

	// Test tanh
	result, err = UQL(`extend "result"=tanh(0)`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice = result.([]interface{})
	m = slice[0].(map[string]interface{})
	if m["result"] != 0.0 {
		t.Errorf("expected tanh(0) to be 0, got %v", m["result"])
	}
}

// Test extended datetime functions
func TestUQLExtendedDatetimeFunctions(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"timestamp": 1609459200000000}, // 2021-01-01 00:00:00 in microseconds
	}

	// Test unixtime_microseconds_todatetime
	result, err := UQL(`extend "dt"=unixtime_microseconds_todatetime("timestamp")`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice := result.([]interface{})
	m := slice[0].(map[string]interface{})
	if m["dt"] == nil {
		t.Errorf("expected datetime value, got nil")
	}

	// Test todatetime
	data2 := []interface{}{
		map[string]interface{}{"datestr": "2021-01-01"},
	}
	result, err = UQL(`extend "dt"=todatetime("datestr")`, &Options{Data: data2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice = result.([]interface{})
	m = slice[0].(map[string]interface{})
	if m["dt"] == nil {
		t.Errorf("expected datetime value, got nil")
	}
}

// Test array functions
func TestUQLArrayFunctions(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"items": []interface{}{1, 2, 2, 3, 3, 3}},
	}

	// Test distinct
	result, err := UQL(`extend "unique"=distinct("items")`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice := result.([]interface{})
	m := slice[0].(map[string]interface{})
	unique := m["unique"].([]interface{})
	if len(unique) != 3 {
		t.Errorf("expected 3 unique items, got %d", len(unique))
	}

	// Test pack
	data2 := []interface{}{
		map[string]interface{}{"a": 1, "b": 2},
	}
	result, err = UQL(`extend "packed"=pack("a", "b", 3)`, &Options{Data: data2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice = result.([]interface{})
	m = slice[0].(map[string]interface{})
	packed := m["packed"].([]interface{})
	if len(packed) != 3 {
		t.Errorf("expected 3 packed items, got %d", len(packed))
	}

	// Test kv
	result, err = UQL(`extend "pair"=kv('name', 'value123')`, &Options{Data: data2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice = result.([]interface{})
	m = slice[0].(map[string]interface{})
	pair := m["pair"].(map[string]interface{})
	if pair["key"] != "name" || pair["value"] != "value123" {
		t.Errorf("expected kv pair with key='name' and value='value123', got %v", pair)
	}
}

// Test URL parsing functions
func TestUQLURLFunctions(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"url": "https://example.com:8080/path/to/resource?foo=bar&baz=qux#section"},
	}

	// Test parse_url
	result, err := UQL(`extend "parsed"=parse_url("url")`, &Options{Data: data})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice := result.([]interface{})
	m := slice[0].(map[string]interface{})
	parsed := m["parsed"].(map[string]interface{})

	if parsed["scheme"] != "https" {
		t.Errorf("expected scheme to be 'https', got %v", parsed["scheme"])
	}
	if parsed["hostname"] != "example.com" {
		t.Errorf("expected hostname to be 'example.com', got %v", parsed["hostname"])
	}
	if parsed["port"] != "8080" {
		t.Errorf("expected port to be '8080', got %v", parsed["port"])
	}

	// Test parse_urlquery
	data2 := []interface{}{
		map[string]interface{}{"query": "foo=bar&baz=qux"},
	}
	result, err = UQL(`extend "params"=parse_urlquery("query")`, &Options{Data: data2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	slice = result.([]interface{})
	m = slice[0].(map[string]interface{})
	params := m["params"].(map[string]interface{})

	if params["foo"] != "bar" {
		t.Errorf("expected foo to be 'bar', got %v", params["foo"])
	}
	if params["baz"] != "qux" {
		t.Errorf("expected baz to be 'qux', got %v", params["baz"])
	}
}
