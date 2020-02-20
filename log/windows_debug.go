package log

import (
	"fmt"
	"syscall"
	"unsafe"
)

type WindowsDebugLog struct {
	outputDebugStringW *syscall.LazyProc
}

func (w *WindowsDebugLog) Fatal(format string, a ...interface{}) {
	w.Log("[fatal]", format, a...)
}

func (w *WindowsDebugLog) Debug(format string, a ...interface{}) {
	w.Log("[debug]", format, a...)
}

func (w *WindowsDebugLog) Info(format string, a ...interface{}) {
	w.Log("[info]", format, a...)
}

func (w *WindowsDebugLog) Error(format string, a ...interface{}) {
	w.Log("[error]", format, a...)
}

func NewWindowsDebugLog() *WindowsDebugLog {
	kernel := syscall.NewLazyDLL("kernel32")
	return &WindowsDebugLog{
		outputDebugStringW: kernel.NewProc("OutputDebugStringW"),
	}
}

func (w *WindowsDebugLog) Log(prefix, format string, a ...interface{}) {
	s := fmt.Sprintf(prefix+format, a...)
	p, err := syscall.UTF16PtrFromString(s)
	if err == nil {
		_, _, _ = w.outputDebugStringW.Call(uintptr(unsafe.Pointer(p)))
	}
}
