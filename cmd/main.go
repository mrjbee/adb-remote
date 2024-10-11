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

type AdbOperation func() bool

type AndroidKey struct {
	action AdbOperation
	title  string
}

func CreateSendKeyOperation(keyEvent string) AdbOperation {
	return func() bool {
		return services.SendKeyEvent(keyEvent)
	}
}

var OPS_MAP = map[string]AndroidKey{
	string(fyne.KeyBackspace): {CreateSendKeyOperation(services.KEYCODE_BACK), "Back"},
	"RightShift":              {CreateSendKeyOperation(services.KEYCODE_HOME), "Home"},
	string(fyne.KeyRight):     {CreateSendKeyOperation(services.KEYCODE_DPAD_RIGHT), "Right"},
	string(fyne.KeyLeft):      {CreateSendKeyOperation(services.KEYCODE_DPAD_LEFT), "Left"},
	string(fyne.KeyUp):        {CreateSendKeyOperation(services.KEYCODE_DPAD_UP), "Up"},
	string(fyne.KeyDown):      {CreateSendKeyOperation(services.KEYCODE_DPAD_DOWN), "Down"},
	string(fyne.KeyReturn):    {CreateSendKeyOperation("66"), "Enter"},
	string(fyne.KeyMinus):     {CreateSendKeyOperation("25"), "Volume Down"},
	string(fyne.KeyEqual):     {CreateSendKeyOperation("24"), "Volume Up"},
	string(fyne.KeyPageDown):  {CreateSendKeyOperation("25"), "Volume Down"},
	string(fyne.KeyPageUp):    {CreateSendKeyOperation("24"), "Volume Up"},
	string(fyne.KeySpace):     {CreateSendKeyOperation("85"), "Play/Pause"},
	string(fyne.KeyA): {func() bool {
		return services.SendCustomShell("input tap 200 400;input tap 200 400")
	}, "Double Tap Left"},
	string(fyne.KeyD): {func() bool {
		return services.SendCustomShell("input tap 2000 500;input tap 2000 500")
	}, "Double Tap Reft"},
	string(fyne.KeyDelete): {func() bool {
		return services.SendKeyEvent("127") && services.SendKeyEvent(services.KEYCODE_HOME)
	}, "Script: Pause & Home"},
}

func redirectKeyPress(k *fyne.KeyEvent, owner *ui.UI) {
	log.Print("Get key - " + k.Name)
	adbOperation, prs := OPS_MAP[string(k.Name)]
	if prs {
		log.Print("Going to send - " + adbOperation.title)
		if adbOperation.action() {
			owner.SetText("<" + adbOperation.title + ">")
		}
	} else {
		log.Printf("No mapping for %s[%d] ", k.Name, k.Physical.ScanCode)
	}
}
