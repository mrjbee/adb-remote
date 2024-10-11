# How to use
## First setup
- install android adb on PC
- enable wireless debuggin on Android device
- pair android device with PC
- assign short ket to execute `releases/cmd showHide` to any keyboard shortcut (Im using `super`+`enter`)

## Each next session
- connect with `adb connect <android host>:<android port>` (Im afraid port is different each time, quick way to check is to hold `wireless debugging` button in notification tray)
- hit your shortcut (my is `super`+`enter`)
- use:
    - arrow key for navigation
    - plus/minus for volume
    - space for pause/play
    - backspace for `Back`
    - right shift for `Home`
    - ... see source code for more

![Key Capture Mode](/screenshot.png?raw=true "Key awaiting window")
