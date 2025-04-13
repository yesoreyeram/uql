import { isArray } from "lodash";
import { Command, CommandResult } from "../../types";

export const projectAway = (pv: CommandResult, cv: Extract<Command, { type: "project-away" }>): CommandResult => {
  let output = pv.output;
  if (isArray(output)) {
    let keys = cv.value.map((c) => c.value || "");
    output = output.map((o) => {
      let oo = o;
      Object.keys(o).forEach((key) => {
        if (keys.includes(key)) {
          delete oo[key];
        }
      });
      return oo;
    });
  }
  return { ...pv, output };
};
