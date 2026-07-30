// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

import { describe, expect, it, vi } from 'vitest';
import { validateAdviceStatus } from './validateStatus.ts';

describe('advice', () => {
  describe('validate-status', () => {
    const logSpy = vi.spyOn(console, 'log');
    it('should exit with != 0 if status not matching', async () => {
      await expect(
        validateAdviceStatus({ environment: 'Global', query: 'mock.response=fail', status: 'Implemented' })
      ).rejects.toThrow('2 of 3 advice did not match the expected status.');
      expect(logSpy).toHaveBeenCalledWith('Fetched 3 of 3 matching advice.');
    });

    it('should exit with 0 if all ok', async () => {
      const logSpy = vi.spyOn(console, 'log');
      await validateAdviceStatus({ environment: 'Global', query: 'mock.response=ok', status: 'Implemented' });
      expect(logSpy).toHaveBeenCalledWith('Fetched 1 of 1 matching advice.');
    });

    // The platform reports IMPLEMENTED while the default for --status is written
    // Implemented. Comparing them exactly meant the command failed even when every
    // piece of advice was implemented, which is the whole point of the check.
    it.each(['Implemented', 'IMPLEMENTED', 'implemented', ' Implemented '])(
      'should accept %s as the expected status',
      async status => {
        await expect(
          validateAdviceStatus({ environment: 'Global', query: 'mock.response=ok', status })
        ).resolves.toBeUndefined();
      }
    );

    it.each([
      ['ACTION_NEEDED', 2],
      ['action needed', 2],
      ['Action needed', 2],
    ])('should treat %s as the same status the platform reports', async (status, expectedFailures) => {
      await expect(
        validateAdviceStatus({ environment: 'Global', query: 'mock.response=fail', status: String(status) })
      ).rejects.toThrow(`${expectedFailures} of 3 advice did not match the expected status.`);
    });

    it('should still reject a status that genuinely differs', async () => {
      await expect(
        validateAdviceStatus({ environment: 'Global', query: 'mock.response=ok', status: 'ACTION_NEEDED' })
      ).rejects.toThrow('1 of 1 advice did not match the expected status.');
    });
  });
});
