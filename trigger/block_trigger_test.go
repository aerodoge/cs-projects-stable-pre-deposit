package trigger

import (
	"context"
	"testing"

	"cs-projects-stable-pre-deposit/rpc"

	"github.com/ethereum/go-ethereum/core/types"
)

func TestBlockTrigger_Start(t *testing.T) {
	rpcMgr := rpc.NewRpcMgr()
	type fields struct {
		targetBlock uint64
		urls        []string
		ctx         context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		want    chan *types.Header
		want1   chan error
		wantErr bool
	}{
		{
			name: "subscribe",
			fields: fields{
				targetBlock: 17632025,
				urls:        []string{"wss://mainnet.infura.io/ws/v3/152c87a0611a4e88ac2b5a8e92a0bba9", "ws://172.16.46.192:8546"},
				ctx:         context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := BlockTrigger{
				targetBlock: tt.fields.targetBlock,
				rpcMgr:      rpcMgr,
				ctx:         tt.fields.ctx,
			}
			header, _, err := b.Listen()
			if err != nil {
				panic(err)
			}
			h := <-header
			if h.Number.Uint64() != tt.fields.targetBlock {
				t.Errorf("block num not match")
			}

		})
	}
}
