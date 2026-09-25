// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import { abortExecution } from '../errors.ts';
import { format, output, readStructuredFile } from '../structuredFiles.ts';
import { resolveExperimentFiles } from '../experiment/files.ts';
import { createTable } from '../table.ts';
import {
  fetchServiceProfile,
  fetchServiceProfiles,
  removeServiceProfile,
  type ServiceProfile,
  type ServiceProfileFilter,
  type UpsertServiceProfile,
  upsertServiceProfile,
} from './api.ts';

export async function listServiceProfiles(options: ServiceProfileFilter) {
  const profiles = await fetchServiceProfiles({ ...options, origin: options.origin?.map(o => o.toUpperCase()) });
  if (profiles.length === 0) {
    console.log('No service profiles found.');
    return;
  }
  const table = createTable({
    columns: [
      { name: 'id', title: 'Id', alignment: 'left' },
      { name: 'name', title: 'Name', alignment: 'left' },
      { name: 'origin', title: 'Origin', alignment: 'left' },
      { name: 'defaultProfile', title: 'Default', alignment: 'left' },
      { name: 'templates', title: 'Templates' },
    ],
  });
  table.addRows(
    profiles.map(p => ({
      id: p.id,
      name: p.name,
      origin: p.origin,
      defaultProfile: String(p.defaultProfile),
      templates: p.templates.reduce((count, category) => count + (category.templateIds?.length ?? 0), 0),
    }))
  );
  table.printTable();
}

export interface GetOptions {
  id: string;
  file?: string;
  type?: string;
}

export async function getServiceProfile(options: GetOptions) {
  await output(toFileContent(await fetchServiceProfile(options.id)), options);
  if (options.file) {
    console.log('Service profile %s written to %s.', options.id, options.file);
  }
}

// Only what can be sent back is kept, so that `get` followed by `apply` is a round trip.
// Whether a profile is the default is not part of it; it is changed in the platform.
function toFileContent(profile: ServiceProfile): UpsertServiceProfile {
  const content: Partial<ServiceProfile> = { ...profile };
  delete content.created;
  delete content.createdBy;
  delete content.edited;
  delete content.editedBy;
  delete content.version;
  delete content.defaultProfile;
  return content as UpsertServiceProfile;
}

export interface ApplyOptions {
  file: string[];
  recursive: boolean;
  deleteExperiments: boolean;
}

export async function applyServiceProfiles(options: ApplyOptions) {
  const files = await resolveExperimentFiles(options.file, options.recursive);
  for (const file of files) {
    const { content, datatype } = await readStructuredFile<ServiceProfile>(file, 'service profile');
    if (!content?.name) {
      throw abortExecution("Service profile file '%s' does not name the profile.", file);
    }
    // Profiles written by hand are the team's own; PROVIDED ones come from Steadybit.
    const upsert = toFileContent(content);
    const profile: UpsertServiceProfile = { ...upsert, origin: upsert.origin ?? 'CUSTOM' };
    const result = await upsertServiceProfile(profile, options.deleteExperiments);
    if (!content.id) {
      const withId = Object.assign({ id: result.profile.id }, content, { id: result.profile.id });
      await fs.writeFile(file, format(withId, datatype), { encoding: 'utf8' });
    }
    console.log(
      'Service profile %s (%s) %s.',
      result.profile.name,
      result.profile.id,
      result.created ? 'created' : 'updated'
    );
  }
}

export interface IdOptions {
  id: string;
}

export async function deleteServiceProfile(options: IdOptions) {
  await removeServiceProfile(options.id);
  console.log('Service profile %s deleted.', options.id);
}
