// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { deleteExperiment } from './delete';

describe('experiment', () => {

  describe('delete', () => {
    it('should delete an experiment', async () => {
      const logSpy = jest.spyOn(console, 'log');

      await (deleteExperiment({ key: 'TST-1' }));

      expect(logSpy).toHaveBeenCalledWith('Experiment %s deleted.', 'TST-1');
    });

    it('should report not found', async () => {
      await expect(deleteExperiment({ key: 'TST-999' })).rejects.toThrow('Experiment TST-999 not found.')
    });
  });
});
