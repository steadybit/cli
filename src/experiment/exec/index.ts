// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Options } from './types';
import { confirm } from '../../prompt/confirm';
import { executeExperiment } from './execution';

export async function execute(options: Options) {
  if (!options.yes) {
    const confirmation = await confirm('Are you sure you want to run the experiment?', {
      defaultYes: false,
      defaultWhenNonInteractive: true,
    });
    if (!confirmation) {
      process.exit(0);
    }
  }

  const successful = await executeExperiment(options);
  process.exit(successful === false ? 1 : 0);
}
