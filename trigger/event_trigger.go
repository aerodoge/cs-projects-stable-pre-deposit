package trigger

import (
	"context"
	"fmt"
	"time"

	"cs-projects-stable-pre-deposit/rpc"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

// 事件监控
type EventTrigger struct {
	ctx      context.Context
	rpcMgr   *rpc.RpcMgr
	contract common.Address
	topic0   string
	onceOnly bool // 触发一次后退出
}

func NewEventTrigger(ctx context.Context, rpcMgr *rpc.RpcMgr, contract common.Address, topic0 string, onceOnly bool) *EventTrigger {
	return &EventTrigger{
		ctx:      ctx,
		rpcMgr:   rpcMgr,
		contract: contract,
		topic0:   topic0,
		onceOnly: onceOnly,
	}
}

func (b *EventTrigger) runWsForever(ws string, output chan types.Log) {
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

// 交易需要连续订阅
func (b *EventTrigger) Listen() (chan types.Log, chan error, error) {
	eventCh := make(chan types.Log)
	errCh := make(chan error)
	for _, ws := range b.rpcMgr.Wss() {
		go b.runWsForever(ws, eventCh)
	}
	ch := make(chan types.Log)
	go func() {
		for event := range eventCh {
			fmt.Println("event: ", event)
			if b.topic0 == "" || (len(event.Topics) > 0 && b.topic0 == event.Topics[0].Hex()) {
				ch <- event
				if b.onceOnly {
					break
				}
			}
		}
	}()

	return ch, errCh, nil
}

func (b *EventTrigger) doSubscribeEvent(ctx context.Context, ws string, ch chan types.Log) (ethereum.Subscription, error) {
	cli, err := ethclient.Dial(ws)
	if err != nil {
		return nil, err
	}
	return cli.SubscribeFilterLogs(ctx, ethereum.FilterQuery{
		Addresses: []common.Address{b.contract},
		Topics:    [][]common.Hash{{common.HexToHash(b.topic0)}},
	}, ch)
}
