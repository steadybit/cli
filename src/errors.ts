// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ApiError } from './api/error.ts';
import colors from './colors.ts';
import { format } from 'node:util';

export interface AbortExecutionOptions {
  colorize?: boolean;
}

// Reaching for `(e as Error)?.message` at each catch site had already let two different
// fallback strings drift apart, so the idiom lives here instead.
export function errorMessage(e: unknown): string {
  return (e as Error)?.message || 'Unknown error';
}

export function abortExecution(msg: string, ...args: unknown[]): Error {
  return abortExecutionWithOpts(undefined, msg, ...args);
}

export function abortExecutionWithOpts(
  { colorize = true }: AbortExecutionOptions = {},
  msg: string,
  ...args: unknown[]
): Error {
  // Make unit-testing easier by only aborting the process outside of the test runner,
  // which sets NODE_ENV to 'test'.
  if (process.env.NODE_ENV !== 'test') {
    if (colorize) {
      msg = colors.red(msg);
    }
    console.error(msg, ...args);
    process.exit(1);
  }

  return new Error(format(msg, ...args));
}

export function abortExecutionWithError(error: unknown, msg: string, ...args: unknown[]): Error {
  let message = errorMessage(error);

  const errorBody = getExecutionErrorBody(error);
  if (errorBody) {
    message = `${message}: ${JSON.stringify(errorBody, undefined, 2)}`;
  }

  return abortExecution(`${msg}: %s`, ...args, message);
}

export function getExecutionErrorBody<T>(error: unknown): T | undefined {
  return error instanceof ApiError ? error.problemBody<T>() : undefined;
}
