import { Parser, Grammar } from "nearley";
import grammar from "./grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

const tests: [string, { query: string; expected: unknown }][] = [
  ["parse json", { query: `parse-json`, expected: [{ type: "parse-json", args: [] }] }],
  ["parse csv", { query: `parse-csv`, expected: [{ type: "parse-csv", args: [] }] }],
  ["parse csv with args", { query: `parse-csv --delimiter ":"`, expected: [{ type: "parse-csv", args: [[{ identifier: "delimiter", value: ":" }]] }] }],
  [
    "parse csv with multiple args",
    {
      query: `parse-csv --delimiter ":" --ignoreError true --headers "foo,bar"`,
      expected: [
        {
          type: "parse-csv",
          args: [
            [
              { identifier: "delimiter", value: ":" },
              { identifier: "ignoreError", value: "true" },
              { identifier: "headers", value: "foo,bar" },
            ],
          ],
        },
      ],
    },
  ],
  ["parse xml", { query: `parse-xml `, expected: [{ type: "parse-xml", args: [] }] }],
];

describe("grammar parse", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
