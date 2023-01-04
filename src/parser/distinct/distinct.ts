import { uniq, isArray } from "lodash";
import { Command, CommandResult } from "../../types";
import { get_single_value } from "../utils";

export const distinct = (pv: CommandResult, cv: Extract<Command, { type: "distinct" }>): CommandResult => {
  let output = pv.output;
  if (cv.value === undefined) {
    output = uniq(pv.output as unknown[]);
  } else {
    if (typeof pv.output === "object" && isArray(pv.output) && cv.value) {
      let a = pv.output.map((o) => get_single_value(o, cv.value?.value || ""));
      output = uniq(a);
    } else {
      let value = get_single_value(pv.output, cv.value.value);
      if (typeof value === "object") {
        output = uniq(value);
      }
    }
  }
  return { ...pv, output };
};
