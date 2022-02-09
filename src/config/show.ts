import yaml from 'js-yaml';

import { getConfiguration } from './index';

export async function show(): Promise<void> {
  const configuration = await getConfiguration();
  console.log(yaml.dump(configuration));
}
