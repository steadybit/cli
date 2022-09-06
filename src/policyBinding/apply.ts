// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Response } from 'node-fetch';
import colors from 'colors/safe';
import yaml from 'js-yaml';

import { abortExecutionWithError } from '../errors';
import { loadPolicyBinding } from './files';
import { executeApiCall } from '../api/http';
import { PolicyBinding } from './types';
import { addVcsTags } from '../vcs';

export interface Options {
  file: string;
}

const finishHelp = `
${colors.green('Done!')} The policy binding was successfully applied.
You may now verify how compliant you are by executing:

                   ${colors.bold('steadybit policy-binding verify')}
`.trim();

export async function apply(options: Options) {
  const policyBinding = await loadPolicyBinding(options.file);
  await addVcsTags(options.file, policyBinding);

  try {
    if (policyBindingIsNew(policyBinding)) {
      await executeApiCall({
        method: 'POST',
        path: '/api/policy-bindings',
        body: policyBinding,
      });
    } else {
      await executeApiCall({
        method: 'PUT',
        path: `/api/policy-bindings/${encodeURIComponent(policyBinding.id)}`,
        body: policyBinding,
      });
    }
    console.log(finishHelp);
  } catch (e: any) {
    const response: Response | undefined = e.response;
    if (response?.status === 409) {
      await handleConflict(policyBinding);
    } else {
      const error = await abortExecutionWithError(
        e,
        'Failed to upload policy binding to Steadybit %s',
        policyBinding.id ?? '(no id)'
      );
      throw error;
    }
  }
}

function policyBindingIsNew(policyBinding: PolicyBinding): boolean {
  return !('id' in policyBinding);
}

const conflictHelp = `
Please resolve the conflict by editing the local policy binding file.
For example, by ensuring that the same ID is used between the local
policy binding and the remote policy binding.

${colors.gray('Tip: Most often these conflicts are caused by different IDs for the same mapping.')}
`;

async function handleConflict(policyBinding: PolicyBinding): Promise<void> {
  console.log(colors.red('Policy binding resulted in a conflict.'));

  try {
    const response = await executeApiCall({
      method: 'POST',
      path: '/api/policy-bindings/by-mapping',
      body: policyBinding.mapping,
    });

    const conflictingPolicyBinding: PolicyBinding = await response.json();
    console.log('The following conflicting policy binding is stored:');
    console.log();
    console.log(yaml.dump(conflictingPolicyBinding));
    console.log();
    console.log(conflictHelp.trim());
    process.exit(1);
  } catch (e) {
    await abortExecutionWithError(
      e,
      'Failed to provide policy binding conflict details for policy binding %s',
      policyBinding.id
    );
  }
}
