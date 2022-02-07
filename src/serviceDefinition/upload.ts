import { executeApiCall } from '../api';
import { abortExecution } from '../errors';
import { ServiceDefinition } from './types';

export async function uploadServiceDefinition(serviceDefinition: ServiceDefinition): Promise<void> {
  try {
    await executeApiCall({
      method: 'PUT',
      path: `/api/service-definition/${encodeURIComponent(serviceDefinition.id)}`,
      body: serviceDefinition
    });
  } catch (e) {
    throw abortExecution(`Failed to upload service definition to Steadybit: %s`, (e as Error)?.message ?? 'Unknown error');
  }
}
