// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { loadServiceDefinition } from './files';
import { abortExecution } from '../errors';
import { executeApiCall } from '../api/http';

export interface Options {
  id?: string;
  file: string;
}

export async function deleteServiceDefinition(options: Options) {
  if(!options.id && !options.file){
    throw abortExecution(
      'Failed to delete service definition, no id nor path to the service definition file was given',
    );
  }
  let id = options.id;
  if(!options.id && options.file){
    const serviceDefinition = await loadServiceDefinition(options.file);
    id = serviceDefinition.id ?? id;
  }
  // assuming that the parameter id is given

  try {
    await executeApiCall({
      method: 'DELETE',
      path: `/api/service-definitions/${encodeURIComponent(String(id))}`,
    });
  } catch (e) {
    throw abortExecution(
      'Failed to delete service definition with id %s: %s',
      id,
      (e as Error)?.message ?? 'Unknown error'
    );
  }
}
