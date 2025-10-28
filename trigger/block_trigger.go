package trigger

import (
	"context"
	"fmt"
	"time"

	"cs-projects-stable-pre-deposit/rpc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type BlockTrigger struct {
	targetBlock uint64
	rpcMgr      *rpc.RpcMgr
	ctx         context.Context
}

func NewBlockTrigger(ctx context.Context, blockNum uint64, rpcMgr *rpc.RpcMgr) *BlockTrigger {
	return &BlockTrigger{
		targetBlock: blockNum,
		rpcMgr:      rpcMgr,
		ctx:         ctx,
	}
}

func (b *BlockTrigger) doSubscribeEvent(ctx context.Context, ws string, ch chan *types.Header) (ethereum.Subscription, error) {
	cli, err := ethclient.Dial(ws)
	if err != nil {
		return nil, err
	}
	return cli.SubscribeNewHead(ctx, ch)
}

func (b *BlockTrigger) runWsForever(ws string, output chan *types.Header) {
	for {
		sub, err := b.doSubscribeEvent(b.ctx, ws, output)
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
func (b *BlockTrigger) Listen() (chan *types.Header, chan error, error) {
	ch := make(chan *types.Header)
	resultCh := make(chan *types.Header)
	errCh := make(chan error)

	for _, url := range b.rpcMgr.Wss() {
		go b.runWsForever(url, ch)
	}

	go func() {
		var lastBlock uint64 = 0
		for header := range ch {
			// 订阅所有区块
			currentBlock := header.Number.Uint64()
			if b.targetBlock == 0 {
				if currentBlock > lastBlock {
					resultCh <- header
					lastBlock = currentBlock
				}
				continue
			}
			if currentBlock == b.targetBlock {
				logrus.Infof("expected head comes: %d", currentBlock)
				resultCh <- header
				close(resultCh)
				close(errCh)
				return
			} else if currentBlock < b.targetBlock {
				logrus.Infof("new head: %d", currentBlock)
			} else {
				errCh <- fmt.Errorf("target block height too low: %d, current is %d", b.targetBlock, currentBlock)
				close(errCh)
				close(resultCh)
				return
			}
		}
	}()

	return resultCh, errCh, nil
}
