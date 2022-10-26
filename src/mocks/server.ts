// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH


import { setupServer } from 'msw/node'
import { handlers } from './handlers'

export const server = setupServer(...handlers)
