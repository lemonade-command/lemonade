package lemon

import "fmt"

var Version string
var Usage = fmt.Sprintf(`Usage: lemonade [options]... SUB_COMMAND [arg]
Sub Commands:
  open [URL]                  Open URL by browser
  copy [text]                 Copy text.
  paste                       Paste text.
  server                      Start lemonade server.

Options:
  --port=2489                 TCP port number
  --line-ending               Convert Line Ending(CR/CRLF)
  --allow="0.0.0.0/0,::/0"    Allow IP Range             [Server only]
  --host="localhost"          Destination hostname       [Client only]
  --trans-loopback=true       Translate loopback address [open subcommand only]
  --trans-localfile=true      Translate local file path  [open subcommand only]
  --help                      Show this message


Version:
  %s`, Version)
