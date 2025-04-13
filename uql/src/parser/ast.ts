const { Parser, Grammar } = require("@yesoreyeram/nearley");
import grammar from "../grammar/grammar";
import { Command } from "../types";

const uqlGrammar = Grammar.fromCompiled(grammar);

export const getAST = (input: string = "hello"): Promise<Command[]> => {
  const uqlParser = new Parser(uqlGrammar);
  uqlParser.feed(input.trim() || "hello");
  const commands = uqlParser?.results || [];
  return new Promise((resolve, reject) => {
    if (commands.length === 0) {
      reject(`failed to parse query. no results found`);
    } else {
      resolve(commands[0]);
    }
  });
};
