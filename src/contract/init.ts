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
  console.log('  ' + colors.bold('steadybit contract establish .steadybit.yml'));
}

async function askForContractInformation(): Promise<Contract> {
  return await inquirer.prompt([
    {
      type: 'input',
      name: 'name',
      message: 'What do you call this service?',
      default: path.basename(process.cwd()),
    },
    {
      type: 'input',
      name: 'healthEndpoint',
      message:
        'What HTTP URL can we ping to idenity whether the service is healthy?',
    },
    {
      type: 'list',
      name: 'desiredResilienceLevel',
      message: 'Desired Resilience Level',
      choices: [
        {
          name: 'F - the bare minimum for any service',
          value: 'F',
        },
        {
          name: 'E - following Kubernetes best practices',
          value: 'E',
        },
        {
          name: 'D - TODO',
          value: 'D',
        },
        {
          name: 'C - TODO',
          value: 'C',
        },
        {
          name: 'B - TODO',
          value: 'B',
        },
        {
          name: 'A - for the most mission-critical systems',
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
