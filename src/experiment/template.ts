// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import type { Schemas } from '../api/schemas.ts';
import { abortExecution } from '../errors.ts';
import { readStructuredFile } from '../structuredFiles.ts';
import { createExperimentFromTemplate, updateExperimentFromTemplate } from './api.ts';

type PlaceholderValue = Schemas['ExperimentTemplatePlaceholderValueAO'];

export interface TemplateOptions {
  template: string;
  team?: string;
  environment?: string;
  externalId?: string;
  placeholder?: Record<string, string>;
  placeholders?: string;
  variable?: Record<string, string>;
  resetProperties: boolean;
}

// A placeholders file is either a map of key to value, which is what someone writes by
// hand, or the platform's own list of {key, value} pairs, which is what they copy out
// of an API call. Placeholders given with -p are applied on top, so a pipeline can keep
// the shared values in a file and override one per stage.
export async function resolvePlaceholders(options: TemplateOptions): Promise<PlaceholderValue[]> {
  const values = new Map<string, unknown>();

  if (options.placeholders) {
    const { content } = await readStructuredFile<unknown>(options.placeholders, 'placeholders');
    if (Array.isArray(content)) {
      for (const entry of content) {
        if (typeof entry?.key !== 'string' || !('value' in entry)) {
          throw abortExecution(
            "Placeholders file '%s' must be a map of key to value or a list of {key, value} entries.",
            options.placeholders
          );
        }
        values.set(entry.key, entry.value);
      }
    } else if (content && typeof content === 'object') {
      Object.entries(content).forEach(([key, value]) => values.set(key, value));
    } else {
      throw abortExecution(
        "Placeholders file '%s' must be a map of key to value or a list of {key, value} entries.",
        options.placeholders
      );
    }
  }

  Object.entries(options.placeholder ?? {}).forEach(([key, value]) => values.set(key, value));

  // The spec types a placeholder value as an object although it is any JSON value.
  return [...values].map(([key, value]) => ({ key, value: value as PlaceholderValue['value'] }));
}

export async function toCreateFromTemplate(
  options: TemplateOptions
): Promise<Schemas['CreateExperimentFromTemplateAO']> {
  if (!options.team) {
    throw abortExecution('--team is required to create an experiment from a template.');
  }
  return {
    team: options.team,
    environment: options.environment,
    externalId: options.externalId,
    placeholders: await resolvePlaceholders(options),
    experimentVariables: options.variable,
  };
}

export interface ApplyFromTemplateOptions extends TemplateOptions {
  key?: string;
}

export async function applyExperimentFromTemplate(options: ApplyFromTemplateOptions): Promise<void> {
  if (options.key) {
    // Updating an experiment only re-renders it with new placeholder values. It keeps
    // its team and environment, so accepting those here would silently do nothing.
    const ignored = (['team', 'environment', 'externalId', 'variable'] as const).filter(o => options[o] !== undefined);
    if (ignored.length > 0) {
      throw abortExecution(
        'Updating experiment %s from a template only takes placeholders; remove %s.',
        options.key,
        ignored.map(o => `--${o === 'externalId' ? 'external-id' : o}`).join(', ')
      );
    }
    await updateExperimentFromTemplate(
      options.template,
      options.key,
      { placeholders: await resolvePlaceholders(options) },
      { resetProperties: options.resetProperties }
    );
    console.log('Experiment %s updated from template %s.', options.key, options.template);
    return;
  }

  const result = await createExperimentFromTemplate(options.template, await toCreateFromTemplate(options), {
    resetProperties: options.resetProperties,
  });
  console.log(
    'Experiment %s %s from template %s.',
    result.key,
    result.created ? 'created' : 'updated',
    options.template
  );
}
