import { isArray } from "lodash";
import { Command, CommandResult } from "../../types";
import { get_single_value } from "../utils";

export const projectReorder = (pv: CommandResult, cv: Extract<Command, { type: "project-reorder" }>): CommandResult => {
  let output = pv.output;
  if (typeof output === "object" && isArray(output)) {
    let newOutItems: any[] = [];
    output.forEach((item) => {
      let newOutItem: Record<string, any> = {};
      (cv.value || []).forEach((ci) => {
        newOutItem[ci.alias || ci.value] = get_single_value(item, ci.value);
      });
      newOutItems.push(newOutItem);
    });
    output = newOutItems;
  }
  return { ...pv, output };
};
