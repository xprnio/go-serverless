import fs from 'node:fs/promises';

// This is where the HTTP request is saved to on the container
const REQUEST_PATH = '/context/request.json';

export async function readRequest(): Promise<string> {
  return await fs.readFile(REQUEST_PATH, 'utf-8');
}
