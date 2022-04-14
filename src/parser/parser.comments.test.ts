import { uql } from "./index";

describe("parser", () => {
  describe("comments", () => {
    it("comments at line", async () => {
      expect(await uql(`hello#world`)).toStrictEqual("hello");
      expect(await uql(`hello #world`)).toStrictEqual("hello");
      expect(await uql(`hello# world `)).toStrictEqual("hello");
      expect(await uql(`hello # world `)).toStrictEqual("hello");
    });
    it("comments at the beginning of new line", async () => {
      expect(await uql("hello#world\n| # count")).toStrictEqual("hello");
      expect(await uql("hello #world\n|# count")).toStrictEqual("hello");
      expect(await uql(`hello# world\n| # count `)).toStrictEqual("hello");
      expect(await uql(`hello # world\n| # count `)).toStrictEqual("hello");
      expect(await uql(`hello # world\n| count `)).toStrictEqual(5);
    });
    it("comments at the beginning in new line windows style", async () => {
      expect(await uql("hello#world\r\n| # count")).toStrictEqual("hello");
      expect(await uql("hello #world\r\n|# count")).toStrictEqual("hello");
      expect(await uql(`hello# world\r\n| # count `)).toStrictEqual("hello");
      expect(await uql(`hello # world\r\n| # count `)).toStrictEqual("hello");
      expect(await uql(`hello # world\r\n|  count `)).toStrictEqual(5);
    });
    it("comment at both the places and in between", async () => {
      expect(
        await uql(`project "newName"=strcat("name",'--#--',"name") # something else | # count | order by "newName" desc`, { data: [{ name: "foo" }, { name: "bar" }, { name: "baz" }] })
      ).toStrictEqual([{ newName: "foo--#--foo" }, { newName: "bar--#--bar" }, { newName: "baz--#--baz" }]);
    });
  });
});
