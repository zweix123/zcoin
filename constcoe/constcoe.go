package constcoe

const (
	Diffculty = 12  // target difficulty
	InitCoin  = 100 // 创建区块**链**时的币数

	TransactionPoolFile = "tmp/transactionpool.data"
	BCPath              = "tmp/blocks"
	BCFile              = "tmp/blocks/MANIFEST" // manifest

	ChecksumLength = 4
	NetworkVersion = byte(0x00)

	Wallets        = "tmp/wallets/"
	WalletsRefList = "tmp/ref_list/"
)
