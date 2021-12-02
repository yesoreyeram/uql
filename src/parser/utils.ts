import { get, set, sum, min, max, mean, uniq } from "lodash";
import { type_function_assignment, type_summarize_assignment, FunctionName } from "../types";

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
    case "random":
    case "dcount":
    case "distinct":
    default:
      throw "not implemented";
  }
};

export const get_extended_object = (o: object, assignment: type_function_assignment): object => {
  let args = assignment.args.map((arg) => {
    if (arg.type === "ref") return get(o, arg.value);
    else if (arg.type === "string") return arg.value;
    else if (arg.type === "number") return +arg.value;
  });
  let value = get_value(assignment.operator, args);
  set(o, assignment.alias || assignment.operator, value);
  return o;
};
