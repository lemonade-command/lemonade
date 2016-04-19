package server

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/pocke/go-iprange"
)

var connCh = make(chan net.Conn, 1)

func Serve(port int, allowIP string) error {
	ra, err := iprange.New(allowIP)
	if err != nil {
		return err
	}

	uri := &URI{}
	rpc.Register(uri)
	clipboard := &Clipboard{}
	rpc.Register(clipboard)

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
		}
		log.Printf("Request from %s", conn.RemoteAddr())
		if !ra.InlucdeConn(conn) {
			continue
		}
		connCh <- conn
		rpc.ServeConn(conn)
	}
}
