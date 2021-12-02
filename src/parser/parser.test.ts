import { uql } from "./index";

const sample_data = {
  users: [
    { patron: "jonny", age: 48 },
    { patron: "kamal", age: 34 },
    { patron: "john", age: 12 },
    { patron: "john", age: 40 },
    { patron: "jonny", age: 36 },
  ],
};

describe("parser", () => {
  describe("hello", () => {
    it("hello", async () => {
      const result = await uql("hello");
      expect(result).toStrictEqual("hello");
    });
  });
  describe("ping", () => {
    it("ping", async () => {
      const result = await uql("ping");
      expect(result).toStrictEqual("pong");
    });
  });
  describe("echo", () => {
    it("echo", async () => {
      const result = await uql(`echo "hi there!"`);
      expect(result).toStrictEqual("hi there!");
    });
  });
  describe("count", () => {
    it("count string", async () => {
      const result = await uql(`count`, { data: "vibgyor" });
      expect(result).toStrictEqual(7);
    });
    it("count number", async () => {
      const result = await uql(`count`, { data: 123 });
      expect(result).toStrictEqual(123);
    });
    it("count array", async () => {
      const result = await uql(`count`, { data: [1, 2, 3] });
      expect(result).toStrictEqual(3);
    });
  });
  describe("limit", () => {
    it("limit string", async () => {
      const result = await uql(`limit 2`, { data: "vibgyor" });
      expect(result).toStrictEqual("vi");
    });
    it("limit number", async () => {
      const result = await uql(`limit 2`, { data: 123 });
      expect(result).toStrictEqual(123);
    });
    it("limit array", async () => {
      const result = await uql(`limit 2`, { data: [1, 2, 3] });
      expect(result).toStrictEqual([1, 2]);
    });
  });
  describe("orderby", () => {
    it("string", async () => {
      const result = await uql(`order by "name" asc`, { data: "hello" });
      expect(result).toStrictEqual("hello");
    });
    it("number", async () => {
      const result = await uql(`order by "name" asc`, { data: 312 });
      expect(result).toStrictEqual(312);
    });
    it("simple asc", async () => {
      const result = await uql(`order by "name" asc`, { data: [{ name: "foo" }, { name: "bar" }] });
      expect(result).toStrictEqual([{ name: "bar" }, { name: "foo" }]);
    });
    it("simple desc", async () => {
      const result = await uql(`order by "name" desc`, { data: [{ name: "foo" }, { name: "bar" }] });
      expect(result).toStrictEqual([{ name: "foo" }, { name: "bar" }]);
    });
    it("multi order", async () => {
      const result = await uql(`order by "patron" desc,"age" asc`, { data: sample_data.users });
      expect(result).toStrictEqual([
        { patron: "kamal", age: 34 },
        { patron: "jonny", age: 36 },
        { patron: "jonny", age: 48 },
        { patron: "john", age: 12 },
        { patron: "john", age: 40 },
      ]);
    });
  });
  describe("project", () => {
    it("basic", async () => {
      const result = await uql(`project "bar"`, {
        data: [
          { foo: "foo1", bar: "bar1", baz: "baz1" },
          { foo: "foo2", bar: "bar2", baz: "baz2" },
        ],
      });
      expect(result).toStrictEqual([{ bar: "bar1" }, { bar: "bar2" }]);
    });
    it("multi", async () => {
      const result = await uql(`project "bar" , "foo" `, {
        data: [
          { foo: "foo1", bar: "bar1", baz: "baz1" },
          { foo: "foo2", bar: "bar2", baz: "baz2" },
        ],
      });
      expect(result).toStrictEqual([
        { bar: "bar1", foo: "foo1" },
        { bar: "bar2", foo: "foo2" },
      ]);
    });
    it("multi with pipe", async () => {
      const result = await uql(`limit 1 | project "bar" , "foo"`, {
        data: [
          { foo: "foo1", bar: "bar1", baz: "baz1" },
          { foo: "foo2", bar: "bar2", baz: "baz2" },
        ],
      });
      expect(result).toStrictEqual([{ foo: "foo1", bar: "bar1" }]);
    });
  });
  describe("project-away", () => {
    it("basic", async () => {
      const result = await uql(`project-away "bar"`, {
        data: [
          { foo: "foo1", bar: "bar1", baz: "baz1" },
          { foo: "foo2", bar: "bar2", baz: "baz2" },
        ],
      });
      expect(result).toStrictEqual([
        { baz: "baz1", foo: "foo1" },
        { baz: "baz2", foo: "foo2" },
      ]);
    });
    it("multi", async () => {
      const result = await uql(`project-away "bar" , "foo" `, {
        data: [
          { foo: "foo1", bar: "bar1", baz: "baz1" },
          { foo: "foo2", bar: "bar2", baz: "baz2" },
        ],
      });
      expect(result).toStrictEqual([{ baz: "baz1" }, { baz: "baz2" }]);
    });
    it("multi with pipe", async () => {
      const result = await uql(`limit 1 | project-away "bar" , "foo"`, {
        data: [
          { foo: "foo1", bar: "bar1", baz: "baz1" },
          { foo: "foo2", bar: "bar2", baz: "baz2" },
        ],
      });
      expect(result).toStrictEqual([{ baz: "baz1" }]);
    });
  });
  describe("extend", () => {
    it("basic", async () => {
      const result = await uql(`extend "three"=sum(1,2)`, { data: [{}] });
      expect(result).toStrictEqual([{ three: 3 }]);
    });
    it("with ref field", async () => {
      const result = await uql(`extend "c"=sum("a","b")`, { data: [{ a: 1, b: 2 }] });
      expect(result).toStrictEqual([{ a: 1, b: 2, c: 3 }]);
    });
    it("mixed without alias", async () => {
      const result = await uql(`extend sum( "a" , "b" , 3 ) `, { data: [{ a: 1, b: 2 }] });
      expect(result).toStrictEqual([{ a: 1, b: 2, sum: 6 }]);
    });
    it("strcat", async () => {
      const result = await uql(`extend "message"=strcat(str("hello"),str(" "),str("world"))`, { data: [{}] });
      expect(result).toStrictEqual([{ message: "hello world" }]);
    });
    it("strcat with single quoted string", async () => {
      const result = await uql(`extend "message"=strcat('hello',' ','world')`, { data: [{}] });
      expect(result).toStrictEqual([{ message: "hello world" }]);
    });
    it("strcat with ref fields", async () => {
      const result = await uql(`extend "message"=strcat("greeting",str(" "),"person",'!') | project-away "greeting", "person"`, { data: [{ greeting: "hello", person: "world" }] });
      expect(result).toStrictEqual([{ message: "hello world!" }]);
    });
    it("tolower", async () => {
      const result = await uql(`extend "message"=tolower(str("HELLO"))`, { data: [{}] });
      expect(result).toStrictEqual([{ message: "hello" }]);
    });
    it("strlen", async () => {
      const result = await uql(`extend "len"=strlen(str("HELLO")), "len1"=strlen("name") | project-away "name" `, { data: [{ name: "foo" }] });
      expect(result).toStrictEqual([{ len: 5, len1: 3 }]);
    });
  });
  describe("summarize", () => {
    it("count", async () => {
      const result = await uql(`summarize count()`, { data: [{}, {}, {}] });
      expect(result).toStrictEqual({ count: 3 });
    });
    it("count with alias", async () => {
      const result = await uql(`summarize "total value"=count() `, { data: [{}, {}, {}] });
      expect(result).toStrictEqual({ "total value": 3 });
    });
    it("multiple metrics", async () => {
      const result = await uql(`summarize "total value"=count() , min("age")`, { data: [{ age: 12 }, { age: 3 }, { age: 30 }] });
      expect(result).toStrictEqual({ "total value": 3, "age (min)": 3 });
    });
    it("multiple metrics by field", async () => {
      const result = await uql(`summarize max("age") by "city"`, {
        data: [
          { age: 12, city: "london" },
          { age: 3, city: "london" },
          { age: 30, city: "tokyo" },
        ],
      });
      expect(result).toStrictEqual([
        { city: "london", "age (max)": 12 },
        { city: "tokyo", "age (max)": 30 },
      ]);
    });
    it("dcount", async () => {
      const result = await uql(`summarize dcount("region")`, {
        data: [
          { cat: "one", region: "EU" },
          { cat: "one", region: "APAC" },
          { cat: "two", region: "EU" },
          { cat: "two", region: "NASA" },
          { cat: "two", region: "NASA" },
        ],
      });
      expect(result).toStrictEqual({ "region (dcount)": 3 });
    });
    it("dcount by", async () => {
      const result = await uql(`summarize "regions"=dcount("region") by "cat"`, {
        data: [
          { cat: "one", region: "EU" },
          { cat: "one", region: "APAC" },
          { cat: "one", region: "NASA" },
          { cat: "two", region: "EU" },
          { cat: "two", region: "NASA" },
          { cat: "two", region: "NASA" },
        ],
      });
      expect(result).toStrictEqual([
        { cat: "one", regions: 3 },
        { cat: "two", regions: 2 },
      ]);
    });
  });
  describe("multi", () => {
    it("limit and count", async () => {
      const result = await uql("limit 3 | count", { data: sample_data.users });
      expect(result).toStrictEqual(3);
    });
    it("limit and count with multi line", async () => {
      const result = await uql("limit 3 \n | count ", { data: sample_data.users });
      expect(result).toStrictEqual(3);
    });
    it("with summarize", async () => {
      const result = await uql(`summarize mean("age") by "city" | order by "age (mean)" asc`, {
        data: [
          { age: 12, city: "london" },
          { age: 3, city: "london" },
          { age: 30, city: "tokyo" },
        ],
      });
      expect(result).toStrictEqual([
        { city: "london", "age (mean)": 7.5 },
        { city: "tokyo", "age (mean)": 30 },
      ]);
    });
    it("with summarize with alias", async () => {
      const result = await uql(`summarize "city avg age"=mean("age") by "city" | order by "city avg age" asc`, {
        data: [
          { age: 12, city: "london" },
          { age: 3, city: "london" },
          { age: 30, city: "tokyo" },
        ],
      });
      expect(result).toStrictEqual([
        { city: "london", "city avg age": 7.5 },
        { city: "tokyo", "city avg age": 30 },
      ]);
    });
  });
});
