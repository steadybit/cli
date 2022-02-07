import colors from 'colors/safe';
import inquirer from 'inquirer';
import fs from 'fs/promises';
import yaml from 'js-yaml';
import path from 'path';

import { Contract } from './types';

export async function init() {
  const contract = await askForContractInformation();

  const outputFile = path.join(process.cwd(), '.steadybit.yml');
  const fileContent = yaml.dump(contract);

  console.log(`About to write to ${colors.bold(outputFile)}:`);
  console.log();
  console.log(fileContent);
  console.log();

  const ok = await confirm();
  if (!ok) {
    process.exit(0);
  }
  console.log();

  fs.writeFile(outputFile, fileContent);

  console.log('File created!');
  console.log('You can now establish the contract by executing');
  console.log(
    '  ' + colors.bold('steadybit contract establish .steadybit.yml')
  );
}

async function askForContractInformation(): Promise<Contract> {
  return await inquirer.prompt([
    {
      type: 'input',
      name: 'name',
      message: 'What do you call this service?',
      default: path.basename(process.cwd()),
      validate: input => input != null && input.trim().length > 0
    },
    {
      type: 'input',
      name: 'healthEndpoint',
      message:
        'What HTTP URL can we ping to idenity whether the service is healthy?',
      validate: input => {
        try {
          new URL(input);
          return true;
        } catch (e) {
          return 'Invalid URL. Please specify an absolute URL, e.g., https://example.com/health';
        }
      },
    },
    {
      type: 'list',
      name: 'desiredResilienceLevel',
      message: 'Desired Resilience Level',
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
}

async function confirm(): Promise<boolean> {
  const answers = await inquirer.prompt([
    {
      type: 'confirm',
      name: 'confirm',
      message: 'Should we create this file?',
    },
  ]);

  return answers.confirm;
}
