package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

// 私钥签名, 公钥认证

func Sign(msg []byte, privKey ecdsa.PrivateKey) []byte {
	r, s, err := ecdsa.Sign(rand.Reader, &privKey, msg)
	Handle(err)
	signature := append(r.Bytes(), s.Bytes()...)
	return signature
}

func Verify(msg []byte, pubkey []byte, signature []byte) bool {
	return ecdsa.Verify(
		&ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(pubkey[:(len(pubkey) / 2)]),
			Y:     new(big.Int).SetBytes(pubkey[(len(pubkey) / 2):]),
		},
		msg,
		new(big.Int).SetBytes(signature[:(len(signature)/2)]),
		new(big.Int).SetBytes(signature[(len(signature)/2):]),
	)
}
