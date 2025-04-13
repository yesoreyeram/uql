import { isArray, uniq } from "lodash";
import { UQLsummarize } from "./../summarize/summarize";
import { Command, CommandResult } from "../../types";

export const pivot = (pv: CommandResult, cv: Extract<Command, { type: "pivot" }>): CommandResult => {
  let input = pv.output;
  if (input == null || !isArray(input)) {
    return { ...pv, output: null };
  }
  let item = cv.value;
  let rows: string[] = [];
  if (item && item.fields && item.fields.length > 0) {
    rows = uniq(input.map((u) => (item && item.fields ? u[item.fields[0].value] : ""))).filter((v) => v !== "");
  }
  let cols: string[] = [];
  if (item && item.fields && item.fields.length > 1) {
    cols = uniq(input.map((u) => (item && item.fields ? u[item.fields[1].value] : ""))).filter((v) => v !== "");
  }
  if (item.fields?.length === 2) {
    let out: any[] = [];
    rows.forEach((r) => {
      let rowName = item && item.fields ? item.fields[0].value : "";
      let colName = item && item.fields ? item.fields[1].value : "";
      let outValue: Record<string, any> = { [rowName]: r };
      cols.forEach((c) => {
        let currentItems = (input as any[]).filter((ins) => ins[rowName] === r && ins[colName] === c) || [];
        if (currentItems.length === 0) {
          outValue[c] = item.metric.operator === "count" || item.metric.operator === "dcount" || item.metric.operator === "sum" ? 0 : null;
        } else {
          let v: any = UQLsummarize({}, [{ ...item.metric, alias: item.metric.operator }], currentItems);
          outValue[c] = v[item.metric.operator];
        }
      });
      out.push(outValue);
    });
    return { ...pv, output: out };
  }
  if (item.fields?.length === 1) {
    let out: any[] = [];
    rows.forEach((r) => {
      let rowName = item && item.fields ? item.fields[0].value : "";
      let outValue = { [rowName]: r };
      let v: any = UQLsummarize(
        {},
        [{ ...item.metric, alias: item.metric.operator }],
        (input as any[]).filter((ins) => ins[rowName] === r)
      );
      outValue["value"] = v[item.metric.operator];
      out.push(outValue);
    });
    return { ...pv, output: out };
  }
  let v: any = UQLsummarize({}, [{ ...item.metric, alias: item.metric.operator }], input);
  return { ...pv, output: v[item.metric.operator] };
};
