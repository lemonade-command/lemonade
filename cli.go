package main

import (
	"fmt"
	"io"
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
	Port  int
	Allow string
	Host  string

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

	switch c.Type {
	case OPEN:
		return c.Open()
	case COPY:
		return c.Copy()
	case PASTE:
		return c.Paste()
	case SERVER:
		return c.Server()
	default:
		panic("")
	}
}

func (c *CLI) writeError(err error) {
	c.Err.Write([]byte(err.Error()))
	c.Err.Write([]byte{'\n'})
}
