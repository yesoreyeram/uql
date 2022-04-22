import { Parser, Grammar } from "nearley";
import grammar from "../grammar/grammar";
import { Command } from "../types";

const oqlGrammar = Grammar.fromCompiled(grammar);

export const getAST = (input: string = "hello"): Promise<Command[]> => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input.trim() || "hello");
  const commands = oqlParser?.results || [];
  return new Promise((resolve, reject) => {
    if (commands.length === 0) {
      reject(`failed to parse query. no results found`);
    } else {
      resolve(commands[0]);
    }
  });
};
