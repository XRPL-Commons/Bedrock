#!/usr/bin/env node

const xrpl = require('@transia/xrpl');
const { readFileSync } = require('fs');

async function main() {
  try {
    // Read configuration from stdin or argument
    const input = process.argv[2] || readFileSync(0, 'utf8');
    const config = JSON.parse(input);

    let result = {};

    switch (config.action) {
      case 'derive_address':
        result = await deriveAddress(config);
        break;
      case 'generate_wallet':
        result = await generateWallet(config);
        break;
      case 'validate_seed':
        result = await validateSeed(config);
        break;
      default:
        throw new Error(`Unknown action: ${config.action}`);
    }

    process.stdout.write(JSON.stringify({
      success: true,
      ...result
    }));

  } catch (error) {
    process.stdout.write(JSON.stringify({
      success: false,
      error: error.message,
      details: error.toString()
    }));
    process.exit(1);
  }
}

async function deriveAddress(config) {
  const { seed, algorithm } = config;
  
  if (!seed) {
    throw new Error('Seed is required');
  }

  // Create wallet from seed using specified algorithm (default: secp256k1)
  const algo = algorithm === 'ed25519' ? undefined : xrpl.ECDSA.secp256k1;
  const wallet = algo ? xrpl.Wallet.fromSeed(seed, { algorithm: algo }) : xrpl.Wallet.fromSeed(seed);
  
  return {
    address: wallet.address,
    public_key: wallet.publicKey
  };
}

async function generateWallet(config) {
  const { algorithm } = config;
  
  // Generate new wallet with specified algorithm (default: secp256k1)
  const algo = algorithm === 'ed25519' ? undefined : xrpl.ECDSA.secp256k1;
  const wallet = algo ? xrpl.Wallet.generate(algo) : xrpl.Wallet.generate();
  
  return {
    seed: wallet.seed,
    address: wallet.address,
    public_key: wallet.publicKey
  };
}

async function validateSeed(config) {
  const { seed, algorithm } = config;
  
  if (!seed) {
    throw new Error('Seed is required');
  }

  try {
    // Try to create wallet from seed - if it fails, seed is invalid
    const algo = algorithm === 'ed25519' ? undefined : xrpl.ECDSA.secp256k1;
    const wallet = algo ? xrpl.Wallet.fromSeed(seed, { algorithm: algo }) : xrpl.Wallet.fromSeed(seed);
    return {
      valid: true,
      address: wallet.address
    };
  } catch (error) {
    return {
      valid: false,
      error: error.message
    };
  }
}

if (require.main === module) {
  main();
}

module.exports = { deriveAddress, generateWallet, validateSeed };