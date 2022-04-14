import { uql } from "./index";

const sample_data = {
  users: [
    { patron: "jonny", patrons: ["jonny", "foo", "bar", "foo"], age: 48, name: { f: "foo1", l: "bar1" } },
    { patron: "kamal", patrons: ["kamal", "foo", "bar", "foo"], age: 34, name: { f: "foo2", l: "bar2" } },
    { patron: "john", patrons: ["john", "foo", "bar", "foo"], age: 48, name: { f: "foo3", l: "bar3" } },
    { patron: "john", patrons: ["john", "foo", "bar", "foo"], age: 40, name: { f: "foo4", l: "bar4" } },
    { patron: "jonny", patrons: ["jonny", "foo", "bar", "foo"], age: 36, name: { f: "foo5", l: "bar5" } },
  ],
};

describe("project", () => {
  describe("distinct", () => {
    it("default", async () => {
      const result = await uql(`project "users"=distinct("patrons")`, { data: sample_data.users });
      expect(result).toStrictEqual([
        { users: ["jonny", "foo", "bar"] },
        { users: ["kamal", "foo", "bar"] },
        { users: ["john", "foo", "bar"] },
        { users: ["john", "foo", "bar"] },
        { users: ["jonny", "foo", "bar"] },
      ]);
    });
  });
});

describe("distinct", () => {
  it("default", async () => {
    const result = await uql(`distinct`, { data: ["a", "b", "c", "a", "d"] });
    expect(result).toStrictEqual(["a", "b", "c", "d"]);
  });
  it("default with args", async () => {
    const result = await uql(`distinct "age"`, { data: sample_data.users });
    expect(result).toStrictEqual([48, 34, 40, 36]);
  });
  it("default with args", async () => {
    const result = await uql(`distinct "patrons"`, { data: sample_data.users[0] });
    expect(result).toStrictEqual(["jonny", "foo", "bar"]);
  });
  it("default with args", async () => {
    const result = await uql(`distinct "patron"`, { data: sample_data.users });
    expect(result).toStrictEqual(["jonny", "kamal", "john"]);
  });
});
