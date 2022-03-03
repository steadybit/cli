/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import fs from 'fs/promises';
import path from 'path';

import { getGitHubRepositoryNameFromGitRemotes, getGitHubRepositoryName, setVersion } from './setVersion';
import { executeShellCommand } from '../shell';

const dummyDefinitionRepository = path.join(__dirname, '__tests__', 'dummyDefinitionRepository');

describe('defRepo/setVersion', () => {
  describe('getGitHubRepositoryNameFromGitRemotes', () => {
    it('must extract from HTTPS URLs', () => {
      expect(getGitHubRepositoryNameFromGitRemotes('https://github.com/steadybit/cli.git'))
        .toEqual('steadybit/cli');
    });

    it('must extract from SSH URLs', () => {
      expect(getGitHubRepositoryNameFromGitRemotes('git@github.com:steadybit/cli.git'))
        .toEqual('steadybit/cli');
    });

    it('must extract from git remote -v output', () => {
      expect(getGitHubRepositoryNameFromGitRemotes(`
origin  git@github.com:steadybit/cli.git (fetch)
origin  git@github.com:steadybit/cli.git (push)
`))
        .toEqual('steadybit/cli');
    });
  });

  describe('getGitHubRepositoryName', () => {
    it('must identify the repository name', async () => {
      const gitHubRepositoryName = await getGitHubRepositoryName(process.cwd());
      expect(gitHubRepositoryName).toEqual('steadybit/cli');
    });
  });

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
        cwd: dummyDefinitionRepository
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
