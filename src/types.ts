export type type_str_type = { type: "string"; value: string };
export type type_ref_type = { type: "ref"; value: string; alias?: string };
export type type_num_type = { type: "number"; value: number };
export type type_any_type = type_str_type | type_num_type | type_ref_type;
export type type_identifier_type = { type: "identifier"; value: string };
export type type_function_arg = type_any_type | type_identifier_type;
export type type_where_arg = type_any_type | { type: "value_array"; value: type_any_type[] } | { type: "operation"; value: Operator };
export type type_function = { alias?: string; operator: FunctionName; args: type_function_arg[]; type: "function" };
export type type_orderby_arg = { field: string; direction: "asc" | "desc" };
export type type_summarize_arg = type_str_type;
export type type_summarize_function = { operator: FunctionName; args: type_summarize_arg[] } | { operator: ConditionalFunctionName; condition: type_where_arg[]; ref: type_ref_type };
export type type_summarize_assignment = { alias?: string } & type_summarize_function;
export type type_summarize_item = { metrics: type_summarize_assignment[]; by: type_summarize_arg[] };
export type type_parse_arg = { identifier: string; value: string };

export type Operator =
  | ">"
  | ">="
  | "<"
  | "<="
  | "=="
  | "!="
  | "=~"
  | "!~"
  | "in"
  | "!in"
  | "in~"
  | "!in~"
  | "between"
  | "inside"
  | "outside"
  | "matches regex"
  | "!matches regex"
  | "contains"
  | "!contains"
  | "contains_cs"
  | "!contains_cs"
  | "startswith"
  | "!startswith"
  | "startswith_cs"
  | "!startswith_cs"
  | "endswith"
  | "!endswith"
  | "endswith_cs"
  | "!endswith_cs";
export type ConditionalFunctionName = "countif" | "sumif" | "minif" | "maxif";
export type FunctionName =
  | "add_datetime"
  | "array_from_entries"
  | "array_to_map"
  | "bag_pack"
  | "ceil"
  | "cos"
  | "count"
  | "dcount"
  | "diff"
  | "distinct"
  | "div"
  | "extract"
  | "first"
  | "floor"
  | "format_datetime"
  | "kv"
  | "last"
  | "latest"
  | "log"
  | "log10"
  | "log2"
  | "max"
  | "mean"
  | "min"
  | "mul"
  | "pack"
  | "parse_url"
  | "parse_urlquery"
  | "percentage"
  | "pow"
  | "random"
  | "replace_string"
  | "reverse"
  | "round"
  | "sign"
  | "sin"
  | "split"
  | "startofday"
  | "startofhour"
  | "startofminute"
  | "startofmonth"
  | "startofweek"
  | "startofyear"
  | "strcat"
  | "strlen"
  | "sum"
  | "tan"
  | "tobool"
  | "todatetime"
  | "todouble"
  | "tofloat"
  | "toint"
  | "tolong"
  | "tolower"
  | "tonumber"
  | "tostring"
  | "tounixtime"
  | "toupper"
  | "trim_end"
  | "trim_start"
  | "trim"
  | "unixtime_microseconds_todatetime"
  | "unixtime_milliseconds_todatetime"
  | "unixtime_nanoseconds_todatetime"
  | "unixtime_seconds_todatetime";

export type CommandType =
  | "comment"
  | "hello"
  | "ping"
  | "echo"
  | "count"
  | "limit"
  | "command"
  | "orderby"
  | "project"
  | "project-away"
  | "project-reorder"
  | "extend"
  | "summarize"
  | "range"
  | "scope"
  | "where"
  | "distinct"
  | "mv-expand"
  | "jsonata"
  | "parse-json"
  | "parse-csv"
  | "parse-xml"
  | "parse-yaml";

type CommandBase<T extends CommandType> = { type: T };

type CommandComment = { value: string } & CommandBase<"comment">;
type CommandHello = {} & CommandBase<"hello">;
type CommandPing = {
  value: string;
} & CommandBase<"ping">;
type CommandEcho = {
  value: string;
} & CommandBase<"echo">;
type CommandScope = { value: type_ref_type } & CommandBase<"scope">;
type CommandWhere = { value: type_where_arg[] } & CommandBase<"where">;
type CommandDistinct = { value: type_ref_type | undefined } & CommandBase<"distinct">;
type CommandMvExpand = { value: type_ref_type & { alias: string } } & CommandBase<"mv-expand">;
type CommandParseJSON = { args: type_parse_arg[][] } & CommandBase<"parse-json">;
type CommandParseCSV = { args: type_parse_arg[][] } & CommandBase<"parse-csv">;
type CommandParseXML = { args: type_parse_arg[][] } & CommandBase<"parse-xml">;
type CommandParseYAML = { args: type_parse_arg[][] } & CommandBase<"parse-yaml">;
type CommandJSONata = { expression: string } & CommandBase<"jsonata">;
type CommandCount = {} & CommandBase<"count">;
type CommandLimit = {
  value: number;
} & CommandBase<"limit">;
type CommandCommand = {
  value: {
    type: "function";
    operator: FunctionName;
    args: type_function_arg[];
  };
} & CommandBase<"command">;
type CommandOrderBy = {
  value: type_orderby_arg[];
} & CommandBase<"orderby">;
type CommandProject = {
  value: (type_function | (type_ref_type & { alias: string }))[];
} & CommandBase<"project">;
type CommandProjectReorder = {
  value: type_ref_type[];
} & CommandBase<"project-reorder">;
type CommandProjectAway = {
  value: type_ref_type[];
} & CommandBase<"project-away">;
type CommandExtend = {
  value: type_function[] | type_ref_type[];
} & CommandBase<"extend">;
type CommandSummarize = {
  value: type_summarize_item;
} & CommandBase<"summarize">;
type CommandRange = {
  value: { start: number; end: number; step: number } | { start: string; end: number; step: string };
} & CommandBase<"range">;
export type Command =
  | CommandComment
  | CommandHello
  | CommandPing
  | CommandEcho
  | CommandCount
  | CommandLimit
  | CommandCommand
  | CommandOrderBy
  | CommandProject
  | CommandProjectAway
  | CommandProjectReorder
  | CommandScope
  | CommandWhere
  | CommandDistinct
  | CommandMvExpand
  | CommandJSONata
  | CommandParseJSON
  | CommandParseCSV
  | CommandParseXML
  | CommandParseYAML
  | CommandExtend
  | CommandSummarize
  | CommandRange;

export type CommandResult = { context: Record<string, unknown>; output: unknown };
