package trigger

//
//import (
//	"arb-scaffold/constant"
//	"bytes"
//	"context"
//	"fmt"
//	"github.com/ethereum/go-ethereum"
//	"github.com/ethereum/go-ethereum/accounts/abi"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/core/types"
//	"github.com/ethereum/go-ethereum/ethclient/gethclient"
//	"github.com/ethereum/go-ethereum/rpc"
//	"github.com/sirupsen/logrus"
//	"math/big"
//	"strings"
//)
//

// 【WIP】: value监控需要用户提供abi和call data
//type ValueTrigger struct {
//	ctx     context.Context
//	chainId int64
//	addr    *common.Address
//	abi     abi.ABI
//	rpcs    map[string]string //map[ws]rpc, 收到新区块则去对应的节点查询value
//}
//
//func NewValueTrigger(ctx context.Context, chainId int64, wss []string, from *common.Address, to *common.Address, selector []byte) *ValueTrigger {
//	if len(wss) == 0 {
//		wss = constant.DefaultWss
//	}
//	return &ValueTrigger{
//		ctx:     ctx,
//		chainId: chainId,
//	}
//}
//
//// 交易需要连续订阅
//func (b *ValueTrigger) Listen() (chan *types.Transaction, chan error, error) {
//	signer := types.NewLondonSigner(big.NewInt(b.chainId))
//	txCh := make(chan *types.Transaction)
//	errCh := make(chan error)
//	var errs []string
//	healthWsCount := len(b.wss)
//
//	for _, ws := range b.wss {
//		sub, err := doSubscribeMempool(b.ctx, ws, txCh)
//		if err != nil {
//			logrus.Warnf("failed connect to rpc websocket: %s,  err: %s", ws, err)
//			errs = append(errs, ws+": "+err.Error())
//			healthWsCount -= 1
//			if healthWsCount == 0 {
//				return nil, nil, fmt.Errorf(strings.Join(errs, "\n"))
//			}
//			continue
//		}
//		go func(url string, s ethereum.Subscription) {
//			e := <-s.Err()
//			if e != nil {
//				logrus.Warnf("websocket err: %s,  err: %s", url, err)
//				errs = append(errs, url+": "+e.Error())
//				healthWsCount -= 1
//				if healthWsCount == 0 {
//					errCh <- fmt.Errorf(strings.Join(errs, "\n"))
//				}
//			}
//		}(ws, sub)
//	}
//
//	go func() {
//		for tx := range txCh {
//			sender, err := signer.Sender(tx)
//			if err != nil {
//				errCh <- err
//				continue
//			}
//			// TODO: 无法判断合约部署: to == nil
//			if (b.from != nil && sender == *b.from) || (b.to != nil && *tx.SetTo() == *b.to) || (b.selector != nil && len(tx.SetData()) >= 4 && bytes.Equal(tx.SetData()[:4], b.selector)) {
//				txCh <- tx
//			}
//		}
//	}()
//	return txCh, errCh, nil
//}
//
//func doSubscribeMempool(ctx context.Context, ws string, ch chan *types.Transaction) (*rpc.ClientSubscription, error) {
//	v2Cli, err := rpc.Dial(ws)
//	if err != nil {
//		return nil, err
//	}
//	subCli := gethclient.New(v2Cli)
//	return subCli.SubscribeFullPendingTransactions(ctx, ch)
//}
