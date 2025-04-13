import { isArray } from "lodash";
import { Command, CommandResult } from "../../types";

export const count = (pv: CommandResult, cv: Extract<Command, { type: "count" }>): CommandResult => {
  let output = pv.output;
  if (typeof output === "string") {
    output = output.length;
  } else if (typeof pv.output === "number") {
    output = output;
  } else if (isArray(output)) {
    output = output.length;
  }
  return { ...pv, output };
};
