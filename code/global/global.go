package global

import (
	"pure-game-kit/graphics"
)

var ScreenLoading, ScreenMainMenu, ScreenWorld, ScreenBattle int
var ThemesGUI, PopupNarrowGUI, PopupWideGUI, DimGUI, XBtnGUI, TitleGUI string
var TimeScale float32 = 1

var Opts *Options

var BattleTileWidth, BattleTileHeight, BattleColumns, BattleRows float32 = 0, 0, 0, 0

func TryShowFPS(camera *graphics.Camera) {
	if Opts.ShowFPS {
		camera.DrawTextDebug(true, false, false, false)
	}
}
