// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

export function toTwoColumns(key: string, expectedMaximumKeyLength: number, value: string): string {
  const targetPadLength = expectedMaximumKeyLength + 3;
  key = `${key}:`;
  return key.padEnd(targetPadLength) + value;
}
