// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { EXPERIMENTS } from '../mocks/handlers';
import { applyExperiments } from './apply';
import * as fileApi from './files';
import { getTempDir, writeFile } from '../mocks/tempFiles';

describe('experiment', () => {
  describe('apply', () => {
    it('should update experiment from file', async () => {
      const file = await writeFile('experiment.yaml', EXPERIMENTS['TST-1']);
      const logSpy = jest.spyOn(console, 'log');

      await applyExperiments({ file: [file], recursive: false });

      expect(logSpy).toHaveBeenCalledWith('Experiment %s updated.', 'TST-1');
    });

    it('should upsert experiment from file', async () => {
      const file = await writeFile('experiment.yaml', EXPERIMENTS['NEW']);
      const logSpy = jest.spyOn(console, 'log');

      await applyExperiments({ file: [file], recursive: false });

      expect(logSpy).toHaveBeenCalledWith('Experiment %s created.', 'NEW-1');
    });

    it('should upsert experiment from directory', async () => {
      await writeFile('experiment-1.yaml', EXPERIMENTS['NEW']);
      await writeFile('experiment-2.yaml', EXPERIMENTS['NEW']);
      const logSpy = jest.spyOn(console, 'log');

      await applyExperiments({ file: [getTempDir()], recursive: false });

      expect(logSpy).toHaveBeenCalledWith('Experiment %s created.', 'NEW-1');
      expect(logSpy).toHaveBeenCalledWith('Experiment %s created.', 'NEW-2');
    });

    it('should throw when key and two or more files are given', async () => {
      await writeFile('experiment-1.yaml', EXPERIMENTS['NEW']);
      await writeFile('experiment-2.yaml', EXPERIMENTS['NEW']);

      await expect(applyExperiments({ key: 'TST-1', file: [getTempDir()], recursive: false })).rejects.toThrow(
        'If --key is specified, at most one --file can be specified.'
      );
    });

    it('should keep file data type YML', async () => {
      const file = await writeFile('experiment.yaml', { ...EXPERIMENTS['TST-1'], key: undefined });
      const logSpy = jest.spyOn(fileApi, 'writeFile');

      await applyExperiments({ file: [file], recursive: false });

      expect(logSpy).toHaveBeenCalledWith(expect.any(String), expect.any(Object), 'yaml');
    });

    it('should keep file data type JSON', async () => {
      const file = await writeFile('experiment.json', { ...EXPERIMENTS['TST-1'], key: undefined }, 'json');
      const logSpy = jest.spyOn(fileApi, 'writeFile');

      await applyExperiments({ file: [file], recursive: false });

      expect(logSpy).toHaveBeenCalledWith(expect.any(String), expect.any(Object), 'json');
    });
  });
});
