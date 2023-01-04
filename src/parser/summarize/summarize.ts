import { set, sum, min, max, mean, uniq, first, last, forEach, groupBy, toString, isArray } from "lodash";
import { filterData, get_single_value } from "../utils";
import { Command, CommandResult, ConditionalFunctionName, type_summarize_assignment, type_summarize_function } from "../../types";

export const summarize = (pv: CommandResult, cv: Extract<Command, { type: "summarize" }>): CommandResult => {
  let output = pv.output;
  let item = cv.value;
  let groupByValues = item.by;
  if (groupByValues.length === 1) {
    const groupByKey = groupByValues[0].value;
    let groupedResult = groupBy(pv.output as unknown[], groupByKey);
    let o: unknown[] = [];
    forEach(groupedResult, (g, key) => {
      let s = UQLsummarize({}, item.metrics, g);
      o.push({ [groupByKey]: key, ...s });
    });
    output = o;
  } else if (groupByValues.length > 1) {
    let groups = groupBy(output as unknown[], (a: any) => {
      return item.by.map((g) => toString(a[g.value])).join("#___#");
    });
    let out: unknown[] = [];
    forEach(groups, (group) => {
      let o: Record<string, unknown> = {};
      item.by.forEach((gi) => {
        o[gi.value] = (first(group) as any)?.[gi.value];
      });
      let s = UQLsummarize(o, item.metrics, group);
      out.push(s);
    });
    output = out;
  } else {
    output = UQLsummarize({}, item.metrics, output as unknown[]);
  }
  return { ...pv, output };
};

const IsConditionalSummaryMetric = (i: type_summarize_assignment): i is Extract<type_summarize_function, { operator: ConditionalFunctionName }> => {
  if (i.operator === "countif" || i.operator === "sumif" || i.operator === "minif" || i.operator === "maxif") {
    return true;
  }
  return false;
};

const UQLsummarize = (o: object, metrics: type_summarize_assignment[], pi: unknown[]): object => {
  metrics.forEach((i) => {
    if (IsConditionalSummaryMetric(i)) {
      const input: any[] = filterData(pi, i.condition);
      let statName = i.alias || `${i.ref.value} (${i.operator})`;
      let val;
      if (typeof input === "object" && isArray(input)) {
        switch (i.operator) {
          case "countif":
            set(o, statName, input.length);
            break;
          case "sumif":
            val = input.length > 0 ? sum(input.map((p) => get_single_value(p, i.ref.value))) : sum(input);
            set(o, statName, val);
            break;
          case "minif":
            val = input.length > 0 ? min(input.map((p) => get_single_value(p, i.ref.value))) : min(input);
            set(o, statName, val);
            break;
          case "maxif":
            val = input.length > 0 ? max(input.map((p) => get_single_value(p, i.ref.value))) : max(input);
            set(o, statName, val);
            break;
          default:
            break;
        }
      }
    } else {
      let args = i.args || [];
      let statName = args.length > 0 ? `${args[0].value} (${i.operator})` : i.operator;
      statName = i.alias ? i.alias : statName;
      let val;
      switch (i.operator) {
        case "count":
          set(o, statName, pi.length);
          break;
        case "dcount":
          val = args.length > 0 ? pi.map((p) => get_single_value(p, args[0].value)) : [pi];
          set(o, statName, uniq(val).length);
          break;
        case "mean":
          val = args.length > 0 ? mean(pi.map((p) => get_single_value(p, args[0].value))) : mean(pi);
          set(o, statName, val);
          break;
        case "sum":
          val = args.length > 0 ? sum(pi.map((p) => get_single_value(p, args[0].value))) : sum(pi);
          set(o, statName, val);
          break;
        case "min":
          val = args.length > 0 ? min(pi.map((p) => get_single_value(p, args[0].value))) : min(pi);
          set(o, statName, val);
          break;
        case "max":
          val = args.length > 0 ? max(pi.map((p) => get_single_value(p, args[0].value))) : max(pi);
          set(o, statName, val);
          break;
        case "first":
          val = args.length > 0 ? first(pi.map((p) => get_single_value(p, args[0].value))) : first(pi);
          set(o, statName, val);
          break;
        case "last":
        case "latest":
          val = args.length > 0 ? last(pi.map((p) => get_single_value(p, args[0].value))) : last(pi);
          set(o, statName, val);
          break;
        default:
          break;
      }
    }
  });
  return o;
};
