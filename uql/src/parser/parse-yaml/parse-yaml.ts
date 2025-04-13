import { load as yaml_loader } from "js-yaml";
import { Command, CommandResult } from "../../types";

export const parseYaml = (pv: CommandResult, cv: Extract<Command, { type: "parse-yaml" }>): CommandResult => {
  let output = pv.output;
  if (typeof output === "string") {
    try {
      output = yaml_loader(output);
    } catch (ex) {
      throw ex;
    }
  }
  return { ...pv, output };
};
