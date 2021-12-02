type type_str_type = { type: "string"; value: string };
type type_ref_type = { type: "ref"; value: string };
type type_num_type = { type: "number"; value: number };
type type_function_arg = type_str_type | type_ref_type | type_num_type;
type type_function = { operator: FunctionName; args: type_function_arg[]; type: "function" };
export type type_function_assignment = { alias: string } & type_function;
type type_orderby_arg = { field: string; direction: "asc" | "desc" };
type type_summarize_arg = type_str_type;
type type_summarize_function = { operator: FunctionName; args: type_summarize_arg[] };
export type type_summarize_assignment = { alias?: string } & type_summarize_function;
type type_summarize_item = { metrics: type_summarize_assignment[]; by: type_summarize_arg[] };

export type FunctionName =
  | "count"
  | "sum"
  | "diff"
  | "mul"
  | "min"
  | "max"
  | "mean"
  | "strcat"
  | "dcount"
  | "distinct"
  | "random"
  | "toupper"
  | "tolower"
  | "strlen"
  | "trim"
  | "trim_start"
  | "trim_end";

type CommandType = "hello" | "ping" | "echo" | "count" | "limit" | "command" | "orderby" | "project" | "project-away" | "extend" | "summarize" | "range";
type CommandBase<T extends CommandType> = { type: T };

type CommandHello = {} & CommandBase<"hello">;
type CommandPing = {
  value: string;
} & CommandBase<"ping">;
type CommandEcho = {
  value: string;
} & CommandBase<"echo">;
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
  value: (type_function | type_function_assignment | type_ref_type)[];
} & CommandBase<"project">;
type CommandProjectAway = {
  value: type_ref_type[];
} & CommandBase<"project-away">;
type CommandExtend = {
  value: type_function_assignment[];
} & CommandBase<"extend">;
type CommandSummarize = {
  value: type_summarize_item;
} & CommandBase<"summarize">;
type CommandRange = {
  value: { start: number; end: number; step: number } | { start: string; end: number; step: string };
} & CommandBase<"range">;
export type Command =
  | CommandHello
  | CommandPing
  | CommandEcho
  | CommandCount
  | CommandLimit
  | CommandCommand
  | CommandOrderBy
  | CommandProject
  | CommandProjectAway
  | CommandExtend
  | CommandSummarize
  | CommandRange;