import { Batch } from '../models';
import { Wallet } from '.';
export declare function signMultiBatch(wallet: Wallet, transaction: Batch, opts?: {
    batchAccount?: string;
    multisign?: boolean | string;
}): void;
export declare function combineBatchSigners(transactions: Array<Batch | string>): string;
//# sourceMappingURL=batchSigner.d.ts.map