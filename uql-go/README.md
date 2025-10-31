# UQL - Unified Query Language (Go Port)

Go port of the UQL (Unified Query Language) package. UQL is a query language to query JSON-like data structures, inspired by Azure Kusto Query Language (KQL).

## Installation

```bash
go get github.com/yesoreyeram/uql/uql-go
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/yesoreyeram/uql/uql-go"
)

func main() {
	users := []interface{}{
		map[string]interface{}{"name": "foo", "age": 2, "location": "uk"},
		map[string]interface{}{"name": "bar", "age": 3, "location": "usa"},
	}

	query := `order by "name" asc | project "name", "location"`

	result, err := uql.UQL(query, &uql.Options{Data: users})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", result)
	// Output: [map[location:usa name:bar] map[location:uk name:foo]]
}
```

## Supported Commands

### Basic Commands
- `hello` - Returns "hello"
- `ping` - Returns "pong"
- `echo <string>` - Returns the given string
- `count` - Returns the count of items
- `limit <number>` - Limits the number of results

### Data Transformation
- `project <fields>` - Select specific fields
- `project-away <fields>` - Remove specific fields
- `project-reorder <fields>` - Reorder fields
- `extend <assignments>` - Add new fields
- `order by <field> [asc|desc]` - Sort results
- `distinct [field]` - Get distinct values
- `scope <field>` - Navigate to a nested field
- `summarize <metrics> [by <fields>]` - Aggregate data by groups
- `pivot <metric>, [<row_field>], [<col_field>]` - Pivot data into cross-tabular format

### Data Parsing
- `parse-json` - Parse JSON string
- `parse-csv` - Parse CSV string
- `parse-xml` - Parse XML string
- `parse-yaml` - Parse YAML string

### Functions

#### String Functions
- `toupper(str)` - Convert to uppercase
- `tolower(str)` - Convert to lowercase
- `trim(str)` - Trim whitespace
- `trim_start(str)` - Trim leading whitespace
- `trim_end(str)` - Trim trailing whitespace
- `reverse(str)` - Reverse string
- `strlen(str)` - String length
- `strcat(str1, str2, ...)` - Concatenate strings
- `replace_string(str, old, new)` - Replace substrings
- `substring(str, start, length)` - Extract substring

#### Math Functions
- `sum(...)` - Sum of values
- `min(...)` - Minimum value
- `max(...)` - Maximum value
- `mean(...)` - Average value
- `abs(num)` - Absolute value
- `floor(num)` - Floor value
- `ceil(num)` - Ceiling value
- `round(num)` - Round value
- `pow(base, exp)` - Power
- `log(num)`, `log2(num)`, `log10(num)` - Logarithms

#### Trigonometric Functions
- `sin(x)`, `cos(x)`, `tan(x)` - Trigonometric functions
- `asin(x)`, `acos(x)`, `atan(x)` - Inverse trigonometric functions

#### Conversion Functions
- `tostring(val)` - Convert to string
- `toint(val)`, `tolong(val)` - Convert to integer
- `todouble(val)`, `tofloat(val)`, `tonumber(val)` - Convert to number
- `tobool(val)` - Convert to boolean

#### Encoding Functions
- `atob(str)` - Base64 decode
- `btoa(str)` - Base64 encode

#### Date/Time Functions
- `unixtime_milliseconds_todatetime(ms)` - Convert Unix timestamp (ms) to datetime
- `unixtime_seconds_todatetime(sec)` - Convert Unix timestamp (sec) to datetime
- `tounixtime(datetime)` - Convert datetime to Unix timestamp
- `startofday(datetime)` - Get start of day
- `startofhour(datetime)` - Get start of hour
- `startofminute(datetime)` - Get start of minute

## Examples

### Example 1: Simple Projection

```go
data := []interface{}{
	map[string]interface{}{"name": "Alice", "age": 30, "city": "NYC"},
	map[string]interface{}{"name": "Bob", "age": 25, "city": "LA"},
}

result, _ := uql.UQL(`project "name", "city"`, &uql.Options{Data: data})
// Result: [map[city:NYC name:Alice] map[city:LA name:Bob]]
```

### Example 2: Ordering and Limiting

```go
data := []interface{}{
	map[string]interface{}{"name": "Charlie", "score": 85},
	map[string]interface{}{"name": "Alice", "score": 92},
	map[string]interface{}{"name": "Bob", "score": 78},
}

result, _ := uql.UQL(`order by "score" desc | limit 2`, &uql.Options{Data: data})
// Result: Top 2 scores
```

### Example 3: Parse JSON

```go
jsonData := `{"users": [{"name": "foo"}, {"name": "bar"}]}`

result, _ := uql.UQL(`parse-json | scope "users"`, &uql.Options{Data: jsonData})
// Result: Array of users
```

### Example 4: Summarize Data

```go
users := []interface{}{
	map[string]interface{}{"patron": "a", "age": 48, "country": "foo"},
	map[string]interface{}{"patron": "b", "age": 34, "country": "foo"},
	map[string]interface{}{"patron": "c", "age": 12, "country": "bar"},
}

result, _ := uql.UQL(`summarize "total_age"=sum("age") by "country"`, &uql.Options{Data: users})
// Result: Grouped aggregation by country
```

### Example 5: Pivot Data (Cross-tabulation)

```go
fruits := []interface{}{
	map[string]interface{}{"fruit": "apple", "size": "sm", "qty": 1},
	map[string]interface{}{"fruit": "apple", "size": "lg", "qty": 3},
	map[string]interface{}{"fruit": "banana", "size": "sm", "qty": 1},
}

result, _ := uql.UQL(`pivot sum("qty"), "fruit", "size"`, &uql.Options{Data: fruits})
// Result: Cross-tabulated data with fruits as rows and sizes as columns
```

### Example 6: Pipeline Processing

```go
jsonData := `[{"name": "foo", "age": 25}, {"name": "bar", "age": 30}]`

query := `parse-json | order by "name" asc | project "name"`
result, _ := uql.UQL(query, &uql.Options{Data: jsonData})
// Result: [map[name:bar] map[name:foo]]
```

## Features

- ✅ Lexer and parser for UQL syntax
- ✅ Support for piped commands
- ✅ Basic commands (hello, ping, echo, count, limit)
- ✅ Data transformation (project, project-away, project-reorder, extend, order by)
- ✅ Data filtering (distinct, scope)
- ✅ Data aggregation (summarize, pivot)
- ✅ Data parsing (JSON, CSV, XML, YAML)
- ✅ String manipulation functions
- ✅ Math and trigonometric functions
- ✅ Type conversion functions
- ✅ Date/time functions
- 🚧 Where clause (partial support)
- 🚧 mv-expand (planned)
- 🚧 JSONata support (planned)

## Development

### Running Tests

```bash
go test -v
```

### Building

```bash
go build
```

## License

Apache-2.0 License - same as the original UQL project

## Credits

This is a Go port of the original [UQL](https://github.com/yesoreyeram/uql) TypeScript/JavaScript implementation by Sriramajeyam Sugumaran.
