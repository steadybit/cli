// @ts-check

/*
 * Copyright 2024 steadybit GmbH. All rights reserved.
 */

import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';
import globals from 'globals';

export default tseslint.config(
  {
    ignores: ['**/node_modules', '**/dist'],
    files: ['**/*.js', '**/*.jsx', '**/*.ts', '**/*.tsx'],
  },
  eslint.configs.recommended,
  ...tseslint.configs.recommended,
  {
    languageOptions: {
      globals: {
        ...globals.node,
      },
    },
    rules: {
      '@typescript-eslint/no-explicit-any': 0,
      // Temporary disable - wait for fix -
      // 'header/header': [
      //   2,
      //   'line',
      //   [
      //     ' SPDX-License-Identifier: MIT',
      //     {
      //       template: ` SPDX-FileCopyrightText: ${new Date().getFullYear()} Steadybit GmbH`,
      //       pattern: /^ SPDX-FileCopyrightText: \d{4} Steadybit GmbH/,
      //     },
      //   ],
      //   1,
      //   { lineEndings: 'unix' },
      // ],
    },
  }
);
