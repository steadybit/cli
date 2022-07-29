// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { v4 as uuidv4 } from 'uuid';
import fs from 'fs/promises';
import { tmpdir } from 'os';
import yaml from 'js-yaml';
import path from 'path';

import { loadServiceDefinition, writeServiceDefinition } from './files';

describe('serviceDefinition/files', () => {
  describe('loadServiceDefinition', () => {
    const fileWithSyntaxErrors = path.join(tmpdir(), `${uuidv4()}-syntax-errors.json`);
    const fileWithValidJson = path.join(tmpdir(), `${uuidv4()}-valid.json`);
    const fileWithValidYaml = path.join(tmpdir(), `${uuidv4()}-valid.yml`);

    beforeAll(async () => {
      await Promise.all([
        fs.writeFile(fileWithSyntaxErrors, '{$$$42:'),
        fs.writeFile(fileWithValidJson, JSON.stringify({ type: 'json' })),
        fs.writeFile(fileWithValidYaml, yaml.dump({ type: 'yaml' })),
      ]);
    });

    afterAll(async () => {
      await Promise.all([fs.unlink(fileWithSyntaxErrors), fs.unlink(fileWithValidJson), fs.unlink(fileWithValidYaml)]);
    });

    it('must reject when file cannot be found', async () => {
      await expect(loadServiceDefinition(path.join(tmpdir(), uuidv4()))).rejects.toThrow(
        /Failed to read service definition file.*/
      );
    });

    it('must reject when file cannot be parsed', async () => {
      await expect(loadServiceDefinition(fileWithSyntaxErrors)).rejects.toThrow(
        /Failed to parse service definition file.*/
      );
    });

    it('must successfully load JSON files', async () => {
      await expect(loadServiceDefinition(fileWithValidJson)).resolves.toMatchInlineSnapshot(`
Object {
  "type": "json",
}
`);
    });

    it('must successfully load YAML files', async () => {
      await expect(loadServiceDefinition(fileWithValidYaml)).resolves.toMatchInlineSnapshot(`
Object {
  "type": "yaml",
}
`);
    });
  });

  describe('writeServiceDefinition', () => {
    const outputfile = path.join(tmpdir(), `${uuidv4()}-valid.yml`);

    afterAll(async () => {
      await Promise.all([fs.unlink(outputfile)]);
    });

    it('must successfully write YAML files', async () => {
      await writeServiceDefinition(outputfile, { name: 'Test ServiceDefinition', mapping: {} });

      expect((await fs.readFile(outputfile)).toString()).toEqual(`name: Test ServiceDefinition
mapping: {}
`);
    });
  });
});
