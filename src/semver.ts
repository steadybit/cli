// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { InvalidArgumentError } from 'commander';
import { validRange } from 'semver';

export function validateSemverRangeCommanderArgument(value: string) {
  if (value?.trim().length > 0 && validRange(value)) {
    return value;
  }

  throw new InvalidArgumentError(`'${value}' is not a valid semver version range.`);
}
