#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

import { Command, Option } from 'commander';
import { requirePlatformAccess } from './requirements.ts';
import { validateAdviceStatus } from '../advice/validateStatus.ts';
import { withExamples } from './help.ts';

const program = new Command();

const validateStatus = program
  .command('validate-status')
  .description('Validates the status of one or multiple advice for a given environment and an optional query.')
  .addOption(new Option('-e, --environment <environment>', 'The environment name.').makeOptionMandatory(true))
  .addOption(new Option('-s, --status <expected-status>', 'The expected status of the advice.').default('Implemented'))
  .addOption(new Option('-q, --query <query>', '(optional) A target query to filter advice by targets.'))
  .action(requirePlatformAccess(validateAdviceStatus));
withExamples(validateStatus, [
  'steadybit advice validate-status -e Global',
  'steadybit advice validate-status -e Global -q "k8s.cluster-name=dev-demo and k8s.namespace=steadybit-demo"',
]);

program.parseAsync(process.argv);
