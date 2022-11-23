// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { server } from './mocks/server';
import { resetExperiments } from './mocks/handlers';
import { createTempDir, removeTempDir } from './mocks/tempFiles';

process.env.STEADYBIT_URL = 'http://test';
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
