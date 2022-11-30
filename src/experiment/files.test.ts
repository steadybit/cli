// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { getTempDir, writeYamlFile } from '../mocks/tempFiles';
import { resolveExperimentFiles } from './files';


describe('experiment/files', () => {
  describe('resolveExperimentFiles', () => {
    let files: string[];
    let nestedFiles: string[];

    beforeAll(async () => {
      files = await Promise.all([
        writeYamlFile('experiment.yaml', {}),
        writeYamlFile('experiment.yml', {})
      ]);
      nestedFiles = await Promise.all([
        writeYamlFile('nested/experiment.yaml', {}),
        writeYamlFile('nested/experiment.yml', {})
      ]);
      await Promise.all([
        writeYamlFile('other.txt', {}),
        writeYamlFile('nested/other.txt', {})
      ]);
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
