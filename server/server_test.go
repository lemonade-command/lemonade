package server

import (
	"net"
	"reflect"
	"testing"
	"time"
)

func TestSplitHostPort(t *testing.T) {
	assert := func(hostport string, expected []string) {
		got := splitHostPort(hostport)
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("Expected: %v, but got %v", expected, got)
		}
	}
	assert("192.168.0.1", []string{"192.168.0.1"})
	assert("192.168.0.1:3000", []string{"192.168.0.1", "3000"})
	assert("me.pocke.me:3000", []string{"me.pocke.me", "3000"})
	assert("[::1]:3000", []string{"::1", "3000"})
	assert("[::1]", []string{"::1"})
}

type ConnMock struct {
	addr *net.TCPAddr
}

func (_ *ConnMock) Read([]byte) (int, error)         { return 0, nil }
func (_ *ConnMock) Write([]byte) (int, error)        { return 0, nil }
func (_ *ConnMock) Close() error                     { return nil }
func (_ *ConnMock) LocalAddr() net.Addr              { return nil }
func (c *ConnMock) RemoteAddr() net.Addr             { return c.addr }
func (_ *ConnMock) SetDeadline(time.Time) error      { return nil }
func (_ *ConnMock) SetReadDeadline(time.Time) error  { return nil }
func (_ *ConnMock) SetWriteDeadline(time.Time) error { return nil }

var _ net.Conn = &ConnMock{}

func TestURItranslateLoopbackIP(t *testing.T) {
	assert := func(uri string, conn net.Conn, expected string) {
		u := &URI{}
		got := u.translateLoopbackIP(uri, conn)
		if got != expected {
			t.Errorf("Expected: %s, but got %s", expected, got)
		}
	}
	conn := &ConnMock{addr: &net.TCPAddr{IP: net.ParseIP("192.168.0.1")}}

	assert("http://192.168.0.1/foo", conn, "http://192.168.0.1/foo")
	assert("http://192.168.0.1:1124/foo", conn, "http://192.168.0.1:1124/foo")
	assert("http://127.0.0.1:1124/foo", conn, "http://192.168.0.1:1124/foo")
	assert("http://127.0.0.1/foo", conn, "http://192.168.0.1/foo")
	assert("http://127.0.0.1", conn, "http://192.168.0.1")
	assert("https://127.0.0.1", conn, "https://192.168.0.1")
	assert("file://127.0.0.1", conn, "file://192.168.0.1")
	assert("http://[::1]/", conn, "http://192.168.0.1/")
}
