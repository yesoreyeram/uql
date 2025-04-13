import { uql } from "../index";

describe("pivot", () => {
  const data = {
    fruits: [
      { fruit: "apple", size: "sm", qty: 1 },
      { fruit: "apple", size: "md", qty: 2 },
      { fruit: "apple", size: "lg", qty: 3 },
      { fruit: "banana", size: "sm", qty: 1 },
      { fruit: "banana", size: "lg", qty: 6 },
      { fruit: "banana", size: "xl", qty: 5 },
    ],
  };
  describe("basic", () => {
    it("with default arguments", async () => {
      expect(await uql(`pivot count()`, { data: data.fruits })).toStrictEqual(6);
    });
    it("with custom operator", async () => {
      expect(await uql(`pivot sum("qty")`, { data: data.fruits })).toStrictEqual(18);
      expect(await uql(`pivot max("qty")`, { data: data.fruits })).toStrictEqual(6);
    });
    it("with custom operator with row", async () => {
      expect(await uql(`pivot sum("qty"), "fruit"`, { data: data.fruits })).toStrictEqual([
        { fruit: "apple", value: 6 },
        { fruit: "banana", value: 12 },
      ]);
      expect(await uql(`pivot count("qty"), "size"`, { data: data.fruits })).toStrictEqual([
        { size: "sm", value: 2 },
        { size: "md", value: 1 },
        { size: "lg", value: 2 },
        { size: "xl", value: 1 },
      ]);
    });
    it("with custom operator with row and col", async () => {
      expect(await uql(`pivot sum("qty"), "fruit", "size"`, { data: data.fruits })).toStrictEqual([
        { fruit: "apple", sm: 1, md: 2, lg: 3, xl: 0 },
        { fruit: "banana", sm: 1, md: 0, lg: 6, xl: 5 },
      ]);
      expect(await uql(`pivot max("qty"), "fruit", "size"`, { data: data.fruits })).toStrictEqual([
        { fruit: "apple", sm: 1, md: 2, lg: 3, xl: null },
        { fruit: "banana", sm: 1, md: null, lg: 6, xl: 5 },
      ]);
    });
  });
});
