// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { getTaskCountByType, printTaskList, printTaskStateByType } from './taskList';
import { loadServiceDefinition } from './files';
import { bold, gray } from 'colors/safe';
import { getState } from './api';

interface Options {
  file: string;
  printParameters?: boolean;
  printMatrixContext?: boolean;
}

const openHelp = gray(
  `
Do you need more information? You can see more details within Steadybit's
user interface. Try executing the following to open it:

                   ${bold('steadybit service open')}
`.trim()
);

export async function verify(options: Options) {
  const serviceDefinition = await loadServiceDefinition(options.file);
  const state = await getState(serviceDefinition);

  printTaskList({
    ...options,
    colorsByTaskState: true,
    renderTaskState: true
  }, state.tasks);

  console.log();

  const countByType = getTaskCountByType(state.tasks);
  printTaskStateByType(countByType);

  console.log();
  console.log(openHelp);

  if (countByType.FAILURE > 0 || countByType.PENDING > 0) {
    process.exit(1);
  }
}
