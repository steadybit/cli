import opn from 'open';

import { loadServiceDefinition } from './loading';
import { baseUrl } from '../api';

export async function open(serviceDefinitionPath: string) {
  const serviceDefinition = await loadServiceDefinition(serviceDefinitionPath);
  await opn(`${baseUrl}/service/${serviceDefinition.id}`);
}
