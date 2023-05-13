package blockchain

import (
	"bytes"
	"crypto/sha256"
	"time"

	"github.com/zweix123/zcoin/utils"
)

type Block struct {
	Timestamp int64  // 时间戳
	Hash      []byte // 区块数据的Hash
	PrevHash  []byte // 上一个区块的Hash
	Target    []byte // PoW, target difficulty
	Nonce     int64  // Pow. nonce
	Data      []byte //
}

func (b *Block) SetHash() {
	// 将多个字节切片按固定分隔符拼接
	information := bytes.Join(
		[][]byte{
			utils.ToHexInt(b.Timestamp),
			b.PrevHash,
			b.Target,
			utils.ToHexInt(b.Nonce),
			b.Data,
		},
		[]byte{},
	)
	hash := sha256.Sum256(information)
	//SHA 256 Secure Hash Algorithm 256, 将任意长的输入数据压缩成固定长度的输出, 256位即32字节
	b.Hash = hash[:]
}

func CreateBlock(prevhash, data []byte) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Hash:      []byte{},
		PrevHash:  prevhash,
		Target:    []byte{},
		Nonce:     0,
		Data:      data,
	}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce() // 注意这里, 寻找nonce时需要区块的hash, 此时只差一个nonce
	block.SetHash()                 // 到这里所有信息已经全了
	return block
}

func GeneslsBlock() *Block {
	genslsWords := "Hello, zcoin!"
	return CreateBlock([]byte{}, []byte(genslsWords))
}
