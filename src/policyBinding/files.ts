// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import fs from 'fs/promises';
import yaml from 'js-yaml';

import { abortExecution } from '../errors';
import { DefinePolicyBinding, PolicyBinding } from './types';

export async function loadPolicyBinding(policyBindingPath: string): Promise<PolicyBinding> {
  let fileContent: string;
  try {
    fileContent = await fs.readFile(policyBindingPath, { encoding: 'utf8' });
  } catch (e) {
    throw abortExecution(
      "Failed to read policy binding file at path '%s': %s",
      policyBindingPath,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }

  try {
    return yaml.load(fileContent) as PolicyBinding;
  } catch (e) {
    throw abortExecution(
      "Failed to parse policy binding file at path '%s' as YAML/JSON: %s",
      policyBindingPath,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }
}

export async function writePolicyBinding(
  policyBindingPath: string,
  policyBinding: DefinePolicyBinding
): Promise<void> {
  const fileContent = yaml.dump(policyBinding);
  await fs.writeFile(policyBindingPath, fileContent, { encoding: 'utf8' });
}
