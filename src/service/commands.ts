// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import fs from 'node:fs/promises';
import { abortExecution } from '../errors.ts';
import { format, output, readStructuredFile } from '../structuredFiles.ts';
import { resolveExperimentFiles } from '../experiment/files.ts';
import { resolvePlaceholders } from '../experiment/template.ts';
import { createTable } from '../table.ts';
import {
  fetchService,
  fetchServiceExperiments,
  fetchServiceRisk,
  fetchServices,
  fetchServiceVariables,
  linkCustomExperiment,
  removeService,
  type Service,
  type ServiceExperimentFilter,
  type ServiceFilter,
  unlinkCustomExperiment,
  type UpsertService,
  upsertProvidedExperiment,
  upsertService,
  type Variables,
  writeServiceVariables,
} from './api.ts';

export async function listServices(options: ServiceFilter) {
  const services = await fetchServices(options);
  if (services.length === 0) {
    console.log('No services found.');
    return;
  }
  const table = createTable({
    columns: [
      { name: 'id', title: 'Id', alignment: 'left' },
      { name: 'name', title: 'Name', alignment: 'left' },
      { name: 'team', title: 'Team', alignment: 'left' },
      { name: 'environment', title: 'Environment', alignment: 'left' },
    ],
  });
  table.addRows(services.map(s => ({ id: s.id, name: s.name, team: s.team, environment: s.environment })));
  table.printTable();
}

export interface GetOptions {
  id: string;
  file?: string;
  type?: string;
}

export async function getService(options: GetOptions) {
  await output(toFileContent(await fetchService(options.id)), options);
  if (options.file) {
    console.log('Service %s written to %s.', options.id, options.file);
  }
}

// Who created and edited the service is not part of what can be sent. The version is
// left out as `experiment get` does: kept in a file, it turns every apply after a change
// made in the UI into a conflict.
function toFileContent(service: Service): UpsertService {
  const content: Partial<Service> = { ...service };
  delete content.created;
  delete content.createdBy;
  delete content.edited;
  delete content.editedBy;
  delete content.version;
  return content as UpsertService;
}

export interface ApplyOptions {
  file: string[];
  recursive: boolean;
  deleteExperiments: boolean;
}

export async function applyServices(options: ApplyOptions) {
  const files = await resolveExperimentFiles(options.file, options.recursive);
  for (const file of files) {
    const { content, datatype } = await readStructuredFile<Service>(file, 'service');
    if (!content?.name) {
      throw abortExecution("Service file '%s' does not name the service.", file);
    }
    const result = await upsertService(toFileContent(content), options.deleteExperiments);
    if (!content.id) {
      // The next apply of the file then updates this service instead of creating another.
      const withId = Object.assign({ id: result.service.id }, content, { id: result.service.id });
      await fs.writeFile(file, format(withId, datatype), { encoding: 'utf8' });
    }
    console.log('Service %s (%s) %s.', result.service.name, result.service.id, result.created ? 'created' : 'updated');
  }
}

export interface IdOptions {
  id: string;
}

export async function deleteService(options: IdOptions) {
  await removeService(options.id);
  console.log('Service %s deleted.', options.id);
}

export interface RiskOptions {
  id: string;
  type?: string;
  failAbove?: number;
}

export async function showServiceRisk(options: RiskOptions) {
  const risk = await fetchServiceRisk(options.id);

  if (options.type) {
    await output(risk, { type: options.type });
  } else {
    console.log(
      'Risk of service %s: %s (calculated %s)',
      options.id,
      risk.risk ?? 'unknown',
      risk.lastCalculated ?? 'never'
    );
    const categories = Object.entries(risk.categoryRisks ?? {});
    if (categories.length > 0) {
      const table = createTable({
        columns: [
          { name: 'category', title: 'Category', alignment: 'left' },
          { name: 'total', title: 'Total' },
          { name: 'experiment', title: 'Experiments' },
          { name: 'advice', title: 'Advice' },
        ],
      });
      table.addRows(
        categories.map(([category, r]) => ({ category, total: r.total, experiment: r.experiment, advice: r.advice }))
      );
      table.printTable();
    }
    if ((risk.experimentRisks ?? []).length > 0) {
      const table = createTable({
        columns: [
          { name: 'experimentKey', title: 'Experiment', alignment: 'left' },
          { name: 'risk', title: 'Risk' },
        ],
      });
      table.addRows(risk.experimentRisks ?? []);
      table.printTable();
    }
  }

  // Lets a pipeline stop a rollout of a service whose risk is too high, the way
  // `advice validate-status` does for advice.
  if (options.failAbove !== undefined && (risk.risk === undefined || risk.risk > options.failAbove)) {
    throw abortExecution(
      'Risk of service %s is %s, above the accepted %d.',
      options.id,
      risk.risk ?? 'unknown',
      options.failAbove
    );
  }
}

