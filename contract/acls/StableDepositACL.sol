// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
import "./BaseReadOnlyACL.sol";
contract StableDepositACL is BaseReadOnlyACL {
    bytes32 public constant NAME = "StableDepositACL";
    uint256 public constant VERSION = 1;
    
    // Safe
    address constant SAFE = 0x28d9464a56129A75f5cEe38651C098A55feB3C11;
    address constant USDT = 0xdAC17F958D2ee523a2206206994597C13D831ec7;
    address constant STABLE_VAULT = 0x6503de9FE77d256d9d823f2D335Ce83EcE9E153f;
    address constant STABLE_TOKEN_RECEIVER = 0x22fd06cD176d0fa701f7aF31AD0E163D1a8Bae61;

    constructor() BaseOwnable(msg.sender) {}

    // approve USDT to StableVault
    function approve(address spender, uint256 value) public pure {
        require(
            (_txn().to == USDT && spender == STABLE_VAULT),
            "Approve: Invalid spender"
        );
    }
    // depost USDT to StableVault, limit receiver Must be SAFE
    function deposit(
        uint256 assets,
        address receiver
    ) public pure onlyContract(STABLE_VAULT) {
        require(receiver == SAFE, "Deposit: Invalid receiver");
    }
    // transfer STABLE_VAULT to STABLE_TOKEN_RECEIVER, limit to Must be STABLE_TOKEN_RECEIVER
    function transfer(
        address to,
        uint256 amount
    ) external pure onlyContract(STABLE_VAULT) {
        require(to == STABLE_TOKEN_RECEIVER, "invalid to account");
    }
}