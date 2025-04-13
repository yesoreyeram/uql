import { isArray, orderBy } from "lodash";
import { Command, CommandResult } from "../../types";

export const orderby = (pv: CommandResult, cv: Extract<Command, { type: "orderby" }>): CommandResult => {
  let output = pv.output;
  if (isArray(output)) {
    output = orderBy(
      output as unknown[],
      cv.value.map((o) => o.field),
      cv.value.map((o) => o.direction || "asc")
    );
  }
  return { ...pv, output };
};
