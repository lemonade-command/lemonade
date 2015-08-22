package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

var dummy = &struct{}{}

func exists(fname string) bool {
	_, err := os.Stat(fname)
	return err == nil
}

func serveFile(fname string) (string, <-chan struct{}, error) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", nil, err
	}
	finished := make(chan struct{})

	go func() {
		http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b, err := ioutil.ReadFile(fname)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)

			w.(http.Flusher).Flush()
			finished <- struct{}{}
		}))
	}()

	return fmt.Sprintf("http://127.0.0.1:%d/%s", l.Addr().(*net.TCPAddr).Port, fname), finished, nil
}

func (c *CLI) Open() int {
	var finished <-chan struct{}
	st := c.withRPCClient(func(client *rpc.Client) error {
		var uri string
		if exists(c.DataSource) {
			var err error
			uri, finished, err = serveFile(c.DataSource)
			if err != nil {
				return err
			}
		} else {
			uri = c.DataSource
		}

		return client.Call("URI.Open", uri, dummy)
	})
	if finished != nil {
		<-finished
	}
	return st
}

func (c *CLI) Paste() int {
	return c.withRPCClient(func(client *rpc.Client) error {
		var resp string
		err := client.Call("Clipboard.Paste", dummy, &resp)
		if err == nil {
			c.Out.Write([]byte(resp))
		}
		return err
	})
}

func (c *CLI) Copy() int {
	return c.withRPCClient(func(client *rpc.Client) error {
		return client.Call("Clipboard.Copy", c.DataSource, dummy)
	})
}

func (c *CLI) withRPCClient(f func(*rpc.Client) error) int {
	client, err := rpc.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
	if err != nil {
		c.writeError(err)
		return RPCError
	}

	err = f(client)
	if err != nil {
		c.writeError(err)
		return RPCError
	}
	return Success
}
