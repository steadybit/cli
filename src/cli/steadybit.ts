#!/usr/bin/env node
/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { Command } from 'commander';

new Command()
  // Prefer to load at runtime directly from the package.json to simplify
  // the TypeScript build. Without this, we would have to make the build
  // more complicated to adapt the root dir accordingly.
  // eslint-disable-next-line
  .version(require('../../package.json').version)
  .command('config', 'Show/modify the CLI configuration and authentication profiles.')
  .command('service', 'Alias for the "service-definition" command.')
  .command('service-definition', 'Configure or verify service definitions.')
  .parseAsync(process.argv);
