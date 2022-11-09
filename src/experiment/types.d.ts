// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

export type ExecuteResult = {
  location?: string;
};

export type ExecutionResult = {
  state: string;
};

export type Experiment = Record<string, any> & {
  externalId?: string;
  key?: string;
}

export type UpsertResult = {
  created?: boolean;
  key?: string;
};

export type UpsertAndExecuteResult = {
  location?: string;
  key?: string;
}
