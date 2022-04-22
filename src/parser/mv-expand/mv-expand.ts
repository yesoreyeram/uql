import { isArray, get, set } from "lodash";
import { Command, CommandResult } from "../../types";

export const mvExpand = (pv: CommandResult, cv: Extract<Command, { type: "mv-expand" }>): CommandResult => {
  let output = pv.output;
  if (typeof output === "object" && isArray(output)) {
    output = output
      .filter((item) => {
        let v = get(item, cv.value.value);
        return v && isArray(v) && v.length > 0;
      })
      .flatMap((item) => {
        const expandingItem = get(item, cv.value.value);
        if (expandingItem && isArray(expandingItem)) {
          return expandingItem
            .map((e) => {
              return set({ ...item }, [cv.value.alias || cv.value.value], e);
            })
            .map((item) => {
              if (cv.value.alias) {
                delete item[cv.value.value];
              }
              return item;
            });
        }
        return item;
      });
  }
  return { ...pv, output };
};
