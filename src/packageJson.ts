// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { createRequire } from 'node:module';

interface PackageJson {
  name: string;
  version: string;
  engines: { node: string };
}

// Read at runtime rather than imported, so that package.json stays outside of the
// TypeScript root directory and the compiled output keeps its flat dist/ layout.
export const packageJson: PackageJson = createRequire(import.meta.url)('../package.json');
