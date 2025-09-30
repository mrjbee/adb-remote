package ui

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/mrjbee/android-adb-shortcuts/pkg/kbc"
)

type UI struct {
	desk                 desktop.App
	application          fyne.App
	grabberWindow        fyne.Window
	grabberWindowVisible bool
	systemMenu           fyne.Menu
	outputLabel          *widget.Label
	cleanupTimer         *time.Timer
	processingLen        int
}

func New() UI {
	return UI{}
}

func (ui *UI) Start() {
	ui.grabberWindowVisible = true
	ui.grabberWindow.ShowAndRun()
}

func (ui *UI) Show() {
	ui.grabberWindowVisible = true
	ui.grabberWindow.Show()
}

func (ui *UI) Hide() {
	ui.grabberWindowVisible = false
	ui.grabberWindow.Hide()
}

func (ui *UI) ShowHide() {
	ui.SetText("<Nothing Pressed>")
	if ui.grabberWindowVisible {
		ui.Hide()
	} else {
		ui.Show()
		ui.grabberWindow.CenterOnScreen()
	}
}

func (ui *UI) ResetText() {
	if ui.cleanupTimer != nil {
		ui.cleanupTimer.Reset(1 * time.Second)
	} else {
		ui.cleanupTimer = time.NewTimer(1 * time.Second)
	}
	go func() {
		<-ui.cleanupTimer.C
		ui.outputLabel.SetText("<Nothing Pressed>")
		ui.cleanupTimer = nil
	}()
}

func (ui *UI) SetText(text string) {
	ui.outputLabel.SetText(text)
}

func (ui *UI) Init(kbc kbc.KeyBindConfig) bool {
	ui.application = app.New()
	ui.application = app.NewWithID("tapdo")
	ui.application.SetIcon(theme.FyneLogo())
	ui.grabberWindow = ui.application.NewWindow("TapDO")
	ui.grabberWindow.SetIcon(theme.FyneLogo())
	ui.processingLen = 0

	var ok bool
	ui.desk, ok = ui.application.(desktop.App)

	if !ok {
		return false
	}

	ui.systemMenu = *fyne.NewMenu("Android Shortcuts",
		fyne.NewMenuItem("Show", func() {
			ui.Show()
		}),
	)
	ui.desk.SetSystemTrayMenu(&ui.systemMenu)

	ui.grabberWindow.Resize(fyne.Size{Height: 100, Width: 500})
	ui.grabberWindow.CenterOnScreen()
	ui.outputLabel = widget.NewLabelWithStyle("<Nothing Pressed>",
		fyne.TextAlignCenter,
		fyne.TextStyle{
			Symbol: false,
		})
	ui.outputLabel.TextStyle.Bold = true
	content := container.New(layout.NewVBoxLayout(),
		widget.NewLabel("Press shortcut key ..."),
		ui.outputLabel,
		layout.NewSpacer())
	ui.grabberWindow.SetContent(
		content,
	)
	ui.grabberWindow.SetCloseIntercept(func() {
		ui.Hide()
	})

	ui.grabberWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			ui.Hide()
		} else {
			log.Printf("Get key: %s - scanCode: %d", k.Name, k.Physical.ScanCode)
			var adbOperation, exists = kbc.GetByKeyOk(string(k.Name))
			if !exists {
				return
			}
			if ui.processingLen < 3 {
				ui.processingLen += 1
				go func() {
					log.Print("Going to send - " + adbOperation.Title)
					fyne.DoAndWait(func() {
						ui.SetText("<" + adbOperation.Title + ">")
					})
					if adbOperation.Action() {
						//not sure about this
					}
					fyne.DoAndWait(func() {
						ui.SetText("<Nothing Pressed>")
					})
					ui.processingLen -= 1
				}()
			}
		}
	})
	return true
}
