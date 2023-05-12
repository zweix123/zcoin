package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"time"
)

type Block struct {
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Data      []byte
}

type BlockChain struct {
	Blocks []*Block
}

func (b *Block) SetHash() {
	information := bytes.Join([][]byte{ToHexInt(b.Timestamp), b.PrevHash, b.Data}, []byte{})
	hash := sha256.Sum256(information)
	b.Hash = hash[:]
}

func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func CreateBlock(prevhash, data []byte) *Block {
	block := &Block{
		Timestamp: time.Now().Unix(),
		Hash:      []byte{},
		PrevHash:  prevhash,
		Data:      data,
	}
	block.SetHash()
	return block
}

func GeneslsBlock() *Block {
	genslsWords := "Hello, zcoin!"
	return CreateBlock([]byte{}, []byte(genslsWords))
}

func (bc *BlockChain) AddBlock(data string) {
	newBlock := CreateBlock(bc.Blocks[len(bc.Blocks)-1].Hash, []byte(data))
	bc.Blocks = append(bc.Blocks, newBlock)
}

func CreateBlockChain() *BlockChain {
	blockchain := BlockChain{}
	blockchain.Blocks = append(blockchain.Blocks, GeneslsBlock())
	return &blockchain
}

func main() {
	blockchain := CreateBlockChain()
	time.Sleep(time.Second)
	blockchain.AddBlock("After genesis, I have something to say.")
	time.Sleep(time.Second)
	blockchain.AddBlock("Leo Cao is awesome!")
	time.Sleep(time.Second)
	blockchain.AddBlock("I can't wait to follow his github!")
	time.Sleep(time.Second)

	for _, block := range blockchain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Printf("---\n")
	}
}
