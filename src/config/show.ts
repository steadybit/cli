/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import yaml from 'js-yaml';

import { getConfiguration } from './index';

export async function show(): Promise<void> {
  const configuration = await getConfiguration();
  console.log(yaml.dump(configuration));
}
