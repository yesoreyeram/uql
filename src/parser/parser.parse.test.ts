import { get } from "lodash";
import { uql } from "./index";

const sample_data = {
  csv_countries_and_population: `country,population
india,300
uk,100
usa,200`,
  yaml_example: `
---
doe: "a deer, a female deer"
ray: "a drop of golden sun"
pi: 3.14159
xmas: true
french-hens: 3
calling-birds:
  - huey
  - dewey
  - louie
  - fred
xmas-fifth-day:
  calling-birds: four
  french-hens: 3
  golden-rings: 5
  partridges:
    count: 1
    location: "a pear tree"
  turtle-doves: two`,
  xml_countries_and_population: `<countries>
<country><name>india</name><population>300</population></country>
<country><name>uk</name><population>100</population></country>
<country><name>usa</name><population>200</population></country>
</countries>`,
  xml_aws_status_sample: `<?xml version="1.0" encoding="UTF-8"?>
  <rss version="2.0">
    <channel>
      <title><![CDATA[Amazon Web Services Service Status]]></title>
      <link>http://status.aws.amazon.com/</link>
      <language>en-us</language>
      <lastBuildDate>Mon, 20 Dec 2021 00:05:42 PST</lastBuildDate>
      <generator>AWS Service Health Dashboard RSS Generator</generator>
      <description><![CDATA[Amazon Web Services Service Status]]></description>
      <ttl>5</ttl>
      <!-- You seem to care about knowing about your events, why not check out https://docs.aws.amazon.com/health/latest/ug/getting-started-api.html -->
      <item>
      <title><![CDATA[Service is operating normally: [RESOLVED] Internet Connectivity]]></title>
      <link>http://status.aws.amazon.com/</link>
      <pubDate>Wed, 15 Dec 2021 12:16:32 PST</pubDate>
      <guid isPermaLink="false">http://status.aws.amazon.com/#internetconnectivity-us-gov-west-1_1639599392</guid>
      <description><![CDATA[Between 7:14 AM PST and 7:59 AM PST, customers experienced elevated network packet loss that impacted connectivity to a subset of Internet destinations. Traffic within AWS Regions, between AWS Regions, and to other destinations on the Internet was not impacted. The issue was caused by network congestion between parts of the AWS Backbone and a subset of Internet Service Providers, which was triggered by AWS traffic engineering, executed in response to congestion outside of our network. This traffic engineering incorrectly moved more traffic than expected to parts of the AWS Backbone that affected connectivity to a subset of Internet destinations. The issue has been resolved, and we do not expect a recurrence.]]></description>
     </item>
           
     <item>
      <title><![CDATA[Service is operating normally: [RESOLVED] Internet Connectivity]]></title>
      <link>http://status.aws.amazon.com/</link>
      <pubDate>Wed, 15 Dec 2021 08:10:17 PST</pubDate>
      <guid isPermaLink="false">http://status.aws.amazon.com/#internetconnectivity-us-gov-west-1_1639584617</guid>
      <description><![CDATA[We have resolved the issue affecting Internet connectivity to the US-GOV-WEST-1 Region. Connectivity within the region was not affected by this event. The issue has been resolved and the service is operating normally.]]></description>
     </item>
  </channel>
  </rss>`,
};

describe("parser", () => {
  describe("parse-json", () => {
    it("default", async () => {
      const result = await uql("parse-json", { data: "[{},{}]" });
      expect(result).toStrictEqual([{}, {}]);
    });
  });
  describe("parse-csv", () => {
    it("default", async () => {
      const result = await uql("parse-csv", { data: sample_data.csv_countries_and_population });
      expect(result).toStrictEqual([
        { country: "india", population: "300" },
        { country: "uk", population: "100" },
        { country: "usa", population: "200" },
      ]);
    });
    it("with options", async () => {
      const result = await uql("parse-csv --delimiter ':'", { data: sample_data.csv_countries_and_population.replace(/\,/g, ":") });
      expect(result).toStrictEqual([
        { country: "india", population: "300" },
        { country: "uk", population: "100" },
        { country: "usa", population: "200" },
      ]);
    });
    it("tsv", async () => {
      const result: any = await uql("parse-csv --delimiter '\t' --columns 'a,b,c,d,e,f,g,h,i,j'", {
        data: `aar_lsd_broad.tsv	aar	ds	dds	Latin		False	Broad	True	715
acw_ard_broad.tsv	acw	asd sas	sds sds	sss		False	Broad	False	1090`,
      });
      expect(result[0]["b"]).toStrictEqual("aar");
      expect(result[0]["j"]).toStrictEqual("715");
    });
  });
  describe("parse-xml", () => {
    it("default", async () => {
      const result = await uql(`parse-xml | scope "countries.country"`, { data: sample_data.xml_countries_and_population });
      expect(result).toStrictEqual([
        { name: "india", population: 300 },
        { name: "uk", population: 100 },
        { name: "usa", population: 200 },
      ]);
    });
    it("with attribs", async () => {
      const result = await uql(`parse-xml`, { data: `<root a="nice" checked><a>wow</a></root>` });
      expect(result).toStrictEqual({ root: { "@_a": "nice", "@_checked": true, a: "wow" } });
    });
    it("parsing with additional parameters", async () => {
      const result = await uql(`parse-xml --allowBooleanAttributes false --attributeNamePrefix "$."`, { data: `<root a="nice" checked><a>wow</a></root>` });
      expect(result).toStrictEqual({ root: { "$.a": "nice", a: "wow" } });
      expect(get(result, "root['$.a']")).toStrictEqual("nice");
    });
    it("rss feed", async () => {
      const result: any = await uql(`parse-xml`, { data: sample_data.xml_aws_status_sample });
      expect(result.rss.channel.item[0].title).toStrictEqual("Service is operating normally: [RESOLVED] Internet Connectivity");
    });
  });
  describe("parse-yaml", () => {
    it("default", async () => {
      const result: any = await uql("parse-yaml", { data: sample_data.yaml_example });
      expect(result).toStrictEqual({
        doe: "a deer, a female deer",
        ray: "a drop of golden sun",
        pi: 3.14159,
        xmas: true,
        "french-hens": 3,
        "calling-birds": ["huey", "dewey", "louie", "fred"],
        "xmas-fifth-day": {
          "calling-birds": "four",
          "french-hens": 3,
          "golden-rings": 5,
          partridges: {
            count: 1,
            location: "a pear tree",
          },
          "turtle-doves": "two",
        },
      });
    });
  });
});
