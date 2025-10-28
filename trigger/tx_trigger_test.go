package trigger

import (
	"context"
	"fmt"
	"testing"

	"cs-projects-stable-pre-deposit/rpc"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestTxTrigger_Listen(t *testing.T) {
	to := common.HexToAddress("0x3fC91A3afd70395Cd496C647d5a6CC9D4B2b7FAD") // uniswap
	selector := common.Hex2Bytes("3593564c")

	type fields struct {
		ctx      context.Context
		chainId  int64
		from     *common.Address
		to       *common.Address
		selector []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    chan *types.Transaction
		want1   chan error
		wantErr bool
	}{
		{
			name: "tx",
			fields: fields{
				ctx:      context.Background(),
				chainId:  1,
				from:     nil,
				to:       &to,
				selector: selector,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &TxTrigger{
				ctx:      tt.fields.ctx,
				from:     tt.fields.from,
				to:       tt.fields.to,
				rpcMgr:   rpc.NewRpcMgr(),
				selector: tt.fields.selector,
			}
			got, _, err := b.Listen()
			if (err != nil) != tt.wantErr {
				t.Errorf("Listen() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for transaction := range got {
				fmt.Println(transaction)
			}
		})
	}
}
