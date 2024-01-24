// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Response } from 'node-fetch';

import { Team, TeamSummary } from './types';
import { executeApiCall } from '../api/http';

export async function getAllTeams(onlyAccessible = true): Promise<Team[]> {
  const response: Response = await executeApiCall({
    method: 'get',
    path: '/api/teams',
    queryParameters: {
      onlyAccessible: String(onlyAccessible)
    }
  });
  const summary: TeamSummary = await response.json();
  return summary.teams;
}
