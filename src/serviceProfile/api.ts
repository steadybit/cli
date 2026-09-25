// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import type { Schemas } from '../api/schemas.ts';
import { ApiError } from '../api/error.ts';
import { executeApiCall } from '../api/http.ts';
import { fetchAllPages } from '../api/paging.ts';
import { abortExecution, abortExecutionWithError } from '../errors.ts';
import { refusedForProvidedExperiments } from '../service/refusal.ts';

export type ServiceProfile = Schemas['ServiceProfileAO'];
export type UpsertServiceProfile = Schemas['UpsertServiceProfileAO'];

const profilePath = (id: string) => `/api/services/profiles/${encodeURIComponent(id)}`;

function notFoundOr(e: unknown, id: string, msg: string): Error {
  if (e instanceof ApiError && e.status === 404) {
    return abortExecution('Service profile %s not found.', id);
  }
  return abortExecutionWithError(e, msg, id);
}

export interface ServiceProfileFilter {
  name?: string;
  origin?: string[];
  default?: boolean;
}

export async function fetchServiceProfiles(filter: ServiceProfileFilter): Promise<ServiceProfile[]> {
  try {
    return await fetchAllPages<ServiceProfile>('/api/services/profiles', {
      name: filter.name,
      origin: filter.origin,
      defaultProfile: filter.default === undefined ? undefined : String(filter.default),
    });
  } catch (e) {
    throw abortExecutionWithError(e, 'Failed to get the service profiles');
  }
}

export async function fetchServiceProfile(id: string): Promise<ServiceProfile> {
  try {
    return (await (await executeApiCall({ method: 'GET', path: profilePath(id) })).json()) as ServiceProfile;
  } catch (e) {
    throw notFoundOr(e, id, 'Failed to get service profile %s');
  }
}

export async function upsertServiceProfile(
  profile: UpsertServiceProfile,
  deleteExperiments: boolean
): Promise<{ created: boolean; profile: ServiceProfile }> {
  try {
    const response = await executeApiCall({
      method: 'POST',
      path: '/api/services/profiles',
      queryParameters: { deleteExperiments: String(deleteExperiments) },
      body: profile,
    });
    return { created: response.status === 201, profile: (await response.json()) as ServiceProfile };
  } catch (e) {
    if (!deleteExperiments && refusedForProvidedExperiments(e)) {
      throw abortExecution(
        'Service profile %s was not saved: the change would remove provided experiments. Pass --delete-experiments to delete them.',
        profile.name
      );
    }
    throw abortExecutionWithError(e, 'Failed to save service profile %s', profile.name);
  }
}

export async function removeServiceProfile(id: string): Promise<void> {
  try {
    await executeApiCall({ method: 'DELETE', path: profilePath(id) });
  } catch (e) {
    if (e instanceof ApiError && e.status === 422) {
      throw abortExecution('Service profile %s is provided by Steadybit and cannot be deleted.', id);
    }
    throw notFoundOr(e, id, 'Failed to delete service profile %s');
  }
}
