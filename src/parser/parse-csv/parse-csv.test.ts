import { uql } from "../index";

describe("parse-csv", () => {
  it("csv with headers", async () => {
    const result = await uql("parse-csv", { data: "a,b\n1,2\n3,4" });
    expect(result).toStrictEqual([
      { a: "1", b: "2" },
      { a: "3", b: "4" },
    ]);
  });
  it("csv with custom headers", async () => {
    const result = await uql("parse-csv --columns 'a,b'", { data: "1,2\n3,4" });
    expect(result).toStrictEqual([
      { a: "1", b: "2" },
      { a: "3", b: "4" },
    ]);
  });
  it("csv without headers", async () => {
    const result = await uql("parse-csv --columns 'false'", { data: "1,2\n3,4" });
    expect(result).toStrictEqual([
      { col_0: "1", col_1: "2" },
      { col_0: "3", col_1: "4" },
    ]);
  });
  it("csv with custom delimiter", async () => {
    const result = await uql("parse-csv --delimiter ';'", { data: "a;b\n1;2\n3;4" });
    expect(result).toStrictEqual([
      { a: "1", b: "2" },
      { a: "3", b: "4" },
    ]);
  });
  it("tsv", async () => {
    const result = await uql("parse-csv --delimiter '\t'", { data: `a	b\n1	2\n3	4` });
    expect(result).toStrictEqual([
      { a: "1", b: "2" },
      { a: "3", b: "4" },
    ]);
  });
});
