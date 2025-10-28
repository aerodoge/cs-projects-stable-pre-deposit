package contract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var ctStableAddress = common.HexToAddress("0x6503de9FE77d256d9d823f2D335Ce83EcE9E153f")
var usdtAddress = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")

type SafeCallData struct {
	Flag  *big.Int       `json:"flag"`  // uint256
	To    common.Address `json:"to"`    // address
	Value *big.Int       `json:"value"` // uint256
	Data  []byte         `json:"data"`  // bytes
	Hint  []byte         `json:"hint"`  // bytes
	Extra []byte         `json:"extra"` // bytes
}

type Delegate struct {
	client     *ethclient.Client
	privateKey *ecdsa.PrivateKey
	address    common.Address
	safe       common.Address
	argus      common.Address
}

func NewDelegate(rpcURL, privateKey, safe, argus string) *Delegate {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	address := crypto.PubkeyToAddress(pk.PublicKey)
	return &Delegate{
		client:     client,
		privateKey: pk,
		address:    address,
		safe:       common.HexToAddress(safe),
		argus:      common.HexToAddress(argus),
	}
}
func (d *Delegate) IssueAndTransferUSDT(to common.Address, amount *big.Int) {
	issueData, err := IssueUSDT(usdtAbi, to, amount)
	if err != nil {
		fmt.Printf("issue USDT error: %v\n", err)
		return
	}
	fmt.Println("issueData", hexutil.Encode(issueData))
	//txHash, err := d.sendTransaction(usdtAddress, issueData)
	//if err != nil {
	//	fmt.Printf("issue USDT2 error: %v\n", err)
	//	return
	//}
	//err = d.waitForConfirmation(txHash)
	//if err != nil {
	//	log.Printf("wait confirmation error: %v\n", err)
	//	return
	//}

	transferData, err := TransferUSDT(usdtAbi, to, amount)
	if err != nil {
		fmt.Printf("transfer USDT error: %v\n", err)
		return
	}
	fmt.Println("transferData", hexutil.Encode(transferData))
	//txHash, err = d.sendTransaction(usdtAddress, transferData)
	//if err != nil {
	//	fmt.Printf("transfer USDT2 error: %v\n", err)
	//	return
	//}
	//err = d.waitForConfirmation(txHash)
	//if err != nil {
	//	log.Printf("wait confirmation error: %v\n", err)
	//	return
	//}
}

func (d *Delegate) sendTransaction(to common.Address, callData []byte) (*common.Hash, error) {
	// 构建Safe交易
	GasLimit := 6000000
	// 获取nonce
	nonce, err := d.client.PendingNonceAt(context.Background(), d.address)
	if err != nil {
		return nil, fmt.Errorf("获取nonce失败: %v", err)
	}

	// 获取gas价格
	gasPrice, err := d.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取gas价格失败: %v", err)
	}

	fmt.Printf("Gas price: %v, gas limit: %v\n", gasPrice, GasLimit)
	// 创建交易
	tx := types.NewTransaction(
		nonce,
		to,
		big.NewInt(0),
		uint64(GasLimit),
		gasPrice,
		callData,
	)

	chainID, err := d.client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取chain ID失败: %v", err)
	}
	fmt.Println("ChainID:", chainID)
	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), d.privateKey)
	if err != nil {
		return nil, fmt.Errorf("签名交易失败: %v", err)
	}

	// 发送交易
	err = d.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("发送交易失败: %v", err)
	}

	hash := signedTx.Hash()
	return &hash, nil
}

func (d *Delegate) StepByStep(amount *big.Int) {
	// approve
	approveData, err := BuildApproveData(usdtAbi, ctStableAddress, amount)
	if err != nil {
		log.Printf("build approve Error: %v\n", err)
		return
	}
	txHash, err := d.executeSafeTransaction(usdtAddress, big.NewInt(0), approveData)
	err = d.waitForConfirmation(txHash)
	if err != nil {
		log.Printf("wait confirmation error: %v\n", err)
		return
	}

	// deposit
	depositData, err := BuildDepositData(ctStableAbi, d.safe, amount)
	if err != nil {
		log.Printf("build deposit Error: %v\n", err)
		return
	}
	txHash, err = d.executeSafeTransaction(ctStableAddress, big.NewInt(0), depositData)
	err = d.waitForConfirmation(txHash)
	if err != nil {
		log.Printf("wait confirmation error: %v\n", err)
		return
	}
}

