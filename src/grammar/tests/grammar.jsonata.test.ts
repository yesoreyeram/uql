import { Parser, Grammar } from "nearley";
import grammar from "../grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

const tests: [string, { query: string; expected: unknown }][] = [
  ["jsonata basic", { query: `jsonata "something"`, expected: [{ type: "jsonata", expression: "something" }] }],
  ["jsonata basic", { query: `jsonata "some other thing"`, expected: [{ type: "jsonata", expression: "some other thing" }] }],
];

describe("grammar basic", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
