package main

import "os"

func main() {
	c := &CLI{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	c.FlagParse(os.Args)
}
