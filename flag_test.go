package main

import (
	"os"
	"reflect"
	"testing"
)

func TestCLIParse(t *testing.T) {
	assert := func(args []string, expected CLI) {
		expected.In = os.Stdin
		c := &CLI{In: os.Stdin}
		c.FlagParse(args)

		if !reflect.DeepEqual(expected, *c) {
			t.Errorf("Expected:\n %+v, but got\n %+v", expected, c)
		}
	}

	defaultPort := 2489
	defaultHost := "localhost"
	defaultAllow := "0.0.0.0/0,::/0"

	assert([]string{"xdg-open", "http://example.com"}, CLI{
		Type:           OPEN,
		Host:           defaultHost,
		Port:           defaultPort,
		Allow:          defaultAllow,
		DataSource:     "http://example.com",
		TransLoopback:  true,
		TransLocalfile: true,
	})

	assert([]string{"xdg-open"}, CLI{
		Type:           OPEN,
		Host:           defaultHost,
		Port:           defaultPort,
		Allow:          defaultAllow,
		TransLoopback:  true,
		TransLocalfile: true,
	})

	assert([]string{"pbpaste", "--port", "1124"}, CLI{
		Type:           PASTE,
		Host:           defaultHost,
		Port:           1124,
		Allow:          defaultAllow,
		TransLoopback:  true,
		TransLocalfile: true,
	})

	assert([]string{"pbcopy", "hogefuga"}, CLI{
		Type:           COPY,
		Host:           defaultHost,
		Port:           defaultPort,
		Allow:          defaultAllow,
		DataSource:     "hogefuga",
		TransLoopback:  true,
		TransLocalfile: true,
	})

	assert([]string{"lemonade", "--host", "192.168.0.1", "--port", "1124", "open", "http://example.com"}, CLI{
		Type:           OPEN,
		Host:           "192.168.0.1",
		Port:           1124,
		Allow:          defaultAllow,
		DataSource:     "http://example.com",
		TransLoopback:  true,
		TransLocalfile: true,
	})

	assert([]string{"lemonade", "copy", "hogefuga"}, CLI{
		Type:           COPY,
		Host:           defaultHost,
		Port:           defaultPort,
		Allow:          defaultAllow,
		DataSource:     "hogefuga",
		TransLoopback:  true,
		TransLocalfile: true,
	})

	assert([]string{"lemonade", "paste"}, CLI{
		Type:           PASTE,
		Host:           defaultHost,
		Port:           defaultPort,
		Allow:          defaultAllow,
		TransLoopback:  true,
		TransLocalfile: true,
	})

	assert([]string{"lemonade", "--allow", "192.168.0.0/24", "server", "--port", "1124"}, CLI{
		Type:           SERVER,
		Host:           defaultHost,
		Port:           1124,
		Allow:          "192.168.0.0/24",
		TransLoopback:  true,
		TransLocalfile: true,
	})

	assert([]string{"lemonade", "open", "--trans-loopback=false"}, CLI{
		Type:           OPEN,
		Host:           defaultHost,
		Port:           defaultPort,
		Allow:          defaultAllow,
		TransLoopback:  false,
		TransLocalfile: true,
	})

	assert([]string{"lemonade", "open", "--trans-loopback=true"}, CLI{
		Type:           OPEN,
		Host:           defaultHost,
		Port:           defaultPort,
		Allow:          defaultAllow,
		TransLoopback:  true,
		TransLocalfile: true,
	})

	assert([]string{"lemonade", "open", "--trans-localfile=false"}, CLI{
		Type:           OPEN,
		Host:           defaultHost,
		Port:           defaultPort,
		Allow:          defaultAllow,
		TransLoopback:  true,
		TransLocalfile: false,
	})

	assert([]string{"lemonade", "open", "--trans-localfile=true"}, CLI{
		Type:           OPEN,
		Host:           defaultHost,
		Port:           defaultPort,
		Allow:          defaultAllow,
		TransLoopback:  true,
		TransLocalfile: true,
	})
}
