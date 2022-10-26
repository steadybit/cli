// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import path from 'path';
import * as os from 'os';
import fs from 'fs/promises';
import yaml from 'js-yaml';
import { mockExperiments } from '../mocks/handlers';
import { applyExperiment } from './apply';

describe('experiment', () => {
  let tmpFolder: string;

  beforeAll(async () => {
    tmpFolder = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-cli-test'));
  });

  afterAll(async () => {
    await fs.rm(tmpFolder, { recursive: true });
  });

  describe('apply', () => {
    it('should update experiment from file', async () => {
      const file = path.join(tmpFolder, 'experiment.yaml');
      await fs.writeFile(file, yaml.dump(mockExperiments['TST-1']));
      const logSpy = jest.spyOn(console, 'log');

      await applyExperiment({ file });

      expect(logSpy).toHaveBeenCalledWith('Experiment %s updated.', 'TST-1');
    });

    it('should upsert experiment from file', async () => {
      const file = path.join(tmpFolder, 'experiment.yaml');
      await fs.writeFile(file, yaml.dump(mockExperiments['NEW']));
      const logSpy = jest.spyOn(console, 'log');

      await applyExperiment({ file });

      expect(logSpy).toHaveBeenCalledWith('Experiment %s created.', 'TST-2');
    });
  });
});
