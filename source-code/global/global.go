package global

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/gui"
	"pure-game-kit/gui/field"
	"pure-game-kit/tiled"
)

const Version = "v0.0.3"

var ScreenLoading, ScreenMainMenu, ScreenWorld, ScreenBattle int
var Project *tiled.Project
var ThemesGUI, PopupDimGUI string
var TimeScale float32 = 1

func TogglePopup(hud, currentPopup, popup *gui.GUI) *gui.GUI {
	currentPopup = condition.If(currentPopup == popup, nil, popup)
	hud.SetField("popup-dim", field.Hidden, condition.If(currentPopup != popup, "1", ""))
	return currentPopup
}
