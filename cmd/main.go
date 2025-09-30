package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"github.com/mrjbee/android-adb-shortcuts/pkg/kbc"
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
	var keyConfig = kbc.NewKeyBindConfig(OPS_MAP)
	if ui.Init(*keyConfig) {
		services.StartRemoteCommandListing("localhost:7272", func(err error) {
			log.Panic("Could not start server " + err.Error())
		}, func(command string) bool {
			switch command {
			case "showHide":
				fyne.DoAndWait(func() {
					ui.ShowHide()
				})
				return true
			}
			return false
		})
		ui.Start()
	} else {
		log.Panic("No desktop found")
	}
}

func CreateSendKeyOperation(keyEvent string) kbc.AdbOperation {
	return func() bool {
		return services.SendKeyEvent(keyEvent)
	}
}

var OPS_MAP = map[string]kbc.AndroidKey{
	string(fyne.KeyBackspace):    {Action: CreateSendKeyOperation(services.KEYCODE_BACK), Title: "Back"},
	"RightShift":                 {Action: CreateSendKeyOperation(services.KEYCODE_HOME), Title: "Home"},
	string(fyne.KeyRight):        {Action: CreateSendKeyOperation(services.KEYCODE_DPAD_RIGHT), Title: "Right"},
	string(fyne.KeyLeft):         {Action: CreateSendKeyOperation(services.KEYCODE_DPAD_LEFT), Title: "Left"},
	string(fyne.KeyUp):           {Action: CreateSendKeyOperation(services.KEYCODE_DPAD_UP), Title: "Up"},
	string(fyne.KeyDown):         {Action: CreateSendKeyOperation(services.KEYCODE_DPAD_DOWN), Title: "Down"},
	string(fyne.KeyL):            {Action: CreateSendKeyOperation("26"), Title: "Lock/Unlock"},
	string(fyne.KeyReturn):       {Action: CreateSendKeyOperation("66"), Title: "Enter"},
	string(fyne.KeyMinus):        {Action: CreateSendKeyOperation("25"), Title: "Volume Down"},
	string(fyne.KeyEqual):        {Action: CreateSendKeyOperation("24"), Title: "Volume Up"},
	string(fyne.KeyPageDown):     {Action: CreateSendKeyOperation("25"), Title: "Volume Down"},
	string(fyne.KeyPageUp):       {Action: CreateSendKeyOperation("24"), Title: "Volume Up"},
	string(fyne.KeySpace):        {Action: CreateSendKeyOperation("85"), Title: "Play/Pause"},
	string(fyne.KeyLeftBracket):  {Action: CreateSendKeyOperation("88"), Title: "Previous"},
	string(fyne.KeyRightBracket): {Action: CreateSendKeyOperation("87"), Title: "Next"},
	string(fyne.KeyA): {Action: func() bool {
		return services.SendCustomShell("input tap 100 200& sleep 0.03; input tap 100 200")
	}, Title: "Double Tap Left"},
	string(fyne.KeyD): {Action: func() bool {
		return services.SendCustomShell("input tap 1800 200& sleep 0.03; input tap 1800 200")
	}, Title: "Double Tap Reft"},
	string(fyne.KeyDelete): {Action: func() bool {
		return services.SendKeyEvent("127") && services.SendKeyEvent(services.KEYCODE_HOME)
	}, Title: "Script: Pause & Home"},
}
