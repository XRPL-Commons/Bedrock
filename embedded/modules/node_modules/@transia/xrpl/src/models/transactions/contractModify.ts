import { Function, InstanceParameter, InstanceParameterValue } from '../common'
import {
  BaseTransaction,
  // isString,
  validateBaseTransaction,
  // validateOptionalField,
} from './common'

/**
 * @category Transaction Models
 */
export interface ContractModify extends BaseTransaction {
  TransactionType: 'ContractModify'

  ContractAccount?: string

  ContractCode?: string

  ContractHash?: string

  Functions?: Function[]

  InstanceParameters?: InstanceParameter[]

  InstanceParameterValues?: InstanceParameterValue[]

  URI?: string
}

/**
 * Verify the form and type of a ContractModify at runtime.
 *
 * @param tx - A ContractModify Transaction.
 * @throws When the ContractModify is malformed.
 */
export function validateContractModify(tx: Record<string, unknown>): void {
  validateBaseTransaction(tx)

  // validateOptionalField(tx, 'ContractAccount', isString)

  // validateOptionalField(tx, 'ContractCode', isString)

  // validateOptionalField(tx, 'ContractHash', isString)

  // validateOptionalField(tx, 'Functions', isany[])

  // validateOptionalField(tx, 'InstanceParameters', isany[])

  // validateOptionalField(tx, 'InstanceParameterValues', isany[])

  // validateOptionalField(tx, 'URI', isString)
}
