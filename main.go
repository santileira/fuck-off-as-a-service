package main

import (
	"github.com/santileira/fuck-off-as-a-service/cmd"
	"os"
)

func main() {
	if err := cmd.Cmds().Execute(); err != nil {
		os.Exit(1)
	}
}
