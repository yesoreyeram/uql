import { get, set, sum, min, max, mean, uniq } from "lodash";
import { type_function, type_summarize_assignment, type_parse_arg, FunctionName } from "../types";
import { Options as csv_parser_Options } from "csv-parse/lib";
import { X2jOptionsOptional } from "fast-xml-parser";

export const summarize = (o: object, metrics: type_summarize_assignment[], pi: unknown[]): object => {
  metrics.forEach((i) => {
    let args = i.args || [];
    let statName = args.length > 0 ? `${args[0].value} (${i.operator})` : i.operator;
    statName = i.alias ? i.alias : statName;
    let val;
    switch (i.operator) {
      case "count":
        set(o, statName, pi.length);
        break;
      case "dcount":
        val = args.length > 0 ? pi.map((p) => get(p, args[0].value)) : [pi];
        set(o, statName, uniq(val).length);
        break;
      case "mean":
        val = args.length > 0 ? mean(pi.map((p) => get(p, args[0].value))) : mean(pi);
        set(o, statName, val);
        break;
      case "sum":
        val = args.length > 0 ? sum(pi.map((p) => get(p, args[0].value))) : sum(pi);
        set(o, statName, val);
        break;
      case "min":
        val = args.length > 0 ? min(pi.map((p) => get(p, args[0].value))) : min(pi);
        set(o, statName, val);
        break;
      case "max":
        val = args.length > 0 ? max(pi.map((p) => get(p, args[0].value))) : max(pi);
        set(o, statName, val);
        break;
      default:
        break;
    }
  });
  return o;
};

export const get_value = (operator: FunctionName, args: any[]): unknown => {
  switch (operator) {
    case "tolower":
      return (args[0] + "").toLowerCase();
    case "toupper":
      return (args[0] + "").toUpperCase();
    case "strlen":
      return (args[0] + "").length;
    case "trim":
      return (args[0] + "").trim();
    case "trim_start":
      return (args[0] + "").trimStart();
    case "trim_end":
      return (args[0] + "").trimEnd();
    case "strcat":
      return args.join("");
    case "sum":
      return sum(args);
    case "diff":
      return args[0] - args[1];
    case "mul":
      return args[0] * args[1];
    case "min":
      return min(args);
    case "max":
      return max(args);
    case "mean":
      return mean(args);
    case "toint":
    case "tolong":
    case "todouble":
    case "tofloat":
    case "tonumber":
      return +(args[0] || "");
    case "tobool":
      //TODO: write a logic to convert to bool
      return args[0];
    case "tostring":
      return args[0] + "";
    case "todatetime":
      if (typeof args[0] === "number") return new Date(args[0] + "");
      return new Date(args[0] || "");
    case "unixtime_seconds_todatetime":
      if (typeof args[0] === "string") return new Date(+args[0] * 1000);
      return new Date(args[0] * 1000);
    case "unixtime_milliseconds_todatetime":
      if (typeof args[0] === "string") return new Date(+args[0]);
      return new Date(args[0]);
    case "unixtime_microseconds_todatetime":
      if (typeof args[0] === "string") return new Date(+args[0] / 1000);
      return new Date(args[0] / 1000);
    case "unixtime_nanoseconds_todatetime":
      if (typeof args[0] === "string") return new Date(+args[0] / 1000 / 1000);
      return new Date(args[0] / 1000 / 1000);
    case "random":
    case "dcount":
    case "distinct":
    default:
      throw "not implemented";
  }
};

export const get_extended_object = (o: object, assignment: type_function): object => {
  let args = assignment.args.map((arg) => {
    if (arg.type === "ref") return get(o, arg.value);
    else if (arg.type === "string") return arg.value;
    else if (arg.type === "number") return +arg.value;
  });
  let value = get_value(assignment.operator, args);
  set(o, assignment.alias || assignment.operator, value);
  return o;
};

export const get_parse_csv_options = (args: type_parse_arg[][]): csv_parser_Options => {
  let options: csv_parser_Options = { columns: true, delimiter: ",", skipEmptyLines: false, skipLinesWithError: false, relaxColumnCount: false };
  if (args[0] && args[0].length > 0) {
    args[0].forEach((arg) => {
      switch (arg.identifier) {
        case "autoParse":
        case "autoParseDate":
        case "bom":
        case "cast":
        case "castDate":
        case "relaxColumnCount":
        case "skipEmptyLines":
        case "skipLinesWithEmptyValues":
        case "trim":
          options[arg.identifier] = arg.value.toLowerCase() === "true";
          break;
        case "columns":
          options[arg.identifier] = arg.value ? arg.value.split(",") : true;
          break;
        default:
          // @ts-ignore
          options[arg.identifier] = arg.value.toLowerCase() === "true" ? true : arg.value.toLowerCase() === "false" ? false : arg.value;
          break;
      }
    });
  }
  return options;
};
export const get_parse_xml_options = (args: type_parse_arg[][]): X2jOptionsOptional => {
  let options: X2jOptionsOptional = { ignoreAttributes: false, allowBooleanAttributes: true, commentPropName: "#comments" };
  if (args[0] && args[0].length > 0) {
    args[0].forEach((arg) => {
      switch (arg.identifier) {
        case "ignoreAttributes":
        case "allowBooleanAttributes":
        case "alwaysCreateTextNode":
        case "preserveOrder":
        case "parseTagValue":
        case "parseAttributeValue":
        case "trimValues":
          options[arg.identifier] = arg.value.toLowerCase() === "true";
          break;
        case "commentPropName":
          options["commentPropName"] = arg.value || "#comment";
          break;
        case "attributeNamePrefix":
          options["attributeNamePrefix"] = arg.value || "@_";
          break;
        default:
          // @ts-ignore
          options[arg.identifier] = arg.value.toLowerCase() === "true" ? true : arg.value.toLowerCase() === "false" ? false : arg.value;
          break;
      }
    });
  }
  return options;
};
