// Generated automatically by nearley, version 2.20.1
// http://github.com/Hardmath123/nearley
// Bypasses TS6133. Allow declared but unused functions.
// @ts-ignore
function id(d: any[]): any { return d[0]; }
declare var comment: any;
declare var lparan: any;
declare var rparan: any;
declare var plus: any;
declare var dash: any;
declare var mul: any;
declare var divide: any;
declare var mod: any;
declare var gt: any;
declare var assignment: any;
declare var lt: any;
declare var eq: any;
declare var ne: any;
declare var tilde: any;
declare var exclaim: any;
declare var identifier: any;
declare var str: any;
declare var sq_string: any;
declare var string: any;
declare var number: any;
declare var nl: any;
declare var pipe: any;
declare var ws: any;


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


interface NearleyToken {
  value: any;
  [key: string]: any;
};

interface NearleyLexer {
  reset: (chunk: string, info: any) => void;
  next: () => NearleyToken | undefined;
  save: () => any;
  formatError: (token: never) => string;
  has: (tokenType: string) => boolean;
};

interface NearleyRule {
  name: string;
  symbols: NearleySymbol[];
  postprocess?: (d: any[], loc?: number, reject?: {}) => any;
};

type NearleySymbol = string | { literal: any } | { test: (token: any) => boolean };

interface Grammar {
  Lexer: NearleyLexer | undefined;
  ParserRules: NearleyRule[];
  ParserStart: string;
};

