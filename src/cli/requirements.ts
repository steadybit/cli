// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { ensurePlatformAccessConfigurationIsAvailable } from '../config/requirePlatformAccess';

export type ActionFn = (...args: any[]) => void | Promise<void>;

export function requirePlatformAccess(fn: ActionFn): ActionFn {
  return async (...args: any[]) => {
    await ensurePlatformAccessConfigurationIsAvailable();
    return fn(...args);
  };
}
