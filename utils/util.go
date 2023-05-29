package utils

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"runtime"
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
