# UQL - Unstructured Query Language

Unstructured query language (UQL) - It is a query language to query the unstructured or array like objects in javascript.

## Installation

Install the UQL from npm/yarn

```sh
npm install uql
## or
yarn add uql
```

Then in your code use this as follows

```ts
import { uql } from "uql";
const users = [
  { name: "foo", age: 2, location: "uk" },
  { name: "bar", age: 3, location: "usa" },
];
const query = `order by "name" asc | project "name", "location"`;
uql(query, { data: users })
  .then((res) => console.log(res))
  .catch((ex) => console.error(ex));
//
// Output
//
// [ { name: 'bar', location: 'usa' }, { name: 'foo', location: 'uk' } ]
```

## Commands

### ping

A very simple `ping` query will respond you with `pong`.

### echo

A very simple `echo` query will echo what you said. Example: A query `echo "hello world"` will print `hello world`.

### count

As the name suggests, `count` query will give you the length of an array. For example, the input `[{},{}]` and then the query `count` will respond with `2`

### limit

The `limit` query, return up to the specified number of rows in an array like object. For example, the input `[1,2,3]` and then the query `limit 2` will respond with `[1,2]`

limit query expects one mandatory argument which have to be a number

### order by

The `order by` query, sorts the items in the array into order by one or more columns. For example: `order by "name" asc` will order the input array by a field "name".

Note: field name should be referred by single quote and the order must be one of "asc" or "desc"

### project

The `project` query, limits the field you want to return from the array of objects.

Example:

With `[ { name: "foo", age: 1, country: "uk" }, { name: "bar", age :2, country: "usa" }]`, the query `project "name", "age"` will only return name and age properties of each element.

Note: field name should be referred by single quote.

### project-away

The `project-away` query is opposite of `project` query. It omits the specified properties.

Example:

With `[ { name: "foo", age: 1, country: "uk" }, { name: "bar", age :2, country: "usa" }]`, the query `project-away "country", "age"` will only return name property of each element.

Note: field name should be referred by single quote.
