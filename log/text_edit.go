package log

import (
	"fmt"
	"github.com/lxn/walk"
	"strings"
)

type TextEditor struct {
	*walk.TextEdit
}

func (w *TextEditor) Fatal(format string, a ...interface{}) {
	w.Log("[F]", format, a...)
}

func (w *TextEditor) Debug(format string, a ...interface{}) {
	w.Log("[D]", format, a...)
}

func (w *TextEditor) Info(format string, a ...interface{}) {
	w.Log("[I]", format, a...)
}

func (w *TextEditor) Error(format string, a ...interface{}) {
	w.Log("[E]", format, a...)
}

func (w *TextEditor) Log(prefix, format string, a ...interface{}) {
	if w.TextEdit == nil {
		return
	}

	msg := strings.Builder{}
	fmt.Fprintf(&msg, format, a...)
	msg.WriteByte('\n')

	last := w.Text()
	l := len(last)
	w.SetTextSelection(l, l)
	w.ReplaceSelectedText(msg.String(), false)

	nl := len(w.Text())
	w.SetTextSelection(nl, nl)
}
