/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

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
		quotes: ['error', 'single'],
		'header/header': [
			2,
			'block',
			[
				'',
				{
					template: ` * Copyright ${new Date().getFullYear()} steadybit GmbH. All rights reserved.`,
					pattern: /^ \* Copyright \d{4} steadybit GmbH\. All rights reserved\./,
				},
				' ',
			],
			2,
			{ lineEndings: 'unix' },
		],
	},
};
