/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { loadServiceDefinition } from './files';
import { abortExecution } from '../errors';
import { executeApiCall } from '../api';

export interface Options {
  file: string;
}

export async function deleteServiceDefinition(id: string, options: Options) {
  if(!id && !options.file){
    throw abortExecution(
      'Failed to delete service definition, no id nor path to the service definition file',
    );
  }

  if(!id && options.file){
    const serviceDefinition = await loadServiceDefinition(options.file);
    id = serviceDefinition.id ?? id;
  }
  // assuming that the parameter id is given

  try {
    await executeApiCall({
      method: 'DELETE',
      path: `/api/service-definitions/${encodeURIComponent(id)}`,
    });
  } catch (e) {
    throw abortExecution(
      'Failed to delete service definition with id %s: %s',
      id,
      (e as Error)?.message ?? 'Unknown error'
    );
  }
}
