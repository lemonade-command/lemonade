package main

import (
	"fmt"
	"os"

	log "github.com/inconshreveable/log15"

	"github.com/lemonade-command/lemonade/lemon"
	"github.com/lemonade-command/lemonade/server"
)

var logLevelMap = map[int]log.Lvl{
	0: log.LvlDebug,
	1: log.LvlInfo,
	2: log.LvlWarn,
	3: log.LvlError,
	4: log.LvlCrit,
}

func main() {

	cli := &lemon.CLI{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	os.Exit(Do(cli, os.Args))
}

func Do(c *lemon.CLI, args []string) int {
	logger := log.New()
	logger.SetHandler(log.LvlFilterHandler(log.LvlError, log.StdoutHandler))

	if err := c.FlagParse(args, false); err != nil {
		writeError(c, err)
		return lemon.FlagParseError
	}

	logLevel := logLevelMap[c.LogLevel]
	logger.SetHandler(log.LvlFilterHandler(logLevel, log.StdoutHandler))

	if c.Help {
		fmt.Fprint(c.Err, lemon.Usage)
		return lemon.Help
	}

	var err error

	switch c.Type {
	case lemon.SERVER:
		logger.Debug("Starting Server")
		err = server.Serve(c, logger)
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
