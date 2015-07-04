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
  --host="localhost"          Destination hostname [Client only]
  --allow="0.0.0.0/0,::0"     Allow IP Range [Server only]
  --help                      Show this message
`

func main() {
	c := &CLI{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	os.Exit(c.Do(os.Args))
}
