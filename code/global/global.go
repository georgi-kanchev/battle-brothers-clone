package global

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/tiled"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
	"pure-game-kit/utility/time"
)

var ScreenLoading, ScreenMainMenu, ScreenWorld, ScreenBattle int
var Project *tiled.Project
var ThemesGUI, PopupNarrowGUI, PopupWideGUI, DimGUI, XBtnGUI, TitleGUI string
var TimeScale float32 = 1

var Opts *Options

var BattleTileWidth, BattleTileHeight, BattleTileColumns, BattleTileRows float32 = 0, 0, 0, 0

var fps = ""

func TryShowFPS(camera *graphics.Camera) {
	if !Opts.ShowFPS {
		return
	}

	if condition.TrueEvery(0.1, "fps") {
		fps = text.New(" FPS: ", number.Round(time.FrameRate(), 1))
	}

	var tlx, tly = camera.PointFromEdge(0, 0)
	var height = 20 * Opts.ScaleUI / camera.Zoom
	camera.DrawText("", fps, tlx, tly, height, 0.95, 0, palette.Black)
	camera.DrawText("", fps, tlx, tly, height, 0.5, 0, palette.White)
}
