import { getAST } from "./ast";
import { evaluate } from "./evaluate";

export const uql = (query: string, options?: { data?: any }): Promise<unknown> => {
  return new Promise((resolve, reject) => {
    if (!query) {
      resolve("hello there! provide a valid query");
    } else {
      getAST(query)
        .then((res) => {
          evaluate(res, options)
            .then((response) => resolve(response))
            .catch((ex) => reject(ex));
        })
        .catch((ex) => reject(ex));
    }
  });
};
