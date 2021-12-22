@preprocessor typescript

@{%

  import * as moo from "moo";
  const oqlLexer = moo.compile({
    ws: /[ \t]+/,
    comment: /\/\/.*?$/,
    nl: {
      match: "\n",
      lineBreaks: true,
    },
    sq_string: {
        match: /'(?:[^\n\\']|\\['"\\ntbfr])*'/,
        value: (s)=>s.slice(1,-1),
    },
    string: {
      match: /"(?:[^\n\\"]|\\['"\\ntbfr])*"/,
      value: (s) => JSON.parse(s),
    },
    number: {
      // @ts-ignore Ignore the error for now until finding a better regex
      match: /-?[0-9]+(?:\.[0-9]+)?/,
      value: (s) => Number(s),
    },
    pipe: "|",
    plus: "+",
    dash:"-",
    mul:"*",
    divide:"/",
    mod:"%",
    eq: "==",
    lparan: "(",
    rparan: ")",
    comma: ",",
    assignment: "=",
    identifier: {
      match: /[a-z_][a-zA-Z_0-9]*/,
      type: moo.keywords({
        str: "str",
      }),
    },
  });
  const pick = (idx: number) => (d: unknown[]) => d[idx];
  const merge = (newIndex:number, oldIndex:number) => (d:unknown[][]) => [d[newIndex], ...d[oldIndex]];
  const as_array = (idx:number) => (d: unknown[]) => [d[idx]];
  const as_string = (d: { value: string}[]) => d[0].value;
  const as_number = (d: { value: string}[]) => +d[0].value;
  const dispose : () => void = () => null;

%}

@lexer oqlLexer

input   
    -> commands                                             {% pick(0) %}
commands
    ->  command __                                          {% as_array(0) %}
    |   command __ nlo __ pipeo __ commands                 {% merge(0,6) %}
command 
    -> "hello"                                              {% d => ({ type: "hello" }) %}
    |  "ping"                                               {% d => ({ type: "ping", value: "pong" }) %}
    |  "echo" __ str                                        {% d => ({ type: "echo", value: d[2] }) %}
    |  "count"                                              {% d => ({ type: "count" }) %}
    |  "limit" __ number                                    {% d => ({ type: "limit", value: d[2] }) %}
    |  function                                             {% d => ({ type: "command" , value: d[0] }) %}
    |  command_orderby                                      {% d => ({ type: "orderby", value: d[0] })%}
    |  command_extend                                       {% d => ({ type: "extend", value: d[0] })%}
    |  command_project_away                                 {% d => ({ type: "project-away", value: d[0] })%}
    |  command_project                                      {% d => ({ type: "project", value: d[0] })%}
    |  command_parse_json                                   {% d => ({ type: "parse-json", args: d[0] })%}
    |  command_parse_csv                                    {% d => ({ type: "parse-csv", args: d[0] })%}
    |  command_parse_xml                                    {% d => ({ type: "parse-xml", args: d[0] })%}
    |  command_scope                                        {% d => ({ type: "scope", value: d[0] })%}
    |  command_summarize                                    {% d => ({ type: "summarize", value: d[0] })%}
    |  command_range                                        {% d => ({ type: "range", value: d[0] })%}
# Command Function
function_assignments
    -> function_assignment                                  {% as_array(0) %}
    |  function_assignment __ "," __ function_assignments   {% merge(0,4) %}
function_assignment
    -> str:* "=":* expression                               {% d => ({ alias: d[0][0], ...d[2] })%}
    |  str:* "=":* function                                 {% d => ({ alias: d[0][0], ...d[2] })%}
    |  str:* "=":* ref_type                                 {% d => ({ alias: d[0][0], ...d[2] })%}
    |  ref_type                                             {% d => d[0] %}
expression
    ->  %lparan __ expression_args:* %rparan                {% d => ({ type: "expression", args: d[2][0]||[] }) %}
expression_args       
    -> expression_arg __                                    {% as_array(0) %}
    |  expression_arg _ expression_args                     {% merge(0,2)  %}
expression_arg
    -> num_type                                             {% d => d[0] %}
    |  str_type                                             {% d => d[0] %}
    |  ref_type                                             {% d => d[0] %}
    |  %plus                                                {% d => ({ type: "operation", value: "+" }) %}
    |  %dash                                                {% d => ({ type: "operation", value: "-" }) %}
    |  %mul                                                 {% d => ({ type: "operation", value: "*" }) %}
    |  %divide                                              {% d => ({ type: "operation", value: "/" }) %}
    |  %mod                                                 {% d => ({ type: "operation", value: "%" }) %}
    |  function                                             {% d => ({ type: "function", value: d[0] }) %}
function
    ->  function_name "(" __ function_args:* ")"            {% d => ({ type: "function", operator: d[0], args: d[3][0]||[] }) %}
function_name
    -> "count"                                              {% as_string %}
    |  "sum"                                                {% as_string %}
    |  "diff"                                               {% as_string %}
    |  "mul"                                                {% as_string %}
    |  "min"                                                {% as_string %}
    |  "max"                                                {% as_string %}
    |  "mean"                                               {% as_string %}
    |  "strcat"                                             {% as_string %}
    |  "dcount"                                             {% as_string %}
    |  "distinct"                                           {% as_string %}
    |  "random"                                             {% as_string %}
    |  "toupper"                                            {% as_string %}
    |  "tolower"                                            {% as_string %}
    |  "strlen"                                             {% as_string %}
    |  "trim"                                               {% as_string %}
    |  "trim_start"                                         {% as_string %}
    |  "trim_end"                                           {% as_string %}
    |  "toint"                                              {% as_string %}
    |  "tonumber"                                           {% as_string %}
    |  "tolong"                                             {% as_string %}
    |  "tobool"                                             {% as_string %}
    |  "tostring"                                           {% as_string %}
    |  "todouble"                                           {% as_string %}
    |  "tofloat"                                            {% as_string %}
    |  "todatetime"                                         {% as_string %}
    |  "unixtime_seconds_todatetime"                        {% as_string %}
    |  "unixtime_nanoseconds_todatetime"                    {% as_string %}
    |  "unixtime_milliseconds_todatetime"                   {% as_string %}
    |  "unixtime_microseconds_todatetime"                   {% as_string %}
