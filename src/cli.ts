import { Command } from 'commander';
import colors from 'colors/safe';

export const program = new Command();

// Prefer to load at runtime directly from the package.json to simplify
// the TypeScript build. Without this, we would have to make the build
// more complicated to adapt the root dir accordingly.
// eslint-disable-next-line
program.version(require('../package.json').version);

program
  .command('establish [config]')
  .description(
    'yo!'
  )
  .action(async configPath => {
    // what?
  });
