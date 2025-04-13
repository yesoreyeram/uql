import { isArray } from "lodash";
import { Command, CommandResult } from "../../types";

export const limit = (pv: CommandResult, cv: Extract<Command, { type: "limit" }>): CommandResult => {
  let output = pv.output;
  if (typeof output === "string") {
    output = output.substr(0, cv.value);
  } else if (typeof output === "number") {
    output = output;
  } else if (isArray(output)) {
    output = output.slice(0, cv.value);
  }
  return { ...pv, output };
};
