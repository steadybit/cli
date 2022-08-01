// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { DefineServiceDefinition, KubernetesMapping, Parameters, PolicyReference } from './types';
import { validateHttpUrl, validateNotBlank } from '../prompt/validation';
import { abortExecutionWithError } from '../errors';
import { writeServiceDefinition } from './files';
import { confirm } from '../prompt/confirm';
import { getAllTeams } from '../team/get';
import { Team } from '../team/types';
import colors from 'colors/safe';
import inquirer from 'inquirer';
import yaml from 'js-yaml';
import path from 'path';

export async function init() {
  const serviceDefinition = await askForServiceDefinitionInformation();

  const outputFile = path.join(process.cwd(), '.steadybit.yml');

  console.log();
  console.log(`About to write to ${colors.bold(outputFile)}:`);
  console.log();
  console.log(yaml.dump(serviceDefinition));
  console.log();

  const ok = await confirm('Should we create this file?');
  if (!ok) {
    process.exit(0);
  }
  console.log();

  await writeServiceDefinition(outputFile, serviceDefinition);

  console.log('File created!');
  console.log('You can now upload the service definition by executing');
  console.log(`  ${colors.bold('steadybit service apply')}`);
}

async function askForServiceDefinitionInformation(): Promise<DefineServiceDefinition> {
  let teams: Team[];
  try {
    teams = await getAllTeams();
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to retrieve accessible teams from Steadybit');
    throw error;
  }

  console.clear();

  const answers = await inquirer.prompt([
    {
      type: 'input',
      name: 'name',
      message: 'Service name:',
      default: path.basename(process.cwd()),
      validate: validateNotBlank,
    },
  ]);

  const policies = await askForPolicies();

  const k8Mapping = await askForMappingInformation();
  const parameters = await askForParameters(teams);

  return {
    name: answers.name,
    policies,
    mapping: {
      kubernetes: k8Mapping,
    },
    parameters,
  };
}

const policiesHelp = `
Services can refer to policies and tasks to describe the desired resilience
level. The following question will give you the option to select commonly
used policies. You may also choose to proceed without selecting a policy.
In this case, you can modify the generated YAML file once this init step
completes. For more information and documentation about tasks and policies,
please refer to the following page:

           ${colors.bold('https://github.com/steadybit/definitions')}
`;

async function askForPolicies(): Promise<PolicyReference[]> {
  console.log(policiesHelp);

  const answers = await inquirer.prompt([
    {
      type: 'checkbox',
      name: 'policies',
      message: 'Policies:',
      choices: [
        {
          name: `${colors.bold(
            'steadybit/.../recovery-pod'
          )}: At a bare minimum, your service needs to restart in case of an outage`,
          value: 'steadybit/definitions/kubernetes/deployments/policies/recovery-pod',
        },
        {
          name: `${colors.bold(
            'steadybit/.../redundancy-pod'
          )}: Verify to be redundant on pod-level to avoid single point of failure`,
          value: 'steadybit/definitions/kubernetes/deployments/policies/redundancy-pod',
        },
        {
          name: `${colors.bold(
            'steadybit/.../redundancy-host'
          )}: Ensure redundancy on host-level to avoid single point of failure`,
          value: 'steadybit/definitions/kubernetes/deployments/policies/redundancy-host',
        },
        {
          name: `${colors.bold('steadybit/.../rolling-update')}: Ensure rolling rollout without service degradation`,
          value: 'steadybit/definitions/kubernetes/deployments/policies/rolling-update',
        },
        {
          name: `${colors.bold(
            'steadybit/.../http-client-fault-tolerance'
          )}: Test your service when a downstream HTTP service has issues`,
          value: 'steadybit/definitions/kubernetes/deployments/policies/http-client-fault-tolerance',
        },

        {
          name: `${colors.bold(
            'steadybit/.../loose-coupling-on-startup'
          )}: Check for coupling to dependent services when restarting your service`,
          value: 'steadybit/definitions/kubernetes/deployments/policies/loose-coupling-on-startup',
        },
        {
          name: `${colors.bold(
            'steadybit/.../loose-coupling'
          )}: Test that your service functions when dependencies are unavailable`,
          value: 'steadybit/definitions/kubernetes/deployments/policies/loose-coupling',
        },
      ],
    },
  ]);

  if (answers.policies.length === 0) {
    const confirm = await inquirer.prompt([
      {
        type: 'confirm',
        name: 'policiesEmpty',
        message: 'You did not select any policies. Are you sure to continue?',
        default: true,
      },
    ]);
    if (!confirm.policiesEmpty) {
      return askForPolicies();
    }
  }

  return answers.policies.map((name: string) => ({
    name,
    version: '0.5.3',
  }));
}

const httpEndpointHelp = `
We need to ensure that the service is still operating as expected when
verifying compliance with the desired resilience level. To do so we require
a load-balanced HTTP endpoint that can be called during task execution.
`;

const teamAndEnvironmentHelp = `
Tasks are evaluated in context of teams and environments. For this reason,
you need to select a team and environment that will be used by Steadybit.
For example, when executing an experiment.
`;

async function askForParameters(teams: Team[]): Promise<Parameters> {
  console.log(httpEndpointHelp);

  const { httpEndpoint } = await inquirer.prompt([
    {
      type: 'input',
      name: 'httpEndpoint',
      message: 'URL:',
      validate: validateHttpUrl,
    },
  ]);

  console.log(teamAndEnvironmentHelp);

  const { teamKey, environmentName } = await inquirer.prompt([
    {
      type: 'list',
      name: 'teamKey',
      message: 'Team:',
      choices: teams
        .slice()
        .sort((a, b) => a.name.localeCompare(b.name))
        .map(({ name, key }) => ({ name, value: key })),
    },
    {
      type: 'list',
      name: 'environmentName',
      message: 'Environment:',
      // We know that the teamKey is within the teams[] because this is how the user
      // could select the team key.
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      choices: answers =>
        teams
          .find(t => t.key === answers.teamKey)!
          .allowedAreas.slice()
          .sort((a, b) => a.localeCompare(b)),
    },
  ]);

  return {
    httpEndpoint,
    teamKey,
    environmentName,
  };
}

const mappingHelp = `
In order to understand which service you are describing, we need
to be able to map it to data collected by the Steadybit agents.
Right now, we only support mapping to Kubernetes monitored deployments.
`;

async function askForMappingInformation(): Promise<KubernetesMapping> {
  console.log(mappingHelp);

  const answers = await inquirer.prompt([
    {
      type: 'input',
      name: 'cluster',
      message: 'Kubernetes cluster name:',
      validate: validateNotBlank,
    },
    {
      type: 'input',
      name: 'namespace',
      message: 'Kubernetes namespace:',
      validate: validateNotBlank,
    },
    {
      type: 'input',
      name: 'deployment',
      message: 'Kubernetes deployment:',
      validate: validateNotBlank,
    },
  ]);

  return answers;
}
