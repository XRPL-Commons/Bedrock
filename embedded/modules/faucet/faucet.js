#!/usr/bin/env node

/**
 * XRPL Faucet Request Module
 *
 * This module handles requesting funds from XRPL faucets.
 *
 * Usage: node faucet.js <config-json-path>
 *
 * Config JSON format:
 * {
 *   "faucet_url": "https://faucet.altnet.rippletest.net/accounts",
 *   "wallet_seed": "sXXX..." (optional),
 *   "wallet_address": "rXXX..." (optional),
 *   "network_url": "wss://..." (optional, for balance check),
 *   "verbose": true (optional)
 * }
 *
 * Output JSON format:
 * {
 *   "success": true,
 *   "data": {
 *     "txHash": "...",
 *     "walletAddress": "...",
 *     "walletSeed": "...",
 *     "balance": "1000",
 *     "faucetAmount": "1000"
 *   }
 * }
 */

const xrpl = require('@transia/xrpl');
const fs = require('fs');
const https = require('https');
const http = require('http');

/**
 * Request funds from XRPL faucet
 */
async function requestFaucet(config) {
  const { faucet_url, wallet_seed, wallet_address, network_url, verbose } =
    config;

  const log = verbose ? console.error.bind(console) : () => {};

  log('Requesting funds from faucet...\n');

  try {
    let wallet;
    let address;

    // Determine wallet/address
    if (wallet_seed) {
      wallet = xrpl.Wallet.fromSeed(seed, { algorithm: xrpl.ECDSA.secp256k1 });
      address = wallet.address;
      log('Using provided wallet seed');
      log('  Address:', address);
    } else if (wallet_address) {
      address = wallet_address;
      log('Using provided address:', address);
    } else {
      wallet = xrpl.Wallet.generate(xrpl.ECDSA.secp256k1);
      address = wallet.address;
      log('Generated new wallet');
      log('  Address:', address);
      log('  Seed:', wallet.seed);
    }

    log('\nRequesting funds from faucet...');
    log('  Faucet URL:', faucet_url);

    // Make faucet request
    const faucetResult = await makeFaucetRequest(faucet_url, address);
    log('✓ Faucet request successful');

    // Get balance if network URL provided
    let balance = null;
    if (network_url) {
      log('\nChecking balance...');
      const client = new xrpl.Client(network_url);
      await client.connect();
      balance = await client.getXrpBalance(address);
      await client.disconnect();
      log('  Balance:', balance, 'XRP');
    }

    // Output result
    const result = {
      success: true,
      data: {
        txHash: faucetResult.txHash,
        walletAddress: address,
        walletSeed: wallet ? wallet.seed : '',
        balance: balance ? String(balance) : '',
        faucetAmount: faucetResult.amount ? String(faucetResult.amount) : '',
      },
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
  }
}

/**
 * Make HTTP request to faucet
 */
function makeFaucetRequest(faucetUrl, address) {
  return new Promise((resolve, reject) => {
    const url = new URL(faucetUrl);
    const protocol = url.protocol === 'https:' ? https : http;

    const postData = JSON.stringify({
      destination: address,
    });

    const options = {
      hostname: url.hostname,
      port: url.port,
      path: url.pathname + url.search,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Content-Length': Buffer.byteLength(postData),
      },
    };

    const req = protocol.request(options, (res) => {
      let data = '';

      res.on('data', (chunk) => {
        data += chunk;
      });

      res.on('end', () => {
        if (res.statusCode === 200 || res.statusCode === 201) {
          try {
            const response = JSON.parse(data);
            // Different faucets return different formats
            // Try to extract relevant info
            const txHash =
              response.hash || response.txHash || response.tx_hash || 'unknown';
            const amount = response.amount || response.balance || 'unknown';

            resolve({
              txHash: String(txHash),
              amount: String(amount),
            });
          } catch (err) {
            reject(
              new Error(`Failed to parse faucet response: ${err.message}`)
            );
          }
        } else {
          reject(
            new Error(
              `Faucet request failed with status ${res.statusCode}: ${data}`
            )
          );
        }
      });
    });

    req.on('error', (err) => {
      reject(new Error(`Faucet request failed: ${err.message}`));
    });

    req.write(postData);
    req.end();
  });
}

// CLI interface
if (require.main === module) {
  const args = process.argv.slice(2);

  if (args.length < 1) {
    console.error(`
Usage: node faucet.js <config-json-path>

The config JSON file should contain:
{
  "faucet_url": "https://faucet.altnet.rippletest.net/accounts",
  "wallet_seed": "sXXX..." (optional),
  "wallet_address": "rXXX..." (optional),
  "network_url": "wss://..." (optional),
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
    requestFaucet(config);
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

module.exports = { requestFaucet };
