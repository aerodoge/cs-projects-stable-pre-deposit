package trigger

import (
	"context"
	"fmt"
	"testing"

	"cs-projects-stable-pre-deposit/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestEventTrigger_Listen(t *testing.T) {

	type fields struct {
		ctx      context.Context
		contract common.Address
		topic0   string
	}
	tests := []struct {
		name    string
		fields  fields
		want    chan types.Log
		want1   chan error
		wantErr bool
	}{
		{
			name: "event",
			fields: fields{
				ctx:      context.Background(),
				contract: common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
				topic0:   "0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &EventTrigger{
				ctx:      tt.fields.ctx,
				rpcMgr:   rpc.NewRpcMgr(),
				contract: tt.fields.contract,
				topic0:   tt.fields.topic0,
			}
			got, _, err := b.Listen()
			if (err != nil) != tt.wantErr {
				t.Errorf("Listen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for log := range got {
				fmt.Println(log.Data)
			}
		})
	}
}
