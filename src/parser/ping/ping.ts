import { Command, CommandResult } from "../../types";

export const ping = (pv: CommandResult, cv: Extract<Command, { type: "ping" }>): CommandResult => {
  const output = cv.value;
  return { ...pv, output };
};
