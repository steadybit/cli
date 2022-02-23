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
}

export type Parameters = Record<string, never>;

export type TaskState = 'PENDING' | 'SUCCESS' | 'FAILURE';

export interface Task {
  definition: ReferenceCoordinate;
  state: TaskState;
  parameters: Parameters;
  matrixContext: never;
}

export interface ServiceState {
  tasks: Task[];
}
