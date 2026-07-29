// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// The status and the problem body used to travel as untyped properties on a plain Error,
// which every caller duck-typed back out. Naming the type lets callers branch on the
// status through the compiler instead of by convention.
//
// This lives apart from http.ts so that errors.ts can recognise it without the two
// modules having to import each other.
export class ApiError extends Error {
  constructor(
    message: string,
    readonly response: Response,
    // The response body can only be read once, so the text read to build the message
    // above is carried here rather than left on the unusable response.
    private readonly responseBody: string
  ) {
    super(message);
    this.name = 'ApiError';
  }

  get status(): number {
    return this.response.status;
  }

  problemBody<T>(): T | undefined {
    try {
      return this.responseBody ? JSON.parse(this.responseBody) : undefined;
    } catch {
      return undefined;
    }
  }
}
