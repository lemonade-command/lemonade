package main

import (
	"fmt"
	"net/rpc"
)

func (c *CLI) Open() int {
	client, err := c.rpcClient()
	if err != nil {
		c.writeError(err)
		return RPCError
	}

	err = client.Call("URI.Open", c.DataSource, &struct{}{})
	if err != nil {
		c.writeError(err)
		return RPCError
	}

	return Success
}

func (c *CLI) Paste() int {
	client, err := c.rpcClient()
	if err != nil {
		c.writeError(err)
		return RPCError
	}

	var resp string
	err = client.Call("Clipboard.Paste", struct{}{}, &resp)
	if err != nil {
		c.writeError(err)
		return RPCError
	}
	c.Out.Write([]byte(resp))

	return Success
}

func (c *CLI) Copy() int {
	client, err := c.rpcClient()
	if err != nil {
		c.writeError(err)
		return RPCError
	}

	err = client.Call("Clipboard.Copy", c.DataSource, &struct{}{})
	if err != nil {
		c.writeError(err)
		return RPCError
	}

	return Success
}

func (c *CLI) rpcClient() (*rpc.Client, error) {
	return rpc.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port))
}
