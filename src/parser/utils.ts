import { sum, min, max, mean, uniq, isArray, random, first, last, forEach, flatten, get } from "lodash";
import * as dayjs from "dayjs";
import { FunctionName, type_where_arg } from "../types";

export const get_single_value = (input: any, query: string): string | number | any[] => {
  return get(input, query);
};

export const get_value = (operator: FunctionName, args: any[], previous_value?: any): unknown => {
  switch (operator) {
    case "tolower":
      return (args[0] + "").toLowerCase();
    case "toupper":
      return (args[0] + "").toUpperCase();
    case "strlen":
      return (args[0] + "").length;
    case "reverse":
      return typeof args[0] === "string" ? args[0].split("").reverse().join("") : null;
    case "split":
      if (args.length >= 1 && typeof args[0] === "string") {
        return (args[0] || "").split(typeof args[1] === "string" ? args[1] : "");
      }
      return [];
    case "substring":
      if (typeof args[0] === "string") {
        return args[0].substring(args[1], args[2]);
      }
      return "";
    case "replace_string":
      if (args.length >= 2 && typeof args[0] === "string" && typeof args[1] === "string") {
        return (args[0] || "").replace(new RegExp(args[1], args[3] || "g"), typeof args[2] === "string" ? args[2] : "");
      }
      return args[0] || "";
    case "atob":
      return Buffer.from(args[0], "base64").toString("ascii");
    case "btoa":
      return Buffer.from(args[0]).toString("base64");
    case "trim":
      return (args[0] + "").trim();
    case "trim_start":
      return (args[0] + "").trimStart();
    case "trim_end":
      return (args[0] + "").trimEnd();
    case "strcat":
      return args.join("");
    case "sum":
      return sum(args);
    case "diff":
      return args[0] - args[1];
    case "mul":
      if (typeof args === "object" && isArray(args)) {
        return args.reduce((a, b) => a * b, 1);
      }
      return args[0] * args[1];
    case "div":
      return args[0] / args[1];
    case "percentage":
      if (args.length === 2) {
        return (args[0] / args[1]) * 100;
      }
      return null;
    case "min":
      return min(args);
    case "max":
      return max(args);
    case "mean":
      return mean(args);
    case "first":
      return first(args);
    case "last":
    case "latest":
      return last(args);
    case "sign":
      if (args.length > 0 && typeof args[0] === "number") {
        return Math.sign(args[0]);
      }
      return null;
    case "sin":
      if (args.length > 0 && typeof args[0] === "number") {
        return Math.sin(args[0]);
      }
      return null;
    case "cos":
      if (args.length > 0 && typeof args[0] === "number") {
        return Math.cos(args[0]);
      }
      return null;
    case "tan":
      if (args.length > 0 && typeof args[0] === "number") {
        return Math.tan(args[0]);
      }
      return null;
    case "pow":
      if (args.length > 1 && typeof args[0] === "number" && typeof args[1] === "number") {
        return Math.pow(args[0], args[1]);
      }
      return null;
    case "round":
      if (args.length > 0 && typeof args[0] === "number") {
        return Math.round(args[0]);
      }
      return null;
    case "ceil":
      if (args.length > 0 && typeof args[0] === "number") {
        return Math.ceil(args[0]);
      }
      return null;
    case "floor":
      if (args.length > 0 && typeof args[0] === "number") {
        return Math.floor(args[0]);
      }
      return null;
    case "log":
      if (args.length > 0 && typeof args[0] === "number") {
        return Math.log(args[0]);
      }
      return null;
    case "log2":
      if (args.length > 0 && typeof args[0] === "number") {
        return Math.log2(args[0]);
      }
      return null;
    case "log10":
      if (args.length > 0 && typeof args[0] === "number") {
        return Math.log10(args[0]);
      }
      return null;
    case "extract":
      if (args.length >= 3) {
        const reg: string = args[0] || "";
        const index = args.length > 1 ? args[1] : 0;
        const input: string = args[2] || "";
        const match = input.match(reg);
        const out = match === null ? null : match[index];
        if (out !== null) {
          switch (args[3]) {
            case "number":
              return +out;
            case "date":
              return new Date(out);
            default:
              return out;
          }
        }
      }
      return null;
    case "parse_url":
      if (typeof args[0] === "string" && args[0] !== "") {
        try {
          const url = new URL(args[0]);
          if (args.length === 3 && (args[1] === "search" || args[1] === "query")) {
            return url.searchParams.get(args[2]) || "";
          } else if (args.length === 2) {
            return (url as any)[args[1]] || "";
          } else {
            return url.toString();
          }
        } catch (ex) {
          return "";
        }
      }
      return "";
    case "parse_urlquery":
      if (typeof args[0] === "string") {
        const searchParams = new URLSearchParams(args[0]);
        if (args.length == 2) {
          return searchParams.get(args[1]) || "";
        } else {
          let o: Record<string, string> = {};
          searchParams.forEach((value, key) => (o[key] = value));
          return o;
        }
      }
      return "";
    case "kv":
      let out: { key: string | number; value: any }[] = [];
      if (args.length > 0 && typeof args[0] === "object") {
        forEach(args[0], (value, key) => {
          out.push({ key, value });
        });
      } else if (args.length === 0 && previous_value && typeof previous_value === "object") {
        forEach(previous_value, (value, key) => {
          out.push({ key, value });
        });
      }
      return out;
    case "array_to_map":
      const outMap: Record<string, any> = {};
      if (typeof args[0] === "object" && isArray(args[0])) {
        args[0].forEach((item, idx) => {
          outMap[args[idx + 1] || idx] = item;
        });
      } else if (typeof args[0] === "object") {
        return args[0];
      }
      return outMap;
    case "array_from_entries":
      const arrayResult: any[] = [];
      if (args.length >= 2 && typeof args[1] === "object" && isArray(args[1])) {
        for (let index = 0; index < args[1].length; index++) {
          let out: Record<string, any> = {};
          args.forEach((arg, idx) => {
            if (typeof arg === "string" && idx % 2 == 0) {
              out[arg] = args[idx + 1]?.[index] ?? null;
            }
          });
          arrayResult.push(out);
        }
      }
      return arrayResult;
    case "pack":
    case "bag_pack":
      const packResult: Record<string, any> = {};
      args.forEach((arg, index) => {
        if (index % 2 === 0) {
          packResult[arg] = args[index + 1];
        }
      });
      return packResult;
    case "toint":
    case "tolong":
    case "todouble":
    case "tofloat":
    case "tonumber":
      return +(args[0] || "");
    case "tobool":
      if (typeof args[0] === "string") return args[0].toLowerCase() === "true" ? true : false;
      else if (typeof args[0] === "number") return args[0] <= 0 ? true : false;
      else return args[0];
    case "tostring":
      return args[0] + "";
    case "todatetime":
      if (typeof args[0] === "number") return new Date(args[0] + "");
      return new Date(args[0] || "");
    case "tounixtime":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return args[0].getTime();
      return args[0];
    case "format_datetime":
      try {
        if (typeof args[0] === "object" && typeof args[0].getTime === "function" && typeof args[1] === "string") return dayjs(args[0]).format(args[1]);
        return args[0];
      } catch (ex) {
        console.error(ex);
        return null;
      }
    case "add_datetime":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function" && typeof args[1] === "string") {
        let numberParts = args[1].match(/[\-\d]+/g) || [];
        let textParts = args[1].match(/[A-Za-z]+/g) || [];
        if (numberParts.length > 0 && textParts.length > 0) {
          let o = dayjs(args[0]).add(numberParts[0] === "-" ? -1 : +(numberParts[0] || ""), textParts[0] as dayjs.ManipulateType);
          return o.toDate();
        }
      }
      return args[0];
    case "startofminute":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("m").toDate();
      return args[0];
    case "startofhour":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("h").toDate();
      return args[0];
    case "startofday":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("d").toDate();
      return args[0];
    case "startofweek":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("w").toDate();
      return args[0];
    case "startofmonth":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("M").toDate();
      return args[0];
    case "startofyear":
      if (typeof args[0] === "object" && typeof args[0].getTime === "function") return dayjs(args[0]).startOf("y").toDate();
      return args[0];
    case "unixtime_seconds_todatetime":
      if (typeof args[0] === "string") return new Date(+args[0] * 1000);
      return new Date(args[0] * 1000);
    case "unixtime_milliseconds_todatetime":
      if (typeof args[0] === "string") return new Date(+args[0]);
      return new Date(args[0]);
    case "unixtime_microseconds_todatetime":
      if (typeof args[0] === "string") return new Date(+args[0] / 1000);
      return new Date(args[0] / 1000);
    case "unixtime_nanoseconds_todatetime":
      if (typeof args[0] === "string") return new Date(+args[0] / 1000 / 1000);
      return new Date(args[0] / 1000 / 1000);
    case "random":
      return random(...args);
    case "distinct":
      if (args.length > 0) {
        return uniq(args[0]);
      } else if (args.length === 0 && previous_value && typeof previous_value === "object" && isArray(previous_value)) {
        return uniq(previous_value);
      }
      return null;
    case "dcount":
    default:
      throw "not implemented";
  }
};

