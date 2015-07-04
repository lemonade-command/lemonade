package main

import (
	"fmt"
	"net/rpc"
)

var dummy = &struct{}{}

func (c *CLI) Open() int {
	return c.withRPCClient(func(client *rpc.Client) error {
		return client.Call("URI.Open", c.DataSource, dummy)
	})
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
