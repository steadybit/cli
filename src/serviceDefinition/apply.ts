// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Response } from 'node-fetch';
import colors from 'colors/safe';
import yaml from 'js-yaml';

import { loadServiceDefinition, writeServiceDefinition } from './files';
import { abortExecutionWithError } from '../errors';
import { ServiceDefinition } from './types';
import { executeApiCall } from '../api';
import { v4 as uuidv4 } from 'uuid';
import { addVcsTags } from '../vcs';

export interface Options {
  file: string;
}

const finishHelp = `
${colors.green('Done!')} The service definition was successfully applied.
You may now verify how compliant the service is by executing:

                   ${colors.bold('steadybit service verify')}
`.trim();

export async function apply(options: Options) {
  let serviceDefinition = await loadServiceDefinition(options.file);
  serviceDefinition = await initializeIdIfNecessary(serviceDefinition, options.file);
  await addVcsTags(options.file, serviceDefinition);

  try {
    await executeApiCall({
      method: 'PUT',
      path: `/api/service-definitions/${encodeURIComponent(serviceDefinition.id)}`,
      body: serviceDefinition,
    });
    console.log(finishHelp);
  } catch (e: any) {
    const response: Response | undefined = e.response;
    if (response?.status === 409) {
      await handleConflict(serviceDefinition);
    } else {
      const error = await abortExecutionWithError(e, 'Failed to upload service definition to Steadybit %s', serviceDefinition.id);
      throw error;
    }
  }
}

async function initializeIdIfNecessary(serviceDefinition: ServiceDefinition, serviceDefinitionPath: string) {
  if (!('id' in serviceDefinition)) {
    //we have to place to place the id first so it appears on top in file
    const updated = { id: uuidv4(), ...(serviceDefinition as Omit<ServiceDefinition, 'id'>)};
    await writeServiceDefinition(serviceDefinitionPath, updated);
    return updated;
  }
  return serviceDefinition;
}

const conflictHelp = `
Please resolve the conflict by editing the local service definition file.
For example, by ensuring that the same ID is used between the local
service definition and the remote service definition.

${colors.gray('Tip: Most often these conflicts are caused by different IDs for the same mapping.')}
`;

async function handleConflict(serviceDefinition: ServiceDefinition): Promise<void> {
  console.log(colors.red('Service definition application resulted in a conflict.'));

  try {
    const response = await executeApiCall({
      method: 'POST',
      path: '/api/service-definitions/by-mapping',
      body: serviceDefinition.mapping,
    });

    const conflictingServiceDefinition: ServiceDefinition = await response.json();
    console.log('The following conflicting service definition is stored:')
    console.log();
    console.log(yaml.dump(conflictingServiceDefinition));
    console.log();
    console.log(conflictHelp.trim());
    process.exit(1);
  } catch (e) {
    await abortExecutionWithError(
      e,
      'Failed to provide service definition conflict details for service definition %s',
      serviceDefinition.id
    );
  }
}
