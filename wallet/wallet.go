package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/zweix123/zcoin/constcoe"
	"github.com/zweix123/zcoin/utils"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w *Wallet) Address() []byte {
	pubHash := utils.PublicKeyHash(w.PublicKey)
	return utils.PubHash2Address(pubHash)
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	// 在曲线(arg1)随机生成(arg2)生成ECDSA秘钥
	utils.Handle(err)
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return *privateKey, publicKey
}

func NewWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}

func (w *Wallet) Save() {
	filename := constcoe.Wallets + string(w.Address()) + ".wlt"

	privKeyBytes, err := x509.MarshalECPrivateKey(&w.PrivateKey)
	utils.Handle(err)
	privKeyFile, err := os.Create(filename)
	utils.Handle(err)
	err = pem.Encode(privKeyFile, &pem.Block{
		// Type:  "EC PRIVATE KEY",
		Bytes: privKeyBytes,
	})
	utils.Handle(err)
	privKeyFile.Close()
}

func LoadWallet(address string) *Wallet {
	filename := constcoe.Wallets + address + ".wlt"
	if !utils.FileExists(filename) {
		utils.Handle(errors.New("no wallet with such address"))
	}

	privKeyFile, err := os.ReadFile(filename)
	utils.Handle(err)
	pemBlock, _ := pem.Decode(privKeyFile)
	utils.Handle(err)
	privKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	utils.Handle(err)
	publicKey := append(privKey.PublicKey.X.Bytes(), privKey.PublicKey.Y.Bytes()...)
	return &Wallet{
		PrivateKey: *privKey,
		PublicKey:  publicKey,
	}
}
