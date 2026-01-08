// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { getTempDir, writeFile } from '../mocks/tempFiles';
import { resolveExperimentFiles } from './files';

describe('experiment/files', () => {
  describe('resolveExperimentFiles', () => {
    let files: string[];
    let nestedFiles: string[];

    beforeAll(async () => {
      files = await Promise.all([writeFile('experiment.yaml', {}), writeFile('experiment.yml', {})]);
      nestedFiles = await Promise.all([
        writeFile('nested/experiment.yaml', {}),
        writeFile('nested/experiment.yml', {}),
      ]);
      await Promise.all([writeFile('other.txt', {}), writeFile('nested/other.txt', {})]);
    });

    it('should resolve all files recursive', async () => {
      expect(await resolveExperimentFiles([getTempDir()], true)).toEqual([...files, ...nestedFiles]);
    });
    it('should resolve files non-recursive', async () => {
      expect(await resolveExperimentFiles([getTempDir()], false)).toEqual(files);
    });
    it('should resolve files directly', async () => {
      expect(await resolveExperimentFiles(files, false)).toEqual(files);
    });
    it('should resolve empty files', async () => {
      expect(await resolveExperimentFiles([], false)).toEqual([]);
    });
  });
});
