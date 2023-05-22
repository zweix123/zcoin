package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

func Handle(err error) {
	if err != nil {
		log.Panic(err)
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
