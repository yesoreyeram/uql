import { uql } from "./index";

describe("parser", () => {
  describe("jsonata", () => {
    it("should return a specific value", async () => {
      expect(await uql(`jsonata "$sum(example.value)"`, { data })).toStrictEqual(24);
      expect(await uql(`jsonata "$sum(library.books.price)"`, { data: data.library })).toStrictEqual(95.87);
    });
    it("should return an array of objects", async () => {
      expect(await uql(`jsonata "example.value"`, { data })).toStrictEqual([4, 7, 13]);
      expect(await uql(`jsonata "library.loans@$L.books@$B[$L.isbn=$B.isbn].customers[$L.customer=id].{ 'customer': name, 'book': $B.title, 'due': $L.return}"`, { data: data.library })).toStrictEqual(
        [
          { book: "Structure and Interpretation of Computer Programs", customer: "Joe Doe", due: "2016-12-05" },
          { book: "Compilers: Principles, Techniques, and Tools", customer: "Jason Arthur", due: "2016-10-22" },
          { book: "Structure and Interpretation of Computer Programs", customer: "Jason Arthur", due: "2016-12-22" },
        ]
      );
    });
    it("should combine with other queries", async () => {
      expect(
        await uql(`scope "library" | jsonata "library.loans@$L.books@$B[$L.isbn=$B.isbn].customers[$L.customer=id].{ 'customer': name, 'book': $B.title, 'due': $L.return}" | count`, { data: data })
      ).toStrictEqual(3);
    });
  });
});

const data = {
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
