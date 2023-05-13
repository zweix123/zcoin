package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"

	"github.com/zweix123/zcoin/constcoe"
	"github.com/zweix123/zcoin/utils"
)

func (b *Block) GetTarget() []byte {
	// 1左移256-diffculty的字节表示
	// 因为hash使用的是SHA256, 左移256则移出, diffculty相当于“退回”的位数
	// 此时我们发现diffculty越大, target越小, 找到小于它的hash就越难
	// Target的选择很复杂, 这里更多是API的保留
	target := big.NewInt(1)
	target.Lsh(target, uint(256-constcoe.Diffculty))
	return target.Bytes()
}

func (b *Block) GetBaseNonce(nonce int64) []byte {
	// 调用这个方法的Block并没有nonce, 不能算出最后SHA, 但是现在不是在选择那个nonce合适嘛, 当做最终的nonce
	data := bytes.Join(
		[][]byte{
			utils.ToHexInt(b.Timestamp),
			b.PrevHash,
			b.Target,
			utils.ToHexInt(int64(nonce)),
			b.BackTrasactionSummary(),
		},
		[]byte{},
	)
	return data
}

func (b *Block) FindNonce() int64 {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	var nonce int64 = 0
	intTarget.SetBytes(b.Target)

	for nonce < math.MaxInt64 { // 目前这个算法不一定能找到
		data := b.GetBaseNonce(nonce)
		hash = sha256.Sum256(data)
		intHash.SetBytes(hash[:])
		if intHash.Cmp(&intTarget) == -1 {
			break
		}
		nonce++
	}
	return nonce
}

func (b *Block) ValidatePoW() bool {
	var intHash big.Int
	var intTarget big.Int
	var hash [32]byte
	intTarget.SetBytes(b.Target)
	data := b.GetBaseNonce(b.Nonce)
	hash = sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(&intTarget) == -1
}
