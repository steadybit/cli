// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import fs from 'fs/promises';
import path from 'path';
import os from 'os';
import yaml from 'js-yaml';

let tempDir: string;

export async function createTempDir() {
  tempDir = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-cli-test'));
}

export async function writeYamlFile(name: string, content: any) {
  const file = path.join(tempDir, name);
  await fs.mkdir(path.dirname(file), { recursive: true });
  await fs.writeFile(file, yaml.dump(content));
  return file;
}

export function getTempDir() {
  return tempDir;
}

export async function removeTempDir() {
  await fs.rm(tempDir, { recursive: true });
}
