// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

import { fetchAllAdvice } from './api';
import { COLOR, Table } from 'console-table-printer';
import { abortExecution } from '../errors';

export interface Options {
  environment: string;
  query?: string;
  status: string;
}

const red_color: COLOR = 'red';
const green_color: COLOR = 'green';

export async function validateAdviceStatus(options: Options) {
  const allAdvice = await fetchAllAdvice(options.environment, options.query);
  if (allAdvice.length === 0) {
    return;
  }

  let errorCount = 0;
  const p = new Table();
  for (const advice of allAdvice) {
    const statusMatch = advice.advice.status === options.status;
    if (!statusMatch) {
      errorCount++;
    }
    p.addRow(
      {
        target: advice.target.reference,
        advice: advice.advice.label,
        status: advice.advice.status,
      },
      { color: statusMatch ? green_color : red_color }
    );
  }
  p.printTable();
  if (errorCount > 0) {
    throw abortExecution('%d of %d advice did not match the expected status.', errorCount, allAdvice.length);
  }
}
