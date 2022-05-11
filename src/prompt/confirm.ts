// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import inquirer from 'inquirer';

export async function confirm(message: string, defaultYes = true): Promise<boolean> {
  const answers = await inquirer.prompt([
    {
      type: 'confirm',
      name: 'confirm',
      message,
      default: defaultYes,
    },
  ]);

  return answers.confirm;
}
