package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"

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

	err = c.parse(args)
	if err == flag.ErrHelp {
		c.Help = true
		return nil
	}
	return err
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

	flags := c.flags()
	w := bytes.NewBuffer([]byte("Unknown SubCommand\n\nUsage of lemonade:\n"))
	flags.SetOutput(w)
	flags.PrintDefaults()
	return s, fmt.Errorf(w.String())
}

func (c *CLI) flags() *flag.FlagSet {
	flags := flag.NewFlagSet("lemonade", flag.ContinueOnError)
	flags.IntVar(&c.Port, "port", 2489, "TCP port number")
	flags.StringVar(&c.Allow, "allow", "0.0.0.0/0,::0", "Allow IP range")
	flags.StringVar(&c.Host, "host", "localhost", "Destination host name.")
	return flags
}

func (c *CLI) parse(args []string) error {
	flags := c.flags()

	confPath, err := homedir.Expand("~/.config/lemonade.toml")
	if err == nil {
		if confArgs, err := conflag.ArgsFrom(confPath); err == nil {
			flags.Parse(confArgs)
		}
	}

	var arg string
	err = flags.Parse(args[1:])
	if err != nil {
		return err
	}
	if c.Type == PASTE || c.Type == SERVER {
		return nil
	}

	for 0 < flags.NArg() {
		arg = flags.Arg(0)
		err := flags.Parse(flags.Args()[1:])
		if err != nil {
			return err
		}

	}

	if arg != "" {
		c.DataSource = arg
	} else {
		b, err := ioutil.ReadAll(c.In)
		if err != nil {
			return err
		}
		c.DataSource = string(b)
	}
	return nil
}
