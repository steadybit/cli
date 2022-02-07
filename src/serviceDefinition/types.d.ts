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

export interface ServiceDefinition {
  id: string;
  name: string;
  desiredResilienceLevel: string;
  mapping: Mapping;
  health: HealthDefinition[];
}
