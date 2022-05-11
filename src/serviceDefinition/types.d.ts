// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

export interface KubernetesMapping {
  cluster: string;
  namespace: string;
  deployment: string;
}

export interface Mapping {
  kubernetes?: KubernetesMapping;
}

export type Iterable = 'dependency' | 'container';

export interface ForEach {
  iterable: Iterable;
  define: Record<string, never>;
}

export interface ReferenceCoordinate {
  name: string;
  version: string;
}

export type TaskDefinition = ReferenceCoordinate;

export interface PolicyDefinition extends ReferenceCoordinate {
  tasks: TaskReference[];
}

export interface PolicyReference extends ReferenceCoordinate {
  parameters?: Parameters;
}

export interface TaskReference extends ReferenceCoordinate {
  parameters?: Parameters;
  forEach?: ForEach[];
}

export interface ServiceDefinition {
  id: string;
  name: string;
  policies?: PolicyReference[];
  tasks?: TaskReference[];
  mapping: Mapping;
  parameters?: Parameters;
  tags?: Record<string, string>;
}

export type Parameters = Record<string, any>;

export type TaskState = 'PENDING' | 'SUCCESS' | 'FAILURE';

export interface TaskHistoryEntry {
  date: string;
  type: string;
  meta: Record<string, any>;
}

export interface Task {
  id: string;
  definition: ReferenceCoordinate;
  state: TaskState;
  type: string;
  history: TaskHistoryEntry[];
  parameters: Parameters;
  matrixContext: never;
}

export interface ServiceState {
  id: string;
  tasks: Task[];
}
