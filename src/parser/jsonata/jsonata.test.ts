import { uql } from "../index";

describe("parser", () => {
  describe("jsonata", () => {
    it("should return a specific value", async () => {
      expect(await uql(`jsonata "$sum(example.value)"`, { data })).toStrictEqual(24);
      expect(await uql(`jsonata "$sum(library.books.price)"`, { data: data.library })).toStrictEqual(95.87);
    });
    it("should return an array of objects", async () => {
      let result = await uql(`jsonata "example.value"`, { data });
      expect(JSON.stringify(result)).toStrictEqual(JSON.stringify([4, 7, 13]));
      result = await uql(`jsonata "library.loans@$L.books@$B[$L.isbn=$B.isbn].customers[$L.customer=id].{ 'customer': name, 'book': $B.title, 'due': $L.return}"`, { data: data.library });
      expect(JSON.stringify(result)).toStrictEqual(
        JSON.stringify([
          { customer: "Joe Doe", book: "Structure and Interpretation of Computer Programs", due: "2016-12-05" },
          { customer: "Jason Arthur", book: "Compilers: Principles, Techniques, and Tools", due: "2016-10-22" },
          { customer: "Jason Arthur", book: "Structure and Interpretation of Computer Programs", due: "2016-12-22" },
        ])
      );
    });
    it("should combine with other queries", async () => {
      let query = `scope "library" 
| jsonata "library.loans@$L.books@$B[$L.isbn=$B.isbn].customers[$L.customer=id].{ 'customer': name, 'book': $B.title, 'due': $L.return}" 
| count`;
      let result = await uql(query, { data });
      expect(result).toStrictEqual(3);
    });
  });
  describe("jsonata filter", () => {
    it("should filter data correctly with single element", async () => {
      let result = await uql(`scope "Countries" | jsonata "*[Country='India'][]"`, { data });
      expect(JSON.stringify(result)).toStrictEqual(JSON.stringify([{ Country: "India" }]));
    });
    it("should filter data correctly with multiple elements", async () => {
      let result = await uql(`scope "Countries" | jsonata "*[Country in ['India','United Kingdom']][]"`, { data });
      expect(JSON.stringify(result)).toStrictEqual(JSON.stringify([{ Country: "India" }, { Country: "United Kingdom" }]));
    });
  });
});

const data = {
  Countries: [{ Country: "India" }, { Country: "America" }, { Country: "United Kingdom" }],
  example: [{ value: 4 }, { value: 7 }, { value: 13 }],
  library: {
    library: {
      books: [
        {
          title: "Structure and Interpretation of Computer Programs",
          authors: ["Abelson", "Sussman"],
          isbn: "9780262510875",
          price: 38.9,
          copies: 2,
        },
        {
          title: "The C Programming Language",
          authors: ["Kernighan", "Richie"],
          isbn: "9780131103627",
          price: 33.59,
          copies: 3,
        },
        {
          title: "The AWK Programming Language",
          authors: ["Aho", "Kernighan", "Weinberger"],
          isbn: "9780201079814",
          copies: 1,
        },
        {
          title: "Compilers: Principles, Techniques, and Tools",
          authors: ["Aho", "Lam", "Sethi", "Ullman"],
          isbn: "9780201100884",
          price: 23.38,
          copies: 1,
        },
      ],
      loans: [
        {
          customer: "10001",
          isbn: "9780262510875",
          return: "2016-12-05",
        },
        {
          customer: "10003",
          isbn: "9780201100884",
          return: "2016-10-22",
        },
        {
          customer: "10003",
          isbn: "9780262510875",
          return: "2016-12-22",
        },
      ],
      customers: [
        {
          id: "10001",
          name: "Joe Doe",
          address: {
            street: "2 Long Road",
            city: "Winchester",
            postcode: "SO22 5PU",
          },
        },
        {
          id: "10002",
          name: "Fred Bloggs",
          address: {
            street: "56 Letsby Avenue",
            city: "Winchester",
            postcode: "SO22 4WD",
          },
        },
        {
          id: "10003",
          name: "Jason Arthur",
          address: {
            street: "1 Preddy Gate",
            city: "Southampton",
            postcode: "SO14 0MG",
          },
        },
      ],
    },
  },
};
