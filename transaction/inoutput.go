package transaction

import (
	"bytes"

	"github.com/zweix123/zcoin/utils"
)

//UTXO（Unspent Transaction Outputs）

type TxOutput struct { // 转出
	Value      int    // 资产值
	HashPubKey []byte // 接受者地址
}

type TxInput struct {
	TxID   []byte // 前置交易信息
	OutIdx int    // 前置交易信息的Output索引
	// 一个区块有多个交易信息, 通过ID找到具体是哪个
	// 一个交易信息有多个output, 通过Index找到具体是那个
	PubKey []byte // 转出者公钥
	Sig    []byte
}

// 这里的地址即不是hash也不是index, 而且"用户名"

func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(in.PubKey, address)
}

func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(out.HashPubKey, utils.PublicKeyHash(address))
}
