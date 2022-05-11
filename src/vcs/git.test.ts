// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import path from 'path';

import {
  getAbsolutePathToRepositoryRoot,
} from './git';

describe('vcs/git', () => {
  describe('getPathToFileInGitRepository', () => {
    it('must identify repo root path', async () => {
      // When
      const actualRootPath = await getAbsolutePathToRepositoryRoot(__dirname);

      // Then
      const expectedRootPath = path.join(__dirname, '..', '..');
      expect(actualRootPath).toEqual(expectedRootPath);
    });
  });
});
