package global

import (
	"game/code/options"
	"pure-game-kit/tiled"
)

var ScreenLoading, ScreenMainMenu, ScreenWorld, ScreenBattle int
var Project *tiled.Project
var ThemesGUI, PopupNarrowGUI, PopupWideGUI, DimGUI, XBtnGUI, TitleGUI string
var TimeScale float32 = 1

var Options *options.Options

var BattleTileWidth, BattleTileHeight, BattleTileColumns, BattleTileRows float32 = 0, 0, 0, 0
