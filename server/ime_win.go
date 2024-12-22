//go:build windows

package server

import (
	"fmt"
	"syscall"
	log "github.com/inconshreveable/log15"
)

type Ime struct{
}

var logger log.Logger

func init() {
	logger = log.New()
	logger.SetHandler(log.LvlFilterHandler(log.LvlDebug, log.StdoutHandler))
}

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	imm32               = syscall.NewLazyDLL("imm32.dll")
	getForegroundWindow = user32.NewProc("GetForegroundWindow")
  sendMessage         = user32.NewProc("SendMessageW")
  immGetDefaultIMEWnd = imm32.NewProc("ImmGetDefaultIMEWnd")
	immGetContext       = imm32.NewProc("ImmGetContext")
	immReleaseContext   = imm32.NewProc("ImmReleaseContext")
	immSetOpenStatus    = imm32.NewProc("ImmSetOpenStatus")
	immGetOpenStatus    = imm32.NewProc("ImmGetOpenStatus")
)
const (
	WM_IME_CONTROL = 0x283 // IME control message
	IMC_GETOPENSTATUS = 0x5 // Retrieve IME open/close status
  IMC_SETOPENSTATUS = 0x6 // Set IME open/close status
)

func getForegroundWindowHandle() syscall.Handle {
	hWnd, _, _ := getForegroundWindow.Call()
	return syscall.Handle(hWnd)
}

func isIMEEnabled() (bool, error) {
	hWnd, _, _ := getForegroundWindow.Call()
  hIME, _, _ := immGetDefaultIMEWnd.Call(hWnd)
	if hIME == 0 {
		return false, fmt.Errorf("failed to get ime handle")
	}

  ret, _, _ := sendMessage.Call(uintptr(hIME), uintptr(WM_IME_CONTROL), IMC_GETOPENSTATUS, 0)
  logger.Debug("sendMessage", "ret", ret)
	return ret != 0, nil
}

func boolToUintptr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

func (_ *Ime) On(_ string, _ *struct{}) error {
	<-connCh
	logger.Debug("Ime.On requested")
	hWnd, _, _ := getForegroundWindow.Call()
  hIME, _, _ := immGetDefaultIMEWnd.Call(hWnd)
	if hIME == 0 {
		return fmt.Errorf("failed to get ime handle")
	}

  sendMessage.Call(uintptr(hIME), uintptr(WM_IME_CONTROL), IMC_SETOPENSTATUS, 1)

	return nil
}

func (_ *Ime) Off(_ string, _ *struct{}) error {
	<-connCh
	logger.Debug("Ime.Off requested")
	hWnd, _, _ := getForegroundWindow.Call()
  hIME, _, _ := immGetDefaultIMEWnd.Call(hWnd)
	if hIME == 0 {
		return fmt.Errorf("failed to get ime handle")
	}

  sendMessage.Call(uintptr(hIME), uintptr(WM_IME_CONTROL), IMC_SETOPENSTATUS, 0)

	return nil
}

func (_ *Ime) Toggle(_ string, _ *struct{}) error {
	<-connCh
	logger.Debug("Ime.Toggle requested")
	enabled, err := isIMEEnabled()
	if err != nil {
		return nil
	}
  logger.Debug("Current IME status", "sutatus", boolToUintptr(enabled))
	hWnd, _, _ := getForegroundWindow.Call()
  hIME, _, _ := immGetDefaultIMEWnd.Call(hWnd)
	if hIME == 0 {
		return fmt.Errorf("failed to get ime handle")
	}

  sendMessage.Call(uintptr(hIME), uintptr(WM_IME_CONTROL),
                   IMC_SETOPENSTATUS, boolToUintptr(!enabled))
	enabled, err = isIMEEnabled()
	if err != nil {
		return nil
	}
	logger.Debug("Current IME status", "status", boolToUintptr(enabled))

	return nil
}
