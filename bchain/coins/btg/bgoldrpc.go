package bch

import (
	"blockbook/bchain"
	"blockbook/bchain/coins/btc"
	"encoding/hex"
	"encoding/json"

	"github.com/cpacia/bchutil"
	"github.com/golang/glog"
	"github.com/juju/errors"
)

// BGoldRPC is an interface to JSON-RPC bitcoind service.
type BGoldRPC struct {
	*btc.BitcoinRPC
}

// NewBCashRPC returns new BGoldRPC instance.
func NewBGoldRPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}

	s := &BGoldRPC{
		b.(*btc.BitcoinRPC),
	}

	return s, nil
}

// Initialize initializes BGoldRPC instance.
func (b *BGoldRPC) Initialize() error {
	chainName, err := b.GetChainInfoAndInitializeMempool(b)
	if err != nil {
		return err
	}

	params := GetChainParams(chainName)

	// always create parser
	b.Parser, err = NewBGoldParser(params, b.ChainConfig)

	if err != nil {
		return err
	}

	// parameters for getInfo request
	if params.Net == MainnetMagic {
		b.Testnet = false
		b.Network = "livenet"
	} else {
		b.Testnet = true
		b.Network = "testnet"
	}

	glog.Info("rpc: block chain ", params.Name)

	return nil
}

//
// // getblock
//
// type cmdGetBlock struct {
// 	Method string `json:"method"`
// 	Params struct {
// 		BlockHash string `json:"blockhash"`
// 		Verbose   bool   `json:"verbose"`
// 	} `json:"params"`
// }
//
// type resGetBlockRaw struct {
// 	Error  *bchain.RPCError `json:"error"`
// 	Result string           `json:"result"`
// }
//
// type resGetBlockThin struct {
// 	Error  *bchain.RPCError `json:"error"`
// 	Result bchain.ThinBlock `json:"result"`
// }
//
// // estimatesmartfee
//
// type cmdEstimateSmartFee struct {
// 	Method string `json:"method"`
// 	Params struct {
// 		Blocks int `json:"nblocks"`
// 	} `json:"params"`
// }
//
// type resEstimateSmartFee struct {
// 	Error  *bchain.RPCError `json:"error"`
// 	Result struct {
// 		Feerate float64 `json:"feerate"`
// 		Blocks  int     `json:"blocks"`
// 	} `json:"result"`
// }
//
// // GetBlock returns block with given hash.
// func (b *BGoldRPC) GetBlock(hash string, height uint32) (*bchain.Block, error) {
// 	var err error
// 	if hash == "" && height > 0 {
// 		hash, err = b.GetBlockHash(height)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	header, err := b.GetBlockHeader(hash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	data, err := b.GetBlockRaw(hash)
// 	if err != nil {
// 		return nil, err
// 	}
// 	block, err := b.Parser.ParseBlock(data)
// 	if err != nil {
// 		return nil, errors.Annotatef(err, "hash %v", hash)
// 	}
// 	block.BlockHeader = *header
// 	return block, nil
// }
//
// // GetBlockRaw returns block with given hash as bytes.
// func (b *BGoldRPC) GetBlockRaw(hash string) ([]byte, error) {
// 	glog.V(1).Info("rpc: getblock (verbose=0) ", hash)
//
// 	res := resGetBlockRaw{}
// 	req := cmdGetBlock{Method: "getblock"}
// 	req.Params.BlockHash = hash
// 	req.Params.Verbose = false
// 	err := b.Call(&req, &res)
//
// 	if err != nil {
// 		return nil, errors.Annotatef(err, "hash %v", hash)
// 	}
// 	if res.Error != nil {
// 		if isErrBlockNotFound(res.Error) {
// 			return nil, bchain.ErrBlockNotFound
// 		}
// 		return nil, errors.Annotatef(res.Error, "hash %v", hash)
// 	}
// 	return hex.DecodeString(res.Result)
// }
//
// // GetBlockFull returns block with given hash.
// func (b *BGoldRPC) GetBlockFull(hash string) (*bchain.Block, error) {
// 	return nil, errors.New("Not implemented")
// }
//
// // EstimateSmartFee returns fee estimation.
// func (b *BGoldRPC) EstimateSmartFee(blocks int, conservative bool) (float64, error) {
// 	glog.V(1).Info("rpc: estimatesmartfee ", blocks)
//
// 	res := resEstimateSmartFee{}
// 	req := cmdEstimateSmartFee{Method: "estimatesmartfee"}
// 	req.Params.Blocks = blocks
// 	// conservative param is omitted
// 	err := b.Call(&req, &res)
//
// 	if err != nil {
// 		return 0, err
// 	}
// 	if res.Error != nil {
// 		return 0, res.Error
// 	}
// 	return res.Result.Feerate, nil
// }
//
// func isErrBlockNotFound(err *bchain.RPCError) bool {
// 	return err.Message == "Block not found" ||
// 		err.Message == "Block height out of range"
// }
