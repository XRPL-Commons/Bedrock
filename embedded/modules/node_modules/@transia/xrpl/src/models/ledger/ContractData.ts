import { BaseLedgerEntry, HasPreviousTxnID } from './BaseLedgerEntry'

/**
 *
 *
 * @category Ledger Entries
 */
export default interface ContractData
  extends BaseLedgerEntry,
    HasPreviousTxnID {
  LedgerEntryType: 'ContractData'
  /** The owner node for this contract data. */
  OwnerNode: string
  /** The account that owns this contract data. */
  Owner: string
  /** The account associated with this contract. */
  ContractAccount: string
  /** The JSON data for the contract. */
  ContractJson: Record<string, any>
}
