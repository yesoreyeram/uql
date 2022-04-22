import { get, uniq, isArray } from "lodash";
import { Command, CommandResult } from "../../types";

export const distinct = (pv: CommandResult, cv: Extract<Command, { type: "distinct" }>): CommandResult => {
  let output = pv.output;
  if (cv.value === undefined) {
    output = uniq(pv.output as unknown[]);
  } else {
    if (typeof pv.output === "object" && isArray(pv.output) && cv.value) {
      let a = pv.output.map((o) => get(o, cv.value?.value || ""));
      output = uniq(a);
    } else {
      let value = get(pv.output, cv.value.value);
      output = uniq(value);
    }
  }
  return { ...pv, output };
};
