#!/usr/bin/env node

/**
 * XRPL Jade Operations Module - Quick one-liners for XRPL operations
 *
 * Usage: node jade.js <config-json-path>
 *
 * Config JSON format:
 * {
 *   "operation": "balance|send|tx|account_info|server_info",
 *   "network_url": "wss://...",
 *   "network_id": 21465,
 *   "params": { ... operation-specific params ... },
 *   "verbose": false
 * }
 *
 * Operations:
 *
 * balance:
 *   params: { "address": "rXXX..." }
 *   returns: { "address": "...", "balance": "100.5", "balance_drops": "100500000" }
 *
 * send:
 *   params: { "wallet_seed": "sXXX...", "destination": "rXXX...", "amount": "10", "algorithm": "secp256k1" }
 *   returns: { "tx_hash": "...", "from": "...", "to": "...", "amount": "10", "result": "tesSUCCESS" }
 *
 * tx:
 *   params: { "hash": "XXXXXXXX..." }
 *   returns: { "hash": "...", "type": "Payment", "result": "tesSUCCESS", "account": "...", ... }
 *
 * account_info:
 *   params: { "address": "rXXX..." }
 *   returns: { "address": "...", "balance": "...", "sequence": 1, "owner_count": 0, ... }
 *
 * server_info:
 *   params: {}
 *   returns: { "build_version": "...", "network_id": 21465, "server_state": "full", ... }
 */

const xrpl = require('@transia/xrpl');
const fs = require('fs');

/**
 * Get account balance
 */
async function getBalance(client, params, log) {
  const { address } = params;

  if (!address) {
    throw new Error('Missing required parameter: address');
  }

  log(`Getting balance for ${address}...`);

  try {
    const balance = await client.getXrpBalance(address);
    const balanceDrops = xrpl.xrpToDrops(balance);

    return {
      address,
      balance: String(balance),
      balance_drops: balanceDrops,
    };
  } catch (error) {
    if (error.message.includes('Account not found')) {
      return {
        address,
        balance: '0',
        balance_drops: '0',
        funded: false,
      };
    }
    throw error;
  }
}

/**
 * Send XRP payment
 */
async function sendPayment(client, params, log) {
  const { wallet_seed, destination, amount, algorithm } = params;

  if (!wallet_seed) {
    throw new Error('Missing required parameter: wallet_seed');
  }
  if (!destination) {
    throw new Error('Missing required parameter: destination');
  }
  if (!amount) {
    throw new Error('Missing required parameter: amount');
  }

  // Parse algorithm
  const algo =
    algorithm === 'ed25519' ? undefined : xrpl.ECDSA.secp256k1;
  const wallet = algo
    ? xrpl.Wallet.fromSeed(wallet_seed, { algorithm: algo })
    : xrpl.Wallet.fromSeed(wallet_seed);

  log(`Sending ${amount} XRP from ${wallet.address} to ${destination}...`);

  // Convert XRP to drops
  const amountDrops = xrpl.xrpToDrops(amount);

  const payment = {
    TransactionType: 'Payment',
    Account: wallet.address,
    Destination: destination,
    Amount: amountDrops,
  };

  const prepared = await client.autofill(payment);
  const signed = wallet.sign(prepared);
  const result = await client.submitAndWait(signed.tx_blob);

  const txResult = result.result.meta.TransactionResult;
  log(`Transaction result: ${txResult}`);

  return {
    tx_hash: result.result.hash,
    from: wallet.address,
    to: destination,
    amount: amount,
    amount_drops: amountDrops,
    result: txResult,
    fee: result.result.Fee,
    sequence: result.result.Sequence,
    validated: result.result.validated,
  };
}

/**
 * Get transaction details
 */
