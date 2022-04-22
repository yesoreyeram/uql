import { uql } from "../index";

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
  it("should respect subsequent commands", async () => {
    const data = {
      server1: {
        name: "server-1",
        ip: "0.0.0.1",
        disks: [
          { drive: "C", size: 1024 },
          { drive: "D", size: 2048 },
        ],
      },
      server2: {
        name: "server-2",
        ip: "0.0.0.2",
        disks: [
          { drive: "C", size: 2048 },
          { drive: "D", size: 3096 },
        ],
      },
    };
    const result = await uql(
      `parse-json
    | project kv() 
    | project "name"="value.name", "ip"="value.ip", "disks"="value.disks"
    | mv-expand "disk"="disks" 
    | project "name", "ip", "driveName"="disk.drive", "driveSize"="disk.size"
    | project-away "disk"`,
      { data }
    );
    expect(result).toStrictEqual([
      { name: "server-1", ip: "0.0.0.1", driveName: "C", driveSize: 1024 },
      { name: "server-1", ip: "0.0.0.1", driveName: "D", driveSize: 2048 },
      { name: "server-2", ip: "0.0.0.2", driveName: "C", driveSize: 2048 },
      { name: "server-2", ip: "0.0.0.2", driveName: "D", driveSize: 3096 },
    ]);
  });
});
