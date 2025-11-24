#!/usr/bin/env node

/**
 * XRPL Smart Contract Call Module
 *
 * This module handles calling functions on deployed XRPL smart contracts.
 * It uses the ABI to construct properly-typed parameter values.
 *
 * Usage: node call.js <config-json-path>
 *
 * Config JSON format:
 * {
 *   "contract_account": "rContract123...",
 *   "function_name": "register",
 *   "network_url": "wss://alphanet.xrpl.org",
 *   "wallet_seed": "sXXX...",
 *   "abi_path": "/path/to/abi.json" (optional),
 *   "parameters": {"name": "test", "duration": 31536000} (optional),
 *   "computation_allowance": "1000000" (optional),
 *   "fee": "1000000" (optional),
 *   "verbose": true (optional)
 * }
 *
 * Output JSON format:
 * {
 *   "success": true,
 *   "data": {
 *     "txHash": "...",
 *     "returnCode": 0,
 *     "returnValue": "...",
 *     "gasUsed": 12345,
 *     "validated": true
 *   }
 * }
 */

const xrpl = require('@transia/xrpl');
const fs = require('fs');

/**
 * Build Parameters array from ABI and provided values
 */
function buildParametersFromABI(functionDef, paramValues) {
  if (!functionDef.parameters || functionDef.parameters.length === 0) {
    return undefined;
  }

  const parameters = [];

  for (let i = 0; i < functionDef.parameters.length; i++) {
    const paramDef = functionDef.parameters[i];
    const value = paramValues[paramDef.name] || paramValues[i];

    if (value === undefined && paramDef.flag === 0) {
      throw new Error(
        `Required parameter '${paramDef.name}' (${paramDef.type}) not provided`
      );
    }

    if (value !== undefined) {
      parameters.push({
        ParameterValue: {
          ParameterFlag: paramDef.flag,
          ParameterValue: formatParameterValue(paramDef.type, value),
        },
      });
    }
  }

  return parameters.length > 0 ? parameters : undefined;
}

/**
 * Format parameter value according to XRPL type
 */
function formatParameterValue(type, value) {
  switch (type) {
    case 'UINT8':
    case 'UINT16':
    case 'UINT32':
    case 'UINT64':
    case 'UINT128':
    case 'UINT160':
    case 'UINT192':
    case 'UINT256':
      return {
        type: type,
        value: value.toString(),
      };

    case 'VL':
      // Variable length - convert string to hex
      return {
        type: 'VL',
        value:
          typeof value === 'string' && !value.startsWith('0x')
            ? Buffer.from(value).toString('hex').toUpperCase()
            : value.replace('0x', '').toUpperCase(),
      };

    case 'ACCOUNT':
      return {
        type: 'ACCOUNT',
        value: value, // r-address format
      };

    case 'AMOUNT':
      // XRP drops or IOU
      if (typeof value === 'string' || typeof value === 'number') {
        return {
          type: 'AMOUNT',
          value: value.toString(),
        };
      }
      return {
        type: 'AMOUNT',
        value: value, // Already formatted object
      };

    case 'NUMBER':
      return {
        type: 'NUMBER',
        value: parseFloat(value),
      };

    case 'CURRENCY':
      return {
        type: 'CURRENCY',
        value: value,
      };

    case 'ISSUE':
      return {
        type: 'ISSUE',
        value: value,
      };

    default:
      throw new Error(`Unsupported parameter type: ${type}`);
  }
}

/**
 * Call a contract function
 */
