// main.go
package main

import (
	"os"

	"github.com/zweix123/zcoin/cli"
)

func main() {
	defer os.Exit(0)

	cmd := cli.CommandLine{}
	cmd.Run()
}
