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
});
