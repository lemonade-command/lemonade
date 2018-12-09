package server

import (
	"fmt"
	"net"
	"net/url"
	"regexp"

	"github.com/lemonade-command/lemonade/param"
	"github.com/skratchdot/open-golang/open"
)

type URI struct{}

func (u *URI) Open(param *param.OpenParam, _ *struct{}) error {
	conn := <-connCh
	uri := param.URI
	if param.TransLoopback {
		uri = u.translateLoopbackIP(param.URI, conn)
	}
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
