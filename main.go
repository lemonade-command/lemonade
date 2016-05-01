package main

import (
	"os"

	"github.com/pocke/lemonade/lemon"
)

func main() {

	cli := &lemon.CLI{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	os.Exit(cli.Do(os.Args))
}
