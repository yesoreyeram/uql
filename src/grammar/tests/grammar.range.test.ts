import { Parser, Grammar } from "nearley";
import grammar from "../grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

const tests: [string, { query: string; expected: unknown }][] = [
  ["range", { query: `range from 1 to 10`, expected: [{ type: "range", value: { start: 1, end: 10, step: 1 } }] }],
  ["range with space", { query: `range from 1 to 20 `, expected: [{ type: "range", value: { start: 1, end: 20, step: 1 } }] }],
  ["range with step", { query: `range from 1 to 20 step 2`, expected: [{ type: "range", value: { start: 1, end: 20, step: 2 } }] }],
  ["range with step and space", { query: `range from 1 to 20 step 0.5 `, expected: [{ type: "range", value: { start: 1, end: 20, step: 0.5 } }] }],
  [
    "range with pipe",
    {
      query: `range from 1 to 20 step 0.5 | limit 1 `,
      expected: [
        { type: "range", value: { start: 1, end: 20, step: 0.5 } },
        { type: "limit", value: 1 },
      ],
    },
  ],
  ["string range", { query: `range from "2010" to "2020"`, expected: [{ type: "range", value: { start: "2010", end: "2020", step: "" } }] }],
  ["string range with space", { query: `range from  "2010" to "2020" `, expected: [{ type: "range", value: { start: "2010", end: "2020", step: "" } }] }],
  ["string range with step", { query: `range from  "2010-03-10" to "2020" step "1M"`, expected: [{ type: "range", value: { start: "2010-03-10", end: "2020", step: "1M" } }] }],
  ["string range with step and space", { query: `range from  "2010" to "2020" step "20d" `, expected: [{ type: "range", value: { start: "2010", end: "2020", step: "20d" } }] }],
  [
    "string range with pipe",
    {
      query: `range from "2010" to "2020" step "90d" | limit 1 `,
      expected: [
        { type: "range", value: { start: "2010", end: "2020", step: "90d" } },
        { type: "limit", value: 1 },
      ],
    },
  ],
];

describe("grammar range", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
