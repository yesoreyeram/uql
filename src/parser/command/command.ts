import { isArray } from "lodash";
import { get_value } from "./../utils";
import { Command, CommandResult } from "../../types";

export const command = (pv: CommandResult, cv: Extract<Command, { type: "command" }>): CommandResult => {
  let output = pv.output;
  switch (cv.value.operator) {
    case "count":
      if (isArray(output)) {
        output = get_value("count", output);
      }
      break;
    case "sum":
      if (isArray(output)) {
        output = get_value("sum", output);
      }
      break;
    case "diff":
      if (isArray(output) && output.length === 2) {
        output = get_value("diff", output);
      }
      break;
    case "mul":
      if (isArray(output) && output.length === 2) {
        output = get_value("mul", output);
      }
      break;
    case "min":
      if (isArray(output)) {
        output = get_value("min", output);
      }
      break;
    case "max":
      if (isArray(output)) {
        output = get_value("max", output);
      }
      break;
    case "mean":
      if (isArray(output)) {
        output = get_value("mean", output);
      }
      break;
    case "first":
      if (isArray(output)) {
        output = get_value("first", output);
      }
      break;
    case "last":
      if (isArray(output)) {
        output = get_value("last", output);
      }
      break;
    case "latest":
      if (isArray(output)) {
        output = get_value("latest", output);
      }
      break;
    case "strcat":
      if (isArray(output)) {
        output = get_value("strcat", output);
      }
      break;
    case "dcount":
      if (isArray(output)) {
        output = get_value("dcount", output);
      }
      break;
    case "distinct":
      if (isArray(output)) {
        output = get_value("distinct", output);
      }
      break;
    case "random":
      output = get_value("random", []);
      break;
    case "toupper":
      if (typeof output === "string") {
        output = get_value("toupper", [output]);
      }
      break;
    case "tolower":
      if (typeof output === "string") {
        output = get_value("tolower", [output]);
      }
      break;
    case "strlen":
      if (typeof output === "string") {
        output = get_value("strlen", [output]);
      }
      break;
    case "trim":
      if (typeof output === "string") {
        output = get_value("trim", [output]);
      }
      break;
    case "trim_start":
      if (typeof output === "string") {
        output = get_value("trim_start", [output]);
      }
      break;
    case "trim_end":
      if (typeof output === "string") {
        output = get_value("trim_end", [output]);
      }
      break;
    default:
      break;
  }
  return { ...pv, output };
};
