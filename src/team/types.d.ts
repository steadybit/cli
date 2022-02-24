/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

export interface TeamSummary {
  teams: Team[];
}

export interface Team {
  key: string;
  name: string;
  allowedAreas: string[]
}
