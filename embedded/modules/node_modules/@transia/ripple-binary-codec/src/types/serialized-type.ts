import { BytesList } from '../serdes/binary-serializer'
import { BinaryParser } from '../serdes/binary-parser'
import { XrplDefinitionsBase } from '../enums'
import { bytesToHex } from '@transia/isomorphic/utils'

/**
 * Enum for SerializedTypeID values used in XRPL
 * These match the C++ implementation's STI_ constants
 */
export enum SerializedTypeID {
  STI_NOTPRESENT = 0,
  STI_UINT16 = 1,
  STI_UINT32 = 2,
  STI_UINT64 = 3,
  STI_UINT128 = 4,
  STI_UINT256 = 5,
  STI_AMOUNT = 6,
  STI_VL = 7,
  STI_ACCOUNT = 8,
  STI_NUMBER = 9,
  STI_INT32 = 10,
  STI_INT64 = 11,

  STI_OBJECT = 14,
  STI_ARRAY = 15,

  STI_UINT8 = 16,
  STI_UINT160 = 17,
  STI_PATHSET = 18,
  STI_VECTOR256 = 19,
  STI_UINT96 = 20,
  STI_UINT192 = 21,
  STI_UINT384 = 22,
  STI_UINT512 = 23,
  STI_ISSUE = 24,
  STI_XCHAIN_BRIDGE = 25,
  STI_CURRENCY = 26,
  STI_DATA = 27,
  STI_DATATYPE = 28,
  STI_JSON = 29,
}

/**
 * Map of type strings to SerializedTypeID values
 */
export const TYPE_STRING_TO_ID: Record<string, SerializedTypeID> = {
  NOTPRESENT: SerializedTypeID.STI_NOTPRESENT,
  UINT16: SerializedTypeID.STI_UINT16,
  UINT32: SerializedTypeID.STI_UINT32,
  UINT64: SerializedTypeID.STI_UINT64,
  UINT128: SerializedTypeID.STI_UINT128,
  UINT256: SerializedTypeID.STI_UINT256,
  AMOUNT: SerializedTypeID.STI_AMOUNT,
  VL: SerializedTypeID.STI_VL,
  ACCOUNT: SerializedTypeID.STI_ACCOUNT,
  NUMBER: SerializedTypeID.STI_NUMBER,
  INT32: SerializedTypeID.STI_INT32,
  INT64: SerializedTypeID.STI_INT64,

  OBJECT: SerializedTypeID.STI_OBJECT,
  ARRAY: SerializedTypeID.STI_ARRAY,

  UINT8: SerializedTypeID.STI_UINT8,
  UINT160: SerializedTypeID.STI_UINT160,
  PATHSET: SerializedTypeID.STI_PATHSET,
  VECTOR256: SerializedTypeID.STI_VECTOR256,
  UINT96: SerializedTypeID.STI_UINT96,
  UINT192: SerializedTypeID.STI_UINT192,
  UINT384: SerializedTypeID.STI_UINT384,
  UINT512: SerializedTypeID.STI_UINT512,
  ISSUE: SerializedTypeID.STI_ISSUE,
  XCHAIN_BRIDGE: SerializedTypeID.STI_XCHAIN_BRIDGE,
  CURRENCY: SerializedTypeID.STI_CURRENCY,
  DATA: SerializedTypeID.STI_DATA,
  DATATYPE: SerializedTypeID.STI_DATATYPE,
  JSON: SerializedTypeID.STI_JSON,
}

/**
 * Map of type strings to SerializedTypeID values
 */
export const TYPE_NUMBER_TO_ID: Record<number, SerializedTypeID> = {
  0: SerializedTypeID.STI_NOTPRESENT,
  1: SerializedTypeID.STI_UINT16,
  2: SerializedTypeID.STI_UINT32,
  3: SerializedTypeID.STI_UINT64,
  4: SerializedTypeID.STI_UINT128,
  5: SerializedTypeID.STI_UINT256,
  6: SerializedTypeID.STI_AMOUNT,
  7: SerializedTypeID.STI_VL,
  8: SerializedTypeID.STI_ACCOUNT,
  9: SerializedTypeID.STI_NUMBER,
  10: SerializedTypeID.STI_INT32,
  11: SerializedTypeID.STI_INT64,

  14: SerializedTypeID.STI_OBJECT,
  15: SerializedTypeID.STI_ARRAY,

  16: SerializedTypeID.STI_UINT8,
  17: SerializedTypeID.STI_UINT160,
  18: SerializedTypeID.STI_PATHSET,
  19: SerializedTypeID.STI_VECTOR256,
  20: SerializedTypeID.STI_UINT96,
  21: SerializedTypeID.STI_UINT192,
  22: SerializedTypeID.STI_UINT384,
  23: SerializedTypeID.STI_UINT512,
  24: SerializedTypeID.STI_ISSUE,
  25: SerializedTypeID.STI_XCHAIN_BRIDGE,
  26: SerializedTypeID.STI_CURRENCY,
  27: SerializedTypeID.STI_DATA,
  28: SerializedTypeID.STI_DATATYPE,
  29: SerializedTypeID.STI_JSON,
}

/**
 * Map of SerializedTypeID values to type strings
 */
