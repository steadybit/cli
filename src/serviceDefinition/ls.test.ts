/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import { loadServiceDefinition } from './files';
import { executeApiCall } from '../api';
import { executeLs } from './ls';

jest.mock('./files');
jest.mock('../api');

const loadServiceDefinitionMock = loadServiceDefinition as jest.Mock;
const executeApiCallMock = executeApiCall as jest.Mock;

describe('ls', () => {
  beforeEach(() => {
    loadServiceDefinitionMock.mockReturnValue(
      Promise.resolve({
        tasks: [
          {
            name: 'steadybit/definitions/weak-spots/k8s-readiness-probe',
            version: '0.1.5',
          },
        ],
        policies: [{ name: 'steadybit/definitions/policies/level-e', version: '0.1.5' }],
      })
    );

    executeApiCallMock.mockReturnValue(
      Promise.resolve({
        json: () =>
          Promise.resolve({
            tasks: [
              {
                name: 'steadybit/definitions/experiments/faultless-redundancy-container',
                version: '0.1.5',
              },
              {
                name: 'steadybit/definitions/weak-spots/k8s-readiness-probe',
                version: '0.1.5',
              },
              {
                name: 'steadybit/definitions/weak-spots/k8s-single-replica',
                version: '0.1.5',
              },
            ],
          }),
      })
    );
  });

  it('must create unfiltered output', async () => {
    expect(executeLs({ file: 'unused' })).resolves.toMatchInlineSnapshot(`
Object {
  "exitCode": 0,
  "output": "
├─┬ steadybit/definitions/policies/level-e@0.1.5
│ ├── steadybit/definitions/experiments/faultless-redundancy-container@0.1.5
│ ├── steadybit/definitions/weak-spots/k8s-readiness-probe@0.1.5
│ └── steadybit/definitions/weak-spots/k8s-single-replica@0.1.5
└── steadybit/definitions/weak-spots/k8s-readiness-probe@0.1.5
",
}
`);
  });

  it('must filter output by name', async () => {
    expect(executeLs({ file: 'unused', name: 'steadybit/definitions/weak-spots/k8s-readiness-probe' })).resolves.toMatchInlineSnapshot(`
Object {
  "exitCode": 0,
  "output": "
├─┬ steadybit/definitions/policies/level-e@0.1.5
│ └── steadybit/definitions/weak-spots/k8s-readiness-probe@0.1.5
└── steadybit/definitions/weak-spots/k8s-readiness-probe@0.1.5
",
}
`);
  });

  it('must filter output by name and version', async () => {
    expect(executeLs({ file: 'unused', name: 'steadybit/definitions/weak-spots/k8s-single-replica', version: '^0.1.0' })).resolves.toMatchInlineSnapshot(`
Object {
  "exitCode": 0,
  "output": "
└─┬ steadybit/definitions/policies/level-e@0.1.5
  └── steadybit/definitions/weak-spots/k8s-single-replica@0.1.5
",
}
`);
  });

  it('must support empty output', async () => {
    expect(executeLs({ file: 'unused', name: 'does-not-exist' })).resolves.toMatchInlineSnapshot(`
Object {
  "exitCode": 1,
  "output": "
└── (empty)
",
}
`);
  });
});
