// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import type { Schemas } from '../api/schemas.ts';
import { ApiError } from '../api/error.ts';
import { executeApiCall } from '../api/http.ts';
import { abortExecution, abortExecutionWithError } from '../errors.ts';

export interface TemplateFilter {
  tag?: string[];
  targetType?: string[];
  action?: string[];
  search?: string[];
}

export async function fetchTemplates(filter: TemplateFilter): Promise<Schemas['ExperimentTemplateSummaryAO'][]> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: '/api/experiments/templates',
      queryParameters: {
        tag: filter.tag,
        targetType: filter.targetType,
        action: filter.action,
        freeTextPhrases: filter.search,
      },
    });
    return ((await response.json()) as Schemas['ExperimentTemplateSummariesAO']).templates ?? [];
  } catch (e) {
    throw abortExecutionWithError(e, 'Failed to get the experiment templates');
  }
}

export async function fetchTemplate(id: string): Promise<Schemas['ExperimentTemplateAO']> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: `/api/experiments/templates/${encodeURIComponent(id)}`,
    });
    return (await response.json()) as Schemas['ExperimentTemplateAO'];
  } catch (e) {
    if (e instanceof ApiError && e.status === 404) {
      throw abortExecution('Experiment template %s not found.', id);
    }
    throw abortExecutionWithError(e, 'Failed to get experiment template %s', id);
  }
}
