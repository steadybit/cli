// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { EXPERIMENTS } from '../mocks/handlers';
import { executeExperiments } from './exec';
import { getTempDir, writeFile } from '../mocks/tempFiles';

describe('experiment', () => {
  describe('exec', () => {
    it('should throw when neither key nor file is given', async () => {
      await expect(executeExperiments({ recursive: false, yes: true })).rejects.toThrow(
        'Either --key or --file must be specified.'
      );
    });

    it('should run experiment by key', async () => {
      const logSpy = jest.spyOn(console, 'log');

      await executeExperiments({ key: 'TST-1', recursive: false, yes: true });

      expect(logSpy).toHaveBeenCalledWith('Executing experiment:', 'TST-1');
      expect(logSpy).toHaveBeenCalledWith('Experiment run API:', 'http://example.com/api/experiments/executions/1');
      expect(logSpy).toHaveBeenCalledWith(
        'Experiment run UI:',
        'http://example.com/experiments/edit/TST-1/executions/1?tenant=example&team=EXAMPLE'
      );
    });

    it('should run experiment by file with update', async () => {
      const file = await writeFile('experiment.yaml', EXPERIMENTS['TST-1']);
      const logSpy = jest.spyOn(console, 'log');

      await executeExperiments({ file: [file], recursive: false, yes: true });

      expect(logSpy).toHaveBeenCalledWith('Executing experiment:', 'TST-1');
      expect(logSpy).toHaveBeenCalledWith('Experiment run API:', 'http://example.com/api/experiments/executions/1');
      expect(logSpy).toHaveBeenCalledWith(
        'Experiment run UI:',
        'http://example.com/experiments/edit/TST-1/executions/1?tenant=example&team=EXAMPLE'
      );
    });

    it('should run experiment by file with upsert', async () => {
      const file = await writeFile('experiment.yaml', EXPERIMENTS['NEW']);
      const logSpy = jest.spyOn(console, 'log');

      await executeExperiments({ file: [file], recursive: false, yes: true });

      expect(logSpy).toHaveBeenCalledWith('Executing experiment:', 'NEW-1');
      expect(logSpy).toHaveBeenCalledWith('Experiment run API:', 'http://example.com/api/experiments/executions/1');
      expect(logSpy).toHaveBeenCalledWith(
        'Experiment run UI:',
        'http://example.com/experiments/edit/TST-1/executions/1?tenant=example&team=EXAMPLE'
      );
    });

    it('should run experiments from directory with upsert', async () => {
      await writeFile('experiment-1.yaml', EXPERIMENTS['NEW']);
      await writeFile('experiment-2.yaml', EXPERIMENTS['NEW']);
      const logSpy = jest.spyOn(console, 'log');

      await executeExperiments({ file: [getTempDir()], recursive: false, yes: true });

      expect(logSpy).toHaveBeenCalledWith('Executing experiment:', 'NEW-1');
      expect(logSpy).toHaveBeenCalledWith('Experiment run API:', 'http://example.com/api/experiments/executions/1');
      expect(logSpy).toHaveBeenCalledWith(
        'Experiment run UI:',
        'http://example.com/experiments/edit/NEW-1/executions/1?tenant=example&team=EXAMPLE'
      );
      expect(logSpy).toHaveBeenCalledWith('Executing experiment:', 'NEW-2');
      expect(logSpy).toHaveBeenCalledWith(
        'Experiment run UI:',
        'http://example.com/experiments/edit/NEW-2/executions/2?tenant=example&team=EXAMPLE'
      );
    });

    it('should throw when key and two or more files are given', async () => {
      await writeFile('experiment-1.yaml', EXPERIMENTS['NEW']);
      await writeFile('experiment-2.yaml', EXPERIMENTS['NEW']);

      await expect(
        executeExperiments({
          key: 'TST-1',
          file: [getTempDir()],
          recursive: false,
          yes: true,
        })
      ).rejects.toThrow('If --key is specified, at most one --file can be specified.');
    });
  });
});
