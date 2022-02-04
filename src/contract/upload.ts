import { executeApiCall } from '../api';
import { abortExecution } from '../errors';
import { Contract } from './types';

export async function uploadContract(contract: Contract): Promise<void> {
  try {
    await executeApiCall({
      method: 'POST',
      path: '/api/contract',
      body: contract
    });
  } catch (e) {
    throw abortExecution(`Failed upload contract to Steadybit: %s`, (e as Error)?.message ?? 'Unknown error');
  }
}
