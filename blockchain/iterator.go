package blockchain

import (
	"bytes"

	"github.com/dgraph-io/badger"
	"github.com/zweix123/zcoin/utils"
)

type BlockChainIterator struct {
	Current []byte
	DB      *badger.DB
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	iterator := BlockChainIterator{bc.Tip, bc.DB}
	return &iterator
}
func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.DB.View(func(txn *badger.Txn) error {
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
	err := bc.DB.View(func(txn *badger.Txn) error {
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

func (chain *BlockChain) BackOgPrevHash() []byte {
	var ogprevhash []byte
	err := chain.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("ogprevhash"))
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
