import { Parameter } from '../common'
import {
  BaseTransaction,
  // isNumber,
  // isString,
  validateBaseTransaction,
  // validateOptionalField,
  // validateRequiredField,
} from './common'

/**
 * @category Transaction Models
 */
export interface ContractUserDelete extends BaseTransaction {
  TransactionType: 'ContractUserDelete'

  ComputationAllowance: number

  ContractAccount: string

  FunctionName: string

  Parameters?: Parameter[]
}

/**
 * Verify the form and type of a ContractUserDelete at runtime.
 *
 * @param tx - A ContractUserDelete Transaction.
 * @throws When the ContractUserDelete is malformed.
 */
export function validateContractUserDelete(tx: Record<string, unknown>): void {
  validateBaseTransaction(tx)

  // validateRequiredField(tx, 'ComputationAllowance', isNumber)

  // validateRequiredField(tx, 'ContractAccount', isString)

  // validateRequiredField(tx, 'FunctionName', isString)

  // validateOptionalField(tx, 'Parameters', isany[])
}
