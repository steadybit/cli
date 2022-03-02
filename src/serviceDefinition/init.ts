/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { v4 as uuidv4 } from 'uuid';
import colors from 'colors/safe';
import inquirer from 'inquirer';
import yaml from 'js-yaml';
import path from 'path';

import { KubernetesMapping, Parameters, ServiceDefinition } from './types';
import { validateHttpUrl, validateNotBlank } from '../prompt/validation';
import { abortExecutionWithError } from '../errors';
import { writeServiceDefinition } from './files';
import { confirm } from '../prompt/confirm';
import { getAllTeams } from '../team/get';
import { Team } from '../team/types';

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

  await writeServiceDefinition(outputFile, serviceDefinition)

  console.log('File created!');
  console.log('You can now upload the service definition by executing');
  console.log(`  ${colors.bold('steadybit service apply .steadybit.yml')}`);
}

async function askForServiceDefinitionInformation(): Promise<ServiceDefinition> {
  let teams: Team[]
  try {
    teams = await getAllTeams();
  } catch (e) {
    const error = await abortExecutionWithError(e, 'Failed to retrieve accessible teams from Steadybit');
    throw error;
  }

  const answers = await inquirer.prompt([
    {
      type: 'input',
      name: 'name',
      message: 'Service name:',
      default: path.basename(process.cwd()),
      validate: validateNotBlank,
    },
    {
      type: 'list',
      name: 'desiredResilienceLevel',
      message: 'Desired Resilience Level:',
      choices: [
        {
          name: 'F - challenges for recovery',
          value: 'F',
        },
        {
          name: 'E - challenges for redundancy',
          value: 'E',
        },
        {
          name: 'D - challenges for redundancy during updates',
          value: 'D',
        },
        {
          name: 'C - challenges for host redundancy',
          value: 'C',
        },
        {
          name: 'B - challenges for loose coupling on startup',
          value: 'B',
        },
        {
          name: 'A - challenges for loose coupling',
          value: 'A',
        },
      ],
    },
  ]);

  const k8Mapping = await askForMappingInformation();
  const parameters = await askForParameters(teams);

  return {
    id: uuidv4(),
    name: answers.name,
    policies: [
      {
        name: `steadybit/definitions/policies/level-${answers.desiredResilienceLevel.toLowerCase()}`,
        version: '0.1.5'
      }
    ],
    mapping: {
      kubernetes: k8Mapping,
    },
    parameters
  };
}

const httpEndpointHelp = `
We need to ensure that the service is still operating as expected when
verifying compliance with the desired resilience level. To do so require
a load-balanced HTTP endpoint that can be called during task execution.
`;

const teamAndEnvironmentHelp = `
Tasks are evaluated in context of teams and environments. For this reason,
you need to select a team and environment that will be used by Steadybit.
For example, when executing an experiment.
`;

async function askForParameters(teams: Team[]): Promise<Parameters> {
  console.log(httpEndpointHelp);

  const {httpEndpoint} = await inquirer.prompt([
    {
      type: 'input',
      name: 'httpEndpoint',
      message: 'URL:',
      validate: validateHttpUrl,
    }]);

    console.log(teamAndEnvironmentHelp);

    const {teamKey, environmentName} = await inquirer.prompt([
      {
        type: 'list',
        name: 'teamKey',
        message: 'Team:',
        choices: teams
          .slice()
          .sort((a, b) => a.name.localeCompare(b.name))
          .map(({name, key}) => ({name, value: key}))
      },
      {
        type: 'list',
        name: 'environmentName',
        message: 'Environment:',
        // We know that the teamKey is within the teams[] because this is how the user
        // could select the team key.
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        choices: answers => teams
          .find(t => t.key === answers.teamKey)!
          .allowedAreas
          .slice()
          .sort((a, b) => a.localeCompare(b))
      }
  ]);

  return {
    httpEndpoint,
    teamKey,
    environmentName
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
