package utils

import (
	"bytes"
	"encoding/binary"
	"log"
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
