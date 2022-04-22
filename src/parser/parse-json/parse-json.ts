import { Command, CommandResult } from "../../types";

export const parseJson = (pv: CommandResult, cv: Extract<Command, { type: "parse-json" }>): CommandResult => {
  let output = pv.output;
  if (typeof pv.output === "string") {
    output = JSON.parse(pv.output);
  }
  return { ...pv, output };
};