export const filterData = (output: any, args: type_where_arg[]): any => {
  if (typeof output === "object" && isArray(output) && args.length >= 3 && args[1].type === "operation") {
    return output.filter((o) => {
      let lhs = args[0].value;
      let rhs = args[2].value;
      let rhs_rgs =
        args[2].type === "value_array"
          ? flatten(
              args[2].value.map((v) => {
                if (v.type === "ref") return get_single_value(o, v.value);
                return v.value;
              })
            )
          : [];
      if (args[0].type === "ref") {
        lhs = get_single_value(o, args[0].value);
      }
      if (args[2].type === "ref") {
        rhs = get_single_value(o, args[2].value);
      }
      if (args[1].type === "operation") {
        switch (args[1].value) {
          case "==":
            return lhs === rhs;
          case "!=":
            return lhs !== rhs;
          case ">":
            return lhs > rhs;
          case ">=":
            return lhs >= rhs;
          case "<":
            return lhs < rhs;
          case "<=":
            return lhs <= rhs;
          case "between":
            if (args[2].type === "value_array" && rhs_rgs.length >= 2) {
              return lhs >= rhs_rgs[0] && lhs <= rhs_rgs[1];
            }
            return false;
          case "inside":
            if (args[2].type === "value_array" && rhs_rgs.length >= 2) {
              return lhs > rhs_rgs[0] && lhs < rhs_rgs[1];
            }
            return false;
          case "outside":
            if (args[2].type === "value_array" && rhs_rgs.length >= 2) {
              return lhs < rhs_rgs[0] || lhs > rhs_rgs[1];
            }
            return false;
          case "in":
            if (args[2].type === "value_array") {
              return rhs_rgs.includes(lhs);
            }
            return false;
          case "!in":
            if (args[2].type === "value_array") {
              return !rhs_rgs.includes(lhs);
            }
            return false;
          case "in~":
            if (args[2].type === "value_array") {
              return rhs_rgs.map((t) => (typeof t === "string" ? t.toLowerCase() : t)).includes(typeof lhs === "string" ? lhs.toLowerCase() : lhs);
            }
            return false;
          case "!in~":
            if (args[2].type === "value_array") {
              return !rhs_rgs.map((t) => (typeof t === "string" ? t.toLowerCase() : t)).includes(typeof lhs === "string" ? lhs.toLowerCase() : lhs);
            }
            return false;
          case "=~":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return lhs.toLowerCase().includes(rhs.toLowerCase());
            }
            return false;
          case "!~":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return !lhs.toLowerCase().includes(rhs.toLowerCase());
            }
            return false;
          case "matches regex":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return lhs.match(rhs) !== null;
            }
            return false;
          case "!matches regex":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return lhs.match(rhs) === null;
            }
            return false;
          case "contains":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return lhs.toLowerCase().includes(rhs.toLowerCase());
            }
            return false;
          case "!contains":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return !lhs.toLowerCase().includes(rhs.toLowerCase());
            }
            return false;
          case "contains_cs":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return lhs.includes(rhs);
            }
            return false;
          case "!contains_cs":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return !lhs.includes(rhs);
            }
            return false;
          case "startswith":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return lhs.toLowerCase().startsWith(rhs.toLowerCase());
            }
            return false;
          case "!startswith":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return !lhs.toLowerCase().startsWith(rhs.toLowerCase());
            }
            return false;
          case "startswith_cs":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return lhs.startsWith(rhs);
            }
            return false;
          case "!startswith_cs":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return !lhs.startsWith(rhs);
            }
            return false;
          case "endswith":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return lhs.toLowerCase().endsWith(rhs.toLowerCase());
            }
            return false;
          case "!endswith":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return !lhs.toLowerCase().endsWith(rhs.toLowerCase());
            }
            return false;
          case "endswith_cs":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return lhs.endsWith(rhs);
            }
            return false;
          case "!endswith_cs":
            if (typeof lhs === "string" && typeof rhs === "string") {
              return !lhs.endsWith(rhs);
            }
            return false;
        }
      }
      return false;
    });
  } else {
    return [];
  }
};
