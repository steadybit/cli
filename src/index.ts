#!/usr/bin/env node

import {program} from './cli';

if (require.main === module) {
  main(process.argv)
    .catch(error => console.error('Unexpected error while processing commands.', error));
}

export async function main(argv: string[]) {
  await program.parseAsync(argv);
}
