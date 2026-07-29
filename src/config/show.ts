// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { dump } from '../yaml.ts';

import { getConfiguration } from './index.ts';

export async function show(): Promise<void> {
  const configuration = await getConfiguration();
  console.log(dump(configuration));
}
