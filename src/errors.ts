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
