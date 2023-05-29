package transaction

import (
	"bytes"
	"crypto/ecdsa"
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
	tx.ID = tx.TxHash() // ID就是整体Hash
}

func BaseTx(toaddress []byte) *Transaction {
	txIn := TxInput{
		TxID:   []byte{},
		OutIdx: -1,
		PubKey: []byte{},
		Sig:    nil,
	}
	txOut := TxOutput{
		Value:      constcoe.InitCoin,
		HashPubKey: toaddress,
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

func (tx *Transaction) PlainCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	for _, txin := range tx.Inputs {
		inputs = append(inputs, TxInput{txin.TxID, txin.OutIdx, nil, nil})
	}
	for _, txout := range tx.Outputs {
		outputs = append(outputs, TxOutput{txout.Value, txout.HashPubKey})
	}
	return Transaction{tx.ID, inputs, outputs}
}

func (tx *Transaction) PlainHash(inidx int, prevPubKey []byte) []byte {
	txCopy := tx.PlainCopy()
	txCopy.Inputs[inidx].PubKey = prevPubKey
	return txCopy.TxHash()
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey) {
	if tx.IsBase() {
		return
	}
	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey) // 分别签名
		signature := utils.Sign(plainhash, privKey)  // 签名是不包括签名的hash的加密
		tx.Inputs[idx].Sig = signature
	}
}

func (tx *Transaction) Verify() bool {
	for idx, input := range tx.Inputs {
		plainhash := tx.PlainHash(idx, input.PubKey)
		if !utils.Verify(plainhash, input.PubKey, input.Sig) {
			return false
		}
	}
	return true
}
