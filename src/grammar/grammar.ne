@preprocessor typescript

@{%

  import * as moo from "moo";
  const oqlLexer = moo.compile({
    ws: /[ \t]+/,
    comment: {
        match: /#[^\n]*/,
        value: s => s.substring(1)
    },
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
      match: /-?[\d.]+(?:e-?\d+)?/,
      value: (s) => Number(s),
    },
    pipe: "|",
    plus: "+",
    dash:"-",
    mul:"*",
    divide:"/",
    mod:"%",
    gt:">",
    lt:"<",
    eq: "==",
    ne: "!=",
    lparan: "(",
    rparan: ")",
    comma: ",",
    tilde: "~",
    exclaim: "!",
    assignment: "=",
    return: "\r\n",
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
    |   command __ %comment:* "\r\n" __ pipeo __ commands   {% merge(0,7) %}
    |   command __ %comment:* nlo __ pipeo __ commands      {% merge(0,7) %}
command 
    -> line_comment                                         {% pick(0) %}
    |  "hello"                                              {% d => ({ type: "hello" }) %}
    |  "ping"                                               {% d => ({ type: "ping", value: "pong" }) %}
    |  "echo" __ str                                        {% d => ({ type: "echo", value: d[2] }) %}
    |  "count"                                              {% d => ({ type: "count" }) %}
    |  "limit" __ number                                    {% d => ({ type: "limit", value: d[2] }) %}
    |  function                                             {% d => ({ type: "command" , value: d[0] }) %}
    |  command_orderby                                      {% d => ({ type: "orderby", value: d[0] })%}
    |  command_extend                                       {% d => ({ type: "extend", value: d[0] })%}
    |  command_project_away                                 {% d => ({ type: "project-away", value: d[0] })%}
    |  command_project_reorder                              {% d => ({ type: "project-reorder", value: d[0] })%}
    |  command_project                                      {% d => ({ type: "project", value: d[0] })%}
    |  command_parse_json                                   {% d => ({ type: "parse-json", args: d[0] })%}
    |  command_parse_csv                                    {% d => ({ type: "parse-csv", args: d[0] })%}
    |  command_parse_xml                                    {% d => ({ type: "parse-xml", args: d[0] })%}
    |  command_parse_yaml                                   {% d => ({ type: "parse-yaml", args: d[0] })%}
    |  command_scope                                        {% d => ({ type: "scope", value: d[0] })%}
    |  command_where                                        {% d => ({ type: "where", value: d[0] })%}
    |  command_distinct                                     {% d => ({ type: "distinct", value: d[0] })%}
    |  command_mv_expand                                    {% d => ({ type: "mv-expand", value: d[0] })%}
    |  command_summarize                                    {% d => ({ type: "summarize", value: d[0] })%}
    |  command_range                                        {% d => ({ type: "range", value: d[0] })%}
    |  "jsonata" __ str                                     {% d => ({ type: "jsonata", expression: d[2] }) %}
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
    |  expression_arg __ expression_args                    {% merge(0,2)  %}
expression_arg
    -> any_type                                             {% d => d[0] %}
    |  %lparan __ function_args __ %rparan                  {% d => ({ type : "value_array", value: d[2] }) %}
    |  %plus                                                {% d => ({ type: "operation", value: "+" }) %}
    |  %dash                                                {% d => ({ type: "operation", value: "-" }) %}
    |  %mul                                                 {% d => ({ type: "operation", value: "*" }) %}
    |  %divide                                              {% d => ({ type: "operation", value: "/" }) %}
    |  %mod                                                 {% d => ({ type: "operation", value: "%" }) %}
    |  %gt %assignment                                      {% d => ({ type: "operation", value: ">=" }) %}
    |  %gt                                                  {% d => ({ type: "operation", value: ">" }) %}
    |  %lt %assignment                                      {% d => ({ type: "operation", value: "<=" }) %}
    |  %lt                                                  {% d => ({ type: "operation", value: "<" }) %}
    |  %eq                                                  {% d => ({ type: "operation", value: "==" }) %}
    |  %ne                                                  {% d => ({ type: "operation", value: "!=" }) %}
    |  %assignment %tilde                                   {% d => ({ type: "operation", value: "=~" }) %}
    |  %exclaim %tilde                                      {% d => ({ type: "operation", value: "!~" }) %}
    |  "between"                                            {% d => ({ type: "operation", value: "between" }) %}
    |  "inside"                                             {% d => ({ type: "operation", value: "inside" }) %}
    |  "outside"                                            {% d => ({ type: "operation", value: "outside" }) %}
    |  "matches" _ "regex"                                  {% d => ({ type: "operation", value: "matches regex" }) %}
    |  %exclaim "matches" _ "regex"                         {% d => ({ type: "operation", value: "!matches regex" }) %}
    |  "not" _ "matches" _ "regex"                          {% d => ({ type: "operation", value: "!matches regex" }) %}
    |  "in"                                                 {% d => ({ type: "operation", value: "in" }) %}
    |  %exclaim "in"                                        {% d => ({ type: "operation", value: "!in" }) %}
    |  "not" _ "in"                                         {% d => ({ type: "operation", value: "!in" }) %}
    |  "in" %tilde                                          {% d => ({ type: "operation", value: "in~" }) %}
    |  %exclaim "in" %tilde                                 {% d => ({ type: "operation", value: "!in~" }) %}
    |  "not" _ "in" %tilde                                  {% d => ({ type: "operation", value: "!in~" }) %}
    |  "contains"                                           {% d => ({ type: "operation", value: "contains" }) %}
    |  %exclaim "contains"                                  {% d => ({ type: "operation", value: "!contains" }) %}
    |  "not" _ "contains"                                   {% d => ({ type: "operation", value: "!contains" }) %}
    |  "contains_cs"                                        {% d => ({ type: "operation", value: "contains_cs" }) %}
    |  %exclaim "contains_cs"                               {% d => ({ type: "operation", value: "!contains_cs" }) %}
    |  "not" _ "contains_cs"                                {% d => ({ type: "operation", value: "!contains_cs" }) %}
    |  "startswith"                                         {% d => ({ type: "operation", value: "startswith" }) %}
    |  %exclaim "startswith"                                {% d => ({ type: "operation", value: "!startswith" }) %}
    |  "not" _ "startswith"                                 {% d => ({ type: "operation", value: "!startswith" }) %}
    |  "startswith_cs"                                      {% d => ({ type: "operation", value: "startswith_cs" }) %}
    |  %exclaim "startswith_cs"                             {% d => ({ type: "operation", value: "!startswith_cs" }) %}
    |  "not" _ "startswith_cs"                              {% d => ({ type: "operation", value: "!startswith_cs" }) %}
    |  "endswith"                                           {% d => ({ type: "operation", value: "endswith" }) %}
    |  %exclaim "endswith"                                  {% d => ({ type: "operation", value: "!endswith" }) %}
    |  "not" _ "endswith"                                   {% d => ({ type: "operation", value: "!endswith" }) %}
    |  "endswith_cs"                                        {% d => ({ type: "operation", value: "endswith_cs" }) %}
    |  %exclaim "endswith_cs"                               {% d => ({ type: "operation", value: "!endswith_cs" }) %}
    |  "not" _ "endswith_cs"                                {% d => ({ type: "operation", value: "!endswith_cs" }) %}
    |  function                                             {% d => ({ type: "function", value: d[0] }) %}
function
    ->  function_name "(" __ function_args:* ")"            {% d => ({ type: "function", operator: d[0], args: d[3][0]||[] }) %}
conditional_function_name
    -> "countif" {% as_string %}
    | "sumif"    {% as_string %}
    | "minif"    {% as_string %}
    | "maxif"    {% as_string %}
function_name
    -> "add_datetime"                                       {% as_string %}
    |  "array_from_entries"                                 {% as_string %}
    |  "array_to_map"                                       {% as_string %}
    |  "atob"                                               {% as_string %}
    |  "bag_pack"                                           {% as_string %}
    |  "btoa"                                               {% as_string %}
    |  "ceil"                                               {% as_string %}
    |  "cos"                                                {% as_string %}
    |  "count"                                              {% as_string %}
    |  "dcount"                                             {% as_string %}
    |  "diff"                                               {% as_string %}
    |  "distinct"                                           {% as_string %}
    |  "div"                                                {% as_string %}
    |  "extract"                                            {% as_string %}
    |  "first"                                              {% as_string %}
    |  "floor"                                              {% as_string %}
    |  "format_datetime"                                    {% as_string %}
    |  "kv"                                                 {% as_string %}
    |  "last"                                               {% as_string %}
    |  "latest"                                             {% as_string %}
    |  "log"                                                {% as_string %}
    |  "log10"                                              {% as_string %}
    |  "log2"                                               {% as_string %}
    |  "max"                                                {% as_string %}
    |  "mean"                                               {% as_string %}
    |  "min"                                                {% as_string %}
    |  "mul"                                                {% as_string %}
    |  "pack"                                               {% as_string %}
    |  "parse_url"                                          {% as_string %}
    |  "parse_urlquery"                                     {% as_string %}
    |  "percentage"                                         {% as_string %}
    |  "pow"                                                {% as_string %}
    |  "random"                                             {% as_string %}
    |  "replace_string"                                     {% as_string %}
    |  "reverse"                                            {% as_string %}
    |  "round"                                              {% as_string %}
    |  "sign"                                               {% as_string %}
    |  "sin"                                                {% as_string %}
    |  "split"                                              {% as_string %}
    |  "startofday"                                         {% as_string %}
    |  "startofhour"                                        {% as_string %}
    |  "startofminute"                                      {% as_string %}
    |  "startofmonth"                                       {% as_string %}
    |  "startofweek"                                        {% as_string %}
    |  "startofyear"                                        {% as_string %}
    |  "strcat"                                             {% as_string %}
    |  "strlen"                                             {% as_string %}
    |  "substring"                                          {% as_string %}
    |  "sum"                                                {% as_string %}
    |  "tan"                                                {% as_string %}
    |  "tobool"                                             {% as_string %}
    |  "todatetime"                                         {% as_string %}
    |  "todouble"                                           {% as_string %}
    |  "tofloat"                                            {% as_string %}
    |  "toint"                                              {% as_string %}
    |  "tolong"                                             {% as_string %}
    |  "tolower"                                            {% as_string %}
    |  "tonumber"                                           {% as_string %}
    |  "tostring"                                           {% as_string %}
    |  "tounixtime"                                         {% as_string %}
    |  "toupper"                                            {% as_string %}
    |  "trim_end"                                           {% as_string %}
    |  "trim_start"                                         {% as_string %}
    |  "trim"                                               {% as_string %}
    |  "unixtime_microseconds_todatetime"                   {% as_string %}
    |  "unixtime_milliseconds_todatetime"                   {% as_string %}
    |  "unixtime_nanoseconds_todatetime"                    {% as_string %}
    |  "unixtime_seconds_todatetime"                        {% as_string %}
function_args       
    -> function_arg __                                      {% as_array(0) %}
    |  function_arg __ "," __ function_args                 {% merge(0,4)  %}
function_arg   
    -> any_type                                             {% pick(0) %}
    |  %identifier                                          {% d => { return { type: "identifier", value: d[0].value } } %}
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
# Command : Project ReOrder
command_project_reorder
    -> "project" %dash "reorder" _ function_assignments     {% d => d[4] %}
# Command : Project Away
command_project_away
    -> "project" %dash "away" _ ref_types                   {% d => d[4] %}
# Command : Scope 
command_scope
    -> "scope" _ ref_type                                   {% d => d[2] %}
# Command : Where 
command_where
    -> "where" _ expression_args                            {% d => d[2] %}
# Command : Distinct 
command_distinct
    -> "distinct" __ ref_type:*                             {% d => d[2] ? d[2][0] : undefined %}
# Command : mv-exapand 
command_mv_expand
    -> "mv" %dash "expand" _ ref_type                       {% d => d[4] %}
    |  "mv" %dash "expand" _ str:* "=":* ref_type           {% d => ({ alias: d[4][0], ...d[6] })%}
# Command : Parse json
command_parse_json
    -> "parse" %dash "json" __ parse_args:*                 {% d => d[4] %}
# Command : Parse csv
command_parse_csv
    -> "parse" %dash "csv" __ parse_args:*                  {% d => d[4] %}
# Command : Parse xml
command_parse_xml
    -> "parse" %dash "xml" __ parse_args:*                  {% d => d[4] %}
# Command : Parse yaml
command_parse_yaml
    -> "parse" %dash "yaml" __ parse_args:*                  {% d => d[4] %}
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
     ->  str:*  "=":* summarize_function                             {% d => {
                return d[2].condition ? 
                { operator: d[2].operator, alias: d[0][0], args: d[2].args, ref: d[2].ref, condition: d[2].condition }:
                { operator: d[2].operator, alias: d[0][0], args: d[2].args }
            }%}
summarize_function
     ->  function_name "(" summarize_args:* ")"                                 {% d => ({ operator: d[0], args: d[2][0]||[] })%}
     |   conditional_function_name "(" ref_type __ "," __ expression_args __ ")"   {% d => ({ operator: d[0], ref: d[2], condition: d[6] })%}
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
any_type
    -> num_type                                             {% pick(0) %}
    |  str_type                                             {% pick(0) %}
    |  ref_type                                             {% pick(0) %}
line_comment    
    -> %comment                                             {% d => ({ type: "comment", value : d[0]?.value || '' }) %}
str             -> %string                                  {% as_string %}                 # string
number          -> %number                                  {% as_number %}                 # number
nlo             -> %nl:*                                    {% dispose %}                   # optional newline
pipeo           -> %pipe:*                                  {% dispose %}                   # optional pipe
__              -> %ws:*                                    {% dispose %}                   # optional whitespace
_               -> %ws:+                                    {% dispose %}                   # mandatory whitespace
#endregion