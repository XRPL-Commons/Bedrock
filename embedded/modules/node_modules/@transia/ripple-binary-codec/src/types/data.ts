/* eslint-disable max-lines */
/* eslint-disable complexity */
import { BinaryParser } from '../serdes/binary-parser'
import {
  JsonObject,
  SerializedType,
  SerializedTypeID,
  TYPE_ID_TO_STRING,
  TYPE_STRING_TO_ID,
  TYPE_NUMBER_TO_ID,
} from './serialized-type'
import { readUInt16BE, writeUInt16BE } from '../utils'
import { bytesToHex, concat } from '@transia/isomorphic/utils'
import { Hash128 } from './hash-128'
import { Hash160 } from './hash-160'
import { Hash192 } from './hash-192'
import { Hash256 } from './hash-256'
import { AccountID } from './account-id'
import { Amount, AmountObject } from './amount'
import { Blob } from './blob'
import { Currency } from './currency'
import { STNumber } from './st-number'
import { Issue, IssueObject } from './issue'
import { UInt8 } from './uint-8'
import { UInt16 } from './uint-16'
import { UInt32 } from './uint-32'
import { UInt64 } from './uint-64'
import { BinarySerializer } from '../binary'

/**
 * Interface for Data JSON representation
 */
interface DataJSON extends JsonObject {
  type: string
  value: string | number | JsonObject
}

/**
 * Type union for all possible data values
 */
type DataValue =
  | number
  | string
  | bigint
  | Uint8Array
  | UInt8
  | UInt16
  | UInt32
  | UInt64
  | Hash128
  | Hash160
  | Hash192
  | Hash256
  | AccountID
  | Amount
  | Blob
  | Currency
  | STNumber
  | Issue

/**
 * STData: Encodes XRPL's "Data" type.
 *
 * This type wraps both a SerializedTypeID and the actual data value.
 * It's encoded as a 2-byte type ID followed by the serialized data.
 *
 * Usage:
 *   Data.from({ type: "AMOUNT", value: "1000000" })
 *   Data.from({ type: "UINT64", value: "123456789" })
 *   Data.fromParser(parser)
 */
class Data extends SerializedType {
  static readonly ZERO_DATA: Data = new Data(
    concat([
      new Uint8Array([0x00, 0x01]), // Type ID for UINT16 (SerializedTypeID.STI_UINT16 = 1) as uint16
      new Uint8Array([0x00, 0x00]), // Value: two zero bytes for UINT16
    ]),
  )

  /**
   * Construct Data from bytes
   * @param bytes - Uint8Array containing type ID and data
   */
  constructor(bytes: Uint8Array) {
    super(bytes ?? Data.ZERO_DATA.bytes)
  }

  /**
   * Create Data from various input types
   *
   * @param value - Can be:
   *   - Data instance (returns as-is)
   *   - DataJSON object with 'type' and 'value' fields
   * @returns Data instance
   * @throws Error if value type is not supported
   */
  static from(value: unknown): Data {
    if (value instanceof Data) {
      return value
    }

    if (
      typeof value === 'object' &&
      value !== null &&
      'type' in value &&
      'value' in value
    ) {
      const json = value as DataJSON
      return Data.fromJSON(json)
    }

    throw new Error('Data.from: value must be Data instance or DataJSON object')
  }

