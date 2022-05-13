// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import fs from 'fs/promises';
import yaml from 'js-yaml';
import path from 'path';

import { forEachFile } from './util';

interface TestFile {
  name: string;
}

describe('defRepo/util', () => {
  describe('forEachFile', () => {
    const cwd = path.join(__dirname, '__tests__', 'dummyDefinitionRepository');
    const temporaryTaskDir = path.join(cwd, 'test');
    const temporaryTaskFile = path.join(temporaryTaskDir, 'test.yml');

    afterEach(async () => {
      await fs.rm(temporaryTaskDir, { recursive: true, force: true });
    });

    it('must iterate over policies', async () => {
      // Given
      const seenFiles = new Set();

      // When
      await forEachFile({
        pattern: '**/task.?(yml|yaml)',
        options: {
          cwd,
          nocase: true,
        },
        iterator: async ({ fileName }) => {
          seenFiles.add(path.relative(cwd, fileName));
        },
      });

      // Then
      expect(Array.from(seenFiles.values()).sort()).toMatchInlineSnapshot(`
Array [
  "experiments/rolling-update/task.yml",
  "weak-spots/rolling-update/task.yAml",
]
`);
    });

    it('must support file modification', async () => {
      // Given
      const callArgs = new Set();
      await fs.mkdir(temporaryTaskDir, { recursive: true });
      await fs.writeFile(
        temporaryTaskFile,
        yaml.dump({
          name: 'test',
        })
      );

      // When
      await forEachFile<TestFile>({
        pattern: '**/test.yml',
        options: {
          cwd,
        },
        iterator: async ({ fileName, fileContent }) => {
          callArgs.add({ fileName: path.relative(cwd, fileName), fileContent: JSON.parse(JSON.stringify(fileContent)) });
          fileContent.name = 'updated';
          return 'update';
        },
      });

      // Then
      expect(callArgs).toMatchInlineSnapshot(`
Set {
  Object {
    "fileContent": Object {
      "name": "test",
    },
    "fileName": "test/test.yml",
  },
}
`);
      const updatedFileContent = await fs.readFile(temporaryTaskFile, 'utf8');
      expect(updatedFileContent).toMatchInlineSnapshot(`
"name: updated
"
`);
    });
  });
});
