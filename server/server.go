package server

import (
	"fmt"
	"net"
	"net/rpc"

	log "github.com/inconshreveable/log15"

	"github.com/lemonade-command/lemonade/lemon"
	"github.com/pocke/go-iprange"
)

var connCh = make(chan net.Conn, 1)

var LineEndingOpt string

func Serve(c *lemon.CLI, logger log.Logger) error {
	port := c.Port
	allowIP := c.Allow
	LineEndingOpt = c.LineEnding
	ra, err := iprange.New(allowIP)
	if err != nil {
		return err
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		logger.Info("Request from " + conn.RemoteAddr().String())
		if !ra.InlucdeConn(conn) {
			continue
		}
		connCh <- conn
		rpc.ServeConn(conn)
	}
}

// ServeLocal is for fall back when lemonade client can't connect to server.
// returns port number, error
func ServeLocal(logger log.Logger) (int, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				logger.Crit(err.Error())
				continue
			}
			connCh <- conn
			rpc.ServeConn(conn)
		}
	}()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func init() {
	uri := &URI{}
	rpc.Register(uri)
	clipboard := &Clipboard{}
	rpc.Register(clipboard)
}