function_args       
    -> function_arg __                                      {% as_array(0) %}
    |  function_arg __ "," __ function_args                 {% merge(0,4)  %}
function_arg   
    -> str_type                                             {% pick(0) %}
    |  ref_type                                             {% pick(0) %}
    |  num_type                                             {% pick(0) %}
# Command : Order by
command_orderby
    -> "order" _ "by" _ orderby_args                        {% pick(4) %}
orderby_args
    -> orderby_arg __                                       {% as_array(0) %}
    |  orderby_arg __ "," __ orderby_args                   {% merge(0,4) %}
orderby_arg
    -> str __ "asc"                                         {% d => ({ field: d[0] , direction: "asc"})%}
    |  str __ "desc"                                        {% d => ({ field: d[0] , direction: "desc"})%}
# Command : Extend
command_extend
    -> "extend" _ function_assignments                      {% d => d[2] %}
# Command : Project
command_project
    -> "project" _ function_assignments                     {% d => d[2] %}
# Command : Project Away
command_project_away
    -> "project" %dash "away" _ ref_types                   {% d => d[4] %}
# Command : Scope 
command_scope
    -> "scope" _ ref_type                                   {% d => d[2] %}
# Command : Parse json
command_parse_json
    -> "parse" %dash "json" __ parse_args:*                 {% d => d[4] %}
# Command : Parse csv
command_parse_csv
    -> "parse" %dash "csv" __ parse_args:*                  {% d => d[4] %}
# Command : Parse xml
command_parse_xml
    -> "parse" %dash "xml" __ parse_args:*                  {% d => d[4] %}
parse_args
    -> parse_arg                                            {% as_array(0) %}
    |  parse_arg __ parse_args                              {% merge(0,2) %}
parse_arg
    -> %dash %dash %identifier __ str                       {% d => ({ identifier: d[2].value, value: d[4] }) %}
    |  %dash %dash %identifier __ str_type                  {% d => ({ identifier: d[2].value, value: d[4].value }) %}
    |  %dash %dash %identifier __ %identifier               {% d => ({ identifier: d[2].value, value: d[4].value }) %}
# Command : Summarize
command_summarize 
     ->  summarize_item                                              {% pick(0) %}
     |   summarize_item _ "by" _ summarize_args:*                    {% d => ({ ...d[0], by: d[4][0] })%}
summarize_item
     -> "summarize" _ summarize_assignments                          {% d => ({ metrics: d[2], by :[] })%}
summarize_assignments
     ->  summarize_assignment                                        {% d => [ d[0] ] %}
     |   summarize_assignment __ "," __ summarize_assignments        {% merge(0,4) %}
summarize_assignment 
     ->  str:*  "=":* summarize_function                             {% d => ({ operator: d[2].operator, alias: d[0][0], args: d[2].args })%}
summarize_function
     ->  function_name "(" summarize_args:* ")"                      {% d => ({ operator: d[0], args: d[2][0]||[] })%}
summarize_args       
     ->  summarize_arg __                                            {% as_array(0) %}   
     |   summarize_arg __ "," __ summarize_args                      {% merge(0,4) %}
summarize_arg   
     -> ref_type                                            {% pick(0) %}
# Command : Range
command_range    
    -> range_item                                           {% pick(0) %}
    |  range_item _ "step" _ number                         {% d =>  ({...d[0], step : d[4]})%}
    |  range_item _ "step" _ str                            {% d =>  ({...d[0], step : d[4]})%}
range_item          
    -> "range" _ "from" _ number _ "to" _ number            {% d => ({ start: d[4], end: d[8], step: 1 })%}
    |  "range" _ "from" _ str _ "to" _ str                  {% d => ({ start: d[4], end: d[8], step: "" })%}
#region Utils
str_type        
    -> %str "(" str ")"                                     {% d => ({ type: "string", value: d[2] })%}
    |  %sq_string                                           {% d => ({ type: "string", value:d[0].value}) %}
ref_type        
    -> str                                                  {% d => ({ type: "ref", value: d[0] })%}
    | "[" str "]"                                           {% d => ({ type: "ref", value: d[1] })%}
    | "[" %sq_string "]"                                    {% d => ({ type: "ref", value: d[1] })%}
num_type        
    -> number                                               {% d => ({ type: "number", value: d[0] })%}
ref_types
    -> ref_type                                             {% as_array(0) %}
    |  ref_type __ "," __ ref_types                         {% merge(0,4) %}
str             -> %string                                  {% as_string %}                 # string
number          -> %number                                  {% as_number %}                 # number
nlo             -> %nl:*                                    {% dispose %}                   # optional newline
pipeo           -> %pipe:*                                  {% dispose %}                   # optional pipe
__              -> %ws:*                                    {% dispose %}                   # optional whitespace
_               -> %ws:+                                    {% dispose %}                   # mandatory whitespace
#endregion