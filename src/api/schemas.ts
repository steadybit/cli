// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import type { components } from './generated/platform-api.ts';

// The platform's request and response bodies, generated from its OpenAPI spec. Commands
// type what they send and receive against these, so that a breaking change in the spec
// fails the type check rather than a customer's pipeline.
export type Schemas = components['schemas'];
