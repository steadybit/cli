import { v4 as uuidv4 } from 'uuid';
import fs from 'fs/promises';
import { tmpdir } from 'os';
import yaml from 'js-yaml';
import path from 'path';

import { loadContract } from './loading';

describe('contract/loading', () => {
  describe('loadContract', () => {
    const fileWithSyntaxErrors = path.join(tmpdir(), `${uuidv4()}-syntax-errors.json`);
    const fileWithValidJson = path.join(tmpdir(), `${uuidv4()}-valid.json`);
    const fileWithValidYaml = path.join(tmpdir(), `${uuidv4()}-valid.yml`);

    beforeAll(async () => {
      await Promise.all([
        fs.writeFile(fileWithSyntaxErrors, '{$$$42:'),
        fs.writeFile(fileWithValidJson, JSON.stringify({type: 'json'})),
        fs.writeFile(fileWithValidYaml, yaml.dump({type: 'yaml'}))
      ]);
    });

    afterAll(async () => {
      await Promise.all([
        fs.unlink(fileWithSyntaxErrors),
        fs.unlink(fileWithValidJson),
        fs.unlink(fileWithValidYaml)
      ]);
    });

    it('must reject when file cannot be found', async () => {
      await expect(loadContract(path.join(tmpdir(), uuidv4())))
        .rejects.toThrow(/Failed to read contract file.*/);
    });

    it('must reject when file cannot be parsed', async () => {
      await expect(loadContract(fileWithSyntaxErrors))
        .rejects.toThrow(/Failed to parse contract file.*/);
    });

    it('must successfully load JSON files', async () => {
      await expect(loadContract(fileWithValidJson)).resolves.toMatchInlineSnapshot(`
Object {
  "type": "json",
}
`);
    });

    it('must successfully load YAML files', async () => {
      await expect(loadContract(fileWithValidYaml)).resolves.toMatchInlineSnapshot(`
Object {
  "type": "yaml",
}
`);
    });
  });

});
