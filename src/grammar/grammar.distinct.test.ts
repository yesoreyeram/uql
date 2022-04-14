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
    "distinct",
    {
      query: `distinct`,
      expected: [{ type: "distinct", value: undefined }],
    },
  ],
  [
    "distinct with args",
    {
      query: `distinct "foo.bar"`,
      expected: [{ type: "distinct", value: { type: "ref", value: "foo.bar" } }],
    },
  ],
];

describe("grammar project", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
