import { Parser, Grammar } from "nearley";
import grammar from "../grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

const tests: [string, { query: string; expected: unknown }][] = [
  [
    "different strings",
    {
      query: `project strcat(str("hello"), "hello", 'hello', 'hel"lo', "hel'lo", 'hell""o', "hello",'"hello"')`,
      expected: [
        {
          type: "project",
          value: [
            {
              alias: undefined,
              args: [
                { type: "string", value: "hello" },
                { type: "ref", value: "hello" },
                { type: "string", value: "hello" },
                { type: "string", value: 'hel"lo' },
                { type: "ref", value: "hel'lo" },
                { type: "string", value: 'hell""o' },
                { type: "ref", value: "hello" },
                { type: "string", value: '"hello"' },
              ],
              operator: "strcat",
              type: "function",
            },
          ],
        },
      ],
    },
  ],
];

describe("grammar string", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
