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
    "extend",
    {
      query: `extend count()`,
      expected: [
        {
          type: "extend",
          value: [
            {
              type: "function",
              alias: undefined,
              operator: "count",
              args: [],
            },
          ],
        },
      ],
    },
  ],
  [
    "extend with space",
    {
      query: `extend count() `,
      expected: [
        {
          type: "extend",
          value: [
            {
              type: "function",
              alias: undefined,
              operator: "count",
              args: [],
            },
          ],
        },
      ],
    },
  ],
  [
    "extend with assignment",
    {
      query: `extend "foo"=count() `,
      expected: [
        {
          type: "extend",
          value: [
            {
              type: "function",
              alias: "foo",
              operator: "count",
              args: [],
            },
          ],
        },
      ],
    },
  ],
  [
    "extend with assignment",
    {
      query: `extend "foo"=count(), "bar"=min("something"), max(2,3) `,
      expected: [
        {
          type: "extend",
          value: [
            {
              type: "function",
              alias: "foo",
              operator: "count",
              args: [],
            },
            {
              type: "function",
              alias: "bar",
              operator: "min",
              args: [{ type: "ref", value: "something" }],
            },
            {
              type: "function",
              alias: undefined,
              operator: "max",
              args: [
                { type: "number", value: 2 },
                { type: "number", value: 3 },
              ],
            },
          ],
        },
      ],
    },
  ],
  [
    "extend with pipe",
    {
      query: `extend "foo"=count(), "bar"=min("something"), max(2,3) | limit 2`,
      expected: [
        {
          type: "extend",
          value: [
            {
              type: "function",
              alias: "foo",
              operator: "count",
              args: [],
            },
            {
              type: "function",
              alias: "bar",
              operator: "min",
              args: [{ type: "ref", value: "something" }],
            },
            {
              type: "function",
              alias: undefined,
              operator: "max",
              args: [
                { type: "number", value: 2 },
                { type: "number", value: 3 },
              ],
            },
          ],
        },
        {
          type: "limit",
          value: 2,
        },
      ],
    },
  ],
  [
    "random",
    {
      query: `extend "foo"=random(1.2,2.4,true)`,
      expected: [
        {
          type: "extend",
          value: [
            {
              type: "function",
              alias: "foo",
              operator: "random",
              args: [
                { value: 1.2, type: "number" },
                { type: "number", value: 2.4 },
                { type: "identifier", value: "true" },
              ],
            },
          ],
        },
      ],
    },
  ],
];

describe("grammar extend", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
