import { Parser, Grammar } from "nearley";
import grammar from "../grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

const tests: [string, { query: string; expected: unknown }][] = [
  ["summarize", { query: `summarize count()`, expected: [{ type: "summarize", value: { by: [], metrics: [{ operator: "count", alias: undefined, args: [] }] } }] }],
  ["summarize with space", { query: `summarize count() `, expected: [{ type: "summarize", value: { by: [], metrics: [{ operator: "count", alias: undefined, args: [] }] } }] }],
  ["summarize with alias", { query: `summarize "total"=count() `, expected: [{ type: "summarize", value: { by: [], metrics: [{ operator: "count", alias: "total", args: [] }] } }] }],
  [
    "summarize with multiple",
    {
      query: `summarize "total"=count(), "avg"=mean(), min("age") `,
      expected: [
        {
          type: "summarize",
          value: {
            by: [],
            metrics: [
              { operator: "count", alias: "total", args: [] },
              { operator: "mean", alias: "avg", args: [] },
              { operator: "min", alias: undefined, args: [{ type: "ref", value: "age" }] },
            ],
          },
        },
      ],
    },
  ],
  [
    "summarize with multiple and group by",
    {
      query: `summarize "total"=count(), "avg"=mean(), min("age") by "foo",  "bar" `,
      expected: [
        {
          type: "summarize",
          value: {
            by: [
              { type: "ref", value: "foo" },
              { type: "ref", value: "bar" },
            ],
            metrics: [
              { operator: "count", alias: "total", args: [] },
              {
                operator: "mean",
                alias: "avg",
                args: [],
              },
              { operator: "min", alias: undefined, args: [{ type: "ref", value: "age" }] },
            ],
          },
        },
      ],
    },
  ],
  [
    "summarize with pipe",
    { query: `summarize count()\n  | count `, expected: [{ type: "summarize", value: { by: [], metrics: [{ operator: "count", alias: undefined, args: [] }] } }, { type: "count" }] },
  ],
  [
    "summarize with pipe",
    { query: `summarize count() \n  | count`, expected: [{ type: "summarize", value: { by: [], metrics: [{ operator: "count", alias: undefined, args: [] }] } }, { type: "count" }] },
  ],
  [
    "summarize with pipe",
    {
      query: `summarize "total"=count(), "avg"=mean(), min("age") by "foo",  "bar" | count`,
      expected: [
        {
          type: "summarize",
          value: {
            by: [
              { type: "ref", value: "foo" },
              { type: "ref", value: "bar" },
            ],
            metrics: [
              { operator: "count", alias: "total", args: [] },
              {
                operator: "mean",
                alias: "avg",
                args: [],
              },
              { operator: "min", alias: undefined, args: [{ type: "ref", value: "age" }] },
            ],
          },
        },
        {
          type: "count",
        },
      ],
    },
  ],
];

describe("grammar summarize", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
  it("summarize with parameters", () => {
    expect(get(`summarize "a"=sumif("a" , "a" > "b")`)[0]).toStrictEqual([
      {
        type: "summarize",
        value: {
          by: [],
          metrics: [
            {
              args: undefined,
              alias: "a",
              operator: "sumif",
              condition: [
                { type: "ref", value: "a" },
                { type: "operation", value: ">" },
                { type: "ref", value: "b" },
              ],
              ref: {
                type: "ref",
                value: "a",
              },
            },
          ],
        },
      },
    ]);
  });
});
