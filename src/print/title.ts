// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { bold } from 'colors/safe';

export function printTitle(title: string) {
  console.log(bold(title));
  console.log(bold(''.padStart(title.length, '=')));
}
