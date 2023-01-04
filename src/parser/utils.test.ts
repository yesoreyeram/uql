import { get_value } from "./utils";

describe("utils", () => {
  describe("get_value", () => {
    it("tolower", () => {
      expect(get_value("tolower", ["HeLLO"])).toStrictEqual("hello");
    });
    it("toupper", () => {
      expect(get_value("toupper", ["HeLLO"])).toStrictEqual("HELLO");
    });
    it("strlen", () => {
      expect(get_value("strlen", ["HeLLO"])).toStrictEqual(5);
    });
    it("reverse", () => {
      expect(get_value("reverse", ["hello"])).toStrictEqual("olleh");
    });
    it("atob", () => {
      expect(get_value("atob", ["dGVzdA=="])).toStrictEqual("test");
    });
    it("btoa", () => {
      expect(get_value("btoa", ["test"])).toStrictEqual("dGVzdA==");
    });
    it("split", () => {
      expect(get_value("split", ["hello"])).toStrictEqual(["h", "e", "l", "l", "o"]);
      expect(get_value("split", ["hello world", " "])).toStrictEqual(["hello", "world"]);
      expect(get_value("split", [true])).toStrictEqual([]);
      expect(get_value("split", [1])).toStrictEqual([]);
    });
    it("replace_string", () => {
      expect(get_value("replace_string", [])).toStrictEqual("");
      expect(get_value("replace_string", ["hello world"])).toStrictEqual("hello world");
      expect(get_value("replace_string", ["hello world", "world"])).toStrictEqual("hello ");
      expect(get_value("replace_string", ["hello world", "world", "moon"])).toStrictEqual("hello moon");
      expect(get_value("replace_string", ["hello world", "world", undefined])).toStrictEqual("hello ");
      expect(get_value("replace_string", ["hello world", undefined, undefined])).toStrictEqual("hello world");
      expect(get_value("replace_string", ["banana", "an", "or"])).toStrictEqual("borora");
    });
    it("trim", () => {
      expect(get_value("trim", [" hello "])).toStrictEqual("hello");
    });
    it("trim_start", () => {
      expect(get_value("trim_start", [" hello "])).toStrictEqual("hello ");
    });
    it("trim_end", () => {
      expect(get_value("trim_end", [" hello "])).toStrictEqual(" hello");
    });
    it("pack", () => {
      expect(get_value("pack", ["one", 1, "two", 2, "three"])).toStrictEqual({ one: 1, two: 2, three: undefined });
      expect(get_value("pack", ["one", 1, "two", 2, "three", [3, 3]])).toStrictEqual({ one: 1, two: 2, three: [3, 3] });
    });
    it("array_from_entries", () => {
      expect(get_value("array_from_entries", [])).toStrictEqual([]);
      expect(get_value("array_from_entries", ["timestamp"])).toStrictEqual([]);
      expect(get_value("array_from_entries", ["timestamp", [2010, 2020, 2030]])).toStrictEqual([{ timestamp: 2010 }, { timestamp: 2020 }, { timestamp: 2030 }]);
      expect(get_value("array_from_entries", ["timestamp", [2010, 2020, 2030], "value"])).toStrictEqual([
        { timestamp: 2010, value: null },
        { timestamp: 2020, value: null },
        { timestamp: 2030, value: null },
      ]);
      expect(get_value("array_from_entries", ["timestamp", [2010, 2020, 2030], "value", [1, 0, 2]])).toStrictEqual([
        { timestamp: 2010, value: 1 },
        { timestamp: 2020, value: 0 },
        { timestamp: 2030, value: 2 },
      ]);
      expect(get_value("array_from_entries", ["timestamp", [2010, 2020, 2030], "value", [1, 0]])).toStrictEqual([
        { timestamp: 2010, value: 1 },
        { timestamp: 2020, value: 0 },
        { timestamp: 2030, value: null },
      ]);
      expect(get_value("array_from_entries", ["timestamp", [2010, 2020, 2030], "value", [1, , 0]])).toStrictEqual([
        { timestamp: 2010, value: 1 },
        { timestamp: 2020, value: null },
        { timestamp: 2030, value: 0 },
      ]);
    });
    it("array_to_map", () => {
      expect(get_value("array_to_map", [["ONE", "TWO"]])).toStrictEqual({ 0: "ONE", 1: "TWO" });
      expect(get_value("array_to_map", [["ONE", "TWO"], "one"])).toStrictEqual({ one: "ONE", 1: "TWO" });
      expect(get_value("array_to_map", [["ONE", "TWO"], "one", "two"])).toStrictEqual({ one: "ONE", two: "TWO" });
    });
    it("percentage", () => {
      expect(get_value("percentage", [10])).toStrictEqual(null);
      expect(get_value("percentage", [10, 2000])).toStrictEqual(0.5);
    });
  });
});
