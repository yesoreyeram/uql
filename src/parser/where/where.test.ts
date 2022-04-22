import { uql } from "../index";

describe("parser", () => {
  describe("where", () => {
    it("numeric constants", async () => {
      const input = [{ a: 1 }, { a: 10 }, { a: 20 }, { a: 30 }];
      expect(await uql(`where "a" == 10`, { data: input })).toStrictEqual([{ a: 10 }]);
      expect(await uql(`where "a" != 10`, { data: input })).toStrictEqual([{ a: 1 }, { a: 20 }, { a: 30 }]);
      expect(await uql(`where "a" > 10`, { data: input })).toStrictEqual([{ a: 20 }, { a: 30 }]);
      expect(await uql(`where "a" >= 10`, { data: input })).toStrictEqual([{ a: 10 }, { a: 20 }, { a: 30 }]);
      expect(await uql(`where "a" < 10`, { data: input })).toStrictEqual([{ a: 1 }]);
      expect(await uql(`where "a" <= 10`, { data: input })).toStrictEqual([{ a: 1 }, { a: 10 }]);
    });
    it("numeric ref", async () => {
      const input = [
        { a: 1, b: 10 },
        { a: 10, b: 10 },
        { a: 20, b: 10 },
        { a: 30, b: 10 },
      ];
      expect(await uql(`where "a" == "b" | project "a"`, { data: input })).toStrictEqual([{ a: 10 }]);
      expect(await uql(`where "a" != "b" | project "a"`, { data: input })).toStrictEqual([{ a: 1 }, { a: 20 }, { a: 30 }]);
      expect(await uql(`where "a" > "b" | project "a"`, { data: input })).toStrictEqual([{ a: 20 }, { a: 30 }]);
      expect(await uql(`where "a" >= "b" | project "a"`, { data: input })).toStrictEqual([{ a: 10 }, { a: 20 }, { a: 30 }]);
      expect(await uql(`where "a" < "b" | project "a"`, { data: input })).toStrictEqual([{ a: 1 }]);
      expect(await uql(`where "a" <= "b" | project "a"`, { data: input })).toStrictEqual([{ a: 1 }, { a: 10 }]);
    });
    it("string constants", async () => {
      const input = [{ a: "foo" }, { a: "Foo" }, { a: "bar" }, { a: "Bar" }];
      expect(await uql(`where "a" == 'foo'`, { data: input })).toStrictEqual([{ a: "foo" }]);
      expect(await uql(`where "a" != 'foo'`, { data: input })).toStrictEqual([{ a: "Foo" }, { a: "bar" }, { a: "Bar" }]);
      expect(await uql(`where "a" =~ 'foo'`, { data: input })).toStrictEqual([{ a: "foo" }, { a: "Foo" }]);
      expect(await uql(`where "a" !~ 'foo'`, { data: input })).toStrictEqual([{ a: "bar" }, { a: "Bar" }]);
    });
    it("string ref", async () => {
      const input = [
        { a: "foo", b: "foo" },
        { a: "Foo", b: "foo" },
        { a: "bar", b: "foo" },
        { a: "Bar", b: "foo" },
      ];
      expect(await uql(`where "a" == "b" | project "a"`, { data: input })).toStrictEqual([{ a: "foo" }]);
      expect(await uql(`where "a" != "b" | project "a"`, { data: input })).toStrictEqual([{ a: "Foo" }, { a: "bar" }, { a: "Bar" }]);
      expect(await uql(`where "a" =~ "b" | project "a"`, { data: input })).toStrictEqual([{ a: "foo" }, { a: "Foo" }]);
      expect(await uql(`where "a" !~ "b" | project "a"`, { data: input })).toStrictEqual([{ a: "bar" }, { a: "Bar" }]);
    });
    it("string constants - contains", async () => {
      const input = [{ a: "FabriKam" }, { a: "banana" }];
      expect(await uql(`where "a" contains 'BRik'`, { data: input })).toStrictEqual([{ a: "FabriKam" }]);
      expect(await uql(`where "a" !contains 'ana'`, { data: input })).toStrictEqual([{ a: "FabriKam" }]);
      expect(await uql(`where "a" not contains 'ana'`, { data: input })).toStrictEqual([{ a: "FabriKam" }]);
      expect(await uql(`where "a" contains_cs 'Kam'`, { data: input })).toStrictEqual([{ a: "FabriKam" }]);
      expect(await uql(`where "a" !contains_cs 'Kam'`, { data: input })).toStrictEqual([{ a: "banana" }]);
      expect(await uql(`where "a" not contains_cs 'Kam'`, { data: input })).toStrictEqual([{ a: "banana" }]);
    });
    it("string constants - startswith", async () => {
      const input = [{ a: "FabriKam" }, { a: "fabrikam" }, { a: "banana" }];
      expect(await uql(`where "a" startswith 'fab'`, { data: input })).toStrictEqual([{ a: "FabriKam" }, { a: "fabrikam" }]);
      expect(await uql(`where "a" !startswith 'kam'`, { data: input })).toStrictEqual([{ a: "FabriKam" }, { a: "fabrikam" }, { a: "banana" }]);
      expect(await uql(`where "a" startswith_cs 'Fab'`, { data: input })).toStrictEqual([{ a: "FabriKam" }]);
      expect(await uql(`where "a" !startswith_cs 'Fab'`, { data: input })).toStrictEqual([{ a: "fabrikam" }, { a: "banana" }]);
    });
    it("string constants - endswith", async () => {
      const input = [{ a: "FabriKam" }, { a: "fabrikam" }, { a: "banana" }];
      expect(await uql(`where "a" endswith 'kam'`, { data: input })).toStrictEqual([{ a: "FabriKam" }, { a: "fabrikam" }]);
      expect(await uql(`where "a" !endswith 'fab'`, { data: input })).toStrictEqual([{ a: "FabriKam" }, { a: "fabrikam" }, { a: "banana" }]);
      expect(await uql(`where "a" endswith_cs 'Kam'`, { data: input })).toStrictEqual([{ a: "FabriKam" }]);
      expect(await uql(`where "a" !endswith_cs 'Kam'`, { data: input })).toStrictEqual([{ a: "fabrikam" }, { a: "banana" }]);
    });
    it("numeric in", async () => {
      const input = [{ a: 1 }, { a: 10 }, { a: 20 }, { a: 30 }];
      expect(await uql(`where "a" in (10,20)`, { data: input })).toStrictEqual([{ a: 10 }, { a: 20 }]);
      expect(await uql(`where "a" !in (10,20)`, { data: input })).toStrictEqual([{ a: 1 }, { a: 30 }]);
    });
    it("string in", async () => {
      const input = [{ a: "foo" }, { a: "Foo" }, { a: "bar" }, { a: "Bar" }];
      expect(await uql(`where "a" in ('foo','bar')`, { data: input })).toStrictEqual([{ a: "foo" }, { a: "bar" }]);
      expect(await uql(`where "a" !in ('foo','bar')`, { data: input })).toStrictEqual([{ a: "Foo" }, { a: "Bar" }]);
      expect(await uql(`where "a" in~ ('foo')`, { data: input })).toStrictEqual([{ a: "foo" }, { a: "Foo" }]);
      expect(await uql(`where "a" !in~ ('foo')`, { data: input })).toStrictEqual([{ a: "bar" }, { a: "Bar" }]);
    });
    it("numeric between", async () => {
      const input = [{ a: 1 }, { a: 10 }, { a: 20 }, { a: 30 }];
      expect(await uql(`where "a" between (10,20)`, { data: input })).toStrictEqual([{ a: 10 }, { a: 20 }]);
      expect(await uql(`where "a" between (10,19)`, { data: input })).toStrictEqual([{ a: 10 }]);
      expect(await uql(`where "a" between (11,19)`, { data: input })).toStrictEqual([]);
      expect(await uql(`where "a" between (11,20)`, { data: input })).toStrictEqual([{ a: 20 }]);
      expect(await uql(`where "a" between (10)`, { data: input })).toStrictEqual([]); // incorrect arguments
    });
    it("numeric inside", async () => {
      const input = [{ a: 1 }, { a: 10 }, { a: 20 }, { a: 30 }];
      expect(await uql(`where "a" inside (1,20)`, { data: input })).toStrictEqual([{ a: 10 }]);
      expect(await uql(`where "a" inside (10,20)`, { data: input })).toStrictEqual([]);
      expect(await uql(`where "a" inside (10)`, { data: input })).toStrictEqual([]); // incorrect arguments
    });
    it("numeric outside", async () => {
      const input = [{ a: 1 }, { a: 10 }, { a: 20 }, { a: 30 }];
      expect(await uql(`where "a" outside (10,20)`, { data: input })).toStrictEqual([{ a: 1 }, { a: 30 }]);
      expect(await uql(`where "a" outside (10,19)`, { data: input })).toStrictEqual([{ a: 1 }, { a: 20 }, { a: 30 }]);
      expect(await uql(`where "a" outside (11,19)`, { data: input })).toStrictEqual([{ a: 1 }, { a: 10 }, { a: 20 }, { a: 30 }]);
      expect(await uql(`where "a" outside (11,20)`, { data: input })).toStrictEqual([{ a: 1 }, { a: 10 }, { a: 30 }]);
      expect(await uql(`where "a" outside (10)`, { data: input })).toStrictEqual([]); // incorrect arguments
    });
    it("matches regex", async () => {
      const input = [{ a: "KANSAS" }, { a: "ARKANSAS" }, { a: "LAKE SUPERIOR" }, { a: "Kansas" }, { a: "alabama" }, { a: "LAKE ST CLAIR" }];
      expect(await uql(`where "a" matches regex 'K.*S'`, { data: input })).toStrictEqual([{ a: "KANSAS" }, { a: "ARKANSAS" }, { a: "LAKE SUPERIOR" }, { a: "LAKE ST CLAIR" }]);
      expect(await uql(`where "a" matches regex 'K.*s'`, { data: input })).toStrictEqual([{ a: "Kansas" }]);
      expect(await uql(`where "a" matches regex 'K.*[sS]'`, { data: input })).toStrictEqual([{ a: "KANSAS" }, { a: "ARKANSAS" }, { a: "LAKE SUPERIOR" }, { a: "Kansas" }, { a: "LAKE ST CLAIR" }]);
      expect(await uql(`where "a" matches regex '^K.*[sS]$'`, { data: input })).toStrictEqual([{ a: "KANSAS" }, { a: "Kansas" }]);
    });
    it("!matches regex", async () => {
      const input = [{ a: "KANSAS" }, { a: "ARKANSAS" }, { a: "LAKE SUPERIOR" }, { a: "Kansas" }, { a: "alabama" }, { a: "LAKE ST CLAIR" }];
      expect(await uql(`where "a" !matches regex 'K.*S'`, { data: input })).toStrictEqual([{ a: "Kansas" }, { a: "alabama" }]);
      expect(await uql(`where "a" !matches regex 'K.*s'`, { data: input })).toStrictEqual([{ a: "KANSAS" }, { a: "ARKANSAS" }, { a: "LAKE SUPERIOR" }, { a: "alabama" }, { a: "LAKE ST CLAIR" }]);
      expect(await uql(`where "a" !matches regex 'K.*[sS]'`, { data: input })).toStrictEqual([{ a: "alabama" }]);
      expect(await uql(`where "a" !matches regex '^K.*[sS]$'`, { data: input })).toStrictEqual([{ a: "ARKANSAS" }, { a: "LAKE SUPERIOR" }, { a: "alabama" }, { a: "LAKE ST CLAIR" }]);
      expect(await uql(`where "a" not matches regex '^K.*[sS]$'`, { data: input })).toStrictEqual([{ a: "ARKANSAS" }, { a: "LAKE SUPERIOR" }, { a: "alabama" }, { a: "LAKE ST CLAIR" }]);
    });
  });
});
