import { uql } from "./index";

const sample_data = {
  users: [
    { patron: "jonny", age: 48, name: { f: "foo1", l: "bar1" } },
    { patron: "kamal", age: 34, name: { f: "foo2", l: "bar2" } },
    { patron: "john", age: 12, name: { f: "foo3", l: "bar3" } },
    { patron: "john", age: 40, name: { f: "foo4", l: "bar4" } },
    { patron: "jonny", age: 36, name: { f: "foo5", l: "bar5" } },
  ],
};

describe("project", () => {
  describe("basic", () => {
    it("default", async () => {
      const result = await uql(`project "age"`, { data: sample_data.users });
      expect(result).toStrictEqual([{ age: 48 }, { age: 34 }, { age: 12 }, { age: 40 }, { age: 36 }]);
    });
    it("nested", async () => {
      const result = await uql(`project "name.f"`, { data: sample_data.users });
      expect(result).toStrictEqual([{ name: { f: "foo1" } }, { name: { f: "foo2" } }, { name: { f: "foo3" } }, { name: { f: "foo4" } }, { name: { f: "foo5" } }]);
    });
    it("nested with alias", async () => {
      const result = await uql(`project "firstName"="name.f"`, { data: sample_data.users });
      expect(result).toStrictEqual([{ firstName: "foo1" }, { firstName: "foo2" }, { firstName: "foo3" }, { firstName: "foo4" }, { firstName: "foo5" }]);
    });
    it("multiple nested with alias", async () => {
      const result = await uql(`project "firstName"="name.f", "lastName"="name.l", "age"`, { data: sample_data.users });
      expect(result).toStrictEqual([
        { firstName: "foo1", lastName: "bar1", age: 48 },
        { firstName: "foo2", lastName: "bar2", age: 34 },
        { firstName: "foo3", lastName: "bar3", age: 12 },
        { firstName: "foo4", lastName: "bar4", age: 40 },
        { firstName: "foo5", lastName: "bar5", age: 36 },
      ]);
    });
  });
  describe("object", () => {
    it("single", async () => {
      const result = await uql(`project "name"`, { data: { name: "sriramajeyam", score: 12, languages: ["tamil", "english", "sourashtra"] } });
      expect(result).toStrictEqual("sriramajeyam");
    });
    it("single with alias", async () => {
      const result = await uql(`project "full name"="name"`, { data: { name: "sriramajeyam", score: 12, languages: ["tamil", "english", "sourashtra"] } });
      expect(result).toStrictEqual("sriramajeyam");
    });
    it("single", async () => {
      const result = await uql(`project "languages"`, { data: { name: "sriramajeyam", score: 12, languages: ["tamil", "english", "sourashtra"] } });
      expect(result).toStrictEqual(["tamil", "english", "sourashtra"]);
    });
    it("multiple", async () => {
      const result = await uql(`project "name", "score"`, { data: { name: "sriramajeyam", score: 12, languages: ["tamil", "english", "sourashtra"] } });
      expect(result).toStrictEqual({ name: "sriramajeyam", score: 12 });
    });
  });
  describe("math functions", () => {
    it("should calculate the floor/ceil/round/sign/pow correctly", async () => {
      const result = await uql(`extend "one"=ceil(1.2), "two"=floor(1.2), "three"=round(1.2), "four"=round(1.6), "five"=sign(-1.2), "six"=pow(2,3)`, { data: [{}] });
      expect(result).toStrictEqual([
        {
          one: 2,
          two: 1,
          three: 1,
          four: 2,
          five: -1,
          six: 8,
        },
      ]);
    });
    it("should calculate the floor/ceil/round/sign/pow correctly with ref fields", async () => {
      const result = await uql(`project "one"=ceil("a"), "two"=floor("a"), "three"=round("a"), "four"=round("b"), "five"=sign("e"), "six"=pow("c","d")`, {
        data: [{ a: 1.2, b: 1.6, c: 2, d: 3, e: -1.2 }],
      });
      expect(result).toStrictEqual([
        {
          one: 2,
          two: 1,
          three: 1,
          four: 2,
          five: -1,
          six: 8,
        },
      ]);
    });
    it("should calculate the log/log2/log10 calculations correctly", async () => {
      const result = await uql(`extend "one"=log(1.2), "two"=log2(1.2), "three"=log10(1.2)`, { data: [{}] });
      expect(result).toStrictEqual([
        {
          one: Math.log(1.2),
          two: Math.log2(1.2),
          three: Math.log10(1.2),
        },
      ]);
    });
    it("should calculate the log/log2/log10 calculations correctly with ref fields", async () => {
      const result = await uql(`project "one"=log("a"), "two"=log2("a"), "three"=log10("a")`, { data: [{ a: 1.2 }] });
      expect(result).toStrictEqual([
        {
          one: Math.log(1.2),
          two: Math.log2(1.2),
          three: Math.log10(1.2),
        },
      ]);
    });
    it("should calculate the sin/cos/tan calculations correctly", async () => {
      const result = await uql(`extend "one"=sin(1.2), "two"=cos(1.2), "three"=tan(1.2)`, { data: [{}] });
      expect(result).toStrictEqual([
        {
          one: Math.sin(1.2),
          two: Math.cos(1.2),
          three: Math.tan(1.2),
        },
      ]);
    });
    it("should calculate the sin/cos/tan calculations correctly with ref fields", async () => {
      const result = await uql(`project "one"=sin("a"), "two"=cos("a"), "three"=tan("a")`, { data: [{ a: 1.2 }] });
      expect(result).toStrictEqual([
        {
          one: Math.sin(1.2),
          two: Math.cos(1.2),
          three: Math.tan(1.2),
        },
      ]);
    });
  });
});
