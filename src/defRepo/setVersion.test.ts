// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import fs from 'fs/promises';
import path from 'path';

import { executeShellCommand } from '../shell';
import { setVersion } from './setVersion';

const dummyDefinitionRepository = path.join(__dirname, '__tests__', 'dummyDefinitionRepository');

describe('defRepo/setVersion', () => {
  describe('setVersion', () => {
    afterEach(async () => {
      await executeShellCommand('git restore __tests__/dummyDefinitionRepository', {
        cwd: __dirname
      });
    });

    it('must change the version numbers', async () => {
      // When
      await setVersion({
        version: '1.0.0',
        directory: dummyDefinitionRepository
      });

      // Then
      const rollingUpdateExperiment = await fs.readFile(path.join(dummyDefinitionRepository, 'experiments', 'rolling-update', 'task.yml'), 'utf8');
      expect(rollingUpdateExperiment).toMatchInlineSnapshot(`
"name: \\"steadybit/cli/experiments/rolling-update\\"
type: experiment
"
`);

      const rollingUpdateWeakSpot = await fs.readFile(path.join(dummyDefinitionRepository, 'weak-spots', 'rolling-update', 'task.yAml'), 'utf8');
      expect(rollingUpdateWeakSpot).toMatchInlineSnapshot(`
"name: steadybit/cli/weak-spots/rolling-update
type: weak_spot
version: 1.0.0
"
`);

      const rollingUpdatePolicy = await fs.readFile(path.join(dummyDefinitionRepository, 'policies', 'rolling-update', 'policy.yml'), 'utf8');
      expect(rollingUpdatePolicy).toMatchInlineSnapshot(`
"name: steadybit/cli/policies/rolling-update
type: policy
tasks:
  - name: steadybit/cli/weak-spots/rolling-update
    version: 1.0.0
  - name: steadybit/cli/experiments/rolling-update
    version: 1.0.0
"
`);
    });
  });
});
