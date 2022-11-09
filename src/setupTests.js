// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { server } from './mocks/server';

process.env.STEADYBIT_URL = 'http://test'
process.env.STEADYBIT_TOKEN = 'abcdefgh'

beforeAll(() => server.listen());

afterEach(() => server.resetHandlers());

afterAll(() => server.close());
