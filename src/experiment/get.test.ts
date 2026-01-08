// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { getExperiment } from './get';
import path from 'path';
import fs from 'fs/promises';
import { getTempDir, writeFile } from '../mocks/tempFiles';
import { EXPERIMENTS } from '../mocks/handlers';

describe('experiment', () => {
  describe('get', () => {
    it('should output experiment to file', async () => {
      const file = path.join(getTempDir(), 'experiment.yaml');

      await getExperiment({ key: 'TST-1', file });

      const content = await fs.readFile(file, 'utf-8');
      expect(content).toMatchSnapshot();
    });

    it('should output experiment to stdout JSON', async () => {
      const logSpy = jest.spyOn(console, 'log');
      logSpy.mockClear();

      await getExperiment({ key: 'TST-1', type: 'json' });

      const [stdout] = logSpy.mock.calls[0];
      expect(stdout).toMatchSnapshot();
    });

    it('should output experiment to stdout YAML', async () => {
      const logSpy = jest.spyOn(console, 'log');
      logSpy.mockClear();

      await getExperiment({ key: 'TST-1' });

      const [stdout] = logSpy.mock.calls[0];
      expect(stdout).toMatchSnapshot();
    });

    it('should report not found', async () => {
      await expect(getExperiment({ key: 'TST-999' })).rejects.toThrow('Experiment TST-999 not found.');
    });

    it('should write experiment to file with same data type YAML', async () => {
      const file = path.join(getTempDir(), 'experiment.yaml');
      await writeFile(file, EXPERIMENTS['TST-1']);

      await getExperiment({ key: 'TST-1', file });

      const content = await fs.readFile(file, 'utf-8');
      expect(content).toMatchSnapshot();
    });

    it('should write experiment to file with same data type JSON', async () => {
      const file = path.join(getTempDir(), 'experiment.json');
      await writeFile(file, EXPERIMENTS['TST-1'], 'json');

      await getExperiment({ key: 'TST-1', file });

      const content = await fs.readFile(file, 'utf-8');
      expect(content).toMatchSnapshot();
    });
  });
});
