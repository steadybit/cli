/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { loadServiceDefinition } from './files';
import { abortExecution } from '../errors';
import { executeApiCall } from '../api';

export async function deleteServiceDefinition(serviceDefinitionPathOrId: string) {
  let id = serviceDefinitionPathOrId;
  try {
    const serviceDefinition = await loadServiceDefinition(serviceDefinitionPathOrId);
    id = serviceDefinition.id ?? id;
  } catch (e) {
    // swallow silently assuming that the parameter is actually the ID and not a file path
  }

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
