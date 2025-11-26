import { AccountRoot } from '../ledger'

import { BaseRequest, BaseResponse, LookupByLedgerRequest } from './baseMethod'

/**
 * The `contract_info` command retrieves information about a contract, its
 * activity, and its XRP balance. All information retrieved is relative to a
 * particular version of the ledger. Returns an {@link ContractInfoResponse}.
 *
 * @category Requests
 */
export interface ContractInfoRequest
  extends BaseRequest,
    LookupByLedgerRequest {
  command: 'contract_info'
  /** A unique identifier for the contract, most commonly the contract's address. */
  contract_account: string
  /** If you include an account we will return the contract data for that account. */
  account?: string
  /** If you include a function we will return the contract data for that function. */
  function?: string
}

interface ContractInfoResponse extends BaseResponse {
  result: {
    // IDEA
    // contract: any // LedgerObject Contract for the contract instance
    // contract_source: any // LedgerObject ContractSource for the contract code
    // contract_account: any LedgerObject AccountRoot for pseudo-account of the contract

    contract_account: string
    code: string
    hash: string
    functions: any[]
    source_code_uri: string

    /**
     * The AccountRoot ledger object with this account's information, as stored
     * in the ledger.
     */
    account_data: AccountRoot
    contract_data: any
    user_data?: any

    /**
     * The ledger index of the current in-progress ledger, which was used when
     * retrieving this information.
     */
    ledger_current_index?: number
    /**
     * The ledger index of the ledger version used when retrieving this
     * information. The information does not contain any changes from ledger
     * versions newer than this one.
     */
    ledger_index?: number
    validated?: boolean
  }
}
