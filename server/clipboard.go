package server

import (
	"github.com/lemonade-command/lemonade/lemon"
)

var buff string

type Clipboard struct{}

func (_ *Clipboard) Copy(text string, _ *struct{}) error {
	<-connCh
	buff = lemon.ConvertLineEnding(text, LineEndingOpt)
    return nil
}

func (_ *Clipboard) Paste(_ struct{}, resp *string) error {
	<-connCh
	*resp = buff 
	return nil
}
