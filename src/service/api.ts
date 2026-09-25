// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import type { Schemas } from '../api/schemas.ts';
import { ApiError } from '../api/error.ts';
import { executeApiCall } from '../api/http.ts';
import { fetchAllPages } from '../api/paging.ts';
import { abortExecution, abortExecutionWithError } from '../errors.ts';
import { refusedForProvidedExperiments } from './refusal.ts';

export type Service = Schemas['ServiceAO'];
export type UpsertService = Schemas['UpsertServiceAO'];
export type ServiceSummary = Schemas['ServiceSummaryAO'];
export type ServiceExperiment = Schemas['ServiceExperimentAO'];
export type ServiceRisk = Schemas['ServiceRiskAO'];
export type Variables = NonNullable<UpsertService['variables']>;

function notFoundOr(e: unknown, id: string, msg: string): Error {
  if (e instanceof ApiError && e.status === 404) {
    return abortExecution('Service %s not found.', id);
  }
  return abortExecutionWithError(e, msg, id);
}

const servicePath = (id: string) => `/api/services/${encodeURIComponent(id)}`;

export interface ServiceFilter {
  team?: string[];
  environment?: string[];
  experiment?: string[];
}

export async function fetchServices(filter: ServiceFilter): Promise<ServiceSummary[]> {
  try {
    return await fetchAllPages<ServiceSummary>('/api/services', {
      teamKey: filter.team,
      environmentName: filter.environment,
      experimentKey: filter.experiment,
    });
  } catch (e) {
    throw abortExecutionWithError(e, 'Failed to get the services');
  }
}

export async function fetchService(id: string): Promise<Service> {
  try {
    return (await (await executeApiCall({ method: 'GET', path: servicePath(id) })).json()) as Service;
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to get service %s');
  }
}

export async function upsertService(
  service: UpsertService,
  deleteExperiments: boolean
): Promise<{ created: boolean; service: Service }> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: '/api/services',
      queryParameters: { deleteExperiments: String(deleteExperiments) },
      body: service,
    });
    return { created: response.status === 201, service: (await response.json()) as Service };
  } catch (e) {
    if (!deleteExperiments && refusedForProvidedExperiments(e)) {
      throw abortExecution(
        'Service %s was not saved: the change would remove provided experiments. Pass --delete-experiments to delete them.',
        service.name
      );
    }
    throw abortExecutionWithError(e, 'Failed to save service %s', service.name);
  }
}

export async function removeService(id: string): Promise<void> {
  try {
    await executeApiCall({ method: 'DELETE', path: servicePath(id) });
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to delete service %s');
  }
}

export async function fetchServiceRisk(id: string): Promise<ServiceRisk> {
  try {
    return (await (await executeApiCall({ method: 'GET', path: `${servicePath(id)}/risk` })).json()) as ServiceRisk;
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) {
      throw abortExecution('Service %s not found, or its risk has not been calculated yet.', id);
    }
    throw abortExecutionWithError(e, 'Failed to get the risk of service %s', id);
  }
}

export interface ServiceExperimentFilter {
  category?: string[];
  type?: string[];
}

export async function fetchServiceExperiments(
  id: string,
  filter: ServiceExperimentFilter
): Promise<ServiceExperiment[]> {
  try {
    return await fetchAllPages<ServiceExperiment>(`${servicePath(id)}/experiments`, {
      category: filter.category,
      type: filter.type,
    });
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to get the experiments of service %s');
  }
}

export async function linkCustomExperiment(id: string, experimentKey: string, category: string): Promise<void> {
  try {
    await executeApiCall({
      method: 'POST',
      path: `${servicePath(id)}/experiments/custom`,
      body: { experimentKey, category } satisfies Schemas['LinkCustomExperimentRequestAO'],
    });
  } catch (e) {
    throw notFoundOr(e, id, `Failed to link experiment ${experimentKey} to service %s`);
  }
}

export async function unlinkCustomExperiment(id: string, experimentKey: string): Promise<void> {
  try {
    await executeApiCall({
      method: 'DELETE',
      path: `${servicePath(id)}/experiments/custom`,
      queryParameters: { experimentKey },
    });
  } catch (e) {
    throw notFoundOr(e, id, `Failed to unlink experiment ${experimentKey} from service %s`);
  }
}

// Like a create from a template, the platform names the experiment in the Location
// header and sends no body.
export async function upsertProvidedExperiment(
  id: string,
  request: Schemas['UpsertProvidedExperimentRequestAO'],
  resetProperties: boolean
): Promise<{ created: boolean; key?: string }> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: `${servicePath(id)}/experiments/provided`,
      queryParameters: { resetProperties: String(resetProperties) },
      body: request,
    });
    const location = response.headers.get('Location');
    return {
      created: response.status === 201,
      key: location ? location.substring(location.lastIndexOf('/') + 1) : (request.experimentKey ?? undefined),
    };
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) {
      throw abortExecution('Service %s or experiment template %s not found.', id, request.templateId);
    }
    throw abortExecutionWithError(e, 'Failed to save the provided experiment of service %s', id);
  }
}

export async function fetchServiceVariables(id: string): Promise<Variables> {
  try {
    return (await (await executeApiCall({ method: 'GET', path: `${servicePath(id)}/variables` })).json()) as Variables;
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to get the variables of service %s');
  }
}

// PUT replaces every variable of the service, PATCH merges the given ones in.
export async function writeServiceVariables(id: string, variables: Variables, replace: boolean): Promise<void> {
  try {
    await executeApiCall({ method: replace ? 'PUT' : 'PATCH', path: `${servicePath(id)}/variables`, body: variables });
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to update the variables of service %s');
  }
}
