import { uql } from "./index";

const sample_data = {
  values: [
    { a: 12, b: 20, c: 2 },
    { a: 6, b: 32, c: 3 },
  ],
};

describe("extend", () => {
  describe("maths", () => {
    it("default", async () => {
      const result = await uql(`project "a", "foo"=mul(2,3,4,5,"c"), "triple"=sum("a","a","a"),"thrice"=mul("a",3), sum("a","b",10,"c"),  diff("a","b"), mul("a","b"), div("b",2)`, {
        data: sample_data.values,
      });
      expect(result).toStrictEqual([
        { a: 12, diff: -8, mul: 240, sum: 44, thrice: 36, triple: 36, div: 10, foo: 240 },
        { a: 6, diff: -26, mul: 192, sum: 51, thrice: 18, triple: 18, div: 16, foo: 360 },
      ]);
    });
    describe("date math", () => {
      it("default", async () => {
        const result = await uql(`extend "start"=todatetime("start"),"end"=todatetime("end"),"elapsed"=diff("end","start")  | project "elapsed"`, {
          data: [{ start: "2021-01-01", end: "2021-01-02" }],
        });
        expect((result as { elapsed: number }[])[0].elapsed).toStrictEqual(86400000);
      });
    });
  });
  describe("random", () => {
    it("default", async () => {
      const result = await uql(`extend "a"=random(1.2,5.4) | project "a"`, { data: sample_data.values });
      expect((result as { a: unknown }[])[0].a).toBeGreaterThanOrEqual(1.2);
      expect((result as { a: unknown }[])[0].a).toBeLessThanOrEqual(5.4);
    });
    it("without args", async () => {
      const result = await uql(`extend "a"=random() | project "a"`, { data: sample_data.values });
      const value = (result as { a: unknown }[])[0].a;
      expect(value === 1 || value === 0).toBeTruthy();
    });
    it("with all args", async () => {
      const result = await uql(`extend "a"=random(5,10,true) | project "a"`, { data: sample_data.values });
      const value = (result as { a: unknown }[])[0].a;
      expect(value).toBeGreaterThanOrEqual(5);
      expect(value).toBeLessThanOrEqual(10);
    });
  });
});
