#!/usr/bin/env node

/**
 * XRPL Smart Contract Modify Module
 *
 * Handles modifying deployed contracts via ContractModify transaction.
 *
 * Usage: node modify.js <config-json-path>
 */

const xrpl = require('@transia/xrpl');
const fs = require('fs');

function buildFunctionsFromABI(abi, exportedFunctions) {
  const functions = [];

  for (const fn of abi.functions) {
    if (!exportedFunctions.includes(fn.name)) {
      continue;
    }

    const parameters = fn.parameters.map((param) => ({
      Parameter: {
        ParameterName: Buffer.from(param.name).toString('hex').toUpperCase(),
        ParameterType: {
          type: param.type,
        },
      },
    }));

    functions.push({
      Function: {
        FunctionName: Buffer.from(fn.name).toString('hex').toUpperCase(),
        Parameters: parameters.length > 0 ? parameters : undefined,
      },
    });
  }

  return functions;
}

async function modifyContract(config) {
  const {
    contract_account,
    network_url,
    wallet_seed,
    algorithm,
    wasm_path,
    abi_path,
    fee,
    verbose,
  } = config;

  const log = verbose ? console.error.bind(console) : () => {};

  log('Modifying contract on XRPL...\n');

  const client = new xrpl.Client(network_url);

  try {
    await client.connect();
    log('Connected to network');

    const algo = algorithm === 'ed25519' ? undefined : xrpl.ECDSA.secp256k1;
    const wallet = algo
      ? xrpl.Wallet.fromSeed(wallet_seed, { algorithm: algo })
      : xrpl.Wallet.fromSeed(wallet_seed);

    const tx = {
      TransactionType: 'ContractModify',
      Account: wallet.address,
      ContractAccount: contract_account,
      Fee: fee || '10000000',
    };

    // Optionally update WASM code
    if (wasm_path && fs.existsSync(wasm_path)) {
      const wasmBytes = fs.readFileSync(wasm_path);
      tx.ContractCode = wasmBytes.toString('hex').toUpperCase();
      log(`Updated WASM: ${wasmBytes.length} bytes`);
    }

    // Optionally update ABI
    if (abi_path && fs.existsSync(abi_path)) {
      const abiContent = fs.readFileSync(abi_path, 'utf8');
      const abi = JSON.parse(abiContent);
      const functionNames = abi.functions.map(f => f.name);
      tx.Functions = buildFunctionsFromABI(abi, functionNames);
      log(`Updated ABI: ${abi.functions.length} functions`);
    }

    const prepared = await client.autofill(tx);
    const signed = wallet.sign(prepared);

    log('Transaction ID:', signed.hash);

    const result = await client.submitAndWait(signed.tx_blob);

    await client.disconnect();

    console.log(JSON.stringify({
      success: true,
      data: {
        txHash: signed.hash,
        validated: result.result.validated,
        meta: result.result.meta,
      },
    }));
  } catch (error) {
    if (client.isConnected()) {
      await client.disconnect();
    }

    console.log(JSON.stringify({
      success: false,
      error: error.message,
      details: error.data ? JSON.stringify(error.data) : error.stack,
    }));
    process.exit(1);
  }
}

if (require.main === module) {
  const args = process.argv.slice(2);
  if (args.length < 1) {
    console.error('Usage: node modify.js <config-json-path>');
    process.exit(1);
  }

  const configPath = args[0];
  if (!fs.existsSync(configPath)) {
    console.log(JSON.stringify({
      success: false,
      error: `Config file not found: ${configPath}`,
      details: 'Please provide a valid config JSON file path',
    }));
    process.exit(1);
  }

  const configContent = fs.readFileSync(configPath, 'utf8');
  modifyContract(JSON.parse(configContent));
}

module.exports = { modifyContract };
