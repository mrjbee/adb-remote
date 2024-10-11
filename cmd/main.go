package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"github.com/mrjbee/android-adb-shortcuts/pkg/services"
	"github.com/mrjbee/android-adb-shortcuts/pkg/ui"
)

func main() {

	isRunning := services.IsRemoteCommandAvaialable("localhost:7272")
	if isRunning {
		if len(os.Args) == 2 {
			commandName := os.Args[1]
			successful := services.SendRemoteCommand("localhost:7272", commandName)
			if successful {
				os.Exit(0)
			} else {
				os.Exit(1)
			}
		} else {
			log.Panic("Another instance is running and no command provided")
		}
	}

	ui := ui.New()
	if ui.Init(redirectKeyPress) {
		services.StartRemoteCommandListing("localhost:7272", func(err error) {
			log.Panic("Could not start server " + err.Error())
		}, func(command string) bool {
			switch command {
			case "showHide":
				ui.ShowHide()
				return true
			}
			return false
		})
		ui.Start()
	} else {
		log.Panic("No desktop found")
	}
}

type AndroidKey struct {
	code  string
	title string
}

var ANDROID_KEY_MAP = map[string]AndroidKey{
	string(fyne.KeyBackspace): {services.KEYCODE_BACK, "Back"},
	"RightShift":              {services.KEYCODE_HOME, "Home"},
	string(fyne.KeyRight):     {services.KEYCODE_DPAD_RIGHT, "Right"},
	string(fyne.KeyLeft):      {services.KEYCODE_DPAD_LEFT, "Left"},
	string(fyne.KeyUp):        {services.KEYCODE_DPAD_UP, "Up"},
	string(fyne.KeyDown):      {services.KEYCODE_DPAD_DOWN, "Down"},
	string(fyne.KeyReturn):    {"66", "Enter"},
	string(fyne.KeyMinus):     {"25", "Volume Down"},
	string(fyne.KeyEqual):     {"24", "Volume Up"},
	string(fyne.KeyPageDown):  {"25", "Volume Down"},
	string(fyne.KeyPageUp):    {"24", "Volume Up"},
	string(fyne.KeySpace):     {"85", "Play/Pause"},
}

func redirectKeyPress(k *fyne.KeyEvent, owner *ui.UI) {
	androidKey, prs := ANDROID_KEY_MAP[string(k.Name)]
	if prs {
		log.Print("Going to send - " + androidKey.title)
		if services.SendKeyEvent(androidKey.code) {
			owner.SetText("<" + androidKey.title + ">")
		}
	} else {
		log.Printf("No mapping for %s[%d] ", k.Name, k.Physical.ScanCode)
	}
}
