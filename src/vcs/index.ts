/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { ServiceDefinition } from '../serviceDefinition/types';
import { getTags as getGitTags } from './git';

export async function addVcsTags(serviceDefinitionPath: string, serviceDefinition: ServiceDefinition): Promise<void> {
  const gitTags = await getGitTags(serviceDefinitionPath);

  serviceDefinition.tags = {
    ...gitTags,
    ...(serviceDefinition.tags ?? {}),
  };
}
