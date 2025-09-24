package ui

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mrjbee/android-adb-shortcuts/pkg/kbc"
)

type UI struct {
	desk          desktop.App
	application   fyne.App
	grabberWindow fyne.Window
	systemMenu    fyne.Menu
	outputLabel   *widget.Label
	cleanupTimer  *time.Timer
	processingLen int
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

func (ui *UI) Init(config kbc.KeyBindConfig) bool {
	ui.application = app.New()
	ui.grabberWindow = ui.application.NewWindow("Android Shortcuts")
	ui.processingLen = 0

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

	config.ForEach(func(key string, ak kbc.AndroidKey) {

		sc := desktop.CustomShortcut{
			KeyName:  fyne.KeyName(key),
			Modifier: fyne.KeyModifierAlt,
		}

		log.Print("Init key - " + sc.ShortcutName())

		ui.grabberWindow.Canvas().AddShortcut(&sc, func(s fyne.Shortcut) {
			go func() {
				log.Print("Get key - " + key)
				log.Print("Going to send - " + ak.Тitle)
				ui.SetText("<" + ak.Тitle + ">")
			}()
		})
	})

	/*
		ui.grabberWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
			if k.Name == fyne.KeyEscape {
				ui.grabberWindow.Hide()
			} else {
				if ui.processingLen < 3 {
					ui.processingLen += 1
					go func() {
						//	onKeyPress(k, ui)
						ui.processingLen -= 1
					}()
				}
			}
		})
	*/

	return true
}
