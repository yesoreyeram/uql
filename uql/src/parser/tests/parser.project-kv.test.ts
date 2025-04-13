import { uql } from "../index";

const sample_data = {
  "0x455a304a1d3342844f4dd36731f8da066efdd30b": {
    Nodes: [
      {
        Address: "86.95.103.245:46044",
        Capacity: 5,
        Load: 0,
        EthereumAddress: "0x455a304a1d3342844f4dd36731f8da066efdd30b",
      },
    ],
    Pending: 215919916032240,
    Payout: 0,
    Region: "EU",
  },
  "0x48b988eb13a4fabddc82d4b632ccb49d30b7c543": {
    Nodes: [
      {
        Address: "78.62.181.179:63447",
        Capacity: 5,
        Load: 0,
        EthereumAddress: "0x48b988eb13a4fabddc82d4b632ccb49d30b7c543",
      },
    ],
    Pending: 1128157165764810,
    Payout: 0,
    Region: "EU",
  },
};

describe("project", () => {
  describe("kv", () => {
    it("simple object", async () => {
      const data = { a: "a1", b: 2, c: true };
      const result = await uql(`project "k"=kv()`, { data });
      expect(result).toStrictEqual([
        { key: "a", value: "a1" },
        { key: "b", value: 2 },
        { key: "c", value: true },
      ]);
    });
    it("object inside array", async () => {
      const data = [
        { a: "foo", b: { b1: "foo_b1", b2: 2.1, b3: true } },
        { a: "bar", b: { b1: "bar_b1", b2: 2.2, b3: false } },
        { a: "baz", b: { b1: "baz_b1", b2: 2.3, b3: true } },
      ];
      const result = await uql(`project "k"=kv("b")`, { data });
      expect(result).toStrictEqual([
        {
          k: [
            { key: "b1", value: "foo_b1" },
            { key: "b2", value: 2.1 },
            { key: "b3", value: true },
          ],
        },
        {
          k: [
            { key: "b1", value: "bar_b1" },
            { key: "b2", value: 2.2 },
            { key: "b3", value: false },
          ],
        },
        {
          k: [
            { key: "b1", value: "baz_b1" },
            { key: "b2", value: 2.3 },
            { key: "b3", value: true },
          ],
        },
      ]);
      const result1 = await uql(`limit 2 | project "a", "k"=kv("b") | mv-expand "b"="k"`, { data });
      expect(result1).toStrictEqual([
        { a: "foo", b: { key: "b1", value: "foo_b1" } },
        { a: "foo", b: { key: "b2", value: 2.1 } },
        { a: "foo", b: { key: "b3", value: true } },
        { a: "bar", b: { key: "b1", value: "bar_b1" } },
        { a: "bar", b: { key: "b2", value: 2.2 } },
        { a: "bar", b: { key: "b3", value: false } },
      ]);
    });
    it("object inside array with extend", async () => {
      const data = [
        { a: "foo", b: { b1: "foo_b1" } },
        { a: "bar", b: { b1: "bar_b1" } },
        { a: "baz", b: { b1: "baz_b1" } },
      ];
      const result = await uql(`extend "k"=kv("b")`, { data });
      expect(result).toStrictEqual([
        { a: "foo", b: { b1: "foo_b1" }, k: [{ key: "b1", value: "foo_b1" }] },
        { a: "bar", b: { b1: "bar_b1" }, k: [{ key: "b1", value: "bar_b1" }] },
        { a: "baz", b: { b1: "baz_b1" }, k: [{ key: "b1", value: "baz_b1" }] },
      ]);
    });
    it("complex object", async () => {
      const result = await uql(
        `project kv() 
        | project "id"="key", "payout"="value.Payout", "pending"="value.Pending", "region"="value.Region",  "value"="value"
        | mv-expand "node"="value.Nodes"
        | project-away "value"`,
        { data: sample_data }
      );
      expect(result).toStrictEqual([
        {
          id: "0x455a304a1d3342844f4dd36731f8da066efdd30b",
          node: {
            Address: "86.95.103.245:46044",
            Capacity: 5,
            Load: 0,
            EthereumAddress: "0x455a304a1d3342844f4dd36731f8da066efdd30b",
          },
          payout: 0,
          pending: 215919916032240,
          region: "EU",
        },
        {
          id: "0x48b988eb13a4fabddc82d4b632ccb49d30b7c543",
          node: {
            Address: "78.62.181.179:63447",
            Capacity: 5,
            Load: 0,
            EthereumAddress: "0x48b988eb13a4fabddc82d4b632ccb49d30b7c543",
          },
          payout: 0,
          pending: 1128157165764810,
          region: "EU",
        },
      ]);
    });
  });
});
