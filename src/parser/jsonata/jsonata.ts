import { isArray } from "lodash";
import * as JSONata from "jsonata";
import { Command, CommandResult } from "../../types";

export const jsonata = (pv: CommandResult, cv: Extract<Command, { type: "jsonata" }>): CommandResult => {
  const expression = JSONata(cv.expression);
  let out = expression.evaluate(pv.output);
  if (out && typeof out === "object" && isArray(out)) {
    delete (out as any).sequence; // https://github.com/jsonata-js/jsonata/issues/296
    delete (out as any).keepSingleton; // https://github.com/jsonata-js/jsonata/issues/469
  }
  return { ...pv, output: out };
};
