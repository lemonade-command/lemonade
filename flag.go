package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/monochromegane/conflag"
)

func (c *CLI) FlagParse(args []string) error {
	style, err := c.getCommandType(args)
	if err != nil {
		return err
	}
	if style == SUBCOMMAND {
		args = args[:len(args)-1]
	}

	c.parse(args)

	return nil
}

func (c *CLI) getCommandType(args []string) (s CommandStyle, err error) {
	s = ALIAS
	switch args[0] {
	case "xdg-open":
		c.Type = OPEN
		return
	case "pbpaste":
		c.Type = PASTE
		return
	case "pbcopy":
		c.Type = COPY
		return
	}

	del := func(i int) {
		copy(args[i+1:], args[i+2:])
		args[len(args)-1] = ""
	}

	s = SUBCOMMAND
	for i, v := range args[1:] {
		switch v {
		case "open":
			c.Type = OPEN
			del(i)
			return
		case "paste":
			c.Type = PASTE
			del(i)
			return
		case "copy":
			c.Type = COPY
			del(i)
			return
		case "server":
			c.Type = SERVER
			del(i)
			return
		}
	}
	return s, fmt.Errorf("Unknown subcommand")
}

func (c *CLI) parse(args []string) {
	flags := flag.NewFlagSet("lemonade", flag.ContinueOnError)
	flags.IntVar(&c.Port, "port", 2489, "TCP port number")
	flags.StringVar(&c.Allow, "allow", "0.0.0.0/0,::0", "Allow IP range")
	flags.StringVar(&c.Host, "host", "localhost", "Destination host name.")

	confPath, err := homedir.Expand("~/.config/lemonade.toml")
	if err == nil {
		if confArgs, err := conflag.ArgsFrom(confPath); err == nil {
			flags.Parse(confArgs)
		}
	}

	var arg string
	flags.Parse(args[1:])
	if c.Type == PASTE || c.Type == SERVER {
		return
	}
	for 0 < flags.NArg() {
		arg = flags.Arg(0)
		flags.Parse(flags.Args()[1:])
	}
	if arg != "" {
		c.DataSource = strings.NewReader(arg)
	} else {
		c.DataSource = c.In
	}
	return
}
