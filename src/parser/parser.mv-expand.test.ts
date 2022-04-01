import { uql } from "./index";

describe("mv-expand", () => {
  it("should work for basic array", async () => {
    const result = await uql(`mv-expand "users"`, {
      data: [
        { group: "A", users: ["user a1", "user a2"] },
        { group: "B", users: ["user b1"] },
      ],
    });
    expect(result).toStrictEqual([
      { group: "A", users: "user a1" },
      { group: "A", users: "user a2" },
      { group: "B", users: "user b1" },
    ]);
  });
  it("should work for basic array with alias", async () => {
    const result = await uql(`mv-expand "user"="users"`, {
      data: [
        { group: "A", users: ["user a1", "user a2"] },
        { group: "B", users: ["user b1"] },
      ],
    });
    expect(result).toStrictEqual([
      { group: "A", user: "user a1" },
      { group: "A", user: "user a2" },
      { group: "B", user: "user b1" },
    ]);
  });
  it("should ignore non array values", async () => {
    const result = await uql(`mv-expand "user"="users"`, {
      data: [
        { group: "A", users: ["user a1", "user a2"] },
        { group: "B", users: [] },
        { group: "C", users: ["user c1"] },
        { group: "D" },
        { group: "E", users: undefined },
        { group: "F", users: null },
        { group: "G", users: "not found" },
        { group: "H", users: ["user h1", undefined] },
      ],
    });
    expect(result).toStrictEqual([
      { group: "A", user: "user a1" },
      { group: "A", user: "user a2" },
      { group: "C", user: "user c1" },
      { group: "H", user: "user h1" },
      { group: "H", user: undefined },
    ]);
  });
  it("should respect array of objects", async () => {
    const result = await uql(`mv-expand "user"="users" | summarize "total_age"=sum("user.age"),count()`, {
      data: [
        {
          group: "A",
          users: [
            { username: "a1", age: 2 },
            { username: "a2", age: 3 },
          ],
        },
        {
          group: "B",
          users: [
            { username: "b1", age: 2 },
            { username: "b1", age: undefined },
            { username: "b2", age: 3 },
          ],
        },
      ],
    });
    expect(result).toStrictEqual({ total_age: 10, count: 5 });
  });
});