  /**
   * Create Data from JSON representation
   *
   * @param json - Object with 'type' and 'value' fields
   * @returns Data instance
   * @throws Error if type is not supported
   */
  static fromJSON(json: DataJSON): Data {
    const typeId = TYPE_STRING_TO_ID[json.type]
    if (typeId === undefined) {
      throw new Error(`Data: unsupported type string: ${json.type}`)
    }

    let dataValue: DataValue
    let dataBytes: Uint8Array

    switch (typeId) {
      case SerializedTypeID.STI_UINT8: {
        const val =
          typeof json.value === 'string'
            ? parseInt(json.value, 10)
            : typeof json.value === 'number'
            ? json.value
            : Number(json.value)
        if (
          typeof val !== 'number' ||
          Number.isNaN(val) ||
          val < 0 ||
          val > 255
        ) {
          throw new Error('UINT8 value out of range')
        }
        dataValue = UInt8.from(val)
        dataBytes = (dataValue as UInt8).toBytes()
        break
      }

      case SerializedTypeID.STI_UINT16: {
        const val =
          typeof json.value === 'string'
            ? parseInt(json.value, 10)
            : typeof json.value === 'number'
            ? json.value
            : Number(json.value)
        if (
          typeof val !== 'number' ||
          Number.isNaN(val) ||
          val < 0 ||
          val > 65535
        ) {
          throw new Error('UINT16 value out of range')
        }
        dataValue = UInt16.from(val)
        dataBytes = (dataValue as UInt16).toBytes()
        break
      }

      case SerializedTypeID.STI_UINT32: {
        const val =
          typeof json.value === 'string'
            ? parseInt(json.value, 10)
            : typeof json.value === 'number'
            ? json.value
            : Number(json.value)
        dataValue = UInt32.from(val)
        dataBytes = (dataValue as UInt32).toBytes()
        break
      }

      case SerializedTypeID.STI_UINT64: {
        const val =
          typeof json.value === 'string' ? json.value : json.value.toString()
        dataValue = UInt64.from(val)
        dataBytes = (dataValue as UInt64).toBytes()
        break
      }

      case SerializedTypeID.STI_UINT128: {
        const val =
          typeof json.value === 'string' ? json.value : json.value.toString()
        dataValue = Hash128.from(val)
        dataBytes = (dataValue as Hash128).toBytes()
        break
      }

      case SerializedTypeID.STI_UINT160: {
        const val =
          typeof json.value === 'string' ? json.value : json.value.toString()
        dataValue = Hash160.from(val)
        dataBytes = (dataValue as Hash160).toBytes()
        break
      }

      case SerializedTypeID.STI_UINT192: {
        const val =
          typeof json.value === 'string' ? json.value : json.value.toString()
        dataValue = Hash192.from(val)
        dataBytes = (dataValue as Hash192).toBytes()
        break
      }

      case SerializedTypeID.STI_UINT256: {
        const val =
          typeof json.value === 'string' ? json.value : json.value.toString()
        dataValue = Hash256.from(val)
        dataBytes = (dataValue as Hash256).toBytes()
        break
      }

      case SerializedTypeID.STI_VL: {
        const val =
          typeof json.value === 'string' ? json.value : json.value.toString()
        dataValue = Blob.from(val)
        dataBytes = dataValue.toBytes()
        const lengthBytes = BinarySerializer.encodeVariableLength(
          dataBytes.length,
        )
        dataBytes = concat([lengthBytes, dataBytes])
        break
      }

      case SerializedTypeID.STI_ACCOUNT: {
        dataValue = AccountID.from(
          typeof json.value === 'string' ? json.value : json.value.toString(),
        )
        dataBytes = (dataValue as AccountID).toBytes()
        dataBytes = concat([new Uint8Array([0x14]), dataBytes])
        break
      }

      case SerializedTypeID.STI_AMOUNT: {
        dataValue = Amount.from(json.value as AmountObject)
        dataBytes = (dataValue as Amount).toBytes()
        break
      }

      case SerializedTypeID.STI_ISSUE: {
        dataValue = Issue.from(json.value as IssueObject)
        dataBytes = (dataValue as Issue).toBytes()
        break
      }

      case SerializedTypeID.STI_CURRENCY: {
        const val =
          typeof json.value === 'string' ? json.value : json.value.toString()
        dataValue = Currency.from(val)
        dataBytes = (dataValue as Currency).toBytes()
        break
      }

      case SerializedTypeID.STI_NUMBER: {
        dataValue = STNumber.from(json.value)
        dataBytes = (dataValue as STNumber).toBytes()
        break
      }

      default:
        throw new Error(`Data.fromJSON(): unsupported type ID: ${typeId}`)
    }

    // Combine type header with data bytes
    const typeBytes = new Uint8Array(2)
    writeUInt16BE(typeBytes, typeId, 0)
    const fullBytes = concat([typeBytes, dataBytes])
    return new Data(fullBytes)
  }

  /**
   * Read Data from a BinaryParser stream
   *
   * @param parser - BinaryParser positioned at the start of Data
   * @returns Data instance
   */
  static fromParser(parser: BinaryParser): Data {
    // Read the 2-byte type ID
    const typeBytes = parser.read(2)
    const typeId = TYPE_NUMBER_TO_ID[readUInt16BE(typeBytes, 0)]

    let dataValue: DataValue
    let dataBytes: Uint8Array

    switch (typeId) {
      case SerializedTypeID.STI_UINT8:
        dataValue = UInt8.fromParser(parser)
        dataBytes = (dataValue as UInt8).toBytes()
        break

      case SerializedTypeID.STI_UINT16:
        dataValue = UInt16.fromParser(parser)
        dataBytes = (dataValue as UInt16).toBytes()
        break

      case SerializedTypeID.STI_UINT32:
        dataValue = UInt32.fromParser(parser)
        dataBytes = (dataValue as UInt32).toBytes()
        break

      case SerializedTypeID.STI_UINT64:
        dataValue = UInt64.fromParser(parser)
        dataBytes = (dataValue as UInt64).toBytes()
        break

      case SerializedTypeID.STI_UINT128:
        dataValue = Hash128.fromParser(parser)
        dataBytes = (dataValue as Hash128).toBytes()
        break

      case SerializedTypeID.STI_UINT160:
        dataValue = Hash160.fromParser(parser)
        dataBytes = (dataValue as Hash160).toBytes()
        break

      case SerializedTypeID.STI_UINT192:
        dataValue = Hash192.fromParser(parser)
        dataBytes = (dataValue as Hash192).toBytes()
        break

      case SerializedTypeID.STI_UINT256:
        dataValue = Hash256.fromParser(parser)
        dataBytes = (dataValue as Hash256).toBytes()
        break

      case SerializedTypeID.STI_VL:
        const valueVL = parser.readVariableLength()
        dataValue = Blob.from(bytesToHex(valueVL))
        dataBytes = concat([
          BinarySerializer.encodeVariableLength(valueVL.length),
          valueVL,
        ])
        break

      case SerializedTypeID.STI_ACCOUNT:
        parser.skip(1)
        dataValue = AccountID.fromParser(parser)
        dataBytes = concat([
          new Uint8Array([0x14]),
          (dataValue as AccountID).toBytes(),
        ])
        break

      case SerializedTypeID.STI_AMOUNT:
        dataValue = Amount.fromParser(parser)
        dataBytes = (dataValue as Amount).toBytes()
        break

      case SerializedTypeID.STI_ISSUE:
        dataValue = Issue.fromParser(parser)
        dataBytes = (dataValue as Issue).toBytes()
        break

      case SerializedTypeID.STI_CURRENCY:
        dataValue = Currency.fromParser(parser)
        dataBytes = (dataValue as Currency).toBytes()
        break

      case SerializedTypeID.STI_NUMBER:
        dataValue = STNumber.fromParser(parser)
        dataBytes = (dataValue as STNumber).toBytes()
        break

      default:
        throw new Error(`Data: unsupported type ID when parsing: ${typeId}`)
    }

    const fullBytes = concat([typeBytes, dataBytes])
    return new Data(fullBytes)
  }

