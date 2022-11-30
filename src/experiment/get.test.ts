// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { getExperiment } from './get';
import path from 'path';
import fs from 'fs/promises';
import { getTempDir } from '../mocks/tempFiles';

describe('experiment', () => {
  describe('get', () => {
    it('should output experiment to file', async () => {
      const file = path.join(getTempDir(), 'experiment.yaml');

      await getExperiment({ key: 'TST-1', file });

      const content = await fs.readFile(file, 'utf-8');
      expect(content).toMatchSnapshot();
    });

    it('should report not found', async () => {
      await expect(getExperiment({ key: 'TST-999' })).rejects.toThrow('Experiment TST-999 not found.');
    });
  });
});
