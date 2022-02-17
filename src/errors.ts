import { Response } from 'node-fetch';
import colors from 'colors/safe';
import { format } from 'util';

export function abortExecution(msg: string, ...args: unknown[]): Error {
  // Make unit-testing easier by only aborting the process outside of Jest
  // https://jestjs.io/docs/environment-variables
  if (process.env.NODE_ENV !== 'test') {
    console.error(colors.red(msg), ...args);
    process.exit(1);
  }

  return new Error(format(msg, ...args));
}

export async function abortExecutionWithError(error: any, msg: string, ...args: unknown[]): Promise<Error> {
  let message = (error as Error)?.message ?? 'Unknown error';

  if (typeof error?.response?.json === 'function') {
    try {
      const body = await (error.response as Response).json();
      message = `${message}: ${JSON.stringify(body, undefined, 2)}`;
    } catch (e) {
      // swallow silently - we just do our best to provide error insights, but do not promise anything.
    }
  }

  return abortExecution(`${msg}: %s`, ...args, message);
}
