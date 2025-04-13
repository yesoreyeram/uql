import { uql } from "../parser/index";
import data from "./data.json";

describe("uql", () => {
  it("default", async () => {
    expect(
      await uql(
        `parse-json
        | scope "result"
        | project-away "dataPointCountRatio", "dimensionCountRatio"
        | mv-expand "data"
        | project "dimensions"=array_to_map("data.dimensions",'host','disk'), "series"=array_from_entries('timestamp',"data.timestamps",'value',"data.values"), "metricId"
        | project "host"="dimensions.host", "disk"="dimensions.disk", "series", "metricId"
        | mv-expand "series"
        | project "timestamp"="series.timestamp", "value"="series.value", "host", "disk", "metricId"
        | extend "timestamp"=unixtime_milliseconds_todatetime("timestamp")
        | order by "timestamp" asc`,
        { data: JSON.stringify(data.system_metrics) }
      )
    ).toStrictEqual([
      { timestamp: new Date("2069-11-11 22:38:20"), disk: "DISK-F1266E1D0AAC2C3F", host: "HOST-F1266E1D0AAC2C3C", value: 11.1, metricId: "builtin:host.disk.avail" },
      { timestamp: new Date("2069-11-11 22:38:20"), disk: "DISK-F1266E1D0AAC2C3D", host: "HOST-F1266E1D0AAC2C3C", value: 111.1, metricId: "builtin:host.disk.avail" },
      { timestamp: new Date("2069-11-11 22:38:20"), disk: undefined, host: "HOST-F1266E1D0AAC2C3C", value: 1.1, metricId: "builtin:host.cpu.idle" },
      { timestamp: new Date("2069-11-11 23:38:20"), disk: "DISK-F1266E1D0AAC2C3F", host: "HOST-F1266E1D0AAC2C3C", value: 22.2, metricId: "builtin:host.disk.avail" },
      { timestamp: new Date("2069-11-11 23:38:20"), disk: "DISK-F1266E1D0AAC2C3D", host: "HOST-F1266E1D0AAC2C3C", value: 222.2, metricId: "builtin:host.disk.avail" },
      { timestamp: new Date("2069-11-11 23:38:20"), disk: undefined, host: "HOST-F1266E1D0AAC2C3C", value: 2.2, metricId: "builtin:host.cpu.idle" },
      { timestamp: new Date("2069-11-12 00:38:20"), disk: "DISK-F1266E1D0AAC2C3F", host: "HOST-F1266E1D0AAC2C3C", value: 33.3, metricId: "builtin:host.disk.avail" },
      { timestamp: new Date("2069-11-12 00:38:20"), disk: "DISK-F1266E1D0AAC2C3D", host: "HOST-F1266E1D0AAC2C3C", value: 333.3, metricId: "builtin:host.disk.avail" },
      { timestamp: new Date("2069-11-12 00:38:20"), disk: undefined, host: "HOST-F1266E1D0AAC2C3C", value: 3.3, metricId: "builtin:host.cpu.idle" },
    ]);
  });
  describe("project", () => {
    it("replace_string", async () => {
      expect(await uql(`parse-json | project "foo"=replace_string("str", 'foo', 'bar')`, { data: { str: "foo foo" } })).toStrictEqual("bar bar");
      expect(await uql(`parse-json | project "foo"=replace_string("str", '(foo)', 'bar')`, { data: { str: "(foo) (foo)" } })).toStrictEqual("bar bar");
      expect(await uql(`parse-json | project "foo"=replace_string("str", '[foo]', 'bar')`, { data: { str: "[foo] [foo]" } })).toStrictEqual("bar bar");
      expect(await uql(`parse-json | project "foo"=replace_string("str", 'foo-baz', 'bar')`, { data: { str: "foo-baz foo-baz" } })).toStrictEqual("bar bar");
      expect(await uql(`parse-json | project "foo"=replace_string("str", 'foo, baz', 'bar')`, { data: { str: "foo, baz foo, baz" } })).toStrictEqual("bar bar");
      expect(await uql(`parse-json | project "foo"=replace_string("str", 'foo\\'s value', 'bar')`, { data: { str: "foo's value foo's value" } })).toStrictEqual("bar bar");
    });
  });
});
