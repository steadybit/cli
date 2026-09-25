// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { createTable } from '../table.ts';
import { output } from '../structuredFiles.ts';
import { fetchTemplate, fetchTemplates, type TemplateFilter } from './api.ts';

export async function listTemplates(options: TemplateFilter) {
  const templates = await fetchTemplates(options);
  if (templates.length === 0) {
    console.log('No experiment templates found.');
    return;
  }
  const table = createTable({
    columns: [
      { name: 'id', title: 'Id', alignment: 'left' },
      { name: 'templateTitle', title: 'Title', alignment: 'left' },
    ],
  });
  table.addRows(templates.map(t => ({ id: t.id, templateTitle: t.templateTitle })));
  table.printTable();
}

export interface GetOptions {
  id: string;
  file?: string;
  type?: string;
  placeholders?: boolean;
}

export async function getTemplate(options: GetOptions) {
  const template = await fetchTemplate(options.id);
  if (options.placeholders) {
    // A starting point for --placeholders: every key the template asks for, with an
    // empty value to fill in.
    await output(Object.fromEntries((template.placeholders ?? []).map(p => [p.key, ''])), options);
  } else {
    await output(template, options);
  }
  if (options.file) {
    console.log('Experiment template %s written to %s.', options.id, options.file);
  }
}
