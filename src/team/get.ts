// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import type { Team, TeamSummary } from './types.ts';
import { executeApiCall } from '../api/http.ts';

export async function getAllTeams(onlyAccessible = true): Promise<Team[]> {
  const response: Response = await executeApiCall({
    method: 'get',
    path: '/api/teams',
    queryParameters: {
      onlyAccessible: String(onlyAccessible),
    },
  });
  const summary = (await response.json()) as TeamSummary;
  return summary.teams;
}
