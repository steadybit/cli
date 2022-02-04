import { loadPolicy } from "./loading";

export async function establish(policyPath: string) {
  const policy = await loadPolicy(policyPath);
  console.log(policy);
}
