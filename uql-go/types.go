package uql

// Operator represents comparison and logical operators
type Operator string

const (
	OpGreaterThan          Operator = ">"
	OpGreaterThanOrEqual   Operator = ">="
	OpLessThan             Operator = "<"
	OpLessThanOrEqual      Operator = "<="
	OpEqual                Operator = "=="
	OpNotEqual             Operator = "!="
	OpRegexMatch           Operator = "=~"
	OpRegexNotMatch        Operator = "!~"
	OpIn                   Operator = "in"
	OpNotIn                Operator = "!in"
	OpInCaseInsensitive    Operator = "in~"
	OpNotInCaseInsensitive Operator = "!in~"
	OpBetween              Operator = "between"
	OpInside               Operator = "inside"
	OpOutside              Operator = "outside"
	OpMatchesRegex         Operator = "matches regex"
	OpNotMatchesRegex      Operator = "!matches regex"
	OpContains             Operator = "contains"
	OpNotContains          Operator = "!contains"
	OpContainsCS           Operator = "contains_cs"
	OpNotContainsCS        Operator = "!contains_cs"
	OpStartsWith           Operator = "startswith"
	OpNotStartsWith        Operator = "!startswith"
	OpStartsWithCS         Operator = "startswith_cs"
	OpNotStartsWithCS      Operator = "!startswith_cs"
	OpEndsWith             Operator = "endswith"
	OpNotEndsWith          Operator = "!endswith"
	OpEndsWithCS           Operator = "endswith_cs"
	OpNotEndsWithCS        Operator = "!endswith_cs"
)

// FunctionName represents available functions
type FunctionName string

// String manipulation functions
const (
	FnReplaceString FunctionName = "replace_string"
	FnSubstring     FunctionName = "substring"
	FnSplit         FunctionName = "split"
	FnStrcat        FunctionName = "strcat"
	FnStrlen        FunctionName = "strlen"
	FnToupper       FunctionName = "toupper"
	FnTrimEnd       FunctionName = "trim_end"
	FnTrimStart     FunctionName = "trim_start"
	FnTrim          FunctionName = "trim"
	FnReverse       FunctionName = "reverse"
	FnExtract       FunctionName = "extract"
)

// Math functions
const (
	FnSum        FunctionName = "sum"
	FnDiff       FunctionName = "diff"
	FnMul        FunctionName = "mul"
	FnDiv        FunctionName = "div"
	FnMin        FunctionName = "min"
	FnMax        FunctionName = "max"
	FnMean       FunctionName = "mean"
	FnFirst      FunctionName = "first"
	FnLast       FunctionName = "last"
	FnLatest     FunctionName = "latest"
	FnCount      FunctionName = "count"
	FnDCount     FunctionName = "dcount"
	FnLog        FunctionName = "log"
	FnLog2       FunctionName = "log2"
	FnLog10      FunctionName = "log10"
	FnAbs        FunctionName = "abs"
	FnFloor      FunctionName = "floor"
	FnCeil       FunctionName = "ceil"
	FnRound      FunctionName = "round"
	FnPow        FunctionName = "pow"
	FnPercentage FunctionName = "percentage"
)

// Array functions
const (
	FnDistinct         FunctionName = "distinct"
	FnPack             FunctionName = "pack"
	FnArrayFromEntries FunctionName = "array_from_entries"
	FnArrayToMap       FunctionName = "array_to_map"
	FnBagPack          FunctionName = "bag_pack"
	FnKV               FunctionName = "kv"
)

// Trigonometric functions
const (
	FnSin   FunctionName = "sin"
	FnCos   FunctionName = "cos"
	FnTan   FunctionName = "tan"
	FnSinh  FunctionName = "sinh"
	FnCosh  FunctionName = "cosh"
	FnTanh  FunctionName = "tanh"
	FnAsin  FunctionName = "asin"
	FnAcos  FunctionName = "acos"
	FnAtan  FunctionName = "atan"
	FnAtan2 FunctionName = "atan2"
	FnAsinh FunctionName = "asinh"
	FnAcosh FunctionName = "acosh"
	FnAtanh FunctionName = "atanh"
)

// URL parsing functions
const (
	FnParseURL      FunctionName = "parse_url"
	FnParseURLQuery FunctionName = "parse_urlquery"
)

