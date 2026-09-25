// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import { dump, load } from './yaml.ts';
import { abortExecution, errorMessage } from './errors.ts';

export type Datatype = 'json' | 'yaml';

// Explicit type first, then the file's extension, then YAML, which is what people put
// into Git repositories.
export function resolveDatatype(type: string | undefined, file?: string): Datatype {
  if (type === 'json' || type === 'yaml') {
    return type;
  }
  if (type) {
    throw abortExecution('Unsupported output format \'%s\'. Use "json" or "yaml".', type);
  }
  return file?.toLowerCase().endsWith('.json') ? 'json' : 'yaml';
}

export function format(content: unknown, datatype: Datatype): string {
  return datatype === 'json' ? JSON.stringify(content, undefined, 2) : dump(content);
}

// Writes to the file when one is given, to stdout otherwise, so that a command's output
// can be piped as readily as saved.
export async function output(content: unknown, options: { file?: string; type?: string }): Promise<void> {
  const datatype = resolveDatatype(options.type, options.file);
  if (options.file) {
    await fs.writeFile(options.file, format(content, datatype), { encoding: 'utf8' });
  } else {
    console.log(format(content, datatype));
  }
}

export async function readStructuredFile<T>(file: string, what: string): Promise<{ content: T; datatype: Datatype }> {
  let text: string;
  try {
    text = await fs.readFile(file, { encoding: 'utf8' });
  } catch (e) {
    throw abortExecution("Failed to read %s file at path '%s': %s", what, file, errorMessage(e));
  }

  try {
    return { content: JSON.parse(text) as T, datatype: 'json' };
  } catch {
    try {
      return { content: load(text) as T, datatype: 'yaml' };
    } catch (e) {
      throw abortExecution("Failed to parse %s file at path '%s' as YAML/JSON: %s", what, file, errorMessage(e));
    }
  }
}