const grammar: Grammar = {
  Lexer: oqlLexer,
  ParserRules: [
    {"name": "input", "symbols": ["commands"], "postprocess": pick(0)},
    {"name": "commands", "symbols": ["command", "__"], "postprocess": as_array(0)},
    {"name": "commands$ebnf$1", "symbols": []},
    {"name": "commands$ebnf$1", "symbols": ["commands$ebnf$1", (oqlLexer.has("comment") ? {type: "comment"} : comment)], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "commands", "symbols": ["command", "__", "commands$ebnf$1", {"literal":"\r\n"}, "__", "pipeo", "__", "commands"], "postprocess": merge(0,7)},
    {"name": "commands$ebnf$2", "symbols": []},
    {"name": "commands$ebnf$2", "symbols": ["commands$ebnf$2", (oqlLexer.has("comment") ? {type: "comment"} : comment)], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "commands", "symbols": ["command", "__", "commands$ebnf$2", "nlo", "__", "pipeo", "__", "commands"], "postprocess": merge(0,7)},
    {"name": "command", "symbols": ["line_comment"], "postprocess": pick(0)},
    {"name": "command", "symbols": [{"literal":"hello"}], "postprocess": d => ({ type: "hello" })},
    {"name": "command", "symbols": [{"literal":"ping"}], "postprocess": d => ({ type: "ping", value: "pong" })},
    {"name": "command", "symbols": [{"literal":"echo"}, "__", "str"], "postprocess": d => ({ type: "echo", value: d[2] })},
    {"name": "command", "symbols": [{"literal":"count"}], "postprocess": d => ({ type: "count" })},
    {"name": "command", "symbols": [{"literal":"limit"}, "__", "number"], "postprocess": d => ({ type: "limit", value: d[2] })},
    {"name": "command", "symbols": ["function"], "postprocess": d => ({ type: "command" , value: d[0] })},
    {"name": "command", "symbols": ["command_orderby"], "postprocess": d => ({ type: "orderby", value: d[0] })},
    {"name": "command", "symbols": ["command_extend"], "postprocess": d => ({ type: "extend", value: d[0] })},
    {"name": "command", "symbols": ["command_project_away"], "postprocess": d => ({ type: "project-away", value: d[0] })},
    {"name": "command", "symbols": ["command_project_reorder"], "postprocess": d => ({ type: "project-reorder", value: d[0] })},
    {"name": "command", "symbols": ["command_project"], "postprocess": d => ({ type: "project", value: d[0] })},
    {"name": "command", "symbols": ["command_parse_json"], "postprocess": d => ({ type: "parse-json", args: d[0] })},
    {"name": "command", "symbols": ["command_parse_csv"], "postprocess": d => ({ type: "parse-csv", args: d[0] })},
    {"name": "command", "symbols": ["command_parse_xml"], "postprocess": d => ({ type: "parse-xml", args: d[0] })},
    {"name": "command", "symbols": ["command_parse_yaml"], "postprocess": d => ({ type: "parse-yaml", args: d[0] })},
    {"name": "command", "symbols": ["command_scope"], "postprocess": d => ({ type: "scope", value: d[0] })},
    {"name": "command", "symbols": ["command_where"], "postprocess": d => ({ type: "where", value: d[0] })},
    {"name": "command", "symbols": ["command_distinct"], "postprocess": d => ({ type: "distinct", value: d[0] })},
    {"name": "command", "symbols": ["command_mv_expand"], "postprocess": d => ({ type: "mv-expand", value: d[0] })},
    {"name": "command", "symbols": ["command_summarize"], "postprocess": d => ({ type: "summarize", value: d[0] })},
    {"name": "command", "symbols": ["command_range"], "postprocess": d => ({ type: "range", value: d[0] })},
    {"name": "command", "symbols": [{"literal":"jsonata"}, "__", "str"], "postprocess": d => ({ type: "jsonata", expression: d[2] })},
    {"name": "function_assignments", "symbols": ["function_assignment"], "postprocess": as_array(0)},
    {"name": "function_assignments", "symbols": ["function_assignment", "__", {"literal":","}, "__", "function_assignments"], "postprocess": merge(0,4)},
    {"name": "function_assignment$ebnf$1", "symbols": []},
    {"name": "function_assignment$ebnf$1", "symbols": ["function_assignment$ebnf$1", "str"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "function_assignment$ebnf$2", "symbols": []},
    {"name": "function_assignment$ebnf$2", "symbols": ["function_assignment$ebnf$2", {"literal":"="}], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "function_assignment", "symbols": ["function_assignment$ebnf$1", "function_assignment$ebnf$2", "expression"], "postprocess": d => ({ alias: d[0][0], ...d[2] })},
    {"name": "function_assignment$ebnf$3", "symbols": []},
    {"name": "function_assignment$ebnf$3", "symbols": ["function_assignment$ebnf$3", "str"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "function_assignment$ebnf$4", "symbols": []},
    {"name": "function_assignment$ebnf$4", "symbols": ["function_assignment$ebnf$4", {"literal":"="}], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "function_assignment", "symbols": ["function_assignment$ebnf$3", "function_assignment$ebnf$4", "function"], "postprocess": d => ({ alias: d[0][0], ...d[2] })},
    {"name": "function_assignment$ebnf$5", "symbols": []},
    {"name": "function_assignment$ebnf$5", "symbols": ["function_assignment$ebnf$5", "str"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "function_assignment$ebnf$6", "symbols": []},
    {"name": "function_assignment$ebnf$6", "symbols": ["function_assignment$ebnf$6", {"literal":"="}], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "function_assignment", "symbols": ["function_assignment$ebnf$5", "function_assignment$ebnf$6", "ref_type"], "postprocess": d => ({ alias: d[0][0], ...d[2] })},
    {"name": "function_assignment", "symbols": ["ref_type"], "postprocess": d => d[0]},
    {"name": "expression$ebnf$1", "symbols": []},
    {"name": "expression$ebnf$1", "symbols": ["expression$ebnf$1", "expression_args"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "expression", "symbols": [(oqlLexer.has("lparan") ? {type: "lparan"} : lparan), "__", "expression$ebnf$1", (oqlLexer.has("rparan") ? {type: "rparan"} : rparan)], "postprocess": d => ({ type: "expression", args: d[2][0]||[] })},
    {"name": "expression_args", "symbols": ["expression_arg", "__"], "postprocess": as_array(0)},
    {"name": "expression_args", "symbols": ["expression_arg", "__", "expression_args"], "postprocess": merge(0,2)},
    {"name": "expression_arg", "symbols": ["any_type"], "postprocess": d => d[0]},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("lparan") ? {type: "lparan"} : lparan), "__", "function_args", "__", (oqlLexer.has("rparan") ? {type: "rparan"} : rparan)], "postprocess": d => ({ type : "value_array", value: d[2] })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("plus") ? {type: "plus"} : plus)], "postprocess": d => ({ type: "operation", value: "+" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("dash") ? {type: "dash"} : dash)], "postprocess": d => ({ type: "operation", value: "-" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("mul") ? {type: "mul"} : mul)], "postprocess": d => ({ type: "operation", value: "*" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("divide") ? {type: "divide"} : divide)], "postprocess": d => ({ type: "operation", value: "/" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("mod") ? {type: "mod"} : mod)], "postprocess": d => ({ type: "operation", value: "%" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("gt") ? {type: "gt"} : gt), (oqlLexer.has("assignment") ? {type: "assignment"} : assignment)], "postprocess": d => ({ type: "operation", value: ">=" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("gt") ? {type: "gt"} : gt)], "postprocess": d => ({ type: "operation", value: ">" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("lt") ? {type: "lt"} : lt), (oqlLexer.has("assignment") ? {type: "assignment"} : assignment)], "postprocess": d => ({ type: "operation", value: "<=" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("lt") ? {type: "lt"} : lt)], "postprocess": d => ({ type: "operation", value: "<" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("eq") ? {type: "eq"} : eq)], "postprocess": d => ({ type: "operation", value: "==" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("ne") ? {type: "ne"} : ne)], "postprocess": d => ({ type: "operation", value: "!=" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("assignment") ? {type: "assignment"} : assignment), (oqlLexer.has("tilde") ? {type: "tilde"} : tilde)], "postprocess": d => ({ type: "operation", value: "=~" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("exclaim") ? {type: "exclaim"} : exclaim), (oqlLexer.has("tilde") ? {type: "tilde"} : tilde)], "postprocess": d => ({ type: "operation", value: "!~" })},
    {"name": "expression_arg", "symbols": [{"literal":"between"}], "postprocess": d => ({ type: "operation", value: "between" })},
    {"name": "expression_arg", "symbols": [{"literal":"inside"}], "postprocess": d => ({ type: "operation", value: "inside" })},
    {"name": "expression_arg", "symbols": [{"literal":"outside"}], "postprocess": d => ({ type: "operation", value: "outside" })},
    {"name": "expression_arg", "symbols": [{"literal":"matches"}, "_", {"literal":"regex"}], "postprocess": d => ({ type: "operation", value: "matches regex" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("exclaim") ? {type: "exclaim"} : exclaim), {"literal":"matches"}, "_", {"literal":"regex"}], "postprocess": d => ({ type: "operation", value: "!matches regex" })},
    {"name": "expression_arg", "symbols": [{"literal":"not"}, "_", {"literal":"matches"}, "_", {"literal":"regex"}], "postprocess": d => ({ type: "operation", value: "!matches regex" })},
    {"name": "expression_arg", "symbols": [{"literal":"in"}], "postprocess": d => ({ type: "operation", value: "in" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("exclaim") ? {type: "exclaim"} : exclaim), {"literal":"in"}], "postprocess": d => ({ type: "operation", value: "!in" })},
    {"name": "expression_arg", "symbols": [{"literal":"not"}, "_", {"literal":"in"}], "postprocess": d => ({ type: "operation", value: "!in" })},
    {"name": "expression_arg", "symbols": [{"literal":"in"}, (oqlLexer.has("tilde") ? {type: "tilde"} : tilde)], "postprocess": d => ({ type: "operation", value: "in~" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("exclaim") ? {type: "exclaim"} : exclaim), {"literal":"in"}, (oqlLexer.has("tilde") ? {type: "tilde"} : tilde)], "postprocess": d => ({ type: "operation", value: "!in~" })},
    {"name": "expression_arg", "symbols": [{"literal":"not"}, "_", {"literal":"in"}, (oqlLexer.has("tilde") ? {type: "tilde"} : tilde)], "postprocess": d => ({ type: "operation", value: "!in~" })},
    {"name": "expression_arg", "symbols": [{"literal":"contains"}], "postprocess": d => ({ type: "operation", value: "contains" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("exclaim") ? {type: "exclaim"} : exclaim), {"literal":"contains"}], "postprocess": d => ({ type: "operation", value: "!contains" })},
    {"name": "expression_arg", "symbols": [{"literal":"not"}, "_", {"literal":"contains"}], "postprocess": d => ({ type: "operation", value: "!contains" })},
    {"name": "expression_arg", "symbols": [{"literal":"contains_cs"}], "postprocess": d => ({ type: "operation", value: "contains_cs" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("exclaim") ? {type: "exclaim"} : exclaim), {"literal":"contains_cs"}], "postprocess": d => ({ type: "operation", value: "!contains_cs" })},
    {"name": "expression_arg", "symbols": [{"literal":"not"}, "_", {"literal":"contains_cs"}], "postprocess": d => ({ type: "operation", value: "!contains_cs" })},
    {"name": "expression_arg", "symbols": [{"literal":"startswith"}], "postprocess": d => ({ type: "operation", value: "startswith" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("exclaim") ? {type: "exclaim"} : exclaim), {"literal":"startswith"}], "postprocess": d => ({ type: "operation", value: "!startswith" })},
    {"name": "expression_arg", "symbols": [{"literal":"not"}, "_", {"literal":"startswith"}], "postprocess": d => ({ type: "operation", value: "!startswith" })},
    {"name": "expression_arg", "symbols": [{"literal":"startswith_cs"}], "postprocess": d => ({ type: "operation", value: "startswith_cs" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("exclaim") ? {type: "exclaim"} : exclaim), {"literal":"startswith_cs"}], "postprocess": d => ({ type: "operation", value: "!startswith_cs" })},
    {"name": "expression_arg", "symbols": [{"literal":"not"}, "_", {"literal":"startswith_cs"}], "postprocess": d => ({ type: "operation", value: "!startswith_cs" })},
    {"name": "expression_arg", "symbols": [{"literal":"endswith"}], "postprocess": d => ({ type: "operation", value: "endswith" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("exclaim") ? {type: "exclaim"} : exclaim), {"literal":"endswith"}], "postprocess": d => ({ type: "operation", value: "!endswith" })},
    {"name": "expression_arg", "symbols": [{"literal":"not"}, "_", {"literal":"endswith"}], "postprocess": d => ({ type: "operation", value: "!endswith" })},
    {"name": "expression_arg", "symbols": [{"literal":"endswith_cs"}], "postprocess": d => ({ type: "operation", value: "endswith_cs" })},
    {"name": "expression_arg", "symbols": [(oqlLexer.has("exclaim") ? {type: "exclaim"} : exclaim), {"literal":"endswith_cs"}], "postprocess": d => ({ type: "operation", value: "!endswith_cs" })},
    {"name": "expression_arg", "symbols": [{"literal":"not"}, "_", {"literal":"endswith_cs"}], "postprocess": d => ({ type: "operation", value: "!endswith_cs" })},
    {"name": "expression_arg", "symbols": ["function"], "postprocess": d => ({ type: "function", value: d[0] })},
    {"name": "function$ebnf$1", "symbols": []},
    {"name": "function$ebnf$1", "symbols": ["function$ebnf$1", "function_args"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "function", "symbols": ["function_name", {"literal":"("}, "__", "function$ebnf$1", {"literal":")"}], "postprocess": d => ({ type: "function", operator: d[0], args: d[3][0]||[] })},
    {"name": "conditional_function_name", "symbols": [{"literal":"countif"}], "postprocess": as_string},
    {"name": "conditional_function_name", "symbols": [{"literal":"sumif"}], "postprocess": as_string},
    {"name": "conditional_function_name", "symbols": [{"literal":"minif"}], "postprocess": as_string},
    {"name": "conditional_function_name", "symbols": [{"literal":"maxif"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"add_datetime"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"array_from_entries"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"array_to_map"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"atob"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"bag_pack"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"btoa"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"ceil"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"cos"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"count"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"dcount"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"diff"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"distinct"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"div"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"extract"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"first"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"floor"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"format_datetime"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"kv"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"last"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"latest"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"log"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"log10"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"log2"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"max"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"mean"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"min"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"mul"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"pack"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"parse_url"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"parse_urlquery"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"percentage"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"pow"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"random"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"replace_string"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"reverse"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"round"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"sign"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"sin"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"split"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"startofday"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"startofhour"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"startofminute"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"startofmonth"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"startofweek"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"startofyear"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"strcat"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"strlen"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"substring"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"sum"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"tan"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"tobool"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"todatetime"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"todouble"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"tofloat"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"toint"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"tolong"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"tolower"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"tonumber"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"tostring"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"tounixtime"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"toupper"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"trim_end"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"trim_start"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"trim"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"unixtime_microseconds_todatetime"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"unixtime_milliseconds_todatetime"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"unixtime_nanoseconds_todatetime"}], "postprocess": as_string},
    {"name": "function_name", "symbols": [{"literal":"unixtime_seconds_todatetime"}], "postprocess": as_string},
    {"name": "function_args", "symbols": ["function_arg", "__"], "postprocess": as_array(0)},
    {"name": "function_args", "symbols": ["function_arg", "__", {"literal":","}, "__", "function_args"], "postprocess": merge(0,4)},
    {"name": "function_arg", "symbols": ["any_type"], "postprocess": pick(0)},
    {"name": "function_arg", "symbols": [(oqlLexer.has("identifier") ? {type: "identifier"} : identifier)], "postprocess": d => { return { type: "identifier", value: d[0].value } }},
    {"name": "command_orderby", "symbols": [{"literal":"order"}, "_", {"literal":"by"}, "_", "orderby_args"], "postprocess": pick(4)},
    {"name": "orderby_args", "symbols": ["orderby_arg", "__"], "postprocess": as_array(0)},
    {"name": "orderby_args", "symbols": ["orderby_arg", "__", {"literal":","}, "__", "orderby_args"], "postprocess": merge(0,4)},
    {"name": "orderby_arg", "symbols": ["str", "__", {"literal":"asc"}], "postprocess": d => ({ field: d[0] , direction: "asc"})},
    {"name": "orderby_arg", "symbols": ["str", "__", {"literal":"desc"}], "postprocess": d => ({ field: d[0] , direction: "desc"})},
    {"name": "command_extend", "symbols": [{"literal":"extend"}, "_", "function_assignments"], "postprocess": d => d[2]},
    {"name": "command_project", "symbols": [{"literal":"project"}, "_", "function_assignments"], "postprocess": d => d[2]},
    {"name": "command_project_reorder", "symbols": [{"literal":"project"}, (oqlLexer.has("dash") ? {type: "dash"} : dash), {"literal":"reorder"}, "_", "function_assignments"], "postprocess": d => d[4]},
    {"name": "command_project_away", "symbols": [{"literal":"project"}, (oqlLexer.has("dash") ? {type: "dash"} : dash), {"literal":"away"}, "_", "ref_types"], "postprocess": d => d[4]},
    {"name": "command_scope", "symbols": [{"literal":"scope"}, "_", "ref_type"], "postprocess": d => d[2]},
    {"name": "command_where", "symbols": [{"literal":"where"}, "_", "expression_args"], "postprocess": d => d[2]},
    {"name": "command_distinct$ebnf$1", "symbols": []},
    {"name": "command_distinct$ebnf$1", "symbols": ["command_distinct$ebnf$1", "ref_type"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "command_distinct", "symbols": [{"literal":"distinct"}, "__", "command_distinct$ebnf$1"], "postprocess": d => d[2] ? d[2][0] : undefined},
    {"name": "command_mv_expand", "symbols": [{"literal":"mv"}, (oqlLexer.has("dash") ? {type: "dash"} : dash), {"literal":"expand"}, "_", "ref_type"], "postprocess": d => d[4]},
    {"name": "command_mv_expand$ebnf$1", "symbols": []},
    {"name": "command_mv_expand$ebnf$1", "symbols": ["command_mv_expand$ebnf$1", "str"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "command_mv_expand$ebnf$2", "symbols": []},
    {"name": "command_mv_expand$ebnf$2", "symbols": ["command_mv_expand$ebnf$2", {"literal":"="}], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "command_mv_expand", "symbols": [{"literal":"mv"}, (oqlLexer.has("dash") ? {type: "dash"} : dash), {"literal":"expand"}, "_", "command_mv_expand$ebnf$1", "command_mv_expand$ebnf$2", "ref_type"], "postprocess": d => ({ alias: d[4][0], ...d[6] })},
    {"name": "command_parse_json$ebnf$1", "symbols": []},
    {"name": "command_parse_json$ebnf$1", "symbols": ["command_parse_json$ebnf$1", "parse_args"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "command_parse_json", "symbols": [{"literal":"parse"}, (oqlLexer.has("dash") ? {type: "dash"} : dash), {"literal":"json"}, "__", "command_parse_json$ebnf$1"], "postprocess": d => d[4]},
    {"name": "command_parse_csv$ebnf$1", "symbols": []},
    {"name": "command_parse_csv$ebnf$1", "symbols": ["command_parse_csv$ebnf$1", "parse_args"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "command_parse_csv", "symbols": [{"literal":"parse"}, (oqlLexer.has("dash") ? {type: "dash"} : dash), {"literal":"csv"}, "__", "command_parse_csv$ebnf$1"], "postprocess": d => d[4]},
    {"name": "command_parse_xml$ebnf$1", "symbols": []},
    {"name": "command_parse_xml$ebnf$1", "symbols": ["command_parse_xml$ebnf$1", "parse_args"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "command_parse_xml", "symbols": [{"literal":"parse"}, (oqlLexer.has("dash") ? {type: "dash"} : dash), {"literal":"xml"}, "__", "command_parse_xml$ebnf$1"], "postprocess": d => d[4]},
    {"name": "command_parse_yaml$ebnf$1", "symbols": []},
    {"name": "command_parse_yaml$ebnf$1", "symbols": ["command_parse_yaml$ebnf$1", "parse_args"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "command_parse_yaml", "symbols": [{"literal":"parse"}, (oqlLexer.has("dash") ? {type: "dash"} : dash), {"literal":"yaml"}, "__", "command_parse_yaml$ebnf$1"], "postprocess": d => d[4]},
    {"name": "parse_args", "symbols": ["parse_arg"], "postprocess": as_array(0)},
    {"name": "parse_args", "symbols": ["parse_arg", "__", "parse_args"], "postprocess": merge(0,2)},
    {"name": "parse_arg", "symbols": [(oqlLexer.has("dash") ? {type: "dash"} : dash), (oqlLexer.has("dash") ? {type: "dash"} : dash), (oqlLexer.has("identifier") ? {type: "identifier"} : identifier), "__", "str"], "postprocess": d => ({ identifier: d[2].value, value: d[4] })},
    {"name": "parse_arg", "symbols": [(oqlLexer.has("dash") ? {type: "dash"} : dash), (oqlLexer.has("dash") ? {type: "dash"} : dash), (oqlLexer.has("identifier") ? {type: "identifier"} : identifier), "__", "str_type"], "postprocess": d => ({ identifier: d[2].value, value: d[4].value })},
    {"name": "parse_arg", "symbols": [(oqlLexer.has("dash") ? {type: "dash"} : dash), (oqlLexer.has("dash") ? {type: "dash"} : dash), (oqlLexer.has("identifier") ? {type: "identifier"} : identifier), "__", (oqlLexer.has("identifier") ? {type: "identifier"} : identifier)], "postprocess": d => ({ identifier: d[2].value, value: d[4].value })},
    {"name": "command_summarize", "symbols": ["summarize_item"], "postprocess": pick(0)},
    {"name": "command_summarize$ebnf$1", "symbols": []},
    {"name": "command_summarize$ebnf$1", "symbols": ["command_summarize$ebnf$1", "summarize_args"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "command_summarize", "symbols": ["summarize_item", "_", {"literal":"by"}, "_", "command_summarize$ebnf$1"], "postprocess": d => ({ ...d[0], by: d[4][0] })},
    {"name": "summarize_item", "symbols": [{"literal":"summarize"}, "_", "summarize_assignments"], "postprocess": d => ({ metrics: d[2], by :[] })},
    {"name": "summarize_assignments", "symbols": ["summarize_assignment"], "postprocess": d => [ d[0] ]},
    {"name": "summarize_assignments", "symbols": ["summarize_assignment", "__", {"literal":","}, "__", "summarize_assignments"], "postprocess": merge(0,4)},
    {"name": "summarize_assignment$ebnf$1", "symbols": []},
    {"name": "summarize_assignment$ebnf$1", "symbols": ["summarize_assignment$ebnf$1", "str"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "summarize_assignment$ebnf$2", "symbols": []},
    {"name": "summarize_assignment$ebnf$2", "symbols": ["summarize_assignment$ebnf$2", {"literal":"="}], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "summarize_assignment", "symbols": ["summarize_assignment$ebnf$1", "summarize_assignment$ebnf$2", "summarize_function"], "postprocess":  d => {
            return d[2].condition ? 
            { operator: d[2].operator, alias: d[0][0], args: d[2].args, ref: d[2].ref, condition: d[2].condition }:
            { operator: d[2].operator, alias: d[0][0], args: d[2].args }
        }},
    {"name": "summarize_function$ebnf$1", "symbols": []},
    {"name": "summarize_function$ebnf$1", "symbols": ["summarize_function$ebnf$1", "summarize_args"], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "summarize_function", "symbols": ["function_name", {"literal":"("}, "summarize_function$ebnf$1", {"literal":")"}], "postprocess": d => ({ operator: d[0], args: d[2][0]||[] })},
    {"name": "summarize_function", "symbols": ["conditional_function_name", {"literal":"("}, "ref_type", "__", {"literal":","}, "__", "expression_args", "__", {"literal":")"}], "postprocess": d => ({ operator: d[0], ref: d[2], condition: d[6] })},
    {"name": "summarize_args", "symbols": ["summarize_arg", "__"], "postprocess": as_array(0)},
    {"name": "summarize_args", "symbols": ["summarize_arg", "__", {"literal":","}, "__", "summarize_args"], "postprocess": merge(0,4)},
    {"name": "summarize_arg", "symbols": ["ref_type"], "postprocess": pick(0)},
    {"name": "command_range", "symbols": ["range_item"], "postprocess": pick(0)},
    {"name": "command_range", "symbols": ["range_item", "_", {"literal":"step"}, "_", "number"], "postprocess": d =>  ({...d[0], step : d[4]})},
    {"name": "command_range", "symbols": ["range_item", "_", {"literal":"step"}, "_", "str"], "postprocess": d =>  ({...d[0], step : d[4]})},
    {"name": "range_item", "symbols": [{"literal":"range"}, "_", {"literal":"from"}, "_", "number", "_", {"literal":"to"}, "_", "number"], "postprocess": d => ({ start: d[4], end: d[8], step: 1 })},
    {"name": "range_item", "symbols": [{"literal":"range"}, "_", {"literal":"from"}, "_", "str", "_", {"literal":"to"}, "_", "str"], "postprocess": d => ({ start: d[4], end: d[8], step: "" })},
    {"name": "str_type", "symbols": [(oqlLexer.has("str") ? {type: "str"} : str), {"literal":"("}, "str", {"literal":")"}], "postprocess": d => ({ type: "string", value: d[2] })},
    {"name": "str_type", "symbols": [(oqlLexer.has("sq_string") ? {type: "sq_string"} : sq_string)], "postprocess": d => ({ type: "string", value:d[0].value})},
    {"name": "ref_type", "symbols": ["str"], "postprocess": d => ({ type: "ref", value: d[0] })},
    {"name": "ref_type", "symbols": [{"literal":"["}, "str", {"literal":"]"}], "postprocess": d => ({ type: "ref", value: d[1] })},
    {"name": "ref_type", "symbols": [{"literal":"["}, (oqlLexer.has("sq_string") ? {type: "sq_string"} : sq_string), {"literal":"]"}], "postprocess": d => ({ type: "ref", value: d[1] })},
    {"name": "num_type", "symbols": ["number"], "postprocess": d => ({ type: "number", value: d[0] })},
    {"name": "ref_types", "symbols": ["ref_type"], "postprocess": as_array(0)},
    {"name": "ref_types", "symbols": ["ref_type", "__", {"literal":","}, "__", "ref_types"], "postprocess": merge(0,4)},
    {"name": "any_type", "symbols": ["num_type"], "postprocess": pick(0)},
    {"name": "any_type", "symbols": ["str_type"], "postprocess": pick(0)},
    {"name": "any_type", "symbols": ["ref_type"], "postprocess": pick(0)},
    {"name": "line_comment", "symbols": [(oqlLexer.has("comment") ? {type: "comment"} : comment)], "postprocess": d => ({ type: "comment", value : d[0]?.value || '' })},
    {"name": "str", "symbols": [(oqlLexer.has("string") ? {type: "string"} : string)], "postprocess": as_string},
    {"name": "number", "symbols": [(oqlLexer.has("number") ? {type: "number"} : number)], "postprocess": as_number},
    {"name": "nlo$ebnf$1", "symbols": []},
    {"name": "nlo$ebnf$1", "symbols": ["nlo$ebnf$1", (oqlLexer.has("nl") ? {type: "nl"} : nl)], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "nlo", "symbols": ["nlo$ebnf$1"], "postprocess": dispose},
    {"name": "pipeo$ebnf$1", "symbols": []},
    {"name": "pipeo$ebnf$1", "symbols": ["pipeo$ebnf$1", (oqlLexer.has("pipe") ? {type: "pipe"} : pipe)], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "pipeo", "symbols": ["pipeo$ebnf$1"], "postprocess": dispose},
    {"name": "__$ebnf$1", "symbols": []},
    {"name": "__$ebnf$1", "symbols": ["__$ebnf$1", (oqlLexer.has("ws") ? {type: "ws"} : ws)], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "__", "symbols": ["__$ebnf$1"], "postprocess": dispose},
    {"name": "_$ebnf$1", "symbols": [(oqlLexer.has("ws") ? {type: "ws"} : ws)]},
    {"name": "_$ebnf$1", "symbols": ["_$ebnf$1", (oqlLexer.has("ws") ? {type: "ws"} : ws)], "postprocess": (d) => d[0].concat([d[1]])},
    {"name": "_", "symbols": ["_$ebnf$1"], "postprocess": dispose}
  ],
  ParserStart: "input",
};

export default grammar;
