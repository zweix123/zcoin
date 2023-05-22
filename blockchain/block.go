package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"

	"github.com/zweix123/zcoin/transaction"
	"github.com/zweix123/zcoin/utils"
)

type Block struct {
	Timestamp    int64  // 时间戳
	Hash         []byte // 区块数据的Hash
	PrevHash     []byte // 上一个区块的Hash
	Target       []byte // PoW, target difficulty
	Nonce        int64  // Pow. nonce
	Transactions []*transaction.Transaction
}

func (b *Block) BackTrasactionSummary() []byte {
	txIDs := make([][]byte, 0)
	for _, tx := range b.Transactions {
		txIDs = append(txIDs, tx.ID)
	}
	summary := bytes.Join(txIDs, []byte{})
	return summary
}

func (b *Block) SetHash() {
	information := bytes.Join(
		[][]byte{
			utils.ToHexInt(b.Timestamp),
			b.PrevHash,
			b.Target,
			utils.ToHexInt(b.Nonce),
			b.BackTrasactionSummary(),
		},
		[]byte{},
	)
	hash := sha256.Sum256(information) // SHA256
	b.Hash = hash[:]                   // Sum246返回的是数组
}

func CreateBlock(prevhash []byte, txs []*transaction.Transaction) *Block {
	block := &Block{
		Timestamp:    time.Now().Unix(),
		Hash:         []byte{},
		PrevHash:     prevhash,
		Target:       []byte{},
		Nonce:        0,
		Transactions: txs,
	}
	block.Target = block.GetTarget()
	block.Nonce = block.FindNonce() // 注意这里, 寻找nonce时需要区块的hash, 此时只差一个nonce
	block.SetHash()                 // 到这里所有信息已经全了
	return block
}

func GenesisBlock(address []byte) *Block {
	tx := transaction.BaseTx([]byte(address))
	return CreateBlock([]byte("zweix is sawesome!"), []*transaction.Transaction{tx})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	utils.Handle(err)
	return res.Bytes()
}

func DeSerializeBlock(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.Handle(err)
	return &block
}
