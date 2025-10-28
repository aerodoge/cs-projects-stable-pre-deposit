package main

import (
	"context"
	"cs-projects-stable-pre-deposit/contract"
	"cs-projects-stable-pre-deposit/rpc"
	"cs-projects-stable-pre-deposit/trigger"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	safe := "0x28d9464a56129A75f5cEe38651C098A55feB3C11"
	argus := "0x13623ee9047162a2658c7406a5ac2093c5f75541"
	localRPC := "http://localhost:8545"
	privateKey := ""
	ctStable := common.HexToAddress("0x6503de9FE77d256d9d823f2D335Ce83EcE9E153f")
	topics0 := "0xb2ad710f2954a5376267a683f9ece9ec46ee7dfb47075163379904ee941df8da"

	delegate := contract.NewDelegate(localRPC, privateKey, safe, argus)
	//amount, _ := new(big.Int).SetString("100000000000000000000", 10)
	//// 执行批量交易
	//delegate.BatchExecute(amount)
	//minDeposit, _ := new(big.Int).SetString("1000000", 10)
	//maxDeposit, _ := new(big.Int).SetString("9000000000000000", 10)
	//contract.BuildSetDepositLimits(minDeposit, maxDeposit)
	// 启动事件监听（在单独的goroutine中运行）
	fmt.Printf("启动事件监听，目标合约: %s, topic0: %s\n", ctStable, topics0)
	//a, _ := new(big.Int).SetString("10000000000000", 10)
	//delegate.IssueAndTransferUSDT(common.HexToAddress(safe), a)
	rpcMgr := rpc.NewRpcMgr()
	eventTrigger := trigger.NewEventTrigger(context.Background(), rpcMgr, ctStable, topics0, false)
	logCh, errCh, err := eventTrigger.Listen()
	if err != nil {
		fmt.Printf("linsten err: %v", err)
		return
	}
	// 设置信号处理
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fmt.Println("程序正在运行，监听事件中... 按 Ctrl+C 退出")

	for {
		select {
		case l := <-logCh:
			if l.Topics[0].Hex() != topics0 {
				fmt.Printf("%v\n", l.Topics[0].Hex())
				continue
			}
			fmt.Printf("检测到目标事件!\n")
			fmt.Printf("交易哈希: %s\n", l.TxHash.Hex())
			fmt.Printf("区块号: %d\n", l.BlockNumber)
			fmt.Printf("Topic0: %s\n", l.Topics[0].Hex())

			// 处理事件
			amount, _ := new(big.Int).SetString("100000000000", 10)
			delegate.BatchExecute(amount)
			//delegate.StepByStep(amount)
		case e := <-errCh:
			fmt.Printf("err: %v", e)

		case <-c:
			fmt.Println("收到退出信号，程序结束")
			return
		}
	}
}
