package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/url"

	"github.com/atotto/clipboard"
	"github.com/pocke/go-iprange"
	"github.com/skratchdot/open-golang/open"
)

var connCh = make(chan net.Conn, 1)

func (c *CLI) Server() int {
	ra, err := iprange.New(c.Allow)
	if err != nil {
		c.writeError(err)
		return RPCError
	}

	uri := &URI{}
	rpc.Register(uri)
	clipboard := &Clipboard{}
	rpc.Register(clipboard)

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", c.Port))
	if err != nil {
		c.writeError(err)
		return RPCError
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		c.writeError(err)
		return RPCError
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
	return Success
}

type URI struct{}

func (_ *URI) Open(u string, _ *struct{}) error {
	conn := <-connCh
	parsed, err := url.Parse(u)
	if err != nil {
		return err
	}
	if ip := net.ParseIP(parsed.Host); ip != nil {
		if ip.IsLoopback() {
			// TODO
		}
	}
	return open.Run(u)
}

type Clipboard struct{}

func (_ *Clipboard) Copy(text string, _ *struct{}) error {
	<-connCh
	return clipboard.WriteAll(text)
}

func (_ *Clipboard) Paste(_ struct{}, resp *string) error {
	<-connCh
	t, err := clipboard.ReadAll()
	*resp = t
	return err
}
