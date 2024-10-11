package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	desk          desktop.App
	application   fyne.App
	grabberWindow fyne.Window
	systemMenu    fyne.Menu
	outputLabel   *widget.Label
	cleanupTimer  *time.Timer
}

func New() UI {
	return UI{}
}

func (ui UI) Start() {
	ui.grabberWindow.ShowAndRun()
}

func (ui UI) ShowHide() {
	ui.SetText("<Nothing Pressed>")
	if ui.grabberWindow.Content().Visible() {
		ui.grabberWindow.Hide()
	} else {
		ui.grabberWindow.Show()
	}
}
func (ui *UI) SetText(text string) {
	ui.outputLabel.SetText(text)
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

func (ui *UI) Init(onKeyPress func(k *fyne.KeyEvent, owner *UI)) bool {
	ui.application = app.New()
	ui.grabberWindow = ui.application.NewWindow("Android Shortcuts")

	var ok bool
	ui.desk, ok = ui.application.(desktop.App)

	if !ok {
		return false
	}

	ui.systemMenu = *fyne.NewMenu("Android Shortcuts",
		fyne.NewMenuItem("Show", func() {
			ui.grabberWindow.Show()
		}),
	)
	ui.desk.SetSystemTrayMenu(&ui.systemMenu)

	ui.grabberWindow.Resize(fyne.Size{Height: 100, Width: 500})
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
		ui.grabberWindow.Hide()
	})

	ui.grabberWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			ui.grabberWindow.Hide()
		} else {
			onKeyPress(k, ui)
		}
	})
	return true
}
