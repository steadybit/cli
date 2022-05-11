// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

/* eslint-env node */

module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  plugins: [
    '@typescript-eslint',
    'header'
  ],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
  ],
  rules: {
    '@typescript-eslint/no-explicit-any': 0,
		quotes: ['error', 'single', {avoidEscape: true}],
		'header/header': [
			2,
			'line',
			[
				' SPDX-License-Identifier: MIT',
        {
					template: ` SPDX-FileCopyrightText: ${new Date().getFullYear()} Steadybit GmbH`,
					pattern: /^ SPDX-FileCopyrightText: \d{4} Steadybit GmbH/,
				},
			],
			1,
			{ lineEndings: 'unix' },
		],
	},
};