// Other functions
const (
	FnAddDatetime                    FunctionName = "add_datetime"
	FnAtob                           FunctionName = "atob"
	FnBtoa                           FunctionName = "btoa"
	FnFormatDatetime                 FunctionName = "format_datetime"
	FnRandom                         FunctionName = "random"
	FnSign                           FunctionName = "sign"
	FnStartOfDay                     FunctionName = "startofday"
	FnStartOfHour                    FunctionName = "startofhour"
	FnStartOfMinute                  FunctionName = "startofminute"
	FnStartOfMonth                   FunctionName = "startofmonth"
	FnStartOfWeek                    FunctionName = "startofweek"
	FnStartOfYear                    FunctionName = "startofyear"
	FnToBool                         FunctionName = "tobool"
	FnToDatetime                     FunctionName = "todatetime"
	FnToDouble                       FunctionName = "todouble"
	FnToFloat                        FunctionName = "tofloat"
	FnToInt                          FunctionName = "toint"
	FnToLong                         FunctionName = "tolong"
	FnToLower                        FunctionName = "tolower"
	FnToNumber                       FunctionName = "tonumber"
	FnToString                       FunctionName = "tostring"
	FnToUnixTime                     FunctionName = "tounixtime"
	FnUnixTimeMicrosecondsToDatetime FunctionName = "unixtime_microseconds_todatetime"
	FnUnixTimeMillisecondsToDatetime FunctionName = "unixtime_milliseconds_todatetime"
	FnUnixTimeNanosecondsToDatetime  FunctionName = "unixtime_nanoseconds_todatetime"
	FnUnixTimeSecondsToDatetime      FunctionName = "unixtime_seconds_todatetime"
)

// Conditional functions
const (
	FnCountIf FunctionName = "countif"
	FnSumIf   FunctionName = "sumif"
	FnMinIf   FunctionName = "minif"
	FnMaxIf   FunctionName = "maxif"
)

// CommandType represents the type of command
type CommandType string

const (
	CmdComment        CommandType = "comment"
	CmdHello          CommandType = "hello"
	CmdPing           CommandType = "ping"
	CmdEcho           CommandType = "echo"
	CmdCount          CommandType = "count"
	CmdLimit          CommandType = "limit"
	CmdCommand        CommandType = "command"
	CmdOrderBy        CommandType = "orderby"
	CmdProject        CommandType = "project"
	CmdProjectAway    CommandType = "project-away"
	CmdProjectReorder CommandType = "project-reorder"
	CmdExtend         CommandType = "extend"
	CmdSummarize      CommandType = "summarize"
	CmdPivot          CommandType = "pivot"
	CmdRange          CommandType = "range"
	CmdScope          CommandType = "scope"
	CmdWhere          CommandType = "where"
	CmdDistinct       CommandType = "distinct"
	CmdMvExpand       CommandType = "mv-expand"
	CmdJSONata        CommandType = "jsonata"
	CmdParseJSON      CommandType = "parse-json"
	CmdParseCSV       CommandType = "parse-csv"
	CmdParseXML       CommandType = "parse-xml"
	CmdParseYAML      CommandType = "parse-yaml"
)

// TypedValue represents a value with its type
type TypedValue struct {
	Type  string
	Value interface{}
	Alias string
}

// StringType creates a string typed value
func StringType(value string) TypedValue {
	return TypedValue{Type: "string", Value: value}
}

// NumberType creates a number typed value
func NumberType(value float64) TypedValue {
	return TypedValue{Type: "number", Value: value}
}

// RefType creates a reference typed value
func RefType(value string, alias ...string) TypedValue {
	tv := TypedValue{Type: "ref", Value: value}
	if len(alias) > 0 {
		tv.Alias = alias[0]
	}
	return tv
}

// IdentifierType creates an identifier typed value
func IdentifierType(value string) TypedValue {
	return TypedValue{Type: "identifier", Value: value}
}

// FunctionCall represents a function call
type FunctionCall struct {
	Type     string
	Operator FunctionName
	Args     []TypedValue
	Alias    string
}

// OrderByArg represents an order by argument
type OrderByArg struct {
	Field     string
	Direction string // "asc" or "desc"
}

// SummarizeAssignment represents a summarize assignment
type SummarizeAssignment struct {
	Alias    string
	Operator FunctionName
	Args     []TypedValue
}

// SummarizeItem represents a summarize item
type SummarizeItem struct {
	Metrics []SummarizeAssignment
	By      []TypedValue
}

// PivotItem represents a pivot item
type PivotItem struct {
	Metric SummarizeAssignment
	Fields []TypedValue
}

// ParseArg represents a parse argument
type ParseArg struct {
	Identifier string
	Value      string
}

// RangeValue represents a range value
type RangeValue struct {
	Start interface{} // number or string
	End   float64
	Step  interface{} // number or string
}

// Command represents a UQL command
type Command struct {
	Type  CommandType
	Value interface{}
}

// CommandResult represents the result of command evaluation
type CommandResult struct {
	Output  interface{}
	Context map[string]interface{}
}
