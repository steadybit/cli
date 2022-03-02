/*
 * Copyright 2022 steadybit GmbH. All rights reserved.
 */

import opn from 'open';

import { abortExecution, abortExecutionWithError } from '../errors';
import { loadServiceDefinition } from './files';
import { executeApiCall } from '../api';

const errorPrefix = 'Failed to identify deep link to the Steadybit UI. ';

export interface Options {
  file: string;
}

export async function open(options: Options) {
  const serviceDefinition = await loadServiceDefinition(options.file);
  let response;
  try {
    response = await executeApiCall({
      method: 'GET',
      path: `/api/service-definitions/${serviceDefinition.id}/deep-link`,
      expect2xx: false,
      redirect: 'manual'
    });
  } catch (e) {
    const error = await abortExecutionWithError(e, errorPrefix + 'HTTP request failed.');
    throw error;
  }

  if (response.status !== 302) {
    let errorTitle: string | undefined;
    try {
      const body = await response.json();
      errorTitle = body.title;
    } catch (e) {
      // swallow silently - we just do our best to provide error insights, but do not promise anything.
    }

    if (errorTitle) {
      throw abortExecution(errorPrefix + errorTitle);
    } else {
      throw abortExecution(errorPrefix + 'Got status code %s', response.status);
    }
  }

  const location = response.headers.get('Location');
  if (typeof location !== 'string') {
    throw abortExecution(
      errorPrefix + 'Did not receive a Location response header.'
    );
  }

  await opn(response.headers.get('Location') as string);
}
