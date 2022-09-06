// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

import { PolicyBinding } from '../policyBinding/types';
import { getTags as getGitTags } from './git';

export async function addVcsTags(policyBindingPath: string, policyBinding: PolicyBinding): Promise<void> {
  const gitTags = await getGitTags(policyBindingPath);

  policyBinding.tags = {
    ...gitTags,
    ...(policyBinding.tags ?? {}),
  };
}
