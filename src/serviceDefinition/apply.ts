import { uploadServiceDefinition } from './upload';
import { loadServiceDefinition } from './loading';

export async function apply(serviceDefinitionPath: string) {
  const serviceDefinition = await loadServiceDefinition(serviceDefinitionPath);
  await uploadServiceDefinition(serviceDefinition);
}
