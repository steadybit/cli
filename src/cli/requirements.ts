/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { ensurePlatformAccessConfigurationIsAvailable } from '../config/requirePlatformAccess';

export type ActionFn = (...args: any[]) => void | Promise<void>;

export function requirePlatformAccess(fn: ActionFn): ActionFn {
  return async (...args: any[]) => {
    await ensurePlatformAccessConfigurationIsAvailable();
    return fn(...args);
  };
}
