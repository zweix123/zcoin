package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"runtime"

	"github.com/dgraph-io/badger"
	"github.com/zweix123/zcoin/constcoe"
	"github.com/zweix123/zcoin/transaction"
	"github.com/zweix123/zcoin/utils"
)

type BlockChain struct {
	LastBlockHash []byte
	Database      *badger.DB
}

func InitBlockChain(address []byte) *BlockChain {
	var lastBlockHash []byte
	if utils.FileExists(constcoe.BCFile) {
		fmt.Println("blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(constcoe.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		genesis := GenesisBlock(address)
		fmt.Println("Genesis Created")
		err = txn.Set(genesis.Hash, genesis.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte("l"), genesis.Hash)
		utils.Handle(err)
		err = txn.Set([]byte("out-of-bounds"), genesis.PrevHash)
		lastBlockHash = genesis.Hash
		return err
	})
	utils.Handle(err)
	return &BlockChain{LastBlockHash: lastBlockHash, Database: db}
}

func ContinueBlockChain() *BlockChain {
	if !utils.FileExists(constcoe.BCFile) {
		fmt.Println("No blockchain found, please create one first")
		runtime.Goexit()
	}

	var lastBlockHash []byte

	opts := badger.DefaultOptions(constcoe.BCPath)
	opts.Logger = nil

	db, err := badger.Open(opts)
	utils.Handle(err)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("l"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastBlockHash = val
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	return &BlockChain{lastBlockHash, db}
}

func (bc *BlockChain) AddBlock(newBlock *Block) {
	var lastBlockHash []byte

	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("l"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			lastBlockHash = val
			return nil
		})
		utils.Handle(err)

		return err
	})
	utils.Handle(err)

	if !bytes.Equal(newBlock.PrevHash, lastBlockHash) { // 测试中不会出现, 现在是单机区块链
		fmt.Println("This block is out of age")
		runtime.Goexit()
	}

	err = bc.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte("l"), newBlock.Hash)
		bc.LastBlockHash = newBlock.Hash
		return err
	})
	utils.Handle(err)
}

func (bc *BlockChain) FindUnspentTransactions(address []byte) []transaction.Transaction {
	var unSpentTxs []transaction.Transaction
	spentTxs := make(map[string][]int)
	// tx id : list[被使用的交易的ID], 这里的key是string而不是[]byte是因为[]byte不能作为key

	iter := bc.Iterator()
	for {
		block := iter.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		IterOutputs:
			for outIdx, out := range tx.Outputs {
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						// 如果这个交易的有个output被通过
						if spentOut == outIdx {
							continue IterOutputs
							// 那这个交易不可能是答案
						}
					}
				}
				// 如果逃过一劫, 就是
				if out.ToAddressRight(address) {
					unSpentTxs = append(unSpentTxs, *tx)
				}
				// 我们注意到这个添加的tx是一个循环外的, 会不会造成多次添加呢?
				// 不会, 因为它一个一个if里, 一个交易的output肯定是不同的
			}

			if !tx.IsBase() { // 记录信息供上面使用
				for _, in := range tx.Inputs {
					if in.FromAddressRight(address) {
						inTxID := hex.EncodeToString(in.TxID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.OutIdx)
					}
				}
			}
		}

		if bc.IsEnd(block) {
			break
		}

	}
	return unSpentTxs
}

func (bc *BlockChain) FindUTXOs(address []byte) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				continue Work // 直接跳出, 因为一个交易, output的人肯定是不同的, 不会再有了
			}
		}
	}
	return accumulated, unspentOuts
}

func (bc *BlockChain) FindSpendableOutputs(address []byte, amount int) (int, map[string]int) {
	unspentOuts := make(map[string]int)
	unspentTxs := bc.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Outputs {
			if out.ToAddressRight(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = outIdx
				if accumulated >= amount {
					break Work // 足够即可
				}
				continue Work // 同上
			}
		}
	}
	return accumulated, unspentOuts
}

func (bc *BlockChain) CreateTransaction(from_PubKey, to_HashPubKey []byte, amount int, privkey ecdsa.PrivateKey) (*transaction.Transaction, bool) {
	var inputs []transaction.TxInput
	var outputs []transaction.TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from_PubKey, amount)
	if acc < amount {
		fmt.Println("Not enough coins!")
		return &transaction.Transaction{}, false
	}
	for txid, outidx := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.Handle(err)
		input := transaction.TxInput{TxID: txID, OutIdx: outidx, PubKey: from_PubKey, Sig: nil}
		inputs = append(inputs, input)
	}

	outputs = append(outputs, transaction.TxOutput{Value: amount, HashPubKey: to_HashPubKey})
	if acc > amount {
		outputs = append(outputs, transaction.TxOutput{Value: acc - amount, HashPubKey: utils.PublicKeyHash(from_PubKey)})
	}
	tx := transaction.Transaction{ID: nil, Inputs: inputs, Outputs: outputs}

	tx.SetID()

	tx.Sign(privkey)
	return &tx, true
}

func (bc *BlockChain) RunMine() {
	transactionPool := CreateTransactionPool()
	if !bc.VerifyTransactions(transactionPool.PubTxs) {
		log.Println("falls in transactions verification")
		err := RemoveTransactionPoolFile() // 不通过也就删除了
		utils.Handle(err)
		return
	}

	candidateBlock := CreateBlock(bc.LastBlockHash, transactionPool.PubTxs) // PoW
	if candidateBlock.ValidatePoW() {
		bc.AddBlock(candidateBlock)
		err := RemoveTransactionPoolFile()
		utils.Handle(err)
		return
	} else {
		fmt.Println("Block has invalid nonce.")
		return
	}
}
