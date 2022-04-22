import * as csv_parser from "csv-parse/lib/sync";
import { Options as csv_parser_Options } from "csv-parse/lib";
import { Command, CommandResult, type_parse_arg } from "../../types";

export const parseCsv = (pv: CommandResult, cv: Extract<Command, { type: "parse-csv" }>): CommandResult => {
  let output = pv.output;
  if (typeof pv.output === "string") {
    let csv_parser_options = get_parse_csv_options(cv.args);
    let result: string[][] = csv_parser(pv.output, csv_parser_options);
    output = result;
  }
  return { ...pv, output };
};

const get_parse_csv_options = (args: type_parse_arg[][]): csv_parser_Options => {
  let options: csv_parser_Options = { columns: true, delimiter: ",", skipEmptyLines: false, skipLinesWithError: false, relaxColumnCount: false };
  if (args[0] && args[0].length > 0) {
    args[0].forEach((arg) => {
      switch (arg.identifier) {
        case "autoParse":
        case "autoParseDate":
        case "bom":
        case "cast":
        case "castDate":
        case "relaxColumnCount":
        case "skipEmptyLines":
        case "skipLinesWithEmptyValues":
        case "trim":
          options[arg.identifier] = arg.value.toLowerCase() === "true";
          break;
        case "columns":
          options[arg.identifier] = arg.value ? arg.value.split(",") : true;
          break;
        default:
          // @ts-ignore
          options[arg.identifier] = arg.value.toLowerCase() === "true" ? true : arg.value.toLowerCase() === "false" ? false : arg.value;
          break;
      }
    });
  }
  return options;
};
