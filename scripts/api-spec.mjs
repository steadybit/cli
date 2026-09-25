#!/usr/bin/env node
// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// The CLI's request and response types are generated from the platform's OpenAPI spec,
// so that a breaking change in the platform shows up as a type error here instead of as
// a failing command in a customer's pipeline.
//
//   fetch     download the live spec into openapi/platform-api.json
//   generate  regenerate src/api/generated/platform-api.ts from the committed spec
//   check     fail if the generated types are not what the committed spec produces
//
// `fetch` followed by `generate` and a type check is what CI runs against the live
// platform. `check` runs as part of `npm run verify`, so the committed types can never
// silently drift from the committed spec.

import fs from 'node:fs/promises';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import openapiTS, { astToString } from 'openapi-typescript';

const root = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');
const specFile = path.join(root, 'openapi', 'platform-api.json');
const typesFile = path.join(root, 'src', 'api', 'generated', 'platform-api.ts');
const specUrl = process.env.STEADYBIT_SPEC_URL || 'https://platform.steadybit.com/api/spec';

const HEADER = `// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Generated from openapi/platform-api.json by \`npm run api:generate\`. Do not edit.

`;

// Examples in the spec are copied into the generated doc comments, and one of them is a
// Slack incoming webhook URL, which GitHub push protection rightly refuses to accept.
// The examples carry no type information, so they are masked before anything is written.
function redactSecrets(text) {
  return text.replace(
    /https:\/\/hooks\.slack\.com\/services\/[^"\s\\]+/g,
    'https://hooks.slack.com/services/<redacted>'
  );
}

async function fetchSpec() {
  const response = await fetch(specUrl, { headers: { Accept: 'application/json' } });
  if (!response.ok) {
    throw new Error(`Fetching the platform spec from ${specUrl} failed with status ${response.status}`);
  }
  const spec = await response.json();
  if (!spec.openapi || !spec.paths) {
    throw new Error(`${specUrl} did not return an OpenAPI document`);
  }
  await fs.mkdir(path.dirname(specFile), { recursive: true });
  await fs.writeFile(specFile, `${redactSecrets(JSON.stringify(spec, undefined, 2))}\n`);
  console.log(`Wrote ${Object.keys(spec.paths).length} paths from ${specUrl} to ${path.relative(root, specFile)}`);
}

// The platform models polymorphism the way springdoc emits it: a base schema lists its
// subtypes under `oneOf`, and every subtype pulls the base back in through `allOf`. That
// is valid OpenAPI but a type that contains itself in TypeScript. Each subtype is given
// the base's own fields inline instead, which keeps them while breaking the cycle.
function breakPolymorphicCycles(spec) {
  const schemas = spec.components?.schemas ?? {};
  const refTo = name => `#/components/schemas/${name}`;
  for (const [baseName, base] of Object.entries(schemas)) {
    const subtypes = (base.oneOf ?? []).map(member => member.$ref?.split('/').pop()).filter(Boolean);
    const ownFields = structuredClone(base);
    delete ownFields.oneOf;
    delete ownFields.discriminator;
    for (const subtypeName of subtypes) {
      const subtype = schemas[subtypeName];
      if (!subtype?.allOf) {
        continue;
      }
      subtype.allOf = subtype.allOf.map(part => (part.$ref === refTo(baseName) ? structuredClone(ownFields) : part));
    }
  }
  return spec;
}

async function render() {
  const spec = breakPolymorphicCycles(JSON.parse(await fs.readFile(specFile, 'utf8')));
  const ast = await openapiTS(spec, { alphabetize: true });
  return HEADER + astToString(ast);
}

async function generate() {
  await fs.mkdir(path.dirname(typesFile), { recursive: true });
  await fs.writeFile(typesFile, await render());
  console.log(`Wrote ${path.relative(root, typesFile)}`);
}

async function check() {
  const [expected, actual] = await Promise.all([render(), fs.readFile(typesFile, 'utf8').catch(() => '')]);
  if (expected !== actual) {
    console.error(
      `${path.relative(root, typesFile)} is out of date with ${path.relative(root, specFile)}. Run \`npm run api:generate\`.`
    );
    process.exit(1);
  }
}

const commands = { fetch: fetchSpec, generate, check };
const command = commands[process.argv[2]];
if (!command) {
  console.error(`Usage: api-spec.mjs <${Object.keys(commands).join('|')}>`);
  process.exit(2);
}
await command();
