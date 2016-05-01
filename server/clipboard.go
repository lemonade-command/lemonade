package server

import (
	"github.com/atotto/clipboard"
	"github.com/pocke/lemonade/lemon"
)

type Clipboard struct{}

func (_ *Clipboard) Copy(text string, _ *struct{}) error {
	<-connCh
	return clipboard.WriteAll(lemon.ConvertLineEnding(text, lineEndingOpt))
}

func (_ *Clipboard) Paste(_ struct{}, resp *string) error {
	<-connCh
	t, err := clipboard.ReadAll()
	*resp = t
	return err
}
