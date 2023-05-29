package constcoe

const (
	Diffculty = 12  // target difficulty
	InitCoin  = 100 // 创建区块**链**时的币数

	TMPDIR = "tmp"

	TransactionPoolFile = "tmp/transactionpool.data"
	BCPath              = "tmp/blocks/"
	BCFile              = "tmp/blocks/MANIFEST" // 和键值数据库badger有关

	ChecksumLength = 4
	NetworkVersion = byte(0x00)

	Wallets = "tmp/wallets/"

	WalletsRefList  = "tmp/reflist/"
	RefListDateFile = "reflist.data"
)

// kv-data std
// l: last block hash
// out-of-bounds: genesis block PreHash
