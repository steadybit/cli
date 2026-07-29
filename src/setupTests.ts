// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { afterAll, afterEach, beforeAll, beforeEach } from 'vitest';
import { server } from './mocks/server.ts';
import { resetExperiments } from './mocks/handlers.ts';
import { createTempDir, removeTempDir } from './mocks/tempFiles.ts';

process.env.STEADYBIT_URL = 'http://example.com';
process.env.STEADYBIT_TOKEN = 'abcdefgh';

beforeAll(async () => {
  await createTempDir();
  server.listen();
});

beforeEach(async () => {
  resetExperiments();
});

afterEach(async () => {
  server.resetHandlers();
});

afterAll(async () => {
  await removeTempDir();
  server.close();
});
