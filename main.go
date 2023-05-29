package main

import (
	"os"

	"github.com/zweix123/zcoin/cli"
	"github.com/zweix123/zcoin/constcoe"
)

func init() {
	os.MkdirAll(constcoe.TMPDIR, 0777)
	os.MkdirAll(constcoe.BCPath, 0777)
	os.MkdirAll(constcoe.Wallets, 0777)
	os.MkdirAll(constcoe.WalletsRefList, 0777)
}

func terminate() {
	// os.RemoveAll(constcoe.TMPDIR)
	os.Exit(0)
}

func main() {
	defer terminate()

	cmd := cli.CommandLine{}
	cmd.Run()
}
