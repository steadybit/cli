#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Command } from 'commander';

import { show } from '../config/show';

const program = new Command();

program.command('profile', 'Configure authentication profiles.');

program.command('show').description('Show the active CLI configuration. Warning: Prints secrets!').action(show);

program.parseAsync(process.argv);
