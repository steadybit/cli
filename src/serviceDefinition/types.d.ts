export type HealthDefinitionTypes = 'HTTP';

export interface HealthDefinition {
  type: HealthDefinitionTypes
}

export interface HttpHealthDefinition extends HealthDefinition {
  type: 'HTTP';
  url: string;
}

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
  parameters?: Record<string, never>;
}

export interface TaskReference extends ReferenceCoordinate {
  parameters?: Record<string, never>;
  forEach?: ForEach[];
}

export interface ServiceDefinition {
  id: string;
  name: string;
  policies?: PolicyReference[];
  tasks?: TaskReference[];
  mapping: Mapping;
  health: HealthDefinition[];
}

export type TaskState = 'PENDING' | 'SUCCESS' | 'FAILURE';

export interface Task {
  definition: ReferenceCoordinate;
  state: TaskState;
}

export interface ServiceState {
  tasks: Task[];
}
