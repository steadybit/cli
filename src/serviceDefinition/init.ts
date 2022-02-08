import { v4 as uuidv4 } from 'uuid';
import colors from 'colors/safe';
import inquirer from 'inquirer';
import fs from 'fs/promises';
import yaml from 'js-yaml';
import path from 'path';

import { ServiceDefinition, HealthDefinition, KubernetesMapping } from './types';

export async function init() {
  const serviceDefinition = await askForServiceDefinitionInformation();

  const outputFile = path.join(process.cwd(), '.steadybit.yml');
  const fileContent = yaml.dump(serviceDefinition);

  console.log(`About to write to ${colors.bold(outputFile)}:`);
  console.log();
  console.log(fileContent);
  console.log();

  const ok = await confirm('Should we create this file?');
  if (!ok) {
    process.exit(0);
  }
  console.log();

  fs.writeFile(outputFile, fileContent);

  console.log('File created!');
  console.log('You can now upload the service definition by executing');
  console.log(
    '  ' + colors.bold('steadybit service apply .steadybit.yml')
  );
}

async function askForServiceDefinitionInformation(): Promise<ServiceDefinition> {
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
          name: 'F - service can restart',
          value: 'F',
        },
        {
          name: 'E - service is redundant',
          value: 'E',
        },
        {
          name: 'D - service supports rolling restarts',
          value: 'D',
        },
        {
          name: 'C - service is host redundant',
          value: 'C',
        },
        {
          name: 'B - service dependency unavailability does not cause restart issues',
          value: 'B',
        },
        {
          name: 'A - HTTP calls are successful when service dependencies are unavailable',
          value: 'A',
        },
      ],
    },

  ]);

  const k8Mapping = await askForMappingInformation();
  const health = await askForHealthInformation();

  return {
    id: uuidv4(),
    ...answers,
    mapping: {
      kubernetes: k8Mapping,
    },
    health
  };
}

const healthHelp = `
We need to ensure that the service is still operating as expected when
verifying compliance with the desired resilience level. To do so we support
various checks. Please define at least one health check that we can use
when running resilience checks.
`;

async function askForHealthInformation(): Promise<HealthDefinition[]> {
  console.log(healthHelp);

  const health = [];
  let continueAsking = true;

  while (continueAsking) {
    const answers = await inquirer.prompt([
      {
        type: 'list',
        name: 'type',
        message: 'Health Check Type:',
        choices: [
          {
            name: 'HTTP',
            value: 'HTTP',
          }
        ],
      },
      {
        type: 'input',
        name: 'url',
        message:
          'URL:',
        validate: input => {
          try {
            const url = new URL(input);
            if (url.protocol !== 'http:' && url.protocol !== 'https:') {
              return `Unsupported protocol ${url.protocol}. Only http: and https: are supported.`;
            }
            return true;
          } catch (e) {
            return 'Invalid URL. Please specify an absolute URL, e.g., https://example.com/health';
          }
        },
      },
    ]);

    health.push(answers);

    console.log();
    continueAsking = await confirm('Would you like to define another health check?', false);
    console.log();
  }

  return health;
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

function validateNotBlank(input: string): boolean {
  return input != null && input.trim().length > 0;
}

async function confirm(message: string, defaultYes=true): Promise<boolean> {
  const answers = await inquirer.prompt([
    {
      type: 'confirm',
      name: 'confirm',
      message,
      default: defaultYes
    },
  ]);

  return answers.confirm;
}