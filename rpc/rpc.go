package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

type RpcMgr struct {
	rpcs    map[string]*RpcConfig
	clis    []*ethclient.Client
	chainId *big.Int
}

func NewRpcMgr() *RpcMgr {
	//env := os.Getenv("FS_ENV")
	//// env := "PROD"
	//if env == "PROD" {
	//	return NewRpcMgrWithConfig(defaultRpc)
	//} else if env == "DEV" {
	//	return NewRpcMgrWithConfig(map[string]*RpcConfig{
	//		"local": {
	//			Ws:  "ws://localhost:8545",
	//			Rpc: "http://localhost:8545",
	//		},
	//	})
	//} else {
	//	logrus.Fatal("FS_ENV must be set to PROD or DEV !")
	//}

	return NewRpcMgrWithConfig(map[string]*RpcConfig{
		"local": {
			Ws:  "ws://localhost:8545",
			Rpc: "http://localhost:8545",
		},
	})
}

func NewRpcMgrWithConfig(rpcConfig map[string]*RpcConfig) *RpcMgr {
	mgr := &RpcMgr{}
	mgr.rpcs = rpcConfig
	for _, r := range rpcConfig {
		c, err := ethclient.Dial(r.Rpc)
		if err != nil {
			logrus.Fatal(err)
		}
		mgr.clis = append(mgr.clis, c)
	}
	cid, err := mgr.GetCli().ChainID(context.Background())
	if err != nil {
		logrus.Fatalf("failed get chain id %s", err)
	}
	mgr.chainId = cid
	return mgr
}

var (
	defaultRpc = map[string]*RpcConfig{
		"eth1": {
			Ws:  "wss://eth-hk1.csnodes.com/ws/v1/973eeba6738a7d8c3bd54f91adcbea89",
			Rpc: "https://eth-hk1.csnodes.com/v1/973eeba6738a7d8c3bd54f91adcbea89",
		},
		// "infura": {
		// 	Ws:  "wss://mainnet.infura.io/ws/v3/81e90c9cd6a0430182e3a2bec37f2ba0",
		// 	Rpc: "https://mainnet.infura.io/v3/81e90c9cd6a0430182e3a2bec37f2ba0",
		// },
		// alchemy节点版本太低， 不支持订阅pendingFullTransaction
		// "alchemy": {
		// 	Ws:  "wss://eth-mainnet.g.alchemy.com/v2/jP0h5UEZoR7Wpww9tnNKPGihmNwEkECH",
		// 	Rpc: "https://eth-mainnet.g.alchemy.com/v2/jP0h5UEZoR7Wpww9tnNKPGihmNwEkECH",
		// },
		"chainstack": {
			Ws:  "wss://ethereum-mainnet.core.chainstack.com/ws/61cfcc9622136d52af7e48f1da2fee54",
			Rpc: "https://ethereum-mainnet.core.chainstack.com/61cfcc9622136d52af7e48f1da2fee54",
		},
	}
)

type RpcConfig struct {
	Ws  string `json:"ws"`
	Rpc string `json:"rpc"`
}

func (m *RpcMgr) ChainId() *big.Int {
	return m.chainId
}

func (m *RpcMgr) GetClis() []*ethclient.Client {
	return m.clis
}

func (m *RpcMgr) GetCli() *ethclient.Client {
	return m.clis[0]
}

func (m *RpcMgr) Wss() []string {
	var rpcWs []string
	for _, c := range m.rpcs {
		rpcWs = append(rpcWs, c.Ws)
	}
	return rpcWs
}

func (m *RpcMgr) Rpcs() []string {
	var urls []string
	for _, rpc := range m.rpcs {
		urls = append(urls, rpc.Rpc)
	}
	return urls
}

func AddDefaultRpc(name, ws, rpc string) {
	defaultRpc[name] = &RpcConfig{
		ws, rpc,
	}
}
