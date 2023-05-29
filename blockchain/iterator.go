package blockchain

import (
	"bytes"

	"github.com/dgraph-io/badger"
	"github.com/zweix123/zcoin/utils"
)

type BlockChainIterator struct {
	Current  []byte
	Database *badger.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{bc.LastBlockHash, bc.Database}
	return &iterator
}

// func (iterator *BlockChainIterator) Current() *Block

func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.Current)
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			block = DeSerializeBlock(val)
			return nil
		})
		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	iterator.Current = block.PrevHash

	return block
}

func (bc *BlockChain) IsEnd(b *Block) bool {
	var isEnd bool
	var genesishash []byte
	err := bc.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("out-of-bounds"))
		utils.Handle(err)
		err = item.Value(func(val []byte) error {
			genesishash = val
			return nil
		})
		utils.Handle(err)
		isEnd = bytes.Equal(b.PrevHash, genesishash)
		return nil
	})
	utils.Handle(err)
	return isEnd
}

func (chain *BlockChain) GetOutBoundHash() []byte {
	var ogprevhash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("out-of-bounds"))
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			ogprevhash = val
			return nil
		})

		utils.Handle(err)
		return err
	})
	utils.Handle(err)

	return ogprevhash
}