export interface ExperimentListOptions extends ServiceExperimentFilter {
  id: string;
}

export async function listServiceExperiments(options: ExperimentListOptions) {
  const experiments = await fetchServiceExperiments(options.id, {
    category: options.category,
    type: options.type?.map(t => t.toUpperCase()),
  });
  if (experiments.length === 0) {
    const filtered = (options.category?.length ?? 0) > 0 || (options.type?.length ?? 0) > 0;
    console.log(filtered ? 'Service %s has no matching experiments.' : 'Service %s has no experiments.', options.id);
    return;
  }
  const table = createTable({
    columns: [
      { name: 'category', title: 'Category', alignment: 'left' },
      { name: 'associationType', title: 'Type', alignment: 'left' },
      // A provided experiment that has not been created yet has no key, only the
      // template it would be created from.
      { name: 'experimentKey', title: 'Experiment', alignment: 'left' },
      { name: 'templateId', title: 'Template', alignment: 'left' },
    ],
  });
  table.addRows(
    experiments.map(e => ({
      category: e.category,
      associationType: e.associationType,
      experimentKey: e.experimentKey ?? '(not created)',
      templateId: e.templateId ?? '',
    }))
  );
  table.printTable();
}

export interface LinkOptions {
  id: string;
  experiment: string;
  category: string;
}

export async function linkExperiment(options: LinkOptions) {
  await linkCustomExperiment(options.id, options.experiment, options.category);
  console.log('Experiment %s linked to service %s in category %s.', options.experiment, options.id, options.category);
}

export interface UnlinkOptions {
  id: string;
  experiment: string;
}

export async function unlinkExperiment(options: UnlinkOptions) {
  await unlinkCustomExperiment(options.id, options.experiment);
  console.log('Experiment %s unlinked from service %s.', options.experiment, options.id);
}

export interface ProvideOptions {
  id: string;
  template: string;
  experiment?: string;
  placeholder?: Record<string, string>;
  placeholders?: string;
  resetProperties: boolean;
}

export async function provideExperiment(options: ProvideOptions) {
  const placeholders = await resolvePlaceholders(options);
  const result = await upsertProvidedExperiment(
    options.id,
    { templateId: options.template, experimentKey: options.experiment, placeholders },
    options.resetProperties
  );
  console.log(
    'Provided experiment %s of service %s %s from template %s.',
    result.key ?? '',
    options.id,
    result.created ? 'created' : 'updated',
    options.template
  );
}

export interface VariableGetOptions {
  id: string;
  type?: string;
}

export async function getServiceVariables(options: VariableGetOptions) {
  await output(await fetchServiceVariables(options.id), { type: options.type });
}

export interface VariableSetOptions {
  id: string;
  file?: string;
  replace: boolean;
}

// Values from a file may be lists or select expressions; KEY=VALUE arguments are always
// plain strings and override entries of the file.
export async function setServiceVariables(pairs: string[], options: VariableSetOptions) {
  const given: Record<string, string> = {};
  for (const pair of pairs) {
    const separator = pair.indexOf('=');
    if (separator <= 0) {
      throw abortExecution("'%s' is not in the form KEY=VALUE.", pair);
    }
    given[pair.slice(0, separator)] = pair.slice(separator + 1);
  }

  let variables: Variables = {};
  if (options.file) {
    const { content } = await readStructuredFile<Variables>(options.file, 'variables');
    if (!content || typeof content !== 'object' || Array.isArray(content)) {
      throw abortExecution("Variables file '%s' must be a map of variable names to values.", options.file);
    }
    variables = content;
  }
  variables = { ...variables, ...given };
  if (Object.keys(variables).length === 0 && !options.replace) {
    throw abortExecution('No variables given. Pass KEY=VALUE arguments or --file.');
  }
  await writeServiceVariables(options.id, variables, options.replace);
  console.log(
    '%d variable(s) of service %s %s.',
    Object.keys(variables).length,
    options.id,
    options.replace ? 'set, all others removed' : 'set'
  );
}
