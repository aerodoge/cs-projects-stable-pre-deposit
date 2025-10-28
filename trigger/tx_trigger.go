package trigger

import (
	"bytes"
	"context"
	"time"

	"cs-projects-stable-pre-deposit/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
)

// 内存池交易监控
type TxTrigger struct {
	ctx      context.Context
	from     *common.Address
	to       *common.Address
	rpcMgr   *rpc.RpcMgr
	selector []byte
	onceOnly bool
}

func NewTxTrigger(ctx context.Context, rpcMgr *rpc.RpcMgr, from *common.Address, to *common.Address, selector []byte, onceOnly bool) *TxTrigger {

	return &TxTrigger{
		ctx:      ctx,
		rpcMgr:   rpcMgr,
		from:     from,
		to:       to,
		selector: selector,
		onceOnly: onceOnly,
	}
}

// 交易需要连续订阅
func (b *TxTrigger) Listen() (chan *types.Transaction, chan error, error) {
	signer := types.NewLondonSigner(b.rpcMgr.ChainId())
	txCh := make(chan *types.Transaction)
	errCh := make(chan error)

	for _, ws := range b.rpcMgr.Wss() {
		go b.runWsForever(ws, txCh)
	}

	ch := make(chan *types.Transaction)
	go func() {
		for tx := range txCh {
			sender, err := signer.Sender(tx)
			if err != nil {
				errCh <- err
				continue
			}
			// TODO: 无法判断合约部署: to == nil
			if (b.from != nil && sender == *b.from) || (b.to != nil && tx.To() != nil && *tx.To() == *b.to) || (b.selector != nil && len(tx.Data()) >= 4 && bytes.Equal(tx.Data()[:4], b.selector)) {
				ch <- tx
				if b.onceOnly {
					break
				}
			}
		}
	}()
	return ch, errCh, nil
}

func doSubscribeMempool(ctx context.Context, ws string, ch chan *types.Transaction) (*ethrpc.ClientSubscription, error) {
	v2Cli, err := ethrpc.Dial(ws)
	if err != nil {
		return nil, err
	}
	subCli := gethclient.New(v2Cli)
	return subCli.SubscribeFullPendingTransactions(ctx, ch)
}

func (b *TxTrigger) runWsForever(ws string, output chan *types.Transaction) {
	for {
		sub, err := doSubscribeMempool(b.ctx, ws, output)
		if err != nil {
			logrus.Warnf("failed connect to rpc websocket: %s,  err: %s", ws, err)
			time.Sleep(1 * time.Second)
			continue
		}
		err = <-sub.Err()
		if err != nil {
			logrus.Warnf("err on websocket: %s,  err: %s", ws, err)
			time.Sleep(1 * time.Second)
			continue
		} else {
			logrus.Infof("websocket close: %s", ws)
			break
		}
	}
}
