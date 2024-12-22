//go:build linux

package server

import (
	//"fmt"
	//log "github.com/inconshreveable/log15"
)

type Ime struct{
}

func (_ *Ime) On(_ string, _ *struct{}) error {
	return nil
}

func (_ *Ime) Off(_ string, _ *struct{}) error {
	return nil
}

func (_ *Ime) Toggle(_ string, _ *struct{}) error {
	return nil
}
