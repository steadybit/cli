import { uploadContract } from './upload';
import { loadContract } from './loading';

export async function establish(contractPath: string) {
  const contract = await loadContract(contractPath);
  await uploadContract(contract);
}