// 通过Safe执行交易x
func (d *Delegate) executeSafeTransaction(to common.Address, value *big.Int, data []byte) (*common.Hash, error) {
	// 构建Safe交易
	GasLimit := 6000000
	// 获取nonce
	nonce, err := d.client.PendingNonceAt(context.Background(), d.address)
	if err != nil {
		return nil, fmt.Errorf("获取nonce失败: %v", err)
	}

	// 获取gas价格
	gasPrice, err := d.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取gas价格失败: %v", err)
	}

	// 构建Safe执行交易的calldata
	safeExecData := d.buildSafeExecTransactionData(to, value, data)
	fmt.Printf("Gas price: %v, gas limit: %v\n", gasPrice, GasLimit)
	// 创建交易
	tx := types.NewTransaction(
		nonce,
		d.argus,
		big.NewInt(0),
		uint64(GasLimit),
		gasPrice,
		safeExecData,
	)

	chainID, err := d.client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取chain ID失败: %v", err)
	}
	fmt.Println("ChainID:", chainID)
	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), d.privateKey)
	if err != nil {
		return nil, fmt.Errorf("签名交易失败: %v", err)
	}

	// 发送交易
	err = d.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("发送交易失败: %v", err)
	}

	hash := signedTx.Hash()
	return &hash, nil
}

func (d *Delegate) buildSafeExecTransactionData(to common.Address, value *big.Int, data []byte) []byte {
	contractABI, err := abi.JSON(strings.NewReader(safeAbi))
	if err != nil {
		return nil
	}
	callData := SafeCallData{
		Flag:  big.NewInt(0),
		To:    to,
		Value: value,
		Data:  data,
		Hint:  []byte{},
		Extra: []byte{},
	}

	// 使用ABI编码
	encodedData, err := contractABI.Pack("execTransaction", callData)
	if err != nil {
		log.Fatal("build safe exec transaction err: ", err)
		return nil
	}
	return encodedData
}

func (d *Delegate) BatchExecute(amount *big.Int) {
	approveData, err := BuildApproveData(usdtAbi, ctStableAddress, amount)
	if err != nil {
		log.Printf("build approve Error: %v\n", err)
		return
	}
	depositData, err := BuildDepositData(ctStableAbi, d.safe, amount)
	if err != nil {
		log.Printf("build deposit eposit Error: %v\n", err)
		return
	}
	var addrs []common.Address
	var values []*big.Int
	var datas [][]byte
	value1 := big.NewInt(0)
	addrs = append(addrs, usdtAddress)
	values = append(values, value1)
	datas = append(datas, approveData)

	value2 := big.NewInt(0)
	addrs = append(addrs, ctStableAddress)
	values = append(values, value2)
	datas = append(datas, depositData)

	txHash, err := d.batchExecuteSafeTransactions(addrs, values, datas)
	if err != nil {
		log.Printf("execute safe error: %v\n", err)
		return
	}

	err = d.waitForConfirmation(txHash)
	if err != nil {
		log.Printf("wait confirmation error: %v\n", err)
		return
	}
}

