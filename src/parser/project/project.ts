import { get, set, isArray } from "lodash";
import { get_value } from "./../utils";
import { Command, CommandResult } from "../../types";

export const project = (pv: CommandResult, cv: Extract<Command, { type: "project" }>): CommandResult => {
  let output = pv.output;
  if (typeof output === "object") {
    if (isArray(output)) {
      let refs = cv.value.filter((v) => v.type === "ref");
      let functions = cv.value.filter((v) => v.type === "function");
      output = output.map((o) => {
        let oo = {};
        refs.forEach((r) => {
          if (r.type === "ref") {
            let key = r.alias || r.value;
            set(oo, key, get(o, r.value));
          }
        });
        functions.forEach((f) => {
          if (f.type === "function") {
            let key = f.alias || f.operator;
            let args = f.args.map((arg) => {
              if (arg.type === "ref") return get(o, arg.value);
              else if (arg.type === "string") return arg.value;
              else if (arg.type === "number") return +arg.value;
            });
            let value = get_value(f.operator, args);
            set(oo, key, value);
          }
        });
        return oo;
      });
    } else {
      let refs = cv.value.filter((v) => v.type === "ref");
      let functions = cv.value.filter((v) => v.type === "function");
      let oo: Record<string, any> = {};
      let isSingle = cv.value.filter((v) => v.type === "ref" || v.type === "function").length === 1;
      refs.forEach((r) => {
        if (r.type === "ref") {
          let key = r.alias || r.value;
          set(oo, key, get(output, r.value));
        }
      });
      functions.forEach((f) => {
        if (f.type === "function") {
          let key = f.alias || f.operator;
          let args = f.args.map((arg) => {
            if (arg.type === "ref") return get(output, arg.value);
            else if (arg.type === "string") return arg.value;
            else if (arg.type === "number") return +arg.value;
          });
          let value = get_value(f.operator, args, output);
          set(oo, key, value);
        }
      });
      if (isSingle && refs.length === 1 && refs[0].type === "ref") {
        output = oo[refs[0].alias || refs[0].value];
      } else if (isSingle && functions.length === 1 && functions[0].type === "function") {
        output = oo[functions[0].alias || functions[0].operator];
      } else {
        output = oo;
      }
    }
  }
  return { ...pv, output };
};
