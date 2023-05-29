# zcoin

+ Reference of [深入底层：Go语言从零构建区块链](https://www.krad.top/zh-cn/categories/goblockchain/)：自底向上的用Go实现类似比特币的本地区块链(教程的第二部分计划实现分布式网络的部署，但是在2023.6.1日之前还没有更新)

这套教程是自底向上完成区块链的实现，我在这里提供自顶向下的视角

首先我们要怎么确定一个人的身份呢？在比特币中使用是加密算法ECCSA，倒不要理解其内涵，只需要明白一个加密算法包括私钥和公钥，在这里，我通过私钥加密一段话，任何人都可以通过公钥进行解密，如果我加密的是一个哈希，众人还原发现和这个哈希还真一样，就可以证明这个哈希是我的了。

+ 一个钱包归属于一个人，包含它的ECC的公钥私钥（有私钥，所以要自己保管好）
+ 通过公钥[生成](https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses)一个Bitcoin address  
    在区块链中通过这个地址表示一个人，这个生成过程有不可逆哈希，所以得不到你的公钥
+ 每个交易信息就是通过私钥对这个交易信息的哈希进行加密来证明这个交易信息是你的

那么什么是区块链呢？它真的是一个链数据结构，我们看它的[数据结构](https://en.bitcoin.it/wiki/Block_hashing_algorithm)
```go
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
```

其中前两个就有将链进行连接的属性

那么这里Pow是什么呢？
...


这里的Transactions即使区块的数据，以一个一个交易信息描述，这就是UTXO
那么什么是UTXO呢？
...


---

除此之外，区块链的信息是怎么持久化的保存到本地的呢？比特币使用的是键值数据库，  
在我这里，是这样的规定
```
// kv-data standard
// l: last block hash
// out-of-bounds: genesis block PreHash
```
而比特币是[这样](https://en.bitcoin.it/wiki/Bitcoin_Core_0.11_(ch_2):_Data_Storage)的