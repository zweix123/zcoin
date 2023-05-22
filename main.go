package main

import (
	"os"

	"github.com/zweix123/zcoin/cli"
	"github.com/zweix123/zcoin/constcoe"
)

func main() {
	os.MkdirAll("tmp", 0777)
	os.MkdirAll(constcoe.BCPath, 0777)
	os.MkdirAll(constcoe.Wallets, 0777)
	os.MkdirAll(constcoe.WalletsRefList, 0777)

	defer os.Exit(0)

	cmd := cli.CommandLine{}
	cmd.Run()
}
