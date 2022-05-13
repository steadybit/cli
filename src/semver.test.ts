// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { validateSemverRangeCommanderArgument } from './semver';

describe('semver', () => {
  describe('validateSemverRangeCommanderArgument', () => {
    it('must identify valid semver ranges', () => {
      expect(validateSemverRangeCommanderArgument('1.0.0')).toEqual('1.0.0');
      expect(validateSemverRangeCommanderArgument('^1.0.0')).toEqual('^1.0.0');
    });

    it('must identify invalid semver ranges', () => {
      expect(() => validateSemverRangeCommanderArgument('')).toThrow();
      expect(() => validateSemverRangeCommanderArgument('a')).toThrow();
    });
  });
});
