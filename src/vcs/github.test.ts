// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import {
  getGitHubRepositoryNameFromGitRemotes,
  getGitHubRepositoryName
} from './github';

describe('defRepo/github', () => {
  describe('getGitHubRepositoryNameFromGitRemotes', () => {
    it('must extract from HTTPS URLs', () => {
      expect(getGitHubRepositoryNameFromGitRemotes('https://github.com/steadybit/cli.git')).toEqual('steadybit/cli');
    });

    it('must extract from SSH URLs', () => {
      expect(getGitHubRepositoryNameFromGitRemotes('git@github.com:steadybit/cli.git')).toEqual('steadybit/cli');
    });

    it('must extract from git remote -v output', () => {
      expect(
        getGitHubRepositoryNameFromGitRemotes(`
origin  git@github.com:steadybit/cli.git (fetch)
origin  git@github.com:steadybit/cli.git (push)
`)
      ).toEqual('steadybit/cli');

      expect(
        getGitHubRepositoryNameFromGitRemotes(`
origin\thttps://github.com/steadybit/cli (fetch)
origin\thttps://github.com/steadybit/cli (push)
origin\thttps://github.com/steadybit/cli (fetch)
origin\thttps://github.com/steadybit/cli (push)
`)
      ).toEqual('steadybit/cli');

    });
  });

  describe('getGitHubRepositoryName', () => {
    it('must identify the repository name', async () => {
      const gitHubRepositoryName = await getGitHubRepositoryName(__dirname);
      expect(gitHubRepositoryName).toEqual('steadybit/cli');
    });
  });
});