async function getTransaction(client, params, log) {
  const { hash } = params;

  if (!hash) {
    throw new Error('Missing required parameter: hash');
  }

  log(`Fetching transaction ${hash}...`);

  const response = await client.request({
    command: 'tx',
    transaction: hash,
    binary: false,
  });

  const tx = response.result;

  // Build a simplified response
  const result = {
    hash: tx.hash,
    type: tx.TransactionType,
    account: tx.Account,
    result: tx.meta?.TransactionResult || tx.validated ? 'tesSUCCESS' : 'unknown',
    fee: tx.Fee,
    sequence: tx.Sequence,
    validated: tx.validated,
    ledger_index: tx.ledger_index,
    date: tx.date,
  };

  // Add type-specific fields
  if (tx.TransactionType === 'Payment') {
    result.destination = tx.Destination;
    result.amount = tx.Amount;
    result.delivered_amount = tx.meta?.delivered_amount;
  } else if (tx.TransactionType === 'ContractCreate') {
    result.contract_account = tx.meta?.ContractAccount;
    result.wasm_size = tx.WasmHex?.length / 2;
  } else if (tx.TransactionType === 'ContractCall') {
    result.contract_account = tx.ContractAccount;
    result.function_name = tx.FunctionName;
    result.return_value = tx.meta?.HookReturnString;
  }

  return result;
}

/**
 * Get account info
 */
async function getAccountInfo(client, params, log) {
  const { address } = params;

  if (!address) {
    throw new Error('Missing required parameter: address');
  }

  log(`Getting account info for ${address}...`);

  try {
    const response = await client.request({
      command: 'account_info',
      account: address,
      ledger_index: 'validated',
    });

    const account = response.result.account_data;

    return {
      address: account.Account,
      balance: xrpl.dropsToXrp(account.Balance),
      balance_drops: account.Balance,
      sequence: account.Sequence,
      owner_count: account.OwnerCount,
      previous_txn_id: account.PreviousTxnID,
      previous_txn_lgr_seq: account.PreviousTxnLgrSeq,
      flags: account.Flags,
      ledger_index: response.result.ledger_index,
    };
  } catch (error) {
    if (error.data?.error === 'actNotFound') {
      return {
        address,
        funded: false,
        error: 'Account not found (not funded)',
      };
    }
    throw error;
  }
}

/**
 * Get server info
 */
async function getServerInfo(client, params, log) {
  log('Getting server info...');

  const response = await client.request({
    command: 'server_info',
  });

  const info = response.result.info;

  return {
    build_version: info.build_version,
    complete_ledgers: info.complete_ledgers,
    hostid: info.hostid,
    network_id: info.network_id,
    peers: info.peers,
    pubkey_node: info.pubkey_node,
    server_state: info.server_state,
    uptime: info.uptime,
    validated_ledger: info.validated_ledger
      ? {
          hash: info.validated_ledger.hash,
          seq: info.validated_ledger.seq,
          age: info.validated_ledger.age,
        }
      : null,
  };
}

/**
 * Main jade operations function - routes to specific operations
 */
async function jadeOps(config) {
  const { operation, network_url, network_id, params, verbose } = config;

  const log = verbose ? console.error.bind(console) : () => {};

  if (!operation) {
    throw new Error('Missing required field: operation');
  }

  if (!network_url) {
    throw new Error('Missing required field: network_url');
  }

  log(`Connecting to ${network_url}...`);

  const client = new xrpl.Client(network_url);
  await client.connect();

  try {
    let data;

    switch (operation) {
      case 'balance':
        data = await getBalance(client, params || {}, log);
        break;
      case 'send':
        data = await sendPayment(client, params || {}, log);
        break;
      case 'tx':
        data = await getTransaction(client, params || {}, log);
        break;
      case 'account_info':
        data = await getAccountInfo(client, params || {}, log);
        break;
      case 'server_info':
        data = await getServerInfo(client, params || {}, log);
        break;
      default:
        throw new Error(`Unknown operation: ${operation}`);
    }

    const result = {
      success: true,
      data,
    };

    console.log(JSON.stringify(result));
    return result;
  } catch (error) {
    const errorResult = {
      success: false,
      error: error.message,
      details: error.stack,
    };

    console.log(JSON.stringify(errorResult));
    process.exit(1);
  } finally {
    await client.disconnect();
  }
}

// CLI interface
if (require.main === module) {
  const args = process.argv.slice(2);

  if (args.length < 1) {
    console.error(`
Usage: node jade.js <config-json-path>

The config JSON file should contain:
{
  "operation": "balance|send|tx|account_info|server_info",
  "network_url": "wss://...",
  "network_id": 21465,
  "params": { ... },
  "verbose": false
}

Operations:
  balance       Get XRP balance for an address
  send          Send XRP to a destination
  tx            Get transaction details by hash
  account_info  Get full account information
  server_info   Get XRPL server information

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
    jadeOps(config);
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

module.exports = { jadeOps };
