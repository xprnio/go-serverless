import fs from 'node:fs/promises';

// This is where the HTTP response must be saved to on the container
const RESPONSE_PATH = '/context/response.json';

export async function writeResponse<T>(response: T) {
  const json = JSON.stringify(response, null, 2);
  await fs.writeFile(RESPONSE_PATH, json, 'utf-8');
}
