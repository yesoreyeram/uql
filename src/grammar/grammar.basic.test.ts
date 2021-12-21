import { Parser, Grammar } from "nearley";
import grammar from "./grammar";

const oqlGrammar = Grammar.fromCompiled(grammar);

const get = (input: string): unknown[] => {
  const oqlParser = new Parser(oqlGrammar);
  oqlParser.feed(input);
  return oqlParser.results;
};

const tests: [string, { query: string; expected: unknown }][] = [
  ["hello", { query: "hello", expected: [{ type: "hello" }] }],
  ["hello with space", { query: "hello ", expected: [{ type: "hello" }] }],
  ["ping", { query: "ping", expected: [{ type: "ping", value: "pong" }] }],
  ["ping with space", { query: "ping ", expected: [{ type: "ping", value: "pong" }] }],
  ["count", { query: "count", expected: [{ type: "count" }] }],
  ["count with space", { query: "count ", expected: [{ type: "count" }] }],
  ["limit", { query: "limit 10", expected: [{ type: "limit", value: 10 }] }],
  ["limit with space", { query: "limit 20  ", expected: [{ type: "limit", value: 20 }] }],
  ["multiple", { query: "hello | hello", expected: [{ type: "hello" }, { type: "hello" }] }],
  ["multiple with space", { query: "limit 10  | hello ", expected: [{ type: "limit", value: 10 }, { type: "hello" }] }],
  [
    "multiple with space",
    {
      query: "hello | sum(1, 3) | hello ",
      expected: [
        { type: "hello" },
        {
          type: "command",
          value: {
            type: "function",
            operator: "sum",
            args: [
              { type: "number", value: 1 },
              { type: "number", value: 3 },
            ],
          },
        },
        { type: "hello" },
      ],
    },
  ],
  ["hello with newline", { query: "hello\n|hello", expected: [{ type: "hello" }, { type: "hello" }] }],
  ["hello with newline", { query: "hello\n | hello ", expected: [{ type: "hello" }, { type: "hello" }] }],
  ["hello with newline", { query: "hello \n | hello ", expected: [{ type: "hello" }, { type: "hello" }] }],
  ["hello with newline", { query: "hello \n |hello", expected: [{ type: "hello" }, { type: "hello" }] }],
  ["scope", { query: `scope "foo.bar"`, expected: [{ type: "scope", value: { type: "ref", value: "foo.bar" } }] }],
];

describe("grammar basic", () => {
  it.each(tests)("%s", (_, test) => {
    const { query, expected } = test as { query: string; expected: unknown };
    const results = get(query as string);
    expect(results[0]).toStrictEqual(expected);
  });
});
