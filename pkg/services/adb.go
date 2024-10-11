package services

import (
	"fmt"
	"os/exec"
)

//adb shell input keyevent 82

const KEYCODE_DPAD_UP = "19"
const KEYCODE_DPAD_DOWN = "20"
const KEYCODE_DPAD_LEFT = "21"
const KEYCODE_DPAD_RIGHT = "22"
const KEYCODE_DPAD_CENTER = "23"
const KEYCODE_HOME = "3"
const KEYCODE_BACK = "4"

// adb shell "input tap 200 400;input tap 200 400"
func SendKeyEvent(keyEvent string) bool {
	cmd := exec.Command("adb", "shell", "input", "keyevent", keyEvent)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func SendCustomShell(shellCommand string) bool {
	cmd := exec.Command("adb", "shell", shellCommand)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
