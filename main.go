package main

import (
	"game/source-code/global"
	"game/source-code/options"
	"game/source-code/screens/battle"
	"game/source-code/screens/loading"
	"game/source-code/screens/menu"
	"game/source-code/screens/world"
	"pure-game-kit/data/assets"
	"pure-game-kit/data/file"
	"pure-game-kit/data/storage"
	"pure-game-kit/execution/screens"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/tiled"
	"pure-game-kit/window"
)

func main() {
	window.Title = "Battle Brothers Clone"
	loadAndApplyOptions()

	window.KeepOpen()
	global.ScreenLoading = screens.Add(loading.New(), true)

	loading.Show("Loading:\nDefault icons...")
	assets.LoadDefaultAtlasIcons()

	loading.Show("Loading:\nReusable GUI...")
	global.ThemesGUI = file.LoadText("data/gui/reusable-themes.xml")
	global.PopupDimGUI = file.LoadText("data/gui/reusable-popup-dim.xml")
	loading.Show("Loading:\nTiled project...")
	global.Project = tiled.NewProject(assets.LoadTiledProject("data/project.tiled-project"))

	global.ScreenMainMenu = screens.Add(menu.New(), true)
	global.ScreenWorld = screens.Add(world.New("data/worlds/test/map.tmx"), true)
	global.ScreenBattle = screens.Add(battle.New("data/battlegrounds/test/map.tmx"), true)

	screens.Enter(global.ScreenMainMenu, false)
	for window.KeepOpen() {
		if keyboard.IsKeyJustPressed(key.F5) {
			var prevScreen = screens.CurrentId()
			loadAndApplyOptions()
			assets.ReloadAll()
			global.Project = tiled.NewProject(assets.LoadTiledProject("data/project.tiled-project"))
			screens.Reload()
			screens.Enter(prevScreen, false)
		}
	}
}

func loadAndApplyOptions() {
	var opts options.Options
	storage.FromYAML(file.LoadText("data/options.yaml"), &opts)
	global.Options = opts

	window.IsVSynced = opts.Graphics.VSync
	window.FrameRateLimit = byte(opts.Graphics.LimitFPS)
	window.ApplyState(opts.Graphics.WindowState)
	window.MoveToMonitor(opts.Graphics.Monitor)
}
