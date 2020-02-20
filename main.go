package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/winxxp/trojanN/log"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var (
	logWidget      = &log.TextEditor{}
	wgtStartButton = &walk.PushButton{}
	wgtStopButton  = &walk.PushButton{}
	app            = &App{}
	nip            = &NotifyIconDisplay{}

	ni *walk.NotifyIcon
)

type App struct {
	*walk.MainWindow
}

func ChangeIcon(ic *walk.Icon) {
	app.SetIcon(ic)
	ni.SetIcon(ic)
}

func ResetIcon() {
	ChangeIcon(nip.Reset().Icon())
}

func NextIcon() {
	ChangeIcon(nip.Next().Icon())
}

func main() {
	var (
		cancel context.CancelFunc
	)

	log.AddLogger(logWidget)

	defer func() {
		if cancel != nil {
			cancel()
			time.Sleep(time.Second)
		}
	}()

	ic, err := walk.NewIconFromImageForDPI(makeDigitImage(0), 96)
	if err != nil {
		log.Error(err.Error())
		return
	}

	mw := MainWindow{
		AssignTo: &app.MainWindow,
		Title:    "TrojanN",
		Size:     Size{Width: 700, Height: 500},
		Icon:     ic,
		Layout: VBox{
			Margins: Margins{
				Left:   15,
				Top:    10,
				Right:  15,
				Bottom: 15,
			},
			Spacing: 15,
		},
		Children: []Widget{
			Composite{
				MinSize: Size{Width: 400},
				Layout: HBox{
					MarginsZero: true,
				},
				Children: []Widget{
					PushButton{
						AssignTo: &wgtStartButton,
						Text:     "开始(&S)",
						OnClicked: func() {
							var ctx context.Context
							ctx, cancel = context.WithCancel(context.Background())

							r, w, err := os.Pipe()
							if err != nil {
								log.Error(err.Error())
								return
							}

							ResetIcon()

							go func() {
								cmd := exec.CommandContext(ctx, "trojan.exe")
								cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
								cmd.Stderr = w
								cmd.Stdout = w

								defer w.Close()

								wgtStartButton.SetEnabled(false)
								if err := cmd.Run(); err != nil {
									wgtStartButton.SetEnabled(true)
									log.Error(err.Error())
									return
								}
								cancel = nil
							}()

							go func() {
								defer func() {
									log.Info("退出服务")
									wgtStartButton.SetEnabled(true)
									ResetIcon()
								}()
								buf := make([]byte, 1024)
								for {
									n, err := r.Read(buf)
									if err != nil {
										if errors.Is(err, io.EOF) {
											return
										}
										log.Error(err.Error())
										return
									} else {
										NextIcon()
										lines := bytes.Split(buf[:n], []byte("\n"))
										for _, line := range lines {
											log.Info(string(string(line)))
										}
									}
								}
							}()
						},
					},
					PushButton{
						AssignTo: &wgtStopButton,
						Text:     "停止(&T)",
						OnClicked: func() {
							defer wgtStartButton.SetEnabled(true)
							if cancel == nil {
								return
							}
							cancel()
							cancel = nil
							log.Info("退出服务")
						},
					},
					PushButton{
						Text: "删除日志(&C)",
						OnClicked: func() {
							l := len(logWidget.Text())
							logWidget.SetTextSelection(0, l)
							logWidget.ReplaceSelectedText("", false)
						},
					},
				},
			},
			TextEdit{
				AssignTo: &logWidget.TextEdit,
				Name:     "Test",
				HScroll:  true,
				VScroll:  true,
				ReadOnly: true,
			},
		},
	}
	if err := mw.Create(); err != nil {
		log.Error("wm.Create: %v", err)
	}

	ni, err = walk.NewNotifyIcon(app.MainWindow)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ni.Dispose()

	// Set the icon and a tool tip text.
	if err := ni.SetIcon(ic); err != nil {
		log.Fatal(err.Error())
	}

	ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		app.SetVisible(true)
	})

	// We put an exit action into the context menu.
	exitAction := walk.NewAction()
	if err := exitAction.SetText("E&xit"); err != nil {
		log.Fatal(err.Error())
	}
	exitAction.Triggered().Attach(func() { walk.App().Exit(0) })
	if err := ni.ContextMenu().Actions().Add(exitAction); err != nil {
		log.Fatal(err.Error())
	}

	// The notify icon is hidden initially, so we have to make it visible.
	if err := ni.SetVisible(true); err != nil {
		log.Fatal(err.Error())
	}

	app.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		*canceled = true
		app.SetVisible(false)
	})

	app.Run()
}

type NotifyIconDisplay struct {
	index int
}

func (n *NotifyIconDisplay) Reset() *NotifyIconDisplay {
	n.index = 0
	return n
}

func (n *NotifyIconDisplay) Next() *NotifyIconDisplay {
	n.index++
	if n.index > 999 {
		n.index = 0
	}

	return n
}

func (n *NotifyIconDisplay) Icon() *walk.Icon {
	ic, err := walk.NewIconFromImageForDPI(makeDigitImage(n.index), 96)
	if err != nil {
		log.Fatal(err.Error())
	}

	return ic
}
