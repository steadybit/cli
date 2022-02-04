#!/usr/bin/env node

import { Command } from 'commander';

new Command()
  // Prefer to load at runtime directly from the package.json to simplify
  // the TypeScript build. Without this, we would have to make the build
  // more complicated to adapt the root dir accordingly.
  // eslint-disable-next-line
  .version(require('../../package.json').version)
  .command('policy', 'Configure or verify policies for an application.', {
    executableFile: 'policy',
  })
  .parseAsync(process.argv);
