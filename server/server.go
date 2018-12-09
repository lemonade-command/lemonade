package server

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/lemonade-command/lemonade/lemon"
	"github.com/pocke/go-iprange"
)

var connCh = make(chan net.Conn, 1)

var LineEndingOpt string

func Serve(c *lemon.CLI) error {
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
			log.Println(err)
			continue
		}
		log.Printf("Request from %s", conn.RemoteAddr())
		if !ra.InlucdeConn(conn) {
			continue
		}
		connCh <- conn
		rpc.ServeConn(conn)
	}
}

// ServeLocal is to fall back when lemonade client can't connect server.
// returns port number, error
func ServeLocal() (int, error) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Panicln(err)
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
