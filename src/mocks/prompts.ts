// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { screen } from '@inquirer/testing/vitest';

// A prompt is not always on screen the moment the call under test starts: confirm()
// reaches its prompt through a dynamic import, and a rejected answer re-renders the
// same question. Polling for the text covers both without the caller having to know
// which case it is in, and reports the screen it did see when the wait runs out.
export async function waitForPrompt(text: string, timeoutMillis = 4000): Promise<void> {
  const deadline = Date.now() + timeoutMillis;
  while (Date.now() < deadline) {
    if (screen.getScreen().includes(text)) {
      return;
    }
    await new Promise(resolve => setTimeout(resolve, 5));
  }
  throw new Error(`Prompt "${text}" never appeared. Last screen was:\n${screen.getScreen()}`);
}

// Enter is a carriage return; readline does not submit on a newline. `replace` clears
// what is already in the field first, which a prompt keeps after rejecting an answer —
// without it a second attempt is appended to the first rather than replacing it.
export async function answerPrompt(prompt: string, answer: string, { replace = false } = {}): Promise<void> {
  await waitForPrompt(prompt);
  if (replace) {
    screen.input.write('\x7f'.repeat(64)); // backspace past anything already typed
  }
  screen.input.write(`${answer}\r`);
}

export function pressCtrlC(): void {
  // Written as an escape rather than the literal byte, which is invisible in an editor
  // and reads as an empty string.
  screen.input.write('\x03');
}
