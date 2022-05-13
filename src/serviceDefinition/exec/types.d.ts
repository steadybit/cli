// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

export interface Options {
  file: string;
  wait: boolean;
  yes: boolean;
  errorOnEmptyTaskSet: boolean;
  errorOnTaskFailure: boolean;
  task?: string[];
  printParameters?: boolean;
  printMatrixContext?: boolean;
}
