// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ServiceDefinition } from '../serviceDefinition/types';
import { getTags as getGitTags } from './git';

export async function addVcsTags(serviceDefinitionPath: string, serviceDefinition: ServiceDefinition): Promise<void> {
  const gitTags = await getGitTags(serviceDefinitionPath);

  serviceDefinition.tags = {
    ...gitTags,
    ...(serviceDefinition.tags ?? {}),
  };
}
