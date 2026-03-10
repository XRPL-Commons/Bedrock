package templates

// Template represents a project template
type Template struct {
	Name        string
	Description string
	LibRS       string
}

// Available returns all available project templates
func Available() map[string]Template {
	return map[string]Template{
		"basic": {
			Name:        "basic",
			Description: "Basic contract with hello function",
			LibRS:       basicLib,
		},
		"token": {
			Name:        "token",
			Description: "Fungible token contract with mint/transfer/balance",
			LibRS:       tokenLib,
		},
		"nft": {
			Name:        "nft",
			Description: "NFT contract with mint/transfer/owner tracking",
			LibRS:       nftLib,
		},
		"escrow": {
			Name:        "escrow",
			Description: "Escrow contract with create/release/cancel",
			LibRS:       escrowLib,
		},
		"counter": {
			Name:        "counter",
			Description: "Simple counter contract for learning",
			LibRS:       counterLib,
		},
	}
}

const basicLib = `#![cfg_attr(target_arch = "wasm32", no_std)]

#[cfg(not(target_arch = "wasm32"))]
extern crate std;

use xrpl_wasm_macros::wasm_export;
use xrpl_wasm_std::host::trace::trace;

/// @xrpl-function hello
#[wasm_export]
fn hello() -> i32 {
    let _ = trace("Hello from XRPL Smart Contract!");
    0
}
`

const tokenLib = `#![cfg_attr(target_arch = "wasm32", no_std)]

#[cfg(not(target_arch = "wasm32"))]
extern crate std;

use xrpl_wasm_macros::wasm_export;
use xrpl_wasm_std::host::trace::trace;

/// @xrpl-function mint
/// @param amount UINT64 - Amount to mint
/// @return UINT64 - New total supply
#[wasm_export]
fn mint(amount: u64) -> u64 {
    let _ = trace("Minting tokens");
    amount
}

/// @xrpl-function transfer
/// @param to ACCOUNT - Recipient address
/// @param amount UINT64 - Amount to transfer
/// @return UINT64 - Success code (0 = success)
#[wasm_export]
fn transfer(to: &[u8], amount: u64) -> u64 {
    let _ = trace("Transferring tokens");
    let _ = (to, amount);
    0
}

/// @xrpl-function balance
/// @param account ACCOUNT - Account to check
/// @return UINT64 - Token balance
#[wasm_export]
fn balance(account: &[u8]) -> u64 {
    let _ = trace("Checking balance");
    let _ = account;
    0
}
`

const nftLib = `#![cfg_attr(target_arch = "wasm32", no_std)]

#[cfg(not(target_arch = "wasm32"))]
extern crate std;

use xrpl_wasm_macros::wasm_export;
use xrpl_wasm_std::host::trace::trace;

/// @xrpl-function mint_nft
/// @param to ACCOUNT - Owner address
/// @return UINT64 - Token ID
#[wasm_export]
fn mint_nft(to: &[u8]) -> u64 {
    let _ = trace("Minting NFT");
    let _ = to;
    1
}

/// @xrpl-function transfer_nft
/// @param token_id UINT64 - Token ID to transfer
/// @param to ACCOUNT - New owner address
/// @return UINT64 - Success code
#[wasm_export]
fn transfer_nft(token_id: u64, to: &[u8]) -> u64 {
    let _ = trace("Transferring NFT");
    let _ = (token_id, to);
    0
}

/// @xrpl-function owner_of
/// @param token_id UINT64 - Token ID to check
/// @return UINT64 - Owner info
#[wasm_export]
fn owner_of(token_id: u64) -> u64 {
    let _ = trace("Checking NFT owner");
    let _ = token_id;
    0
}
`

const escrowLib = `#![cfg_attr(target_arch = "wasm32", no_std)]

#[cfg(not(target_arch = "wasm32"))]
extern crate std;

use xrpl_wasm_macros::wasm_export;
use xrpl_wasm_std::host::trace::trace;

/// @xrpl-function create_escrow
/// @param beneficiary ACCOUNT - Escrow beneficiary
/// @param amount UINT64 - Escrow amount
/// @return UINT64 - Escrow ID
#[wasm_export]
fn create_escrow(beneficiary: &[u8], amount: u64) -> u64 {
    let _ = trace("Creating escrow");
    let _ = (beneficiary, amount);
    1
}

/// @xrpl-function release_escrow
/// @param escrow_id UINT64 - Escrow to release
/// @return UINT64 - Success code
#[wasm_export]
fn release_escrow(escrow_id: u64) -> u64 {
    let _ = trace("Releasing escrow");
    let _ = escrow_id;
    0
}

/// @xrpl-function cancel_escrow
/// @param escrow_id UINT64 - Escrow to cancel
/// @return UINT64 - Success code
#[wasm_export]
fn cancel_escrow(escrow_id: u64) -> u64 {
    let _ = trace("Cancelling escrow");
    let _ = escrow_id;
    0
}
`

const counterLib = `#![cfg_attr(target_arch = "wasm32", no_std)]

#[cfg(not(target_arch = "wasm32"))]
extern crate std;

use xrpl_wasm_macros::wasm_export;
use xrpl_wasm_std::host::trace::trace;

/// @xrpl-function increment
/// @return UINT64 - New counter value
#[wasm_export]
fn increment() -> u64 {
    let _ = trace("Incrementing counter");
    1
}

/// @xrpl-function decrement
/// @return UINT64 - New counter value
#[wasm_export]
fn decrement() -> u64 {
    let _ = trace("Decrementing counter");
    0
}

/// @xrpl-function get_count
/// @return UINT64 - Current counter value
#[wasm_export]
fn get_count() -> u64 {
    let _ = trace("Getting counter value");
    0
}
`
