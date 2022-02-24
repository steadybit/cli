/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { Response } from 'node-fetch';

import { Team, TeamSummary } from './types';
import { executeApiCall } from '../api';

export async function getAllTeams(): Promise<Team[]> {
  const response: Response = await executeApiCall({
    method: 'get',
    path: '/api/teams',
    queryParameters: {
      onlyAccessible: 'true'
    }
  });
  const summary: TeamSummary = await response.json();
  return summary.teams;
}
