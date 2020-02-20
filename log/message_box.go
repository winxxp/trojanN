package log

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
)

type WindowsMessageBoxLog struct {
	caption *uint16
}

func (w *WindowsMessageBoxLog) Fatal(format string, a ...interface{}) {
	w.Log(format, a...)
}

func (w *WindowsMessageBoxLog) Debug(format string, a ...interface{}) {

}

func (w *WindowsMessageBoxLog) Info(format string, a ...interface{}) {

}

func (w *WindowsMessageBoxLog) Error(format string, a ...interface{}) {
	w.Log(format, a...)
}

func NewWindowMessageBoxLog(caption string) *WindowsMessageBoxLog {
	caption_, err := syscall.UTF16PtrFromString(caption)
	if err != nil {
		panic(err)
	}
	return &WindowsMessageBoxLog{caption: caption_}
}

func (w *WindowsMessageBoxLog) Log(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	p, err := syscall.UTF16PtrFromString(s)
	if err == nil {
		_, _ = windows.MessageBox(0, p, w.caption, windows.MB_OK)
	}

}
