import { Parser, Grammar } from "nearley";
import grammar from "./grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

const tests: [string, { query: string; expected: unknown }][] = [
  ["order by", { query: `order by "foo" asc`, expected: [{ type: "orderby", value: [{ field: "foo", direction: "asc" }] }] }],
  ["order by with space", { query: `order by "foo" desc `, expected: [{ type: "orderby", value: [{ field: "foo", direction: "desc" }] }] }],
  [
    "order by multiple",
    {
      query: `order by "foo" asc, "bar" desc`,
      expected: [
        {
          type: "orderby",
          value: [
            { field: "foo", direction: "asc" },
            { field: "bar", direction: "desc" },
          ],
        },
      ],
    },
  ],
  [
    "order by multiple with space",
    {
      query: `order by "foo" asc , "bar" desc `,
      expected: [
        {
          type: "orderby",
          value: [
            { field: "foo", direction: "asc" },
            { field: "bar", direction: "desc" },
          ],
        },
      ],
    },
  ],
  [
    "order by with pipe",
    {
      query: `order by "foo" asc, "bar" asc | limit 10`,
      expected: [
        {
          type: "orderby",
          value: [
            { field: "foo", direction: "asc" },
            { field: "bar", direction: "asc" },
          ],
        },
        { type: "limit", value: 10 },
      ],
    },
  ],
];

describe("grammar order by", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