  /**
   * Get the inner SerializedTypeID
   *
   * @returns The inner type ID
   */
  getInnerType(): SerializedTypeID {
    return TYPE_NUMBER_TO_ID[readUInt16BE(this.bytes, 0)]
  }

  /**
   * Get the string representation of the inner type
   *
   * @returns String name of the type
   */
  getInnerTypeString(): string {
    const innerType = this.getInnerType()
    return TYPE_ID_TO_STRING[innerType] || innerType.toString()
  }

  /**
   * Get the data value
   *
   * @returns The stored data value
   */
  getValue(): DataValue {
    const innerType = this.getInnerType()
    const parser = new BinaryParser(bytesToHex(this.bytes.slice(2)))

    switch (innerType) {
      case SerializedTypeID.STI_UINT8:
        return UInt8.fromParser(parser)
      case SerializedTypeID.STI_UINT16:
        return UInt16.fromParser(parser)
      case SerializedTypeID.STI_UINT32:
        return UInt32.fromParser(parser)
      case SerializedTypeID.STI_UINT64:
        return UInt64.fromParser(parser)
      case SerializedTypeID.STI_UINT128:
        return Hash128.fromParser(parser)
      case SerializedTypeID.STI_UINT160:
        return Hash160.fromParser(parser)
      case SerializedTypeID.STI_UINT192:
        return Hash192.fromParser(parser)
      case SerializedTypeID.STI_UINT256:
        return Hash256.fromParser(parser)
      case SerializedTypeID.STI_VL:
        const vlLength = parser.readVariableLengthLength()
        return Blob.fromParser(parser, vlLength)
      case SerializedTypeID.STI_ACCOUNT:
        parser.skip(1)
        return AccountID.fromParser(parser)
      case SerializedTypeID.STI_AMOUNT:
        return Amount.fromParser(parser)
      case SerializedTypeID.STI_ISSUE:
        return Issue.fromParser(parser)
      case SerializedTypeID.STI_CURRENCY:
        return Currency.fromParser(parser)
      case SerializedTypeID.STI_NUMBER:
        return STNumber.fromParser(parser)
      default:
        throw new Error(
          `Data.getValue(): unsupported type ID: ${typeof innerType}`,
        )
    }
  }

  /**
   * Convert to JSON representation
   *
   * @returns JSON object with 'type' and 'value' fields
   */
  toJSON(): DataJSON {
    const data = this.getValue()
    let jsonValue: string | number | JsonObject

    // Convert the data value to its JSON representation
    if (data instanceof SerializedType) {
      jsonValue = data.toJSON() as JsonObject
    } else if (data instanceof Uint8Array) {
      jsonValue = bytesToHex(data)
    } else if (typeof data === 'bigint') {
      jsonValue = data.toString()
    } else {
      jsonValue = data
    }

    return {
      type: this.getInnerTypeString(),
      value: jsonValue,
    }
  }

  /**
   * Compare with another Data for equality
   *
   * @param other - Another Data to compare with
   * @returns true if both have the same inner type and data
   */
  equals(other: Data): boolean {
    if (!(other instanceof Data)) {
      return false
    }

    // Compare bytes directly
    if (this.bytes.length !== other.bytes.length) {
      return false
    }

    for (let i = 0; i < this.bytes.length; i++) {
      if (this.bytes[i] !== other.bytes[i]) {
        return false
      }
    }

    return true
  }

  getSType(): SerializedTypeID {
    return SerializedTypeID.STI_DATA
  }
}

export { Data }
