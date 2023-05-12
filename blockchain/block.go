package blockchain

import (
	"bytes"
	"crypto/sha256"
	"time"

	"github.com/zweix123/zcoin/utils"
)

type Block struct {
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Target    []byte
	Nonce     int64
	Data      []byte
}

func (b *Block) SetHash() {
	information := bytes.Join(
		[][]byte{
			utils.ToHexInt(b.Timestamp),
			b.PrevHash,
			utils.ToHexInt(b.Nonce),
			b.Data,
		},
		[]byte{},
	)
	hash := sha256.Sum256(information)
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
	block.Nonce = block.FindNonce()
	block.SetHash()
	return block
}

func GeneslsBlock() *Block {
	genslsWords := "Hello, zcoin!"
	return CreateBlock([]byte{}, []byte(genslsWords))
}
