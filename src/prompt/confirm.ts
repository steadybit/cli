// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { cancelable } from './cancellation.ts';

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

  // Loaded lazily so that non-interactive runs, the common case in CI, never pay
  // the cost of importing the prompt implementation.
  const { default: confirmPrompt } = await import('@inquirer/confirm');
  return await cancelable(confirmPrompt({ message, default: defaultYes }));
}
