#!/usr/bin/env node

/**
 * Patch @transia/ripple-binary-codec definitions.json to match
 * the field codes in Transia-RnD/rippled alphanet-develop (f5d78179).
 *
 * The npm package has stale nth values for contract-related fields.
 * This script corrects them to match the server's sfields.macro.
 */

const fs = require('fs');
const path = require('path');

const defsPath = path.join(
  __dirname,
  'node_modules',
  '@transia',
  'ripple-binary-codec',
  'dist',
  'enums',
  'definitions.json'
);

if (!fs.existsSync(defsPath)) {
  console.log('postinstall: definitions.json not found, skipping patch');
  process.exit(0);
}

const defs = JSON.parse(fs.readFileSync(defsPath, 'utf8'));

// Correct nth values from rippled include/xrpl/protocol/detail/sfields.macro
const fixes = {
  Function: 38,               // OBJECT, was 37 in codec
  Parameter: 41,              // OBJECT, was 40
  InstanceParameter: 39,      // OBJECT, was 38
  InstanceParameterValue: 40, // OBJECT, was 39
  Functions: 32,              // ARRAY, was 33
  InstanceParameters: 33,     // ARRAY, was 34
  InstanceParameterValues: 34,// ARRAY, was 35
  Parameters: 35,             // ARRAY, was 36
  ParameterFlag: 74,          // UINT32, was 59
  ContractHash: 39,           // UINT256, was 37
  ContractID: 40,             // UINT256, was 38
  ContractAccount: 27,        // ACCOUNT, was 25
  ComputationAllowance: 72,   // UINT32, was 57
};

// Correct transaction type codes from rippled include/xrpl/protocol/detail/transactions.macro
const txTypeFixes = {
  ContractCreate: 85,     // was 72 in codec
  ContractModify: 86,     // was 73
  ContractDelete: 87,     // was 74
  ContractClawback: 88,   // was 75
  ContractUserDelete: 89, // was 76
  ContractCall: 90,       // was 77
};

let patched = 0;

// Patch field definitions
for (const item of defs.FIELDS) {
  if (Array.isArray(item) && item.length === 2) {
    const [name, info] = item;
    if (name in fixes && info.nth !== fixes[name]) {
      info.nth = fixes[name];
      patched++;
    }
  }
}

// Patch transaction types
if (defs.TRANSACTION_TYPES) {
  for (const [name, correctCode] of Object.entries(txTypeFixes)) {
    if (name in defs.TRANSACTION_TYPES && defs.TRANSACTION_TYPES[name] !== correctCode) {
      defs.TRANSACTION_TYPES[name] = correctCode;
      patched++;
    }
  }
}

if (patched > 0) {
  fs.writeFileSync(defsPath, JSON.stringify(defs));
  console.log(`postinstall: patched ${patched} definitions in binary codec`);
} else {
  console.log('postinstall: definitions already correct, no patch needed');
}
