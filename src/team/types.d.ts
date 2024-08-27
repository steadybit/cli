// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

export interface TeamSummary {
  teams: Team[];
}

export interface Team {
  key: string;
  name: string;
  allowedEnvironments: string[];
}
