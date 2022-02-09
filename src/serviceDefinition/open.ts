import opn from 'open';

import { loadServiceDefinition } from './loading';
import { executeApiCall } from '../api';
import { abortExecution, abortExecutionWithError } from '../errors';

const errorPrefix = 'Failed to identify deep link to the Steadybit UI. ';

export async function open(serviceDefinitionPath: string) {
  const serviceDefinition = await loadServiceDefinition(serviceDefinitionPath);
  let response;
  try {
    response = await executeApiCall({
      method: 'GET',
      path: `/api/service-definitions/${serviceDefinition.id}/deep-link`,
      expect2xx: false,
    });
  } catch (e) {
    throw abortExecutionWithError(e, errorPrefix + 'HTTP request failed.');
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
      throw abortExecution(errorPrefix + 'Got status code %s: %s', response.status);
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