async function callContract(config) {
  const {
    contract_account,
    function_name,
    network_url,
    wallet_seed,
    abi_path,
    parameters,
    computation_allowance,
    fee,
    verbose,
  } = config;

  const log = verbose ? console.error.bind(console) : () => {};

  log('Calling smart contract on XRPL...\n');

  const client = new xrpl.Client(network_url);

  try {
    await client.connect();
    log('✓ Connected to network');

    // Create or restore wallet
    const wallet = wallet_seed
      ? xrpl.Wallet.fromSeed(wallet_seed, { algorithm: xrpl.ECDSA.secp256k1 })
      : xrpl.Wallet.generate();

    log('\nWallet:');
    log('  Address:', wallet.address);

    log(`\nContract: ${contract_account}`);
    log(`Function: ${function_name}`);

    // Load ABI if provided
    let functionDef = null;
    let Parameters = undefined;

    if (abi_path && fs.existsSync(abi_path)) {
      log(`\nLoading ABI from: ${abi_path}`);
      const abiContent = fs.readFileSync(abi_path, 'utf8');
      const abi = JSON.parse(abiContent);

      // Find function in ABI
      functionDef = abi.functions.find((f) => f.name === function_name);

      if (functionDef) {
        log(`\nFunction signature:`);
        log(`  ${function_name}(`);
        if (functionDef.parameters) {
          functionDef.parameters.forEach((p) => {
            const required = p.flag === 0 ? 'required' : 'optional';
            log(`    ${p.name}: ${p.type} (${required})`);
          });
        }
        log(`  )`);
        if (functionDef.returns) {
          log(`  -> ${functionDef.returns.type}`);
        }

        // Build parameters from ABI
        if (parameters) {
          log(`\nProvided parameters:`);
          log(JSON.stringify(parameters, null, 2));

          Parameters = buildParametersFromABI(functionDef, parameters);

          if (Parameters) {
            log(`\nFormatted parameters:`);
            log(JSON.stringify(Parameters, null, 2));
          }
        }
      } else {
        log(`\nWarning: Function "${function_name}" not found in ABI`);
      }
    }

    // Check balance
    const balance = await client.getXrpBalance(wallet.address);
    log(`\nWallet balance: ${balance} XRP`);

    if (parseFloat(balance) === 0) {
      log('Warning: Wallet not funded, call will likely fail');
    }

    // Create ContractCall transaction
    log('\nSubmitting contract call transaction...');

    const functionNameHex = Buffer.from(function_name)
      .toString('hex')
      .toUpperCase();

    const tx = {
      TransactionType: 'ContractCall',
      Account: wallet.address,
      ContractAccount: contract_account,
      FunctionName: functionNameHex,
      Parameters: Parameters,
      ComputationAllowance: computation_allowance || '1000000',
      Fee: fee || '1000000', // 1 XRP default
    };

    const prepared = await client.autofill(tx);
    const signed = wallet.sign(prepared);

    log('Transaction ID:', signed.hash);

    const result = await client.submitAndWait(signed.tx_blob);

    log('\n✓ Contract function called successfully!');

    // Extract result information
    const txResult = result.result;
    const meta = txResult.meta;

    log('\nTransaction Status:');
    log(`  Result: ${meta?.TransactionResult || 'N/A'}`);
    log(`  Validated: ${txResult.validated ? 'Yes' : 'No'}`);

    if (meta?.WasmReturnCode !== undefined) {
      log('\nContract Return Code:');
      log(`  Value: ${meta.WasmReturnCode}`);
      log(`  Status: ${meta.WasmReturnCode === 0 ? 'SUCCESS' : 'ERROR'}`);
    }

    if (meta?.ReturnValue) {
      log('\nContract Return Data:');
      log(`  Hex: ${meta.ReturnValue}`);
      log(`  Decimal: ${parseInt(meta.ReturnValue, 16)}`);

      // Try to decode based on ABI return type
      if (functionDef?.returns) {
        log(`  Type: ${functionDef.returns.type}`);
        log(`  Description: ${functionDef.returns.description || 'N/A'}`);
      }
    }

    if (meta?.GasUsed !== undefined) {
      log('\nGas/Computation Used:');
      log(`  Gas Used: ${meta.GasUsed}`);
      log(`  Allowance: ${prepared.ComputationAllowance}`);
      const percentage = (
        (meta.GasUsed / parseInt(prepared.ComputationAllowance)) *
        100
      ).toFixed(2);
      log(`  Percentage: ${percentage}%`);
    }

    await client.disconnect();
    log('\n✓ Disconnected');

    // Output call result as clean JSON to stdout
    const callResult = {
      success: true,
      data: {
        txHash: signed.hash,
        returnCode: meta?.WasmReturnCode,
        returnValue: meta?.ReturnValue,
        gasUsed: meta?.GasUsed,
        validated: txResult.validated,
        transactionResult: meta?.TransactionResult,
        meta: meta,
      },
    };

    console.log(JSON.stringify(callResult));
    return callResult;
  } catch (error) {
    if (client.isConnected()) {
      await client.disconnect();
    }

    // Output error as JSON
    const errorResult = {
      success: false,
      error: error.message,
      details: error.data ? JSON.stringify(error.data) : error.stack,
    };

    console.log(JSON.stringify(errorResult));
    process.exit(1);
  }
}

// CLI interface
if (require.main === module) {
  const args = process.argv.slice(2);

  if (args.length < 1) {
    console.error(`
Usage: node call.js <config-json-path>

The config JSON file should contain:
{
  "contract_account": "rContract123...",
  "function_name": "register",
  "network_url": "wss://alphanet.xrpl.org",
  "wallet_seed": "sXXX...",
  "abi_path": "/path/to/abi.json" (optional),
  "parameters": {"name": "test", "duration": 31536000} (optional),
  "computation_allowance": "1000000" (optional),
  "fee": "1000000" (optional),
  "verbose": true (optional)
}

Output is pure JSON to stdout.
`);
    process.exit(1);
  }

  const configPath = args[0];

  if (!fs.existsSync(configPath)) {
    const errorResult = {
      success: false,
      error: `Config file not found: ${configPath}`,
      details: 'Please provide a valid config JSON file path',
    };
    console.log(JSON.stringify(errorResult));
    process.exit(1);
  }

  try {
    const configContent = fs.readFileSync(configPath, 'utf8');
    const config = JSON.parse(configContent);
    callContract(config);
  } catch (error) {
    const errorResult = {
      success: false,
      error: 'Failed to load config',
      details: error.message,
    };
    console.log(JSON.stringify(errorResult));
    process.exit(1);
  }
}

module.exports = { callContract, buildParametersFromABI, formatParameterValue };