export const TYPE_ID_TO_STRING: Record<SerializedTypeID, string> = {
  [SerializedTypeID.STI_NOTPRESENT]: '',
  [SerializedTypeID.STI_UINT16]: 'UINT16',
  [SerializedTypeID.STI_UINT32]: 'UINT32',
  [SerializedTypeID.STI_UINT64]: 'UINT64',
  [SerializedTypeID.STI_UINT128]: 'UINT128',
  [SerializedTypeID.STI_UINT256]: 'UINT256',
  [SerializedTypeID.STI_AMOUNT]: 'AMOUNT',
  [SerializedTypeID.STI_VL]: 'VL',
  [SerializedTypeID.STI_ACCOUNT]: 'ACCOUNT',
  [SerializedTypeID.STI_NUMBER]: 'NUMBER',
  [SerializedTypeID.STI_INT32]: 'INT32',
  [SerializedTypeID.STI_INT64]: 'INT64',

  [SerializedTypeID.STI_OBJECT]: 'OBJECT',
  [SerializedTypeID.STI_ARRAY]: 'ARRAY',

  [SerializedTypeID.STI_UINT8]: 'UINT8',
  [SerializedTypeID.STI_UINT160]: 'UINT160',
  [SerializedTypeID.STI_PATHSET]: 'PATHSET',
  [SerializedTypeID.STI_VECTOR256]: 'VECTOR256',
  [SerializedTypeID.STI_UINT96]: 'UINT96',
  [SerializedTypeID.STI_UINT192]: 'UINT192',
  [SerializedTypeID.STI_UINT384]: 'UINT384',
  [SerializedTypeID.STI_UINT512]: 'UINT512',
  [SerializedTypeID.STI_ISSUE]: 'ISSUE',
  [SerializedTypeID.STI_XCHAIN_BRIDGE]: 'XCHAIN_BRIDGE',
  [SerializedTypeID.STI_CURRENCY]: 'CURRENCY',
  [SerializedTypeID.STI_DATA]: 'DATA',
  [SerializedTypeID.STI_DATATYPE]: 'DATATYPE',
  [SerializedTypeID.STI_JSON]: 'JSON',
}

type JSON = string | number | boolean | null | undefined | JSON[] | JsonObject

type JsonObject = { [key: string]: JSON }

/**
 * The base class for all binary-codec types
 */
class SerializedType {
  protected readonly bytes: Uint8Array = new Uint8Array(0)

  constructor(bytes?: Uint8Array) {
    this.bytes = bytes ?? new Uint8Array(0)
  }

  static fromParser(parser: BinaryParser, hint?: number): SerializedType {
    throw new Error('fromParser not implemented')
    return this.fromParser(parser, hint)
  }

  static from(value: SerializedType | JSON | bigint): SerializedType {
    throw new Error('from not implemented')
    return this.from(value)
  }

  /**
   * Write the bytes representation of a SerializedType to a BytesList
   *
   * @param list The BytesList to write SerializedType bytes to
   */
  toBytesSink(list: BytesList): void {
    list.put(this.bytes)
  }

  /**
   * Get the hex representation of a SerializedType's bytes
   *
   * @returns hex String of this.bytes
   */
  toHex(): string {
    return bytesToHex(this.toBytes())
  }

  /**
   * Get the bytes representation of a SerializedType
   *
   * @returns A Uint8Array of the bytes
   */
  toBytes(): Uint8Array {
    if (this.bytes) {
      return this.bytes
    }
    const bytes = new BytesList()
    this.toBytesSink(bytes)
    return bytes.toBytes()
  }

  /**
   * Return the JSON representation of a SerializedType
   *
   * @param _definitions rippled definitions used to parse the values of transaction types and such.
   *                          Unused in default, but used in STObject, STArray
   *                          Can be customized for sidechains and amendments.
   * @returns any type, if not overloaded returns hexString representation of bytes
   */
  toJSON(_definitions?: XrplDefinitionsBase, _fieldName?: string): JSON {
    return this.toHex()
  }

  /**
   * @returns hexString representation of this.bytes
   */
  toString(): string {
    return this.toHex()
  }

  getSType(): SerializedTypeID {
    return this.getSType()
  }
}

/**
 * Base class for SerializedTypes that are comparable.
 *
 * @template T - What types you want to allow comparisons between. You must specify all types. Primarily used to allow
 * comparisons between built-in types (like `string`) and SerializedType subclasses (like `Hash`).
 *
 * Ex. `class Hash extends Comparable<Hash | string>`
 */
class Comparable<T extends Object> extends SerializedType {
  lt(other: T): boolean {
    return this.compareTo(other) < 0
  }

  eq(other: T): boolean {
    return this.compareTo(other) === 0
  }

  gt(other: T): boolean {
    return this.compareTo(other) > 0
  }

  gte(other: T): boolean {
    return this.compareTo(other) > -1
  }

  lte(other: T): boolean {
    return this.compareTo(other) < 1
  }

  /**
   * Overload this method to define how two Comparable SerializedTypes are compared
   *
   * @param other The comparable object to compare this to
   * @returns A number denoting the relationship of this and other
   */
  compareTo(other: T): number {
    throw new Error(`cannot compare ${this.toString()} and ${other.toString()}`)
  }
}

export { SerializedType, Comparable, JSON, JsonObject }
