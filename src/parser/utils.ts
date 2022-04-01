import { get, set, sum, min, max, mean, uniq, isArray, random, first, last, forEach } from "lodash";
import { type_function, type_summarize_assignment, type_parse_arg, FunctionName } from "../types";
import { Options as csv_parser_Options } from "csv-parse/lib";
import { X2jOptionsOptional } from "fast-xml-parser";
import * as dayjs from "dayjs";

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
      case "first":
        val = args.length > 0 ? first(pi.map((p) => get(p, args[0].value))) : first(pi);
        set(o, statName, val);
        break;
      case "last":
      case "latest":
        val = args.length > 0 ? last(pi.map((p) => get(p, args[0].value))) : last(pi);
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
      if (typeof args === "object" && isArray(args)) {
        return args.reduce((a, b) => a * b, 1);
      }
      return args[0] * args[1];
    case "div":
      return args[0] / args[1];
    case "min":
      return min(args);
    case "max":
      return max(args);
    case "mean":
      return mean(args);
    case "first":
      return first(args);
    case "last":
    case "latest":
      return last(args);
    case "parse_url":
      if (typeof args[0] === "string") {
        const url = new URL(args[0]);
        if (args.length === 3 && (args[1] === "search" || args[1] === "query")) {
          return url.searchParams.get(args[2]) || "";
        } else if (args.length === 2) {
          return (url as any)[args[1]] || "";
        } else {
          return url.toString();
        }
      }
      return "";
    case "parse_urlquery":
      if (typeof args[0] === "string") {
        const searchParams = new URLSearchParams(args[0]);
        if (args.length == 2) {
          return searchParams.get(args[1]) || "";
        } else {
          let o: Record<string, string> = {};
          searchParams.forEach((value, key) => (o[key] = value));
          return o;
        }
      }
      return "";
    case "toint":
    case "tolong":
    case "todouble":
    case "tofloat":
    case "tonumber":
      return +(args[0] || "");
    case "tobool":
      if (typeof args[0] === "string") return args[0].toLowerCase() === "true" ? true : false;
      else if (typeof args[0] === "number") return args[0] <= 0 ? true : false;
      else return args[0];
    case "tostring":
      return args[0] + "";
    case "todatetime":
      if (typeof args[0] === "number") return new Date(args[0] + "");
      return new Date(args[0] || "");
    case "tounixtime":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return args[0].getTime();
      return args[0];
    case "format_datetime":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function" && typeof args[1] === "string") return dayjs(args[0]).format(args[1]);
      return args[0];
    case "add_datetime":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function" && typeof args[1] === "string") {
        let numberParts = args[1].match(/[\-\d]+/g) || [];
        let textParts = args[1].match(/[A-Za-z]+/g) || [];
        if (numberParts.length > 0 && textParts.length > 0) {
          let o = dayjs(args[0]).add(numberParts[0] === "-" ? -1 : +numberParts[0], textParts[0]);
          return o.toDate();
        }
      }
      return args[0];
    case "startofminute":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("m").toDate();
      return args[0];
    case "startofhour":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("h").toDate();
      return args[0];
    case "startofday":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("d").toDate();
      return args[0];
    case "startofweek":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("w").toDate();
      return args[0];
    case "startofmonth":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("M").toDate();
      return args[0];
    case "startofyear":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("y").toDate();
      return args[0];
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
      return random(...args);
    case "dcount":
    case "distinct":
    default:
      throw "not implemented";
  }
};

export const get_extended_object = (o: object, assignment: type_function): object => {
  let args = assignment.args.map((arg) => {
    if (arg.type === "ref") return get(o, arg.value);
    else if (arg.type === "string" || arg.type === "identifier") {
      if (arg.value.toLowerCase() === "true") return true;
      else if (arg.value.toLowerCase() === "false") return false;
      return arg.value;
    } else if (arg.type === "number") return +arg.value;
    else return arg;
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
