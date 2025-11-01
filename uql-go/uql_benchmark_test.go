package uql

import (
	"testing"
)

// Benchmark basic commands
func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = UQL("hello", nil)
	}
}

func BenchmarkPing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = UQL("ping", nil)
	}
}

func BenchmarkEcho(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`echo "test message"`, nil)
	}
}

func BenchmarkCount(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"id": 1},
		map[string]interface{}{"id": 2},
		map[string]interface{}{"id": 3},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL("count", opts)
	}
}

func BenchmarkLimit(b *testing.B) {
	data := make([]interface{}, 100)
	for i := 0; i < 100; i++ {
		data[i] = map[string]interface{}{"id": i}
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL("limit 10", opts)
	}
}

// Benchmark data transformation commands
func BenchmarkProject(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"name": "John", "age": 30, "city": "NYC"},
		map[string]interface{}{"name": "Jane", "age": 25, "city": "LA"},
		map[string]interface{}{"name": "Bob", "age": 35, "city": "Chicago"},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`project "name", "age"`, opts)
	}
}

func BenchmarkProjectAway(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"name": "John", "age": 30, "city": "NYC"},
		map[string]interface{}{"name": "Jane", "age": 25, "city": "LA"},
		map[string]interface{}{"name": "Bob", "age": 35, "city": "Chicago"},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`project-away "city"`, opts)
	}
}

func BenchmarkExtend(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"name": "John", "age": 30},
		map[string]interface{}{"name": "Jane", "age": 25},
		map[string]interface{}{"name": "Bob", "age": 35},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`extend "info"=strcat("name", " is ", "age")`, opts)
	}
}

func BenchmarkDistinct(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"category": "A"},
		map[string]interface{}{"category": "B"},
		map[string]interface{}{"category": "A"},
		map[string]interface{}{"category": "C"},
		map[string]interface{}{"category": "B"},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`distinct "category"`, opts)
	}
}

func BenchmarkOrderBy(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"name": "John", "age": 30},
		map[string]interface{}{"name": "Alice", "age": 25},
		map[string]interface{}{"name": "Bob", "age": 35},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`order by "name" asc`, opts)
	}
}

// Benchmark filtering
func BenchmarkWhereSimple(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"age": 30},
		map[string]interface{}{"age": 25},
		map[string]interface{}{"age": 35},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`where "age" > 28`, opts)
	}
}

func BenchmarkWhereContains(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"name": "John Smith"},
		map[string]interface{}{"name": "Jane Doe"},
		map[string]interface{}{"name": "Bob Johnson"},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`where "name" contains "john"`, opts)
	}
}

func BenchmarkWhereIn(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"status": "active"},
		map[string]interface{}{"status": "pending"},
		map[string]interface{}{"status": "inactive"},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`where "status" in ('active', 'pending')`, opts)
	}
}

// Benchmark aggregation
func BenchmarkSummarize(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"country": "USA", "age": 30},
		map[string]interface{}{"country": "USA", "age": 25},
		map[string]interface{}{"country": "UK", "age": 35},
		map[string]interface{}{"country": "UK", "age": 40},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`summarize "total"=sum("age") by "country"`, opts)
	}
}

func BenchmarkPivot(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"fruit": "apple", "size": "sm", "qty": 1},
		map[string]interface{}{"fruit": "apple", "size": "lg", "qty": 3},
		map[string]interface{}{"fruit": "banana", "size": "sm", "qty": 2},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`pivot sum("qty"), "fruit", "size"`, opts)
	}
}

// Benchmark parsing
func BenchmarkParseJSON(b *testing.B) {
	jsonData := `{"name":"John","age":30,"city":"NYC"}`
	opts := &Options{Data: jsonData}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL("parse-json", opts)
	}
}

func BenchmarkParseCSVBasic(b *testing.B) {
	csvData := "name,age,city\nJohn,30,NYC\nJane,25,LA\nBob,35,Chicago"
	opts := &Options{Data: csvData}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL("parse-csv", opts)
	}
}

func BenchmarkParseCSVWithOptions(b *testing.B) {
	csvData := "name;age;city\nJohn;30;NYC\nJane;25;LA\nBob;35;Chicago"
	opts := &Options{Data: csvData}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`parse-csv --delimiter ";"`, opts)
	}
}

func BenchmarkParseYAML(b *testing.B) {
	yamlData := "name: John\nage: 30\ncity: NYC"
	opts := &Options{Data: yamlData}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL("parse-yaml", opts)
	}
}

