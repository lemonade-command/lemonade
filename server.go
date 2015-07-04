package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"

	"github.com/atotto/clipboard"
	"github.com/skratchdot/open-golang/open"
)

func (c *CLI) Server() int {
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
		rpc.ServeConn(conn)
	}
	return Success
}

type URI struct{}

func (_ *URI) Open(url string, _ *struct{}) error {
	return open.Run(url)
}

type Clipboard struct{}

func (_ *Clipboard) Copy(text string, _ *struct{}) error {
	return clipboard.WriteAll(text)
}

func (_ *Clipboard) Paste(_ struct{}, resp *string) error {
	t, err := clipboard.ReadAll()
	*resp = t
	return err
}
