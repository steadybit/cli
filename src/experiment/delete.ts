// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { removeExperiment } from './api';

export interface Options {
  key: string;
}

export async function deleteExperiment(options: Options) {
  await removeExperiment(options.key);
  console.log('Experiment %s deleted.', options.key);
}
