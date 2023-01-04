import { Command, CommandResult } from "../../types";
import { get_single_value } from "../utils";

export const scope = (pv: CommandResult, cv: Extract<Command, { type: "scope" }>): CommandResult => {
  const output = get_single_value(pv.output, cv.value.value);
  return { ...pv, output };
};
