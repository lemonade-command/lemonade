package server

import "github.com/atotto/clipboard"

type Clipboard struct{}

func (_ *Clipboard) Copy(text string, _ *struct{}) error {
	<-connCh
	return clipboard.WriteAll(text)
}

func (_ *Clipboard) Paste(_ struct{}, resp *string) error {
	<-connCh
	t, err := clipboard.ReadAll()
	*resp = t
	return err
}
