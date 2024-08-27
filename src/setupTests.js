// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { server } from './mocks/server';
import { resetExperiments } from './mocks/handlers';
import { createTempDir, removeTempDir } from './mocks/tempFiles';

process.env.STEADYBIT_URL = 'http://example.com';
process.env.STEADYBIT_TOKEN = 'abcdefgh';

// eslint-disable-next-line no-undef
beforeAll(async () => {
  await createTempDir();
  server.listen();
});

// eslint-disable-next-line no-undef
beforeEach(async () => {
  resetExperiments();
});

// eslint-disable-next-line no-undef
afterEach(async () => {
  server.resetHandlers();
});

// eslint-disable-next-line no-undef
afterAll(async () => {
  await removeTempDir();
  server.close();
});
