/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import fs from 'fs/promises';
import yaml from 'js-yaml';

import { abortExecution } from '../errors';
import { ServiceDefinition } from './types';

export async function loadServiceDefinition(serviceDefinitionPath: string): Promise<ServiceDefinition> {
  let fileContent: string;
  try {
    fileContent = await fs.readFile(serviceDefinitionPath, { encoding: 'utf8' });
  } catch (e) {
    throw abortExecution(
      `Failed to read service definition file at path '%s': %s`,
      serviceDefinitionPath,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }

  try {
    return yaml.load(fileContent) as ServiceDefinition;
  } catch (e) {
    throw abortExecution(
      `Failed to parse service definition file at path '%s' as YAML/JSON: %s`,
      serviceDefinitionPath,
      (e as Error)?.message ?? 'Unknown Cause'
    );
  }
}

export async function writeServiceDefinition(serviceDefinitionPath: string, serviceDefinition: ServiceDefinition): Promise<void> {
  const fileContent = yaml.dump(serviceDefinition);
  await fs.writeFile(serviceDefinitionPath, fileContent, { encoding: 'utf8' });
}
