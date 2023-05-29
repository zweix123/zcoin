package blockchain

import (
	"bytes"
	"encoding/hex"

	"github.com/zweix123/zcoin/transaction"
)

func (bc *BlockChain) VerifyTransactions(txs []*transaction.Transaction) bool {
	if len(txs) == 0 {
		return true
	}
	spentOutputs := make(map[string]int)
	for _, tx := range txs {
		pubKey := tx.Inputs[0].PubKey
		unspentOutputs := bc.FindUnspentTransactions(pubKey)
		inputAmount := 0
		OutputAmount := 0

		for _, input := range tx.Inputs {
			if outidx, ok := spentOutputs[hex.EncodeToString(input.TxID)]; ok && outidx == input.OutIdx {
				return false
			}
			ok, amount := isInputRight(unspentOutputs, input)
			if !ok {
				return false
			}
			inputAmount += amount
			spentOutputs[hex.EncodeToString(input.TxID)] = input.OutIdx
		}

		for _, output := range tx.Outputs {
			OutputAmount += output.Value
		}
		if inputAmount != OutputAmount {
			return false
		}

		if !tx.Verify() {
			return false
		}
	}
	return true
}

func isInputRight(txs []transaction.Transaction, in transaction.TxInput) (bool, int) {
	for _, tx := range txs {
		if bytes.Equal(tx.ID, in.TxID) {
			return true, tx.Outputs[in.OutIdx].Value
		}
	}
	return false, 0
}
