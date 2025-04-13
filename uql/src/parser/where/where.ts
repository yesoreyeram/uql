import { filterData } from "./../utils";
import { Command, CommandResult } from "../../types";

export const where = (pv: CommandResult, cv: Extract<Command, { type: "where" }>): CommandResult => {
  let output = filterData(pv.output, cv.value);
  return { ...pv, output };
};
