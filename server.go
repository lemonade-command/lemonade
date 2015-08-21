package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/url"
	"regexp"

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

func (u *URI) Open(uri string, _ *struct{}) error {
	conn := <-connCh
	uri = u.translateLoopbackIP(uri, conn)
	return open.Run(uri)
}

func IPv6RemoveBrackets(ip string) string {
	if regexp.MustCompile(`^\[.+\]$`).MatchString(ip) {
		return ip[1 : len(ip)-1]
	}
	return ip
}

func splitHostPort(hostPort string) []string {
	portRe := regexp.MustCompile(`:(\d+)$`)
	portSlice := portRe.FindStringSubmatch(hostPort)
	if len(portSlice) == 0 {
		return []string{IPv6RemoveBrackets(hostPort)}
	}
	port := portSlice[1]
	host := hostPort[:len(hostPort)-len(port)-1]
	return []string{IPv6RemoveBrackets(host), port}
}

func (_ *URI) translateLoopbackIP(uri string, conn net.Conn) string {
	parsed, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	// 0: addr, 1: port
	host := splitHostPort(parsed.Host)

	ip := net.ParseIP(host[0])
	if ip == nil || !ip.IsLoopback() {
		return uri
	}

	addr := conn.RemoteAddr().(*net.TCPAddr).IP.String()

	if len(host) == 1 {
		parsed.Host = addr
	} else {
		parsed.Host = fmt.Sprintf("%s:%s", addr, host[1])
	}

	return parsed.String()
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
