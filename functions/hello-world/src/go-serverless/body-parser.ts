import z from "zod";

const jsonSchema = z.string().transform((json, ctx) => {
  try {
    return JSON.parse(json);
  } catch (error) {
    if (error instanceof Error) {
      ctx.addIssue({
        code: 'custom',
        message: error.message,
      });
      return;
    }

    console.error(error);
    ctx.addIssue({
      code: 'custom',
      message: 'Unknown Error',
    });
  }
})

export function parseBody<T>(schema: z.ZodSchema<T>, json: string): T {
  const body = jsonSchema.parse(json);
  return schema.parse(body);
}