func (d *Delegate) batchExecuteSafeTransactions(addrs []common.Address, values []*big.Int, datas [][]byte) (*common.Hash, error) {
	if len(addrs) != len(datas) {
		return nil, fmt.Errorf("len(addrs) != len(datas)")
	}
	callDatas := []SafeCallData{}
	for i := range addrs {
		callData := SafeCallData{
			Flag:  big.NewInt(0),
			To:    addrs[i],
			Value: values[i],
			Data:  datas[i],
			Hint:  []byte{},
			Extra: []byte{},
		}
		callDatas = append(callDatas, callData)
	}

	// 使用Safe ABI构建execTransactions调用数据
	contractABI, err := abi.JSON(strings.NewReader(safeAbi))
	if err != nil {
		return nil, fmt.Errorf("解析Safe ABI失败: %v", err)
	}

	// 使用ABI编码execTransactions调用
	batchCallData, err := contractABI.Pack("execTransactions", callDatas)
	if err != nil {
		return nil, fmt.Errorf("编码批量交易数据失败: %v", err)
	}

	// 构建Safe交易
	GasLimit := 6000000
	// 获取nonce
	nonce, err := d.client.PendingNonceAt(context.Background(), d.address)
	if err != nil {
		return nil, fmt.Errorf("获取nonce失败: %v", err)
	}

	// 获取gas价格
	gasPrice, err := d.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取gas价格失败: %v", err)
	}

	tx := types.NewTransaction(
		nonce,
		d.argus,
		big.NewInt(0),
		uint64(GasLimit),
		gasPrice,
		batchCallData,
	)

	chainID, err := d.client.NetworkID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取chain ID失败: %v", err)
	}
	fmt.Println("ChainID:", chainID)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), d.privateKey)
	if err != nil {
		return nil, fmt.Errorf("签名交易失败: %v", err)
	}

	// 发送交易
	err = d.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("发送交易失败: %v", err)
	}

	hash := signedTx.Hash()

	return &hash, nil
}

func (d *Delegate) waitForConfirmation(txHash *common.Hash) error {
	fmt.Printf("等待交易确认: %s\n", txHash.Hex())

	for i := 0; i < 60; i++ { // 最多等待5分钟
		receipt, err := d.client.TransactionReceipt(context.Background(), *txHash)
		if err == nil {
			if receipt.Status == types.ReceiptStatusSuccessful {
				fmt.Printf("交易确认成功! Gas使用: %d\n", receipt.GasUsed)
				return nil
			} else {
				return fmt.Errorf("交易执行失败\n")
			}
		}
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("交易确认超时")
}

// ListenToEvents 监听ctStable合约的事件
func (d *Delegate) ListenToEvents(ctStableAddress common.Address, targetTopic0 string) {
	// 创建过滤器查询
	query := ethereum.FilterQuery{
		Addresses: []common.Address{ctStableAddress},
		Topics: [][]common.Hash{
			{common.HexToHash(targetTopic0)}, // 指定的topic0
		},
	}

	// 订阅事件日志
	logs := make(chan types.Log)
	sub, err := d.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Printf("订阅事件失败: %v", err)
		return
	}
	defer sub.Unsubscribe()

	fmt.Printf("开始监听合约 %s 的事件，目标topic0: %s\n", ctStableAddress.Hex(), targetTopic0)

	for {
		select {
		case err := <-sub.Err():
			log.Printf("事件订阅错误: %v", err)
			return
		case vLog := <-logs:
			d.handleEvent(vLog, targetTopic0)
		}
	}
}

// handleEvent 处理接收到的事件
func (d *Delegate) handleEvent(vLog types.Log, targetTopic0 string) {
	// 检查topic0是否匹配
	if len(vLog.Topics) > 0 && vLog.Topics[0].Hex() == targetTopic0 {
		fmt.Printf("检测到目标事件!\n")
		fmt.Printf("交易哈希: %s\n", vLog.TxHash.Hex())
		fmt.Printf("区块号: %d\n", vLog.BlockNumber)
		fmt.Printf("Topic0: %s\n", vLog.Topics[0].Hex())

		// 在这里添加你想要执行的逻辑
		d.onTargetEventDetected(vLog)
	}
}

// onTargetEventDetected 当检测到目标事件时执行的操作
func (d *Delegate) onTargetEventDetected(vLog types.Log) {
	fmt.Printf("执行目标事件处理逻辑...\n")

	// 解析事件数据（根据具体的事件ABI结构）
	// 这里可以添加具体的业务逻辑
	// 比如：
	// 1. 解析事件参数
	// 2. 执行相应的合约调用
	// 3. 记录日志等

	fmt.Printf("事件处理完成\n")
}
