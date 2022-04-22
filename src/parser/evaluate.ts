import * as parsers from "./parsers";
import { Command, CommandResult } from "../types";

export const evaluate = (commands: Command[], options?: { data?: any }): Promise<unknown> => {
  return new Promise((resolve) => {
    const output = options?.data || null;
    const currentTime = new Date();
    let previousValue: CommandResult = { output, context: { currentTime } };
    commands.forEach(async (currentCommand) => {
      switch (currentCommand.type) {
        case "comment":
          previousValue = parsers.comment(previousValue, currentCommand);
          break;
        case "hello":
          previousValue = parsers.hello(previousValue, currentCommand);
          break;
        case "ping":
          previousValue = parsers.ping(previousValue, currentCommand);
          break;
        case "echo":
          previousValue = parsers.echo(previousValue, currentCommand);
          break;
        case "scope":
          previousValue = parsers.scope(previousValue, currentCommand);
          break;
        case "where":
          previousValue = parsers.where(previousValue, currentCommand);
          break;
        case "jsonata":
          previousValue = parsers.jsonata(previousValue, currentCommand);
          break;
        case "distinct":
          previousValue = parsers.distinct(previousValue, currentCommand);
          break;
        case "mv-expand":
          previousValue = parsers.mvExpand(previousValue, currentCommand);
          break;
        case "count":
          previousValue = parsers.count(previousValue, currentCommand);
          break;
        case "limit":
          previousValue = parsers.limit(previousValue, currentCommand);
          break;
        case "command":
          previousValue = parsers.command(previousValue, currentCommand);
          break;
        case "orderby":
          previousValue = parsers.orderby(previousValue, currentCommand);
          break;
        case "extend":
          previousValue = parsers.extend(previousValue, currentCommand);
          break;
        case "project":
          previousValue = parsers.project(previousValue, currentCommand);
          break;
        case "project-away":
          previousValue = parsers.projectAway(previousValue, currentCommand);
          break;
        case "project-reorder":
          previousValue = parsers.projectReorder(previousValue, currentCommand);
          break;
        case "summarize":
          previousValue = parsers.summarize(previousValue, currentCommand);
          break;
        case "parse-json":
          previousValue = parsers.parseJson(previousValue, currentCommand);
          break;
        case "parse-csv":
          previousValue = parsers.parseCsv(previousValue, currentCommand);
          break;
        case "parse-xml":
          previousValue = parsers.parseXML(previousValue, currentCommand);
          break;
        case "parse-yaml":
          previousValue = parsers.parseYaml(previousValue, currentCommand);
          break;
        default:
          break;
      }
    });
    resolve(previousValue.output);
  });
};
