package main

import "os"

const Usage = `Usage: lemonade [options]... SUB_COMMAND [arg]
Sub Commands:
  open [URL]                  Open URL by browser
  copy [text]                 Copy text.
  paste                       Paste text.
  server                      Start lemonade server.

Options:
  --port=2489                 TCP port number
  --host="localhost"          Destination hostname       [Client only]
  --allow="0.0.0.0/0,::0"     Allow IP Range             [Server only]
  --trans-loopback=true       Translate loopback address [Server only]
  --trans-localfile=true      Translate local file path  [Client only]
  --help                      Show this message
`

// XXX: rpc 内で使うためにGlobalにしている
var cli *CLI

func main() {
	cli = &CLI{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	os.Exit(cli.Do(os.Args))
}
