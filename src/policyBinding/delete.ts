// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { loadPolicyBinding } from './files';
import { abortExecution } from '../errors';
import { executeApiCall } from '../api/http';

export interface Options {
  id?: string;
  file: string;
}

export async function deletePolicyBinding(options: Options) {
  if(!options.id && !options.file){
    throw abortExecution(
      'Failed to delete policy binding, no id nor path to the policy binding file was given',
    );
  }
  let id = options.id;
  if(!options.id && options.file){
    const policyBinding = await loadPolicyBinding(options.file);
    id = policyBinding.id ?? id;
  }
  // assuming that the parameter id is given

  try {
    await executeApiCall({
      method: 'DELETE',
      path: `/api/policy-bindings/${encodeURIComponent(String(id))}`,
    });
  } catch (e) {
    throw abortExecution(
      'Failed to delete policy binding with id %s: %s',
      id,
      (e as Error)?.message ?? 'Unknown error'
    );
  }
}
