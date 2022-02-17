import { abortExecutionWithError } from '../errors';
import { ServiceDefinition } from './types';
import { executeApiCall } from '../api';

export async function uploadServiceDefinition(serviceDefinition: ServiceDefinition): Promise<void> {
  try {
    await executeApiCall({
      method: 'PUT',
      path: `/api/service-definitions/${encodeURIComponent(serviceDefinition.id)}`,
      body: serviceDefinition
    });
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to upload service definition to Steadybit %s', serviceDefinition.id);
    throw error;
  }
}
