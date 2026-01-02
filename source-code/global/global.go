package global

import (
	"game/source-code/options"
	"pure-game-kit/execution/condition"
	"pure-game-kit/gui"
	"pure-game-kit/gui/field"
	"pure-game-kit/tiled"
)

var ScreenLoading, ScreenMainMenu, ScreenWorld, ScreenBattle int
var Project *tiled.Project
var ThemesGUI, PopupDimGUI string
var TimeScale float32 = 1

var Options options.Options

func TogglePopup(hud, currentPopup, popup *gui.GUI) *gui.GUI {
	currentPopup = condition.If(currentPopup == popup, nil, popup)
	hud.SetField("popup-dim", field.Hidden, condition.If(currentPopup != popup, "1", ""))
	return currentPopup
}
