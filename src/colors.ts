// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

import { createColors } from 'picocolors';

// picocolors enables colouring whenever CI is set, and unconditionally on win32, even
// when stdout is a pipe. The CLI's output is routinely parsed by GitOps pipelines, so
// colouring is gated on stdout actually being a terminal instead, which is what the
// previously used `colors` package did.
const colorsSupported = !process.env.NO_COLOR && Boolean(process.env.FORCE_COLOR || process.stdout.isTTY);

export default createColors(colorsSupported);
