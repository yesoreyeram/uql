import { uql } from "../index";

describe("project-reorder", () => {
  it("default", async () => {
    const result = await uql(`project-reorder "b", "c", "a"`, { data: [{ a: "a1", b: "b1", c: "c1" }] });
    expect(JSON.stringify(result)).toStrictEqual('[{"b":"b1","c":"c1","a":"a1"}]');
    const result1 = await uql(`project-reorder "a", "c", "b"`, { data: [{ a: "a1", b: "b1", c: "c1" }] });
    expect(JSON.stringify(result1)).toStrictEqual('[{"a":"a1","c":"c1","b":"b1"}]');
  });
  it("reorder with alias", async () => {
    const result = await uql(`project-reorder "bee"="b", "c", "a"`, { data: [{ a: "a1", b: "b1", c: "c1" }] });
    expect(JSON.stringify(result)).toStrictEqual('[{"bee":"b1","c":"c1","a":"a1"}]');
    const result1 = await uql(`project-reorder "a", "cee"="c", "b"`, { data: [{ a: "a1", b: "b1", c: "c1" }] });
    expect(JSON.stringify(result1)).toStrictEqual('[{"a":"a1","cee":"c1","b":"b1"}]');
  });
});
