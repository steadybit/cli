// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import path from 'path';
import * as os from 'os';
import fs from 'fs/promises';
import yaml from 'js-yaml';
import { mockExperiments } from '../mocks/handlers';
import { executeExperiment } from './exec';

describe('experiment', () => {
  let tmpFolder: string;

  beforeAll(async () => {
    tmpFolder = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-cli-test'));
  });

  afterAll(async () => {
    await fs.rm(tmpFolder, { recursive: true });
  });

  describe('exec', () => {
    it('should throw when neither key nor file is given', async () => {
      await expect(executeExperiment({})).rejects.toThrow('Either --key or --file must be specified.');
    });

    it('should run experiment by key', async () => {
      const logSpy = jest.spyOn(console, 'log');

      await executeExperiment({ key: 'TST-1' });

      expect(logSpy).toHaveBeenCalledWith('Executing experiment:', 'TST-1');
      expect(logSpy).toHaveBeenCalledWith('Experiment run:', 'http://test/api/experiments/executions/1');
    });

    it('should run experiment by file with update', async () => {
      const file = path.join(tmpFolder, 'experiment.yaml');
      await fs.writeFile(file, yaml.dump(mockExperiments['TST-1']));
      const logSpy = jest.spyOn(console, 'log');

      await executeExperiment({ file });

      expect(logSpy).toHaveBeenCalledWith('Executing experiment:', 'TST-1');
      expect(logSpy).toHaveBeenCalledWith('Experiment run:', 'http://test/api/experiments/executions/1');
    });

    it('should run experiment by file with upsert', async () => {
      const file = path.join(tmpFolder, 'experiment.yaml');
      await fs.writeFile(file, yaml.dump(mockExperiments['NEW']));
      const logSpy = jest.spyOn(console, 'log');

      await executeExperiment({ file });

      expect(logSpy).toHaveBeenCalledWith('Executing experiment:', 'TST-2');
      expect(logSpy).toHaveBeenCalledWith('Experiment run:', 'http://test/api/experiments/executions/2');
    });
  });
});
