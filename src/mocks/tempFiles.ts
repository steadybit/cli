// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import fs from 'node:fs/promises';
import path from 'node:path';
import os from 'node:os';
import { dump } from '../yaml.ts';
import type { Datatype } from '../experiment/files.ts';

let tempDir: string;

export async function createTempDir() {
  tempDir = await fs.mkdtemp(path.join(os.tmpdir(), 'steadybit-cli-test'));
}

export async function writeFile(name: string, content: any, datatype: Datatype = 'yaml') {
  const file = path.join(tempDir, name);
  await fs.mkdir(path.dirname(file), { recursive: true });
  await fs.writeFile(file, datatype === 'json' ? JSON.stringify(content) : dump(content));
  return file;
}

export function getTempDir() {
  return tempDir;
}

export async function removeTempDir() {
  await fs.rm(tempDir, { recursive: true });
}
