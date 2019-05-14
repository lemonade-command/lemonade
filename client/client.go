package client

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	log "github.com/inconshreveable/log15"
	"github.com/lemonade-command/lemonade/lemon"
)

type client struct {
	host               string
	port               int
	addr               string
	lineEnding         string
	noFallbackMessages bool
	logger             log.Logger
}

func New(c *lemon.CLI, logger log.Logger) *client {
	return &client{
		host:               c.Host,
		port:               c.Port,
		addr:               fmt.Sprintf("http://%s:%d", c.Host, c.Port),
		lineEnding:         c.LineEnding,
		noFallbackMessages: c.NoFallbackMessages,
		logger:             logger,
	}
}

func (c *client) Copy(text string) error {
	c.logger.Debug("Sending: " + text)
	url := fmt.Sprintf("%s/copy", c.addr)
	_, err := http.Post(url, "text/plain", strings.NewReader(text))
	if err != nil {
		clipboard.WriteAll(text)
		return err
	}
	return nil
}

func (c *client) Paste() (string, error) {
	url := fmt.Sprintf("%s/paste", c.addr)
	r, err := http.Get(url)
	if err != nil {
		c.logger.Error("http error.", "err", err.Error())
		return "", err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.logger.Error("http read body error.", "err", err.Error())
		return "", err
	}

	return lemon.ConvertLineEnding(string(body), c.lineEnding), nil
}

func (c *client) Open(uri string, transLoopback bool) error {
	url := fmt.Sprintf("%s/open?uri=%s&transLoopback=%s&base64=true", c.addr, base64.URLEncoding.EncodeToString([]byte(uri)), strconv.FormatBool(transLoopback))
	c.logger.Info("Opening: " + uri)

	_, err := http.Get(url)
	if err != nil {
		c.logger.Error("http error.", "err", err.Error())
		return err
	}
	return nil
}
