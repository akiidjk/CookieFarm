import path from 'node:path';
import { createOpenAPI } from 'fumadocs-openapi/server';

const schemaPath = path.resolve(process.cwd(), '../cookiefarm/server/api/docs/swagger.json');

export const openapi = createOpenAPI({
  input: [schemaPath],
});
