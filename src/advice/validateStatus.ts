// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

import { fetchAllAdvice } from './api.ts';
import { COLOR, Table } from 'console-table-printer';
import { abortExecution } from '../errors.ts';

export interface Options {
  environment: string;
  query?: string;
  status: string;
}

const red_color: COLOR = 'red';
const green_color: COLOR = 'green';

// The platform reports IMPLEMENTED, ACTION_NEEDED and VALIDATION_NEEDED, while the
// default for --status is written Implemented, so an exact comparison never matched and
// the command failed even when every piece of advice was implemented. Case and the
// separator are both ignored, so either spelling works whichever way round it is given.
function sameStatus(reported: string, expected: string): boolean {
  const normalise = (status: string) => status.trim().toLowerCase().replace(/[\s_-]+/g, '_');
  return normalise(reported) === normalise(expected);
}

export async function validateAdviceStatus(options: Options) {
  const allAdvice = await fetchAllAdvice(options.environment, options.query);
  if (allAdvice.length === 0) {
    return;
  }

  let errorCount = 0;
  const p = new Table();
  for (const advice of allAdvice) {
    const statusMatch = sameStatus(advice.advice.status, options.status);
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