func BenchmarkParseXML(b *testing.B) {
	xmlData := `<root><name>John</name><age>30</age></root>`
	opts := &Options{Data: xmlData}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL("parse-xml", opts)
	}
}

// Benchmark range
func BenchmarkRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = UQL("range from 1 to 100", nil)
	}
}

func BenchmarkRangeWithStep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = UQL("range from 1 to 100 step 5", nil)
	}
}

// Benchmark mv-expand
func BenchmarkMvExpand(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"group": "A", "users": []interface{}{"user1", "user2", "user3"}},
		map[string]interface{}{"group": "B", "users": []interface{}{"user4", "user5"}},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`mv-expand "users"`, opts)
	}
}

// Benchmark pipeline operations
func BenchmarkPipelineSimple(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"name": "John", "age": 30},
		map[string]interface{}{"name": "Jane", "age": 25},
		map[string]interface{}{"name": "Bob", "age": 35},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`project "name", "age" | limit 2`, opts)
	}
}

func BenchmarkPipelineComplex(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"name": "John", "age": 30, "city": "NYC"},
		map[string]interface{}{"name": "Jane", "age": 25, "city": "LA"},
		map[string]interface{}{"name": "Bob", "age": 35, "city": "NYC"},
		map[string]interface{}{"name": "Alice", "age": 28, "city": "Chicago"},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`where "age" > 25 | order by "name" asc | project "name", "city" | limit 2`, opts)
	}
}

// Benchmark string functions
func BenchmarkFunctionToupper(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"name": "john"},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`extend "upper"=toupper("name")`, opts)
	}
}

func BenchmarkFunctionStrcat(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"first": "John", "last": "Doe"},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`extend "full"=strcat("first", " ", "last")`, opts)
	}
}

func BenchmarkFunctionSplit(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"text": "a,b,c,d,e"},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`extend "parts"=split("text", ',')`, opts)
	}
}

// Benchmark math functions
func BenchmarkFunctionSum(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"values": []interface{}{1, 2, 3, 4, 5}},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`extend "total"=sum("values")`, opts)
	}
}

func BenchmarkFunctionMean(b *testing.B) {
	data := []interface{}{
		map[string]interface{}{"values": []interface{}{10, 20, 30, 40, 50}},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`extend "avg"=mean("values")`, opts)
	}
}

// Benchmark JSONata
func BenchmarkJSONataSimple(b *testing.B) {
	data := map[string]interface{}{
		"example": []interface{}{
			map[string]interface{}{"value": 4},
			map[string]interface{}{"value": 7},
			map[string]interface{}{"value": 13},
		},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`jsonata "$sum(example.value)"`, opts)
	}
}

func BenchmarkJSONataFilter(b *testing.B) {
	data := map[string]interface{}{
		"Countries": []interface{}{
			map[string]interface{}{"Country": "India"},
			map[string]interface{}{"Country": "America"},
			map[string]interface{}{"Country": "United Kingdom"},
		},
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`scope "Countries" | jsonata "$[Country='India']"`, opts)
	}
}

// Benchmark large datasets
func BenchmarkLargeDatasetProject(b *testing.B) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = map[string]interface{}{
			"id":   i,
			"name": "User" + string(rune(i)),
			"age":  20 + (i % 50),
			"city": "City" + string(rune(i%10)),
		}
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`project "name", "age"`, opts)
	}
}

func BenchmarkLargeDatasetWhere(b *testing.B) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = map[string]interface{}{
			"id":  i,
			"age": 20 + (i % 50),
		}
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`where "age" > 40`, opts)
	}
}

func BenchmarkLargeDatasetSummarize(b *testing.B) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = map[string]interface{}{
			"country": "Country" + string(rune(i%10)),
			"value":   i,
		}
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`summarize "total"=sum("value") by "country"`, opts)
	}
}

func BenchmarkLargeDatasetOrderBy(b *testing.B) {
	data := make([]interface{}, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = map[string]interface{}{
			"id":   1000 - i,
			"name": "User" + string(rune(1000-i)),
		}
	}
	opts := &Options{Data: data}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UQL(`order by "id" asc`, opts)
	}
}

// Benchmark parser
func BenchmarkParserSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(`project "name", "age"`)
		tokens, _ := lexer.Tokenize()
		parser := NewParser(tokens)
		_, _ = parser.ParseCommands()
	}
}

func BenchmarkParserComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lexer := NewLexer(`where "age" > 25 | order by "name" asc | project "name", "city" | limit 10`)
		tokens, _ := lexer.Tokenize()
		parser := NewParser(tokens)
		_, _ = parser.ParseCommands()
	}
}
