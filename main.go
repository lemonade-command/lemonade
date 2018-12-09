package main

import (
	"fmt"
	"os"

	"github.com/lemonade-command/lemonade/client"
	"github.com/lemonade-command/lemonade/lemon"
	"github.com/lemonade-command/lemonade/server"
)

func main() {

	cli := &lemon.CLI{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	os.Exit(Do(cli, os.Args))
}

func Do(c *lemon.CLI, args []string) int {
	if err := c.FlagParse(args, false); err != nil {
		writeError(c, err)
		return lemon.FlagParseError
	}
	if c.Help {
		fmt.Fprint(c.Err, lemon.Usage)
		return lemon.Help
	}

	lc := client.New(c)
	var err error

	switch c.Type {
	case lemon.OPEN:
		err = lc.Open(c.DataSource, c.TransLocalfile, c.TransLoopback)
	case lemon.COPY:
		err = lc.Copy(c.DataSource)
	case lemon.PASTE:
		var text string
		text, err = lc.Paste()
		c.Out.Write([]byte(text))
	case lemon.SERVER:
		err = server.Serve(c)
	default:
		panic("Unreachable code")
	}

	if err != nil {
		writeError(c, err)
		return lemon.RPCError
	}
	return lemon.Success
}

func writeError(c *lemon.CLI, err error) {
	fmt.Fprintln(c.Err, err.Error())
}
