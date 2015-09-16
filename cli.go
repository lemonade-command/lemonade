package main

import (
	"fmt"
	"io"

	"github.com/pocke/lemonade/client"
	"github.com/pocke/lemonade/server"
)

type CommandType int

// Commands
const (
	OPEN CommandType = iota + 1
	COPY
	PASTE
	SERVER
)

const (
	Success        = 0
	FlagParseError = iota + 10
	RPCError
	Help
)

type CommandStyle int

const (
	ALIAS CommandStyle = iota + 1
	SUBCOMMAND
)

type CLI struct {
	In       io.Reader
	Out, Err io.Writer

	Type       CommandType
	DataSource string

	// options
	Port           int
	Allow          string
	Host           string
	TransLoopback  bool
	TransLocalfile bool

	Help bool
}

func (c *CLI) Do(args []string) int {
	if err := c.FlagParse(args); err != nil {
		c.writeError(err)
		return FlagParseError
	}
	if c.Help {
		fmt.Fprint(c.Err, Usage)
		return Help
	}

	lc := client.New(c.Host, c.Port)
	var err error

	switch c.Type {
	case OPEN:
		err = lc.Open(c.DataSource, c.TransLocalfile, c.TransLoopback)
	case COPY:
		err = lc.Copy(c.DataSource)
	case PASTE:
		var text string
		text, err = lc.Paste()
		c.Out.Write([]byte(text))
	case SERVER:
		err = server.Serve(c.Port, c.Allow)
	default:
		panic("Unreachable code")
	}

	if err != nil {
		c.writeError(err)
		return RPCError
	}
	return Success
}

func (c *CLI) writeError(err error) {
	fmt.Fprintln(c.Err, err.Error())
}
