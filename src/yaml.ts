// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import {
  CORE_SCHEMA,
  binaryTag,
  dump as dumpYaml,
  load as loadYaml,
  mergeTag,
  omapTag,
  pairsTag,
  setTag,
  timestampTag,
} from 'js-yaml';

// js-yaml 5 narrowed its default schema to the YAML core schema, which silently drops
// merge keys (`<<:`) and timestamps: an experiment factoring shared step fields into an
// anchor would parse into a literal "<<" property. This restores the tag set that
// js-yaml 4 enabled by default, so experiment files keep round-tripping unchanged.
const schema = CORE_SCHEMA.withTags(mergeTag, timestampTag, binaryTag, omapTag, pairsTag, setTag);

export function load(input: string): unknown {
  return loadYaml(input, { schema });
}

export function dump(input: unknown): string {
  return dumpYaml(input, { schema });
}
