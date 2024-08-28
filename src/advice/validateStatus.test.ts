// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

import { validateAdviceStatus } from './validateStatus';

describe('advice', () => {
  describe('validate-status', () => {
    const logSpy = jest.spyOn(console, 'log');
    it('should exit with != 0 if status not matching', async () => {
      await expect(
        validateAdviceStatus({ environment: 'Global', query: 'mock.response=fail', status: 'Implemented' })
      ).rejects.toThrow('2 of 3 advice did not match the expected status.');
      expect(logSpy).toHaveBeenCalledWith('Fetched 3 of 3 matching advice.');
    });

    it('should exit with 0 if all ok', async () => {
      const logSpy = jest.spyOn(console, 'log');
      await validateAdviceStatus({ environment: 'Global', query: 'mock.response=ok', status: 'Implemented' });
      expect(logSpy).toHaveBeenCalledWith('Fetched 1 of 1 matching advice.');
    });
  });
});
