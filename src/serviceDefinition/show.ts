/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { intersects } from 'semver';
import archy from 'archy';

import { PolicyDefinition, PolicyReference, ReferenceCoordinate, ServiceDefinition } from './types';
import { abortExecutionWithError } from '../errors';
import { loadServiceDefinition } from './files';
import { executeApiCall } from '../api';

export interface Options {
  file: string;
  name?: string;
  version?: string;
}

interface TaskReferenceThroughPolicy {
  policyReference: PolicyReference;
  policyDefinition: PolicyDefinition;
}

interface ExecutionResult {
  exitCode: number;
  output: string;
}

export async function show(options: Options) {
  const result = await executeShow(options);
  console.log(result.output);
  process.exit(result.exitCode);
}

export async function executeShow(options: Options): Promise<ExecutionResult> {
  const serviceDefinition = await loadServiceDefinition(options.file);
  const policyTaskReferences = await getPolicyTaskReferences(serviceDefinition);

  let exitCode = 0;

  const tree: archy.Data = {
    label: serviceDefinition.name,
    nodes: [],
  };

  policyTaskReferences
    .slice()
    .sort((a, b) => a.policyDefinition.name.localeCompare(b.policyDefinition.name))
    .forEach(policyTaskReference => {
      const node: archy.Data = {
        label: `${policyTaskReference.policyReference.name}@${policyTaskReference.policyReference.version}`,
        nodes: [],
      };

      policyTaskReference.policyDefinition.tasks.forEach(task => {
        if (shouldAddToOutput(task, options)) {
          node.nodes?.push(`${task.name}@${task.version}`);
        }
      });

      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      if (node.nodes!.length > 0) {
        tree.nodes?.push(node);
      }
    });

  serviceDefinition.tasks
    ?.slice()
    .sort((a, b) => a.name.localeCompare(b.name))
    .forEach(task => {
      if (shouldAddToOutput(task, options)) {
        tree.nodes?.push(`${task.name}@${task.version}`);
      }
    });

  if (tree.nodes?.length === 0) {
    tree.nodes?.push('(empty)');
    exitCode = 1;
  }

  return {
    exitCode,
    output: archy(tree),
  };
}

function shouldAddToOutput(referenceCoordinate: ReferenceCoordinate, options: Options): boolean {
  if (!options.name) {
    return true;
  }

  if (options.name !== referenceCoordinate.name) {
    return false;
  }

  if (!options.version) {
    return true;
  }

  return intersects(referenceCoordinate.version, options.version);
}

async function getPolicyTaskReferences(serviceDefinition: ServiceDefinition): Promise<TaskReferenceThroughPolicy[]> {
  if (!serviceDefinition.policies) {
    return [];
  }

  return await Promise.all(serviceDefinition.policies.map(getTaskReferenceThroughPolicy));
}

async function getTaskReferenceThroughPolicy(policyReference: PolicyReference): Promise<TaskReferenceThroughPolicy> {
  const policyDefinition = await loadPolicy(policyReference.name, policyReference.version);
  return {
    policyReference,
    policyDefinition,
  };
}

async function loadPolicy(name: string, version: string): Promise<PolicyDefinition> {
  try {
    const response = await executeApiCall({
      method: 'GET',
      path: '/api/policy-definitions',
      queryParameters: {
        name,
        version,
      },
    });

    return (await response.json()) as PolicyDefinition;
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to find policy definition for %s@%s', name, version);
    throw error;
  }
}
