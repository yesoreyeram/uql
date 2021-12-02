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
    "project",
    {
      query: `project count()`,
      expected: [
        {
          type: "project",
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
    "project with space",
    {
      query: `project count() `,
      expected: [
        {
          type: "project",
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
    "project with assignment",
    {
      query: `project "foo"=count() `,
      expected: [
        {
          type: "project",
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
    "project with assignment",
    {
      query: `project "foo"=count(), "bar"=min("something"), max(2,3) `,
      expected: [
        {
          type: "project",
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
    "project with pipe",
    {
      query: `project "foo"=count(), "bar"=min("something"), max(2,3) | limit 2`,
      expected: [
        {
          type: "project",
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
];

describe("grammar project", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
