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
    "fn count",
    {
      query: "count()",
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "count",
            args: [],
          },
        },
      ],
    },
  ],
  [
    "fn count with space",
    {
      query: "count() ",
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "count",
            args: [],
          },
        },
      ],
    },
  ],
  [
    "fn count with inner space",
    {
      query: "count( ) ",
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "count",
            args: [],
          },
        },
      ],
    },
  ],
  [
    "fn sum",
    {
      query: "sum()",
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "sum",
            args: [],
          },
        },
      ],
    },
  ],
  [
    "fn sum with space",
    {
      query: "sum() ",
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "sum",
            args: [],
          },
        },
      ],
    },
  ],
  [
    "fn sum with inner space",
    {
      query: "sum( ) ",
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "sum",
            args: [],
          },
        },
      ],
    },
  ],
  [
    "fn sum with args",
    {
      query: `sum("foo")`,
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "sum",
            args: [{ type: "ref", value: "foo" }],
          },
        },
      ],
    },
  ],
  [
    "fn sum with multiple args",
    {
      query: `sum("foo", "bar")`,
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "sum",
            args: [
              { type: "ref", value: "foo" },
              { type: "ref", value: "bar" },
            ],
          },
        },
      ],
    },
  ],
  [
    "fn with number args",
    {
      query: `sum(1,5)`,
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "sum",
            args: [
              { type: "number", value: 1 },
              { type: "number", value: 5 },
            ],
          },
        },
      ],
    },
  ],
  [
    "fn with string args",
    {
      query: `strcat(str("foo"), str("-"), str("bar"))`,
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "strcat",
            args: [
              { type: "string", value: "foo" },
              { type: "string", value: "-" },
              { type: "string", value: "bar" },
            ],
          },
        },
      ],
    },
  ],
  [
    "fn with mixed args",
    {
      query: `count("foo", str("-"), -12.34 )`,
      expected: [
        {
          type: "command",
          value: {
            type: "function",
            operator: "count",
            args: [
              { type: "ref", value: "foo" },
              { type: "string", value: "-" },
              { type: "number", value: -12.34 },
            ],
          },
        },
      ],
    },
  ],
];

describe("grammar functions", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
