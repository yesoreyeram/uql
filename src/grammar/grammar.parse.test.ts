import { Parser, Grammar } from "nearley";
import grammar from "./grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

const tests: [string, { query: string; expected: unknown }][] = [
  ["parse json", { query: `parse-json`, expected: [{ type: "parse-json" }] }],
  ["parse csv", { query: `parse-csv`, expected: [{ type: "parse-csv" }] }],
  ["parse xml", { query: `parse-xml `, expected: [{ type: "parse-xml" }] }],
];

describe("grammar parse", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
