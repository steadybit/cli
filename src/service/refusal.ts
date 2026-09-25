// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { ApiError } from '../api/error.ts';
import { getExecutionErrorBody } from '../errors.ts';

interface Violations {
  violations?: { message?: string }[];
}

// The platform refuses a change to a service or profile that would orphan provided
// experiments, and names its query parameter in the reason. Recognised here so the
// user is pointed at the CLI flag instead.
export function refusedForProvidedExperiments(e: unknown): boolean {
  return (
    e instanceof ApiError &&
    e.status === 422 &&
    (getExecutionErrorBody<Violations>(e)?.violations ?? []).some(v => v.message?.includes('deleteExperiments'))
  );
}
