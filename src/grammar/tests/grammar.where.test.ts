import { Parser, Grammar } from "nearley";
import grammar from "../grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

describe("grammar where - basic", () => {
  const tests: [string, { query: string; expected: unknown }][] = [
    [
      "where 1",
      {
        query: `where "a" > "b"`,
        expected: [
          {
            type: "where",
            value: [
              { type: "ref", value: "a" },
              { type: "operation", value: ">" },
              { type: "ref", value: "b" },
            ],
          },
        ],
      },
    ],
    [
      "where 2",
      {
        query: `where "a" != 'foo'`,
        expected: [
          {
            type: "where",
            value: [
              { type: "ref", value: "a" },
              { type: "operation", value: "!=" },
              { type: "string", value: "foo" },
            ],
          },
        ],
      },
    ],
    [
      "where 3",
      {
        query: `where "a" >= 32.12`,
        expected: [
          {
            type: "where",
            value: [
              { type: "ref", value: "a" },
              { type: "operation", value: ">=" },
              { type: "number", value: 32.12 },
            ],
          },
        ],
      },
    ],
    [
      "where 4",
      {
        query: `where "a" =~ 'foo'`,
        expected: [
          {
            type: "where",
            value: [
              { type: "ref", value: "a" },
              { type: "operation", value: "=~" },
              { type: "string", value: "foo" },
            ],
          },
        ],
      },
    ],
    [
      "where 5",
      {
        query: `where "a" !~ 'foo'`,
        expected: [
          {
            type: "where",
            value: [
              { type: "ref", value: "a" },
              { type: "operation", value: "!~" },
              { type: "string", value: "foo" },
            ],
          },
        ],
      },
    ],
    [
      "where 6",
      {
        query: `where "a" contains 'foo'`,
        expected: [
          {
            type: "where",
            value: [
              { type: "ref", value: "a" },
              { type: "operation", value: "contains" },
              { type: "string", value: "foo" },
            ],
          },
        ],
      },
    ],
    [
      "where 7",
      {
        query: `where "a" !contains 'foo'`,
        expected: [
          {
            type: "where",
            value: [
              { type: "ref", value: "a" },
              { type: "operation", value: "!contains" },
              { type: "string", value: "foo" },
            ],
          },
        ],
      },
    ],
    [
      "where 8",
      {
        query: `where "a" contains_cs 'foo'`,
        expected: [
          {
            type: "where",
            value: [
              { type: "ref", value: "a" },
              { type: "operation", value: "contains_cs" },
              { type: "string", value: "foo" },
            ],
          },
        ],
      },
    ],
    [
      "where 9",
      {
        query: `where "a" !contains_cs 'foo'`,
        expected: [
          {
            type: "where",
            value: [
              { type: "ref", value: "a" },
              { type: "operation", value: "!contains_cs" },
              { type: "string", value: "foo" },
            ],
          },
        ],
      },
    ],
  ];
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
describe("grammar where - multi", () => {
  expect(get(`where "a" in ("a",'b',3,"d")`)[0]).toStrictEqual([
    {
      type: "where",
      value: [
        { type: "ref", value: "a" },
        { type: "operation", value: "in" },
        {
          type: "value_array",
          value: [
            { type: "ref", value: "a" },
            { type: "string", value: "b" },
            { type: "number", value: 3 },
            { type: "ref", value: "d" },
          ],
        },
      ],
    },
  ]);
  expect(get(`where "a" !in ( "a", 'b', 3 ,"d" )`)[0]).toStrictEqual([
    {
      type: "where",
      value: [
        { type: "ref", value: "a" },
        { type: "operation", value: "!in" },
        {
          type: "value_array",
          value: [
            { type: "ref", value: "a" },
            { type: "string", value: "b" },
            { type: "number", value: 3 },
            { type: "ref", value: "d" },
          ],
        },
      ],
    },
  ]);
});
