import { Command, CommandResult } from "../../types";

export const hello = (pv: CommandResult, cv: Extract<Command, { type: "hello" }>): CommandResult => {
  const output = "hello";
  return { ...pv, output };
};
