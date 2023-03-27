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
    "pivot - with default arguments",
    {
      query: `pivot count()`,
      expected: [
        {
          type: "pivot",
          value: {
            metric: { alias: undefined, args: [], operator: "count" },
            fields: [],
          },
        },
      ],
    },
  ],
  [
    "pivot - with just aggregation",
    {
      query: `pivot sum("quantity")`,
      expected: [
        {
          type: "pivot",
          value: {
            metric: { alias: undefined, args: [{ type: "ref", value: "quantity" }], operator: "sum" },
            fields: [],
          },
        },
      ],
    },
  ],
  [
    "pivot - with aggregation and col",
    {
      query: `pivot sum("quantity"), "fruit"`,
      expected: [
        {
          type: "pivot",
          value: {
            metric: { alias: undefined, args: [{ type: "ref", value: "quantity" }], operator: "sum" },
            fields: [{ type: "ref", value: "fruit" }],
          },
        },
      ],
    },
  ],
  [
    "pivot - with aggregation and col and row",
    {
      query: `pivot sum("quantity"), "fruit", "size"`,
      expected: [
        {
          type: "pivot",
          value: {
            metric: { alias: undefined, args: [{ type: "ref", value: "quantity" }], operator: "sum" },
            fields: [
              { type: "ref", value: "fruit" },
              { type: "ref", value: "size" },
            ],
          },
        },
      ],
    },
  ],
  [
    "pivot - with aggregation and col and row with alias",
    {
      query: `pivot "qty"=sum("quantity"), "fruit", "size"`,
      expected: [
        {
          type: "pivot",
          value: {
            metric: { alias: "qty", args: [{ type: "ref", value: "quantity" }], operator: "sum" },
            fields: [
              { type: "ref", value: "fruit" },
              { type: "ref", value: "size" },
            ],
          },
        },
      ],
    },
  ],
];

describe("grammar pivot", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
