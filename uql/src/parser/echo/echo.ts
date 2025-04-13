import { Command, CommandResult } from "../../types";

export const echo = (pv: CommandResult, cv: Extract<Command, { type: "echo" }>): CommandResult => {
  const output = cv.value;
  return { ...pv, output };
};
