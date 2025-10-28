// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
import "./BaseReadOnlyACL.sol";
contract SparkSavingWETHACL is BaseReadOnlyACL {
    bytes32 public constant NAME = "SparkSavingWETHACL";
    uint256 public constant VERSION = 1;
    address public constant SAFE = 0x969b37A287bBFb4080E6cb293fB6E21995fd1f83;
    address public constant SparkSavingWETHAddress = 0xfE6eb3b609a7C8352A241f7F3A21CEA4e9209B8f;
    address public constant ERC20WETH = 0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2;

    constructor() BaseOwnable(msg.sender) {}

    //1. SAFE APPROVE WETH For SparkSavingWETHAddress
    function approve(address spender, uint256 value) external view {
        require(
            (_txn().to == ERC20WETH && spender == SparkSavingWETHAddress),
            "Approve: Invalid approve");
    }
    //WETH deposit to spark saving, receiver MUST be SAFE
    function deposit(uint256 assets, address receiver, uint16 referral) public view onlyContract(SparkSavingWETHAddress) {
        require(receiver == SAFE, "invalid receiver");
    }
    //WETH redeem from spark, receiver MUST be SAFE
    function redeem(uint256 shares, address receiver, address owner) public view onlyContract(SparkSavingWETHAddress) {
        require(receiver == SAFE, "invalid receiver");
        require(owner == SAFE, "invalid owner");
    }
    //WETH - deposit
    function deposit()  public view onlyContract(ERC20WETH)  {}
    //WETH - withdraw
    function withdraw(uint wad) public view onlyContract(ERC20WETH) {}
}