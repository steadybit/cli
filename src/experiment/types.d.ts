// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

export type ExecuteResult = {
  location?: string;
};

export type ExecutionResult = {
  state: string;
  key: string;
  id: string;
  reason?: string;
};

export type Experiment = Record<string, any> & {
  externalId?: string;
  key?: string;
};

export type UpsertResult = {
  created?: boolean;
  key?: string;
};

export type UpsertAndExecuteResult = ExecuteResult & {
  key?: string;
};

export type ExperimentList = {
  experiments: Record<string, any> &
    {
      key: string;
      name: string;
    }[];
};

export type ExecutionList = {
  executions: Record<string, any> &
    {
      id: number;
      key: string;
      name: string;
      created: string;
      ended: string;
      state: string;
    }[];
};
