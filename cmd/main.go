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

func CreateSendKeyOperation(keyEvent string) kbc.AdbOperation {
	return func() bool {
		return services.SendKeyEvent(keyEvent)
	}
}

var OPS_MAP = map[string]kbc.AndroidKey{
	string(fyne.KeyBackspace):    {Аction: CreateSendKeyOperation(services.KEYCODE_BACK), Тitle: "Back"},
	"RightShift":                 {Аction: CreateSendKeyOperation(services.KEYCODE_HOME), Тitle: "Home"},
	string(fyne.KeyRight):        {Аction: CreateSendKeyOperation(services.KEYCODE_DPAD_RIGHT), Тitle: "Right"},
	string(fyne.KeyLeft):         {Аction: CreateSendKeyOperation(services.KEYCODE_DPAD_LEFT), Тitle: "Left"},
	string(fyne.KeyUp):           {Аction: CreateSendKeyOperation(services.KEYCODE_DPAD_UP), Тitle: "Up"},
	string(fyne.KeyDown):         {Аction: CreateSendKeyOperation(services.KEYCODE_DPAD_DOWN), Тitle: "Down"},
	string(fyne.KeyReturn):       {Аction: CreateSendKeyOperation("66"), Тitle: "Enter"},
	string(fyne.KeyMinus):        {Аction: CreateSendKeyOperation("25"), Тitle: "Volume Down"},
	string(fyne.KeyEqual):        {Аction: CreateSendKeyOperation("24"), Тitle: "Volume Up"},
	string(fyne.KeyPageDown):     {Аction: CreateSendKeyOperation("25"), Тitle: "Volume Down"},
	string(fyne.KeyPageUp):       {Аction: CreateSendKeyOperation("24"), Тitle: "Volume Up"},
	string(fyne.KeySpace):        {Аction: CreateSendKeyOperation("85"), Тitle: "Play/Pause"},
	string(fyne.KeyLeftBracket):  {Аction: CreateSendKeyOperation("88"), Тitle: "Previous"},
	string(fyne.KeyRightBracket): {Аction: CreateSendKeyOperation("87"), Тitle: "Next"},
	string(fyne.KeyA): {Аction: func() bool {
		return services.SendCustomShell("input tap 100 200& sleep 0.03; input tap 100 200")
	}, Тitle: "Double Tap Left"},
	string(fyne.KeyD): {Аction: func() bool {
		return services.SendCustomShell("input tap 1800 200& sleep 0.03; input tap 1800 200")
	}, Тitle: "Double Tap Reft"},
	string(fyne.KeyDelete): {Аction: func() bool {
		return services.SendKeyEvent("127") && services.SendKeyEvent(services.KEYCODE_HOME)
	}, Тitle: "Script: Pause & Home"},
}

func redirectKeyPress(k *fyne.KeyEvent, owner *ui.UI) {
	log.Print("Get key - " + k.Name)
	adbOperation, prs := OPS_MAP[string(k.Name)]
	if prs {
		log.Print("Going to send - " + adbOperation.Тitle)
		owner.SetText("<" + adbOperation.Тitle + ">")
		if adbOperation.Аction() {
		}
		owner.SetText("<Nothing Pressed>")
	} else {
		log.Printf("No mapping for %s[%d] ", k.Name, k.Physical.ScanCode)
	}
}
