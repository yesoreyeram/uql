import { uql } from "../index";

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
    describe("summarize conditional", () => {
      const users = [
        { username: "a", age: 48, country: "foo" },
        { username: "b", age: 34, country: "foo" },
        { username: "c", age: 12, country: "bar" },
        { username: "d", age: 10, country: "bar" },
        { username: "e", age: 36, country: "baz" },
      ];
      it("default", async () => {
        expect(await uql(`summarize countif("age","age" >= 18)`, { data: users })).toStrictEqual({ "age (countif)": 3 });
        expect(await uql(`summarize "eligible_voters"=countif("age","age" >= 18)`, { data: users })).toStrictEqual({ eligible_voters: 3 });
        expect(await uql(`summarize "eligible_voters"=countif("age","age" >= 18)`, { data: users })).toStrictEqual({ eligible_voters: 3 });
        expect(await uql(`summarize "eligible_voters"=countif("age","age" >= 18), "non_eligible_voters"=countif("age","age" < 18)`, { data: users })).toStrictEqual({
          eligible_voters: 3,
          non_eligible_voters: 2,
        });
        expect(await uql(`summarize "eligible_voters"=countif("age","age" >= 18), "non_eligible_voters"=countif("age","age" < 18) by "country"`, { data: users })).toStrictEqual([
          { country: "foo", eligible_voters: 2, non_eligible_voters: 0 },
          { country: "bar", eligible_voters: 0, non_eligible_voters: 2 },
          { country: "baz", eligible_voters: 1, non_eligible_voters: 0 },
        ]);
      });
    });
  });
});
