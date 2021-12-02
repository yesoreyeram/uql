import { Parser, Grammar } from "nearley";
import grammar from "./grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

const tests: [string, { query: string; expected: unknown }][] = [
  [
    "expression",
    {
      query: `project "foo"=( 1 + "two" - str("three") * 1 / 1.2 % -3.14 + 1.2 + sum(1,2) - 12)`,
      expected: [
        {
          type: "project",
          value: [
            {
              type: "expression",
              alias: "foo",
              args: [
                { type: "number", value: 1 },
                { type: "operation", value: "+" },
                { type: "ref", value: "two" },
                { type: "operation", value: "-" },
                { type: "string", value: "three" },
                { type: "operation", value: "*" },
                { type: "number", value: 1 },
                { type: "operation", value: "/" },
                { type: "number", value: 1.2 },
                { type: "operation", value: "%" },
                { type: "number", value: -3.14 },
                { type: "operation", value: "+" },
                { type: "number", value: 1.2 },
                { type: "operation", value: "+" },
                {
                  type: "function",
                  value: {
                    operator: "sum",
                    args: [
                      { type: "number", value: 1 },
                      { type: "number", value: 2 },
                    ],
                    type: "function",
                  },
                },
                { type: "operation", value: "-" },
                { type: "number", value: 12 },
              ],
            },
          ],
        },
      ],
    },
  ],
];

describe("grammar expression", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
