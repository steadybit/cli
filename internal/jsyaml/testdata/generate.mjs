// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Regenerates cases.json: values and what the TypeScript CLI rendered for them with
// js-yaml and JSON.stringify. Run from the repository root after `npm run build`:
//   node internal/jsyaml/testdata/generate.mjs
import fs from 'node:fs';
const { dump } = await import(new URL('../../../dist/yaml.js', import.meta.url));

const long = 'When a single container from steadybit-demo/toys-bestseller fails then within 2m all pods are ready.';
const strings = [
  '', ' ', '  lead', 'trail ', 'plain', 'with: colon', 'a:b', 'key #not comment', 'a #comment', '#start',
  '*', '&anchor', '!tag', '|pipe', '>gt', "'q'", '"dq"', '%pct', '@at', '`tick`', '-', '- dash', '-x', '?', '? q', ':x',
  '[br', ']', '{x', '}', ',comma', 'a,b', 'a[b]', 'x{y}',
  'true', 'True', 'TRUE', 'false', 'yes', 'no', 'on', 'off', 'null', 'Null', '~', 'y', 'n',
  '0', '42', '-7', '+3', '007', '0x1F', '0o17', '0b101', '1_000', '1.5', '.5', '1.', '1e3', '1E-3', '.inf', '-.Inf', '.nan', '1e999',
  '2026-02-25', '2026-02-30', '2026-02-25T06:35:52.676257Z', '2026-02-25 06:35:52', '2026-2-5T1:02:03Z', '12:30', '1:20:30',
  '<<', '=', '---', '--- x', '...', 'a---', 'tab\there', 'ctrl\u0001', 'del\u007f', 'nbsp x', 'bom﻿', 'line sep',
  'multi\nline', 'multi\nline\n', 'trailing\n\n', '\n', '\nlead', ' \nx', 'x\n ', 'a\n\n\nb', '  indented\nnext',
  long, long + '\n' + long, 'x'.repeat(100), 'nospaceshere'.repeat(9) + ' tail', 'a  b ' + 'word '.repeat(20),
  'emoji 😀 ' + 'word '.repeat(18), 'ünïcödé text ' + 'wörd '.repeat(16), 'quote\'s', 'back\\slash', 'http://x?a=b&c=<d>',
  'k8s.namespace="steadybit-demo" AND k8s.deployment="hot-deals" AND host.hostname="very-long-host-name-value"',
];
// -0 is left out: JSON.stringify writes it as 0, so it cannot round-trip through this file.
const numbers = [0, 1, -1, 1.5, 0.1, 1e21, 1e-7, 123456789012345680000, 5e-324, 2 ** 53, 1.0, 100, 3.14159];
const values = [
  ...strings.map(s => ({ s })),
  ...numbers.map(n => ({ n })),
  { b: true, f: false, z: null },
  { empty: {}, list: [], nested: { deep: { deeper: [1, 'two', { three: 3 }] } } },
  { list: [[1, 2], [], [{}], { a: [] }] },
  { '10': 'ten', '2': 'two', b: 'b', '01': 'zero-one', a: 'a', '-1': 'minus' },
  { 'key with: colon': 1, 'true': 2, '123': 3, '': 4, 'multi\nkey': 5, '# hash': 6 },
  { lanes: [{ steps: [{ type: 'action', parameters: { duration: '30s', note: long }, radius: { predicate: { operator: 'AND', predicates: [{ key: 'k8s.namespace', operator: 'EQUALS', values: ['steadybit-demo'] }] } } }] }] },
];
const cases = values.map(value => ({
  value,
  yaml: dump(value),
  json: JSON.stringify(value, undefined, 2),
  compact: JSON.stringify(value),
}));
fs.writeFileSync(new URL('./cases.json', import.meta.url), JSON.stringify(cases, undefined, 2) + '\n');
console.log(`wrote ${cases.length} cases`);
