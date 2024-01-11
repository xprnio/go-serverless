import z from "zod";

import { parseBody } from "./body-parser";
import { readRequest } from "./request-reader";
import { writeResponse } from "./response-writer";

interface LambdaParams<T, TResult> {
  body: z.ZodSchema<T>;
  handler(body: T): TResult | Promise<TResult>;
}

export function lambda<T, TResult>(params: LambdaParams<T, TResult>) {
  run(params).then(
    (result) => writeResponse(result),
    (err) => writeResponse(err),
  );
}

async function run<T, TResult>(params: LambdaParams<T, TResult>) {
  const request = await readRequest();
  const body = parseBody(params.body, request);
  return params.handler(body);
}
