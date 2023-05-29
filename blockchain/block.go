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
	// header
	// version  // 这里没有
	PrevHash  []byte // 上一个区块的Hash
	Hash      []byte // 区块的Hash  // 实际上是MerkleTree Root的哈希
	Timestamp int64  // 时间戳
	Target    []byte // PoW, target difficulty
	Nonce     int64  // Pow. nonce
	// body
	Transactions []*transaction.Transaction // UTXO
}

func (b *Block) BackTrasactionSummary() []byte {
	// 将body date转换成字节, 这里取各个交易ID
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
			b.PrevHash,
			utils.ToHexInt(b.Timestamp),
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
		PrevHash:     prevhash,
		Hash:         []byte{}, // 占位
		Timestamp:    time.Now().Unix(),
		Target:       []byte{}, // 占位
		Nonce:        0,        // 占位
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
