package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/zweix123/zcoin/constcoe"
	"github.com/zweix123/zcoin/utils"
)

type Transaction struct {
	ID      []byte     // 每个交易信息的唯一ID, 是整个结构体的hash
	Inputs  []TxInput  // 该交易信息维护的inputs
	Outputs []TxOutput // 该交易信息维护的outputs
}

func (tx *Transaction) TxHash() []byte {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	utils.Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	return hash[:]
}

func (tx *Transaction) SetID() {
	tx.ID = tx.TxHash()
}

func BaseTx(toaddress []byte) *Transaction {
	txIn := TxInput{
		TxID:        []byte{},
		OutIdx:      -1,
		FromAddress: []byte{},
	}
	txOut := TxOutput{
		Value:     constcoe.InitCoin,
		ToAddress: toaddress,
	}
	tx := Transaction{
		ID:      []byte("This is the Base Transaction!"),
		Inputs:  []TxInput{txIn},
		Outputs: []TxOutput{txOut},
	}
	return &tx
}

func (tx *Transaction) IsBase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].OutIdx == -1
}
