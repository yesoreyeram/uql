import { get } from "lodash";
import { Command, CommandResult } from "../../types";

export const scope = (pv: CommandResult, cv: Extract<Command, { type: "scope" }>): CommandResult => {
  const output = get(pv.output, cv.value.value);
  return { ...pv, output };
};
