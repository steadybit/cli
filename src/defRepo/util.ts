/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import globWithoutPromises from 'glob';
import { promisify } from 'util';
import fs from 'fs/promises';
import yaml from 'js-yaml';
import path from 'path';

import { abortExecutionWithError } from '../errors';

const glob = promisify(globWithoutPromises);

export type ForEachAction = 'update' | void;

export interface IteratorArgs<T> {
  fileName: string;
  /**
   * yaml-parsed file content. Can be mutated by the iterator. Returning ForEachAction.update
   * will cause this object to be yaml serialized and written back to disk.
   */
  fileContent: T;
}

/**
 * Return a non undefined/null value to update the file content.
 */
export type Iterator<T> = (args: IteratorArgs<T>) => Promise<ForEachAction>;

export interface ForEachFileArgs<T> {
  iterator: Iterator<T>;

  /**
   * https://www.npmjs.com/package/glob
   */
  pattern: string;
  options?: globWithoutPromises.IOptions;
}

export async function forEachFile<T>(args: ForEachFileArgs<T>): Promise<void> {
  let fileNames: string[];
  try {
    fileNames = await glob(args.pattern, args.options ?? {});
    fileNames = fileNames.map(fileName => path.join(args.options?.cwd ?? process.cwd(), fileName));
  } catch (e) {
    const error = await abortExecutionWithError(
      e,
      "Failed to identify definition files. CWD '%s'. Pattern '%s'.",
      __dirname,
      args.pattern
    );
    throw error;
  }

  await Promise.all(fileNames.map(fileName => iterate(args, fileName)));
}

async function iterate<T>(args: ForEachFileArgs<T>, fileName: string): Promise<void> {
  try {
    const fileContent = await fs.readFile(fileName, 'utf8');
    const parsedFileContent: T = yaml.load(fileContent) as T;
    const action = await args.iterator({ fileName, fileContent: parsedFileContent });
    if (action === 'update') {
      console.log(`Updating file ${path.relative(args.options?.cwd ?? process.cwd(), fileName)}`);
      await fs.writeFile(fileName, yaml.dump(parsedFileContent));
    }
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to iterate over file %s', fileName);
    throw error;
  }
}
