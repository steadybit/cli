// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { getExperiment } from './get';
import path from 'path';
import * as os from 'os';
import fs from 'fs/promises';

describe('experiment', () => {
  let tmpFolder: string;

  beforeAll(async () => {
    tmpFolder = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-cli-test'));
  });

  afterAll(async () => {
    await fs.rm(tmpFolder, { recursive: true });
  });

  describe('get', () => {
    it('should output experiment to file', async () => {
      const file = path.join(tmpFolder, 'experiment.yaml');
      await getExperiment({ key: 'TST-1', file });

      const content = await fs.readFile(file, 'utf-8');
      expect(content).toMatchSnapshot();
    });

    it('should report not found', async () => {
        await expect(getExperiment({ key: 'TST-999' })).rejects.toThrow('Experiment TST-999 not found.');
    });
  });
});
