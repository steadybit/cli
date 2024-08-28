#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2024 Steadybit GmbH

import { Command, Option } from 'commander';
import { requirePlatformAccess } from './requirements';
import { validateAdviceStatus } from '../advice/validateStatus';

const program = new Command();

program
  .command('validate-status')
  .description('Validates the status of one or multiple advice for a given environment and an optional query.')
  .addOption(new Option('-e, --environment <environment>', 'The environment name.').makeOptionMandatory(true))
  .addOption(new Option('-s, --status <expected-status>', 'The expected status of the advice.').default('Implemented'))
  .addOption(
    new Option('-q, --query <query>', '(optional) A target query to filter advice by targets.').makeOptionMandatory(
      true
    )
  )
  .action(requirePlatformAccess(validateAdviceStatus));

program.parseAsync(process.argv);
