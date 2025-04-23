<!-- markdownlint-configure-file {
  "MD013": false,
  "MD033": false
} -->

<h1 align="center">
  UQL - Unified Query Language
</h1>

<p align="center">Unified query language (UQL) - It is a query language to query the JSON like data in Javascript. Inspired by azure Kusto query language (KQL).</p>

<p align="center">
    <a href="https://yesoreyeram.github.io/uql">
      <img src="https://raw.githubusercontent.com/yesoreyeram/uql/refs/heads/main/public/logo.svg" alt="UQL" width="200" height="200">
    </a>
</p>

## Installation

Install the UQL from npm / yarn

```sh
## With npm
npm install uql

## With yarn
yarn add uql
```

Then in your code use this as follows

```ts
import { uql } from "uql";

const users = [
  { name: "foo", age: 2, location: "uk" },
  { name: "bar", age: 3, location: "usa" },
];

const query = `parse-json
| order by "name" asc
| project "name", "location"`;

uql(query, { data: users })
  .then((res) => console.log(res))
  .catch((ex) => console.error(ex));

//
// Output
//
// [ { name: 'bar', location: 'usa' }, { name: 'foo', location: 'uk' } ]
```

## Basic UQL Commands

### project

The `project` query, limits the field you want to return from the array of objects.

Example:

With `[ { name: "foo", age: 1, country: "uk" }, { name: "bar", age :2, country: "usa" }]`, the query `project "name", "age"` will only return name and age properties of each element.

**Note**: All the field name should be wrapped with double quotes

### project-away

The `project-away` query is opposite of `project` query. It omits the specified properties.

Example:

With `[ { name: "foo", age: 1, country: "uk" }, { name: "bar", age :2, country: "usa" }]`, the query `project-away "country", "age"` will only return name property of each element.
