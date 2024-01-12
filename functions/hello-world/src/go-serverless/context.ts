const { REQUEST_ID, CONTEXT_PATH } = process.env;

if (!REQUEST_ID) throw new Error('REQUEST_ID not set');
if (!CONTEXT_PATH) throw new Error('CONTEXT_PATH not set');

export const RequestContext = {
  RequestId: REQUEST_ID,
  ContextPath: CONTEXT_PATH,
} as const;
