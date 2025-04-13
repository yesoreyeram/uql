import { uql } from "../index";

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
  describe("date", () => {
    it("default", async () => {
      const result = await uql(`extend "in"=todatetime("in") | extend "in"=format_datetime("in",'DD/MM/YYYY')  | project "in"`, { data: [{ in: "1990-02-27" }] });
      expect(result).toStrictEqual([{ in: "27/02/1990" }]);
    });
    it("format_datetime", async () => {
      const result = await uql(`extend "in"=todatetime("in") | extend "in"=format_datetime("in",'dddd')  | project "in"`, { data: [{ in: "1990-02-27" }] });
      expect(result).toStrictEqual([{ in: "Tuesday" }]);
    });
    it("add_datetime", async () => {
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=add_datetime("in",'1d')  | project "in"`, { data: [{ in: "1990-02-27" }] })).toStrictEqual([{ in: new Date("1990-02-28") }]);
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=add_datetime("in",'-1d')  | project "in"`, { data: [{ in: "1990-02-27" }] })).toStrictEqual([{ in: new Date("1990-02-26") }]);
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=add_datetime("in",'-1h')  | project "in"`, { data: [{ in: "1990-02-27" }] })).toStrictEqual([
        { in: new Date("1990-02-26 23:00:00") },
      ]);
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=add_datetime("in",'-1y')  | project "in"`, { data: [{ in: "1990-02-27" }] })).toStrictEqual([
        { in: new Date("1989-02-27 00:00:00") },
      ]);
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=add_datetime("in",'10y')  | project "in"`, { data: [{ in: "1990-02-27" }] })).toStrictEqual([
        { in: new Date("2000-02-27 00:00:00") },
      ]);
    });
    it("start of", async () => {
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=startofyear("in")  | project "in"`, { data: [{ in: "1990-02-27 13:23:56.876" }] })).toStrictEqual([
        { in: new Date("1990-01-01 00:00:00.000") },
      ]);
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=startofmonth("in")  | project "in"`, { data: [{ in: "1990-02-27 13:23:56.876" }] })).toStrictEqual([
        { in: new Date("1990-02-01 00:00:00.000") },
      ]);
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=startofweek("in")  | project "in"`, { data: [{ in: "1990-02-27 13:23:56.876" }] })).toStrictEqual([
        { in: new Date("1990-02-25 00:00:00.000") },
      ]);
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=startofday("in")  | project "in"`, { data: [{ in: "1990-02-27 13:23:56.876" }] })).toStrictEqual([
        { in: new Date("1990-02-27 00:00:00.000") },
      ]);
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=startofhour("in")  | project "in"`, { data: [{ in: "1990-02-27 13:23:56.876" }] })).toStrictEqual([
        { in: new Date("1990-02-27 13:00:00.000") },
      ]);
      expect(await uql(`extend "in"=todatetime("in") | extend "in"=startofminute("in")  | project "in"`, { data: [{ in: "1990-02-27 13:23:56.876" }] })).toStrictEqual([
        { in: new Date("1990-02-27 13:23:00.000") },
      ]);
    });
  });
  describe("basic", () => {
    it("default", async () => {
      const result3 = await uql(`project-away "b" | extend "a1"="a", "c"`, {
        data: [
          { a: 12, b: 20, c: 2 },
          { a: 6, b: 32, c: 3 },
        ],
      });
      expect(result3).toStrictEqual([
        { a: 12, a1: 12, c: 2 },
        { a: 6, a1: 6, c: 3 },
      ]);
    });
    it("advanced", async () => {
      const xml_data = `<nmapresult>
      <hosts>
          <host name="host1">
              <ports>
                  <port name="http">80</port>
                  <port name="https">443</port>
              </ports>
           </host>
          <host name="host2">
              <ports>
                  <port name="http">80</port>
                  <port name="ssh">22</port>
              </ports>
          </host>
      </hosts>
      </nmapresult>`;
      const result4 = await uql(
        `parse-xml 
      | scope "nmapresult.hosts.host" 
      | project "host"="@_name", "ports"="ports.port"
      | mv-expand "port"="ports" 
      | extend "portName"="port.@_name", "portNumber"="port.#text"
      | project-away "port"`,
        { data: xml_data }
      );
      expect(result4).toStrictEqual([
        { host: "host1", portNumber: 80, portName: "http" },
        { host: "host1", portNumber: 443, portName: "https" },
        { host: "host2", portNumber: 80, portName: "http" },
        { host: "host2", portNumber: 22, portName: "ssh" },
      ]);
    });
  });
});
