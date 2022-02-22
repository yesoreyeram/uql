import { uql } from "./index";

describe("parser", () => {
  describe("summarize", () => {
    describe("summarize - single", () => {
      const data = {
        users: [
          { patron: "a", age: 48, country: "foo" },
          { patron: "b", age: 34, country: "foo" },
          { patron: "c", age: 12, country: "bar" },
          { patron: "d", age: 40, country: "bar" },
          { patron: "e", age: 36, country: "baz" },
        ],
      };
      it("first", async () => {
        const result = await uql(`summarize "user"=first("patron") by "country"`, { data: data.users });
        expect(result).toStrictEqual([
          { country: "foo", user: "a" },
          { country: "bar", user: "c" },
          { country: "baz", user: "e" },
        ]);
      });
      it("last", async () => {
        const result = await uql(`summarize "user"=last("patron") by "country"`, { data: data.users });
        expect(result).toStrictEqual([
          { country: "foo", user: "b" },
          { country: "bar", user: "d" },
          { country: "baz", user: "e" },
        ]);
      });
      it("latest", async () => {
        const result = await uql(`summarize "user"=latest("patron") by "country"`, { data: data.users });
        expect(result).toStrictEqual([
          { country: "foo", user: "b" },
          { country: "bar", user: "d" },
          { country: "baz", user: "e" },
        ]);
      });
    });
    describe("summarize - multi", () => {
      it("summarize by two props", async () => {
        const result = await uql(`summarize "age"=sum("age") by "country", "city"`, {
          data: [
            { age: 1, name: "foo1", city: "chennai", country: "india" },
            { age: 2, name: "foo1", city: "chennai", country: "india" },
            { age: 3, name: "foo1", city: "mumbai", country: "india" },
            { age: 4, name: "foo1", city: "london", country: "england" },
          ],
        });
        expect(result).toStrictEqual([
          { country: "india", city: "chennai", age: 3 },
          { country: "india", city: "mumbai", age: 3 },
          { country: "england", city: "london", age: 4 },
        ]);
      });
      it("summarize by two props and two values", async () => {
        const result = await uql(`summarize "age"=sum("age"),"minAge"=min("age") by "country", "city"`, {
          data: [
            { age: 1, name: "foo1", city: "chennai", country: "india" },
            { age: 2, name: "foo1", city: "chennai", country: "india" },
            { age: 3, name: "foo1", city: "mumbai", country: "india" },
            { age: 4, name: "foo1", city: "london", country: "england" },
          ],
        });
        expect(result).toStrictEqual([
          { country: "india", city: "chennai", age: 3, minAge: 1 },
          { country: "india", city: "mumbai", age: 3, minAge: 3 },
          { country: "england", city: "london", age: 4, minAge: 4 },
        ]);
      });
    });
  });
});
