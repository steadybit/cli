// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { abortExecution, errorMessage } from '../errors.ts';
import { changeExecutionProperty, type PropertyOperation } from './api.ts';

export interface Options {
  id: number;
  key: string;
  value: string[];
  json?: boolean;
}

// Values are sent as strings unless --json asks otherwise. Guessing from the text would
// turn a ticket number such as "0042" into the number 42.
function parseValue(value: string, json: boolean | undefined): unknown {
  if (!json) {
    return value;
  }
  try {
    return JSON.parse(value);
  } catch (e) {
    throw abortExecution("'%s' is not valid JSON: %s", value, errorMessage(e));
  }
}

export async function setProperty(options: Options) {
  const values = options.value.map(v => parseValue(v, options.json));
  // Several values set a list property; a single one stays a scalar.
  await change('set', options, values.length === 1 ? values[0] : values);
}

export async function addProperty(options: Options) {
  if (options.value.length !== 1) {
    throw abortExecution('Adding to a list property takes exactly one --value.');
  }
  await change('add', options, parseValue(options.value[0], options.json));
}

async function change(operation: PropertyOperation, options: Options, value: unknown) {
  await changeExecutionProperty(options.id, options.key, operation, value);
  console.log('Property %s of experiment run %s updated.', options.key, options.id);
}
