package transaction

import "bytes"

//UTXO（Unspent Transaction Outputs）

type TxOutput struct { // 转出
	Value     int    // 资产值
	ToAddress []byte // 接收者地址
}

type TxInput struct {
	TxID   []byte // 前置交易信息
	OutIdx int    // 前置交易信息的Output索引
	// 一个区块有多个交易信息, 通过ID找到具体是哪个
	// 一个交易信息有多个output, 通过Index找到具体是那个
	FromAddress []byte // 转出者地址
}

// 这里的地址即不是hash也不是index, 而且"用户名"

func (in *TxInput) FromAddressRight(address []byte) bool {
	return bytes.Equal(in.FromAddress, address)
}

func (out *TxOutput) ToAddressRight(address []byte) bool {
	return bytes.Equal(out.ToAddress, address)
}
