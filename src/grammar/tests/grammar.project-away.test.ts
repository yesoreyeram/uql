import { Parser, Grammar } from "nearley";
import grammar from "../grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

const tests: [string, { query: string; expected: unknown }][] = [
  ["project-away", { query: `project-away "foo"`, expected: [{ type: "project-away", value: [{ type: "ref", value: "foo" }] }] }],
  [
    "project-away multi",
    {
      query: `project-away "foo" , "bar"`,
      expected: [
        {
          type: "project-away",
          value: [
            { type: "ref", value: "foo" },
            { type: "ref", value: "bar" },
          ],
        },
      ],
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
