import fs from 'fs/promises';
import yaml from 'js-yaml';

import { abortExecution } from '../errors';
import { Contract } from './types';

export async function loadContract(contractPath: string): Promise<Contract> {
  let fileContent: string;
  try {
    fileContent = await fs.readFile(contractPath, { encoding: 'utf8' });
  } catch (e) {
    throw abortExecution(
      `Failed to read contract file at path '%s': %s`,
      contractPath,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }

  try {
    return yaml.load(fileContent) as Contract;
  } catch (e) {
    throw abortExecution(
      `Failed to parse contract file at path '%s' as YAML/JSON: %s`,
      contractPath,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }
}
