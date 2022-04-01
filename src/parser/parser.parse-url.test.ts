import { uql } from "./index";

describe("parse_url", () => {
  it("should work for default query", async () => {
    const result = await uql(`project "url"=parse_url("url")`, {
      data: [{ url: "https://foo.com?k1=v1&k2=v2#something" }, { url: "https://bar.com?k2=v2" }, { url: "http://user:pass@baz.com?k2=v2" }],
    });
    expect(result).toStrictEqual([{ url: "https://foo.com/?k1=v1&k2=v2#something" }, { url: "https://bar.com/?k2=v2" }, { url: "http://user:pass@baz.com/?k2=v2" }]);
  });
  it("should pick chosen prop", async () => {
    const result = await uql(`project "host_name"=parse_url("url",'origin')`, {
      data: [{ url: "https://foo.com?k1=v1&k2=v2#something" }, { url: "https://bar.com?k2=v2" }, { url: "http://user:pass@baz.com?k2=v2" }],
    });
    expect(result).toStrictEqual([{ host_name: "https://foo.com" }, { host_name: "https://bar.com" }, { host_name: "http://baz.com" }]);
  });
  it("should pick chosen prop", async () => {
    const result = await uql(`project "k2"=parse_url("url",'search','k2')`, {
      data: [{ url: "https://foo.com?k1=v1&k2=v2#something" }, { url: "https://bar.com?k3=v3" }, { url: "http://user:pass@baz.com?k2=v2" }],
    });
    expect(result).toStrictEqual([{ k2: "v2" }, { k2: "" }, { k2: "v2" }]);
  });
});
describe("parse_urlquery", () => {
  it("should work for default query", async () => {
    const result = await uql(`project "search"=parse_urlquery("search")`, {
      data: [{ search: "?k1=v1&k2=v2" }, { search: "k2=v2" }, { search: "k2=v2&k4=v4" }],
    });
    expect(result).toStrictEqual([{ search: { k1: "v1", k2: "v2" } }, { search: { k2: "v2" } }, { search: { k2: "v2", k4: "v4" } }]);
  });
  it("should work for default query", async () => {
    const result = await uql(`project "search"=parse_urlquery("search",'k1')`, {
      data: [{ search: "?k1=v1&k2=v2" }, { search: "k2=v2" }, { search: "k2=v2&k4=v4" }],
    });
    expect(result).toStrictEqual([{ search: "v1" }, { search: "" }, { search: "" }]);
  });
});
