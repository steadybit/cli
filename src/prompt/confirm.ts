// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import inquirer from 'inquirer';

export interface ConfirmOptions {
  defaultYes?: boolean;
  defaultWhenNonInteractive?: boolean;
}

export async function confirm(
  message: string,
  { defaultYes = true, defaultWhenNonInteractive = true }: ConfirmOptions = {}
): Promise<boolean> {
  if (!process.stdout.isTTY) {
    return defaultWhenNonInteractive;
  }

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
