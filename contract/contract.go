package contract

import (
	_ "embed"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//go:embed abis/ctStableUSDT.abi.json
var ctStableAbi string

//go:embed abis/safe.abi.json
var safeAbi string

//go:embed abis/usdt.abi.json
var usdtAbi string

func IssueUSDT(contractABI string, to common.Address, amount *big.Int) ([]byte, error) {
	data, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}
	encodedData, err := data.Pack("issue", amount)
	if err != nil {
		return nil, err
	}
	return encodedData, nil
}

func TransferUSDT(contractABI string, to common.Address, amount *big.Int) ([]byte, error) {
	data, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}
	encodedData, err := data.Pack("transfer", to, amount)
	if err != nil {
		return nil, err
	}
	return encodedData, nil
}

func BuildSetDepositLimits(minDeposit, maxDeposit *big.Int) ([]byte, error) {
	data, err := abi.JSON(strings.NewReader(ctStableAbi))
	if err != nil {
		return nil, err
	}
	encodedData, err := data.Pack("setDepositLimits", minDeposit, maxDeposit)
	if err != nil {
		return nil, err
	}
	fmt.Println("setDepositLimits: ", hexutil.Encode(encodedData))
	return encodedData, nil
}

func BuildApproveData(contractABI string, spender common.Address, amount *big.Int) ([]byte, error) {
	data, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}

	encodedData, err := data.Pack("approve", spender, amount)
	if err != nil {
		return nil, err
	}

	return encodedData, nil
}

func BuildDepositData(contractABI string, receiver common.Address, amount *big.Int) ([]byte, error) {
	data, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		return nil, err
	}

	encodedData, err := data.Pack("deposit", amount, receiver)
	if err != nil {
		return nil, err
	}

	return encodedData, nil
}
