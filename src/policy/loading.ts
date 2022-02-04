import fs from 'fs/promises';
import yaml from 'js-yaml';

import { abortExecution } from '../errors';
import { Policy } from './types';

export async function loadPolicy(policyPath: string): Promise<Policy> {
  let fileContent: string;
  try {
    fileContent = await fs.readFile(policyPath, { encoding: 'utf8' });
  } catch (e) {
    throw abortExecution(
      `Failed to read policy file at path '%s': %s`,
      policyPath,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }

  try {
    return yaml.load(fileContent) as Policy;
  } catch (e) {
    throw abortExecution(
      `Failed to parse policy file at path '%s' as YAML/JSON: %s`,
      policyPath,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }
}
