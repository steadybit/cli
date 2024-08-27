// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

export function indent(s: string, numberOfCharacters: number): string {
  return s
    .split('\n')
    .map(s => ' '.repeat(numberOfCharacters) + s)
    .join('\n');
}
