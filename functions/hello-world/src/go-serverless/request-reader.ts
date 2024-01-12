import fs from 'node:fs/promises';
import { RequestContext } from './context';

// This is where the HTTP request is saved to on the container
const REQUEST_PATH = `${RequestContext.ContextPath}/request.json`;

export async function readRequest(): Promise<string> {
  return await fs.readFile(REQUEST_PATH, 'utf-8');
}
