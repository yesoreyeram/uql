import { get, set, isArray } from "lodash";
import { get_value } from "./../utils";
import { Command, CommandResult, type_function } from "../../types";

export const extend = (pv: CommandResult, cv: Extract<Command, { type: "extend" }>): CommandResult => {
  let output = pv.output;
  if (isArray(output)) {
    output = output.map((o) => {
      cv.value.forEach((ci) => {
        if (ci.type === "function") {
          o = get_extended_object(o, ci, output);
        } else if (ci.type === "ref") {
          set(o, ci.alias || ci.value, get(o, ci.value));
        }
      });
      return o;
    });
  }
  return { ...pv, output };
};

const get_extended_object = (o: object, assignment: type_function, previous_value?: any): object => {
  let args = assignment.args.map((arg) => {
    if (arg.type === "ref") return get(o, arg.value);
    else if (arg.type === "string" || arg.type === "identifier") {
      if (arg.value.toLowerCase() === "true") return true;
      else if (arg.value.toLowerCase() === "false") return false;
      return arg.value;
    } else if (arg.type === "number") return +arg.value;
    else return arg;
  });
  let value = get_value(assignment.operator, args, previous_value);
  set(o, assignment.alias || assignment.operator, value);
  return o;
};
