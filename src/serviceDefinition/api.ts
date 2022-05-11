// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ServiceDefinition, ServiceState } from './types';

import { abortExecutionWithError } from '../errors';
import { executeApiCall } from '../api';

export async function getState(serviceDefinition: ServiceDefinition): Promise<ServiceState> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/service-states/${serviceDefinition.id}`,
    });
    return (await response.json()) as ServiceState;
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to read state for service %s', serviceDefinition.id);
    throw error;
  }
}
