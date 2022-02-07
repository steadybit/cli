#!/usr/bin/env node

import { Command } from 'commander';

new Command()
  // Prefer to load at runtime directly from the package.json to simplify
  // the TypeScript build. Without this, we would have to make the build
  // more complicated to adapt the root dir accordingly.
  // eslint-disable-next-line
  .version(require('../../package.json').version)
  .command('service', 'Alias for the "service-definition" command.')
  .command('service-definition', 'Configure or verify service definitions.')
  .parseAsync(process.argv);
