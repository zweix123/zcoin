package main

import (
	"fmt"
	"time"

	"github.com/zweix123/zcoin/blockchain"
)

func main() {
	chain := blockchain.CreateBlockChain()
	time.Sleep(time.Second)
	chain.AddBlock("After genesis, I have something to say.")
	time.Sleep(time.Second)
	chain.AddBlock("Leo Cao is awesome!")
	time.Sleep(time.Second)
	chain.AddBlock("I can't wait to follow his github!")
	time.Sleep(time.Second)

	for _, block := range chain.Blocks {
		fmt.Printf("Timestamp: %d\n", block.Timestamp)
		fmt.Printf("hash: %x\n", block.Hash)
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Printf("data: %s\n", block.Data)
		fmt.Println("Proof of Work validation:", block.ValidatePoW())
	}
}
