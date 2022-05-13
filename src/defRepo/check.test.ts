// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { Response } from 'node-fetch';
import path from 'path';

import { getAbsolutePathToRepositoryRoot, getGitSha } from '../vcs/git';
import { executeApiCall, ApiCallArguments } from '../api/http';
import { getGitHubRepositoryName } from '../vcs/github';
import { executeCheck } from './check';

jest.mock('../vcs/github');
jest.mock('../vcs/git');
jest.mock('../api/http');

const executeApiCallMock = executeApiCall as jest.Mock;
const getGitShaMock = getGitSha as jest.Mock;
const getAbsolutePathToRepositoryRootMock = getAbsolutePathToRepositoryRoot as jest.Mock;
const getGitHubRepositoryNameMock = getGitHubRepositoryName as jest.Mock;

const dummyDefinitionRepository = path.join(__dirname, '__tests__', 'dummyDefinitionRepository');
const version = '4.2.0';

getAbsolutePathToRepositoryRootMock.mockResolvedValue(dummyDefinitionRepository);
getGitShaMock.mockResolvedValue(version);
getGitHubRepositoryNameMock.mockResolvedValue('steadybit/cli');

describe('defRepo/check', () => {
  it('must not fail the check', async () => {
    // Given
    executeApiCallMock.mockImplementation((args: ApiCallArguments) => Promise.resolve(getMockResponse({
      name: args.queryParameters.name,
      version: args.queryParameters.version,
    })));

    // When
    const result = await executeCheck({
      directory: dummyDefinitionRepository,
      version
    });

    // Then
    expect(result).toEqual('');
  });

  it('must fail when the names do not match the repository path', async () => {
    // Given
    getGitHubRepositoryNameMock.mockResolvedValue('steadybit/definitions');
    executeApiCallMock.mockImplementation((args: ApiCallArguments) => Promise.resolve(getMockResponse({
      name: args.queryParameters.name.replace('steadybit/definitions', 'steadybit/cli'),
      version: args.queryParameters.version,
    })));

    // When
    const result = await executeCheck({
      directory: dummyDefinitionRepository,
      version
    });

    // Then
    expect(result.split('\n\n').sort().join('\n\n').trim()).toMatchInlineSnapshot(`
"policies/rolling-update/policy.yml:
 - Name in policy definition file does not match the fully qualified policy name.
   Got:      steadybit/cli/policies/rolling-update
   Expected: steadybit/definitions/policies/rolling-update


weak-spots/rolling-update/task.yAml:
 - Name in task definition file does not match the fully qualified task name.
   Got:      steadybit/cli/weak-spots/rolling-update
   Expected: steadybit/definitions/weak-spots/rolling-update

experiments/rolling-update/task.yml:
 - Name in task definition file does not match the fully qualified task name.
   Got:      steadybit/cli/experiments/rolling-update
   Expected: steadybit/definitions/experiments/rolling-update"
`);
  });

  it('must declare as problematic when the backend says so', async () => {
    // Given
    const e: any = new Error('Failed to load task definition');
    e.response = {
      json: () => Promise.resolve({
        reason: 'Could not locate file.',
      })
    };
    executeApiCallMock.mockRejectedValue(e);

    // When
    const result = await executeCheck({
      directory: dummyDefinitionRepository,
      version
    });

    // Then
    expect(result.split('\n\n').sort().join('\n\n').trim()).toMatchInlineSnapshot(`
"experiments/rolling-update/task.yml:
 - Task retrieval/validation failed. Details:
{
  \\"reason\\": \\"Could not locate file.\\"
}

policies/rolling-update/policy.yml:
 - Policy retrieval/validation failed. Details:
{
  \\"reason\\": \\"Could not locate file.\\"
}

weak-spots/rolling-update/task.yAml:
 - Task retrieval/validation failed. Details:
{
  \\"reason\\": \\"Could not locate file.\\"
}"
`);
  });
});

function getMockResponse<T>(json: T): Response {
  return {
    json: () => Promise.resolve(json)
  } as Response;
}
