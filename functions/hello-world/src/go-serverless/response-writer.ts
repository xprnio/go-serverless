import fs from 'node:fs/promises';
import { RequestContext } from './context';

// This is where the HTTP response must be saved to on the container
const RESPONSE_PATH = `${RequestContext.ContextPath}/response.json`;

export async function writeResponse<T>(response: T) {
  const json = JSON.stringify(response, null, 2);
  await fs.writeFile(RESPONSE_PATH, json, 'utf-8');
}
