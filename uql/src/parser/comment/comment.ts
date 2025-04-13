import { Command, CommandResult } from "../../types";

export const comment = (pv: CommandResult, cv: Extract<Command, { type: "comment" }>): CommandResult => {
  return { ...pv };
};
