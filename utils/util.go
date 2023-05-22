package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"os"
	"runtime"

	"github.com/mr-tron/base58"
	"github.com/zweix123/zcoin/constcoe"
	"golang.org/x/crypto/ripemd160"
)

func Handle(err error) {
	log.SetFlags(0)
	if err != nil {
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			funcName := runtime.FuncForPC(pc).Name()
			log.Fatalf("Error in %s\nAt %s:%d: \n\t%v", funcName, file, line, err)
		} else {
			log.Fatalf("Error: %v", err)
		}
	}
}

func ToHexInt(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	// binary.BigEndian 大端字节序
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func IsInSlice(value int, slice []int) bool {
	for _, val := range slice {
		if val == value {
			return true
		}
	}
	return false
}

func FileExists(fileAddr string) bool {
	if _, err := os.Stat(fileAddr); os.IsNotExist(err) {
		return false
	}
	return true
}

func PublicKeyHash(publicKey []byte) []byte {
	// HASH160(PUBLICKEY): RIPEMD160(SHA256(PUBLICKEY))
	hashedPublicKey := sha256.Sum256(publicKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPublicKey[:])
	Handle(err)
	publicRipeMd := hasher.Sum(nil)
	return publicRipeMd
}

func CheckSum(ripeMdHash []byte) []byte {
	// HASH256(PUBLICKEY): HASH256(HASH256(PUBLICKEY))
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:constcoe.ChecksumLength]
}

func Base58Encode(input []byte) []byte {
	// Base58
	encode := base58.Encode(input)
	return []byte(encode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input[:]))
	Handle(err)
	return decode
}

func PubHash2Address(pubKeyHash []byte) []byte {
	networkVersionedHash := append([]byte{constcoe.NetworkVersion}, pubKeyHash...)
	checkSum := CheckSum(networkVersionedHash)
	finalHash := append(networkVersionedHash, checkSum...)
	address := Base58Encode(finalHash)
	return address
}

func Address2PubHash(address []byte) []byte {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-constcoe.ChecksumLength]
	return pubKeyHash
}
