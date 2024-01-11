import { z } from 'zod';
import { lambda } from './go-serverless';

lambda({
  body: z.object({
    message: z.string(),
  }),
  handler(body) {
    return {
      success: true,
      data: { message: body.message }
    };
  },
});
