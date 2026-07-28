// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// @ts-check

import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import globals from 'globals';
import { spdxHeader } from './eslint-rules/spdx-header.mjs';

export default tseslint.config(
  // A config object containing only `ignores` acts as a global ignore. Combining it
  // with `files` would scope the ignores to that single object instead.
  {
    ignores: ['**/node_modules/**', '**/dist/**', '**/coverage/**'],
  },
  {
    files: ['**/*.js', '**/*.mjs', '**/*.ts'],
    extends: [eslint.configs.recommended, tseslint.configs.recommended],
    languageOptions: {
      globals: {
        ...globals.node,
      },
    },
    plugins: {
      steadybit: {
        rules: { 'spdx-header': spdxHeader },
      },
    },
    rules: {
      '@typescript-eslint/no-explicit-any': 0,
      'steadybit/spdx-header': 'error',
      // Both wrappers exist to correct a dependency default, and reaching past them
      // fails silently: raw js-yaml drops merge keys and timestamps from experiment
      // files, raw picocolors puts ANSI escapes into piped output.
      'no-restricted-imports': [
        'error',
        {
          paths: [
            { name: 'js-yaml', message: "Use '../yaml.ts', which restores the tag set experiment files rely on." },
            { name: 'picocolors', message: "Use '../colors.ts', which only colours a real terminal." },
          ],
        },
      ],
    },
  },
  {
    files: ['src/yaml.ts', 'src/colors.ts'],
    rules: { 'no-restricted-imports': 0 },
  }
);
